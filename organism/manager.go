package organism

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	c "github.com/Zebbeni/protozoa/constants"
	d "github.com/Zebbeni/protozoa/decisions"
	"github.com/Zebbeni/protozoa/food"
	"github.com/Zebbeni/protozoa/utils"
)

var descendantsPrintThreshold = 10

// Manager contains 2D array of booleans showing if organism present
type Manager struct {
	worldAPI WorldAPI

	organisms             map[int]*Organism
	organismIDGrid        [][]int
	totalOrganismsCreated int

	organismUpdateOrder []int
	newOrganismIDs      []int

	MostReproductiveAllTime  *organismInfo
	MostReproductiveCurrent  *organismInfo
	AncestorDescendantsCount map[int]int

	UpdateDuration, ResolveDuration time.Duration
}

type organismInfo struct {
	id           int
	size         float64
	health       float64
	ancestorID   int
	age          int
	children     int
	decisionTree string
	traits       *Traits
}

func (o *organismInfo) ID() int {
	return o.id
}

// NewManager creates all Organisms and updates grid
func NewManager(api WorldAPI) *Manager {
	grid := initializeGrid()
	organisms := make(map[int]*Organism)
	manager := &Manager{
		worldAPI:                 api,
		organismIDGrid:           grid,
		organisms:                organisms,
		organismUpdateOrder:      make([]int, 0, c.MaxOrganisms),
		newOrganismIDs:           make([]int, 0, 100),
		AncestorDescendantsCount: make(map[int]int),
		MostReproductiveAllTime:  &organismInfo{traits: &Traits{}},
		MostReproductiveCurrent:  &organismInfo{traits: &Traits{}},
	}
	return manager
}

// Update walks through decision tree of each organism and applies the
// chosen action to the organism, the grid, and the environment
func (m *Manager) Update() {
	m.MostReproductiveCurrent = &organismInfo{traits: &Traits{}}
	// Periodically add new random organisms if population below a certain amount
	if len(m.organisms) < c.MaxOrganisms && rand.Float64() < c.ChanceToAddOrganism {
		m.SpawnRandomOrganism()
	}
	// FUTURE: do this multi-threaded
	start := time.Now()
	for _, id := range m.organismUpdateOrder {
		m.updateOrganism(m.organisms[id])
	}
	m.UpdateDuration = time.Since(start)
	start = time.Now()
	for _, id := range m.organismUpdateOrder {
		m.resolveOrganismAction(m.organisms[id])
	}
	m.ResolveDuration = time.Since(start)
	m.updateOrganismOrder()
}

// updateOrganismOrder creates a new ordered list of all organismIDs that are
// alive after the current cycle, appending any newly spawned organisms.
// This means iterating the full list of organisms again, but this should be
// faster than just deleting the dead IDs and shifting all others to the left
func (m *Manager) updateOrganismOrder() {
	orderedIDs := append(m.organismUpdateOrder, m.newOrganismIDs...)
	organismUpdateOrder := make([]int, 0, len(orderedIDs))
	for _, id := range orderedIDs {
		if _, ok := m.organisms[id]; ok {
			organismUpdateOrder = append(organismUpdateOrder, id)
		}
	}
	m.organismUpdateOrder = organismUpdateOrder
	m.newOrganismIDs = make([]int, 0, 100)
}

func initializeGrid() [][]int {
	grid := make([][]int, c.GridWidth)
	for r := 0; r < c.GridWidth; r++ {
		grid[r] = make([]int, c.GridHeight)
	}
	for x := 0; x < c.GridWidth; x++ {
		for y := 0; y < c.GridHeight; y++ {
			grid[x][y] = -1
		}
	}
	return grid
}

func (m *Manager) updateOrganism(o *Organism) {
	o.UpdateStats()
	o.UpdateAction()
}

func (m *Manager) resolveOrganismAction(o *Organism) {
	if o == nil {
		return
	}
	m.updateHealth(o)
	m.applyAction(o)
	m.evaluateBest(o)
	m.removeIfDead(o)
}

func (m *Manager) evaluateBest(o *Organism) {
	if o.Children > m.MostReproductiveCurrent.children {
		decisionTree := o.getBestDecisionTree()
		if decisionTree == nil {
			decisionTree = o.decisionTree
		}
		organismInfo := &organismInfo{
			id:           o.ID,
			size:         o.Size,
			health:       o.Health,
			ancestorID:   o.OriginalAncestorID,
			decisionTree: decisionTree.Print("", true, false),
			age:          o.Age,
			children:     o.Children,
			traits:       o.traits.copy(),
		}
		m.MostReproductiveCurrent = organismInfo

		if o.Children > m.MostReproductiveAllTime.children {
			m.MostReproductiveAllTime = organismInfo
		}
	}
}

// SpawnRandomOrganism creates an Organism with random position.
//
// Checks random positions on the grid until it finds an empty one. Calls
// NewOrganism to initialize decision tree, other random attributes.
func (m *Manager) SpawnRandomOrganism() {
	if spawnPoint, found := m.getRandomSpawnLocation(); found {
		index := m.totalOrganismsCreated
		organism := NewRandom(index, spawnPoint, m.worldAPI)
		m.registerNewOrganism(organism, index)
	}
}

