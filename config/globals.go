package config

import (
	"encoding/json"
	"io"
	"os"
)

var defaultFilePath = "settings/default.json"
var constants *Globals

// SetGlobals allows a one-time initialization of all globally-referenced constants
func SetGlobals(g *Globals) {
	if constants != nil {
		return
	}
	constants = g
}

func GridUnitSize() int                        { return constants.GridUnitSize }
func GridWidth() int                           { return constants.GridWidth }
func GridHeight() int                          { return constants.GridHeight }
func GridUnitsWide() int                       { return constants.GridUnitsWide }
func GridUnitsHigh() int                       { return constants.GridUnitsHigh }
func ScreenWidth() int                         { return constants.ScreenWidth }
func ScreenHeight() int                        { return constants.ScreenHeight }
func PopulationUpdateInterval() int            { return constants.PopulationUpdateInterval }
func InitialOrganisms() int                    { return constants.InitialOrganisms }
func InitialFood() int                         { return constants.InitialFood }
func ChanceToAddOrganism() float64             { return constants.ChanceToAddOrganism }
func ChanceToAddFoodItem() float64             { return constants.ChanceToAddFoodItem }
func MaxFoodValue() int                        { return constants.MaxFoodValue }
func MinFoodValue() int                        { return constants.MinFoodValue }
func MinPh() float64                           { return constants.MinPh }
func MaxPh() float64                           { return constants.MaxPh }
func InitialPh() float64                       { return constants.InitialPh }
func MaxCyclesBetweenSpawns() int              { return constants.MaxCyclesBetweenSpawns }
func MinSpawnHealth() float64                  { return constants.MinSpawnHealth }
func MaxSpawnHealthPercent() float64           { return constants.MaxSpawnHealthPercent }
func MinChanceToMutateDecisionTree() float64   { return constants.MinChanceToMutateDecisionTree }
func MaxChanceToMutateDecisionTree() float64   { return constants.MaxChanceToMutateDecisionTree }
func MaxOrganisms() int                        { return constants.MaxOrganisms }
func GrowthFactor() float64                    { return constants.GrowthFactor }
func MaximumMaxSize() float64                  { return constants.MaximumMaxSize }
func MinimumMaxSize() float64                  { return constants.MinimumMaxSize }
func MinPhTolerance() float64                  { return constants.MinPhTolerance }
func MaxPhTolerance() float64                  { return constants.MaxPhTolerance }
func MinPhToleranceRange() float64             { return constants.MinPhToleranceRange }
func MaxPhToleranceRange() float64             { return constants.MaxPhToleranceRange }
func MinChangeToPh() float64                   { return constants.MinChangeToPh }
func MaxChangeToPh() float64                   { return constants.MaxChangeToPh }
func PhIncrementToDisplay() float64            { return constants.PhIncrementToDisplay }
func PhDiffuseFactor() float64                 { return constants.PhDiffuseFactor }
func HealthChangePerCycle() float64            { return constants.HealthChangePerCycle }
func HealthChangeFromBeingIdle() float64       { return constants.HealthChangeFromBeingIdle }
func HealthChangeFromTurning() float64         { return constants.HealthChangeFromTurning }
func HealthChangeFromMoving() float64          { return constants.HealthChangeFromMoving }
func HealthChangeFromEatingAttempt() float64   { return constants.HealthChangeFromEatingAttempt }
func HealthChangeFromAttacking() float64       { return constants.HealthChangeFromAttacking }
func HealthChangeInflictedByAttack() float64   { return constants.HealthChangeInflictedByAttack }
func HealthChangeFromFeeding() float64         { return constants.HealthChangeFromFeeding }
func HealthChangePerDecisionTreeNode() float64 { return constants.HealthChangePerDecisionTreeNode }
func HealthChangeFromUnhealthyPh() float64     { return constants.HealthChangeFromUnhealthyPh }
func MaxDecisionTreeSize() int                 { return constants.MaxDecisionTreeSize }

