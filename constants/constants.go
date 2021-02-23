package constants

// Constants
const (
	// Drawing constants
	GridUnitSize  = 5
	GridWidth     = 1400
	GridHeight    = 1000
	GridUnitsWide = 280
	GridUnitsHigh = 200
	ScreenWidth   = 1800
	ScreenHeight  = 1000

	// Statistics constants
	MinFamilyPopulationToLog = 5
	PopulationUpdateInterval = 100

	// Environment constants
	ChanceToAddOrganism = 0.05
	ChanceToAddFoodItem = 0.5
	MaxFoodValue        = 100
	MinFoodValue        = 2

	// Organism constants
	MaxCyclesBetweenSpawns          = 100
	MinSpawnHealth                  = 1.0
	MaxSpawnHealthPercent           = 0.5
	MinChanceToMutateDecisionTree   = 0.01
	MaxChanceToMutateDecisionTree   = 1.00
	MaxCyclesToEvaluateDecisionTree = 100
	MaxOrganisms                    = 20000
	GrowthFactor                    = 0.5
	MaximumMaxSize                  = 100.0
	MinimumMaxSize                  = 10.0

	// Health change constants (as a percent of an organism's size)
	HealthChangePerCycle          = -0.001
	HealthChangeFromBeingIdle     = +0.003
	HealthChangeFromTurning       = -0.001
	HealthChangeFromMoving        = -0.03
	HealthChangeFromEatingAttempt = -0.01
	HealthChangeFromAttacking     = -0.05
	HealthChangeInflictedByAttack = -0.5
	HealthChangeFromFeeding       = -0.005

	// Decision Tree constants
	HealthPercentToChangeDecisionTree = 0.10
)