// SpawnChildOrganism creates a new organism near an existing 'parent' organism
// with a copy of its parent's node library. (No organism created if no room)
// Returns true / false depending on whether a child was actually spawned.
func (m *Manager) SpawnChildOrganism(parent *Organism) bool {
	if spawnPoint, found := m.getChildSpawnLocation(parent); found {
		index := m.totalOrganismsCreated
		organism := NewChild(parent, index, spawnPoint, m.worldAPI)
		m.registerNewOrganism(organism, index)
		return true
	}
	return false
}

func (m *Manager) registerNewOrganism(o *Organism, index int) {
	m.organisms[index] = o
	m.totalOrganismsCreated++
	m.organismIDGrid[o.X()][o.Y()] = index
	m.newOrganismIDs = append(m.newOrganismIDs, index)

	// update ancestors
	ancestorID := o.OriginalAncestorID
	if ancestorID != o.ID {
		if _, ok := m.AncestorDescendantsCount[ancestorID]; !ok {
			m.AncestorDescendantsCount[ancestorID] = 0
		}
		m.AncestorDescendantsCount[ancestorID]++
	}
}

// returns a random point and whether it is empty
func (m *Manager) getRandomSpawnLocation() (utils.Point, bool) {
	point := utils.GetRandomPoint()
	return point, m.isGridLocationEmpty(point)
}

func (m *Manager) getChildSpawnLocation(parent *Organism) (utils.Point, bool) {
	direction := utils.GetRandomDirection()
	point := parent.Location.Add(direction)
	for i := 0; i < 4; i++ {
		if m.isGridLocationEmpty(point) {
			return point, true
		}
		direction = direction.Left()
		point = parent.Location.Add(direction)
	}
	return point, false
}

func (m *Manager) isGridLocationEmpty(point utils.Point) bool {
	return !m.isFoodAtLocation(point) && !m.isOrganismAtLocation(point)
}

func (m *Manager) isFoodAtLocation(point utils.Point) bool {
	return m.worldAPI.CheckFoodAtPoint(point, func(item *food.Item) bool {
		return item != nil
	})
}

func (m *Manager) isOrganismAtLocation(point utils.Point) bool {
	return m.organismIDGrid[point.X][point.Y] != -1
}

func (m *Manager) getOrganismAt(point utils.Point) *Organism {
	if id, exists := m.getOrganismIDAt(point); exists {
		index := id
		return m.organisms[index]
	}
	return nil
}

func (m *Manager) getOrganismIDAt(point utils.Point) (int, bool) {
	id := m.organismIDGrid[point.X][point.Y]
	if id != -1 {
		return id, true
	}
	return -1, false
}

// CheckOrganismAtPoint returns the result of running a check against any Organism
// found at a given Point.
func (m *Manager) CheckOrganismAtPoint(point utils.Point, checkFunc OrgCheck) bool {
	return checkFunc(m.getOrganismAt(point))
}

// OrganismCount returns the current number of organisms alive in the simulation
func (m *Manager) OrganismCount() int {
	return len(m.organisms)
}

func (m *Manager) applyAction(o *Organism) {
	switch o.Action() {
	case d.ActIdle:
		m.applyIdle(o)
		break
	case d.ActAttack:
		m.applyAttack(o)
		break
	case d.ActFeed:
		m.applyFeed(o)
		break
	case d.ActEat:
		m.applyEat(o)
		break
	case d.ActMove:
		m.applyMove(o)
		break
	case d.ActTurnLeft:
		m.applyLeftTurn(o)
		break
	case d.ActTurnRight:
		m.applyRightTurn(o)
		break
	case d.ActSpawn:
		m.applySpawn(o)
		break
	}
}

func (m *Manager) updateHealth(o *Organism) {
	o.ApplyHealthChange(c.HealthChangePerCycle * o.Size)
}

func (m *Manager) applyIdle(o *Organism) {
	o.ApplyHealthChange(c.HealthChangeFromBeingIdle * o.Size)
}

func (m *Manager) applyAttack(o *Organism) {
	o.ApplyHealthChange(c.HealthChangeFromAttacking * o.Size)
	targetPoint := o.Location.Add(o.Direction)
	if m.isOrganismAtLocation(targetPoint) {
		targetOrganismIndex := m.organismIDGrid[targetPoint.X][targetPoint.Y]
		targetOrganism := m.organisms[targetOrganismIndex]
		targetOrganism.ApplyHealthChange(c.HealthChangeInflictedByAttack * o.Size)
		m.removeIfDead(targetOrganism)
	}
}

func (m *Manager) removeIfDead(o *Organism) bool {
	if o.Health > 0.0 {
		return false
	}

	m.organismIDGrid[o.Location.X][o.Location.Y] = -1
	m.worldAPI.AddFoodAtPoint(o.Location, int(o.Size))
	delete(m.organisms, o.ID)
	return true
}

