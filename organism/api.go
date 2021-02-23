package organism

import (
	"github.com/Zebbeni/protozoa/food"
	"github.com/Zebbeni/protozoa/utils"
)

// FoodCheck is a true/false test to run on a given food Item
type FoodCheck func(item *food.Item) bool

// OrgCheck is a true/false test to run on a given food Organism
type OrgCheck func(item *Organism) bool

// LookupAPI provides functions to look up items and organisms
type LookupAPI interface {
	CheckFoodAtPoint(point utils.Point, checkFunc FoodCheck) bool
	CheckOrganismAtPoint(point utils.Point, checkFunc OrgCheck) bool
	GetFoodAtPoint(point utils.Point) *food.Item
	OrganismCount() int
	Cycle() int
}

// ChangeAPI provides callback functions to make changes to the simulation
type ChangeAPI interface {
	// AddFoodAtPoint requests adding some amount of food at a Point
	// returns how much food was actually added
	AddFoodAtPoint(point utils.Point, value int) int
	// RemoveFoodAtPoint requests removing some amount of food at a Point
	// returns how much food was actually removed
	RemoveFoodAtPoint(point utils.Point, value int) int
	// AddGridPointToUpdate indicates a point on the grid has been updated
	// and needs to be re-rendered
	AddUpdatedGridPoint(point utils.Point)
}

// API provides functions needed to lookup and make changes to world objects
type API interface {
	LookupAPI
	ChangeAPI
}