type Globals struct {
	// Drawing parameters
	GridUnitSize  int `json:"grid_unit_size"`
	GridWidth     int `json:"grid_width"`
	GridHeight    int `json:"grid_height"`
	GridUnitsWide int `json:"grid_units_wide"`
	GridUnitsHigh int `json:"grid_units_high"`
	ScreenWidth   int `json:"screen_width"`
	ScreenHeight  int `json:"screen_height"`

	// Statistics parameters
	PopulationUpdateInterval int `json:"population_update_interval"`

	// Environment parameters
	InitialOrganisms    int     `json:"initial_organisms"`
	InitialFood         int     `json:"initial_food"`
	ChanceToAddOrganism float64 `json:"chance_to_add_organism"`
	ChanceToAddFoodItem float64 `json:"chance_to_add_food_item"`
	MaxFoodValue        int     `json:"max_food_value"`
	MinFoodValue        int     `json:"min_food_value"`
	MinPh               float64 `json:"min_ph"`
	MaxPh               float64 `json:"max_ph"`
	InitialPh           float64 `json:"initial_ph"`

	// Organism parameters
	MaxCyclesBetweenSpawns        int     `json:"max_cycles_between_spawns"`
	MinSpawnHealth                float64 `json:"min_spawn_health"`
	MaxSpawnHealthPercent         float64 `json:"max_spawn_health_percent"`
	MaxOrganisms                  int     `json:"max_organisms"`
	GrowthFactor                  float64 `json:"growth_factor"`
	MaximumMaxSize                float64 `json:"maximum_max_size"`
	MinimumMaxSize                float64 `json:"minimum_max_size"`
	MinChanceToMutateDecisionTree float64 `json:"min_chance_to_mutate_decision_tree"`
	MaxChanceToMutateDecisionTree float64 `json:"max_chance_to_mutate_decision_tree"`
	MaxDecisionTreeSize           int     `json:"max_decision_tree_size"`
	MinPhTolerance                float64 `json:"min_ph_tolerance"`
	MaxPhTolerance                float64 `json:"max_ph_tolerance"`
	MinPhToleranceRange           float64 `json:"min_ph_tolerance_range"`
	MaxPhToleranceRange           float64 `json:"max_ph_tolerance_range"`
	MinChangeToPh                 float64 `json:"min_change_to_ph"`
	MaxChangeToPh                 float64 `json:"max_change_to_ph"`
	PhIncrementToDisplay          float64 `json:"ph_increment_to_display"`
	PhDiffuseFactor               float64 `json:"ph_diffuse_factor"`

	// Health parameters (percent of organism size)
	HealthChangePerCycle            float64 `json:"health_change_per_cycle"`
	HealthChangeFromBeingIdle       float64 `json:"health_change_from_being_idle"`
	HealthChangeFromTurning         float64 `json:"health_change_from_turning"`
	HealthChangeFromMoving          float64 `json:"health_change_from_moving"`
	HealthChangeFromEatingAttempt   float64 `json:"health_change_from_eating_attempt"`
	HealthChangeFromAttacking       float64 `json:"health_change_from_attacking"`
	HealthChangeInflictedByAttack   float64 `json:"health_change_inflicted_by_attack"`
	HealthChangeFromFeeding         float64 `json:"health_change_from_feeding"`
	HealthChangePerDecisionTreeNode float64 `json:"health_change_per_decision_tree_node"`
	HealthChangeFromUnhealthyPh     float64 `json:"health_change_from_unhealthy_ph"`
}

func LoadFile(filePath string) io.Reader {
	file, err := os.Open(filePath)
	if err != nil {
		panic("failed to read config file")
	}
	return file
}

func GetDefaultGlobals() Globals {
	defaultFile := LoadFile(defaultFilePath)
	g := applyGlobalsFromJson(defaultFile, Globals{})
	return *g
}

func LoadGlobals(file io.Reader) *Globals {
	defaults := GetDefaultGlobals()
	g := applyGlobalsFromJson(file, defaults)
	return g
}

func applyGlobalsFromJson(file io.Reader, globals Globals) *Globals {
	g := globals
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&g)
	if err != nil {
		panic("failed to read globals from file")
	}
	return &g
}

func DumpGlobals(g *Globals, file io.Writer) {
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		panic("failed to convert globals to json")
	}

	_, err = file.Write(data)
	if err != nil {
		panic("failed to write globals to file")
	}
}