func (m *Manager) applySpawn(o *Organism) {
	if success := m.SpawnChildOrganism(o); success {
		o.ApplyHealthChange(o.HealthCostToReproduce())
		o.CyclesSinceLastSpawn = 0
		o.Children++
	}
}

func (m *Manager) applyFeed(o *Organism) {
	o.ApplyHealthChange(c.HealthChangeFromFeeding * o.Size)

	amountToFeed := c.HealthChangeFromFeeding * o.Size
	targetPoint := o.Location.Add(o.Direction)
	if m.isOrganismAtLocation(targetPoint) {
		targetOrganismIndex := m.organismIDGrid[targetPoint.X][targetPoint.Y]
		targetOrganism := m.organisms[targetOrganismIndex]
		targetOrganism.ApplyHealthChange(amountToFeed)
	} else {
		m.worldAPI.AddFoodAtPoint(targetPoint, int(amountToFeed))
	}
}

func (m *Manager) applyEat(o *Organism) {
	o.ApplyHealthChange(c.HealthChangeFromEatingAttempt * o.Size)

	targetPoint := o.Location.Add(o.Direction)
	if value, exists := m.worldAPI.GetFoodAtPoint(targetPoint); exists {
		maxCanEat := o.Size
		amountToEat := math.Min(float64(value), maxCanEat)
		m.worldAPI.RemoveFoodAtPoint(targetPoint, int(amountToEat))
		o.ApplyHealthChange(amountToEat)
	}
}

func (m *Manager) applyMove(o *Organism) {
	o.ApplyHealthChange(c.HealthChangeFromMoving * o.Size)

	targetPoint := o.Location.Add(o.Direction)
	if m.isGridLocationEmpty(targetPoint) {
		m.organismIDGrid[o.Location.X][o.Location.Y] = -1
		o.Location = targetPoint
		m.organismIDGrid[o.Location.X][o.Location.Y] = o.ID
	}
}

func (m *Manager) applyRightTurn(o *Organism) {
	o.ApplyHealthChange(c.HealthChangeFromTurning * o.Size)

	o.Direction = o.Direction.Right()
}

func (m *Manager) applyLeftTurn(o *Organism) {
	o.ApplyHealthChange(c.HealthChangeFromTurning * o.Size)

	o.Direction = o.Direction.Left()
}

// GetOrganisms returns an array of all Organisms from organism manager
func (m *Manager) GetOrganisms() map[int]*Organism {
	return m.organisms
}

// PrintBest prints the highest current score of any Organism (and their index)
func (m *Manager) PrintBest() {
	m.printBestAncestors()
	fmt.Print("\n\n")
	m.printBestCurrent()
	fmt.Print("\n\n")
	m.printBestAllTime()
}

func (m *Manager) printBestCurrent() {
	fmt.Printf("\n  - Best Organism Current - \n")
	m.printOrganismInfo(m.MostReproductiveCurrent)
}

func (m *Manager) printBestAllTime() {
	fmt.Printf("\n  - Best Organism All Time - \n")
	m.printOrganismInfo(m.MostReproductiveAllTime)
}

func (m *Manager) printOrganismInfo(info *organismInfo) {
	fmt.Printf("\n      ID: %10d   |         InitialHealth: %4d", info.id, int(info.traits.spawnHealth))
	fmt.Printf("\n     Age: %10d   |      MinHealthToSpawn: %4d", info.age, int(info.traits.minHealthToSpawn))
	fmt.Printf("\nChildren: %10d   |      MinCyclesToSpawn: %4d", info.children, info.traits.minCyclesBetweenSpawns)
	fmt.Printf("\nAncestor: %10d   |  CyclesToEvaluateTree: %4d", info.ancestorID, info.traits.cyclesToEvaluateDecisionTree)
	fmt.Printf("\n  Health: %10.2f   |   ChanceToMutateTree:  %4.2f", info.health, info.traits.chanceToMutateDecisionTree)
	fmt.Printf("\n    Size: %10.2f   |              MaxSize:  %4.2f", info.size, info.traits.maxSize)
	fmt.Printf("\n  DecisionTree:\n%s", info.decisionTree)
}

func (m *Manager) printBestAncestors() {
	fmt.Printf("\n - Original Ancestors: %d\n", len(m.AncestorDescendantsCount))
	fmt.Printf("   Best (%d descendants or more) -\n", descendantsPrintThreshold)
	fmt.Print("  Ancestor ID  | Descendants\n")

	// updateThreshold := false
	for ancestorID, descendants := range m.AncestorDescendantsCount {
		if descendants >= descendantsPrintThreshold {
			fmt.Printf("\n%13d  |%12d", ancestorID, descendants)
			// if descendants > descendantsPrintThreshold*2 {
			// 	updateThreshold = true
			// }
		}
	}
	// if updateThreshold {
	// 	descendantsPrintThreshold = int(math.Ceil(float64(descendantsPrintThreshold) * 1.1))
	// }
}
