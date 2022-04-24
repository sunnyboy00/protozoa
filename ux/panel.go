package ux

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"

	r "github.com/Zebbeni/protozoa/resources"
	s "github.com/Zebbeni/protozoa/simulation"
)

const (
	padding     = 15
	panelWidth  = 400
	panelHeight = 1000

	titleXOffset = padding
	titleYOffset = padding
	playXOffset  = padding
	playYOffset  = padding

	statsXOffset = padding
	statsYOffset = 69

	selectedXOffset = padding
	selectedYOffset = 300

	graphXOffset = padding
	graphYOffset = 130
	graphWidth   = 370
	graphHeight  = 120
)

type Panel struct {
	simulation         *s.Simulation
	previousPanelImage *ebiten.Image
	graph              *Graph
}

func NewPanel(sim *s.Simulation) *Panel {
	return &Panel{
		simulation: sim,
		graph:      NewGraph(sim),
	}
}

func (p *Panel) Render() *ebiten.Image {
	panelImage := ebiten.NewImage(panelWidth, panelHeight)

	if p.shouldRefresh() {
		p.renderDividingLine(panelImage)
		p.renderTitle(panelImage)
		p.renderPlayPauseText(panelImage)
		p.renderStats(panelImage)
		p.renderGraph(panelImage)
		p.renderSelected(panelImage)

		p.previousPanelImage = ebiten.NewImage(panelWidth, panelHeight)
		p.previousPanelImage.DrawImage(panelImage, nil)
	} else {
		panelImage.DrawImage(p.previousPanelImage, nil)
	}

	return panelImage
}

func (p *Panel) shouldRefresh() bool {
	return true
}

func (p *Panel) renderDividingLine(panelImage *ebiten.Image) {
	ebitenutil.DrawRect(panelImage, float64(panelWidth)-1, 0, float64(panelWidth), float64(panelHeight), color.White)
}

func (p *Panel) renderTitle(panelImage *ebiten.Image) {
	bounds := text.BoundString(r.FontInversionz40, "protozoa")
	text.Draw(panelImage, "protozoa", r.FontInversionz40, titleXOffset, titleYOffset+bounds.Dy(), color.White)
}

func (p *Panel) renderPlayPauseText(panelImage *ebiten.Image) {
	message := "[Space] to Pause"
	if p.simulation.IsPaused() {
		message = "[Space] to Resume"
	}

	bounds := text.BoundString(r.FontSourceCodePro12, message)
	xOffset := panelWidth - playXOffset - bounds.Dx()
	text.Draw(panelImage, message, r.FontSourceCodePro12, xOffset, playYOffset+bounds.Dy(), color.White)
}

func (p *Panel) renderStats(panelImage *ebiten.Image) {
	statsString := fmt.Sprintf("CYCLE: %9d\nORGANISMS: %5d\nDEAD: %10d",
		p.simulation.Cycle(), p.simulation.OrganismCount(), p.simulation.GetDeadCount())
	text.Draw(panelImage, statsString, r.FontSourceCodePro12, statsXOffset, statsYOffset, color.White)
}

func (p *Panel) renderGraph(panelImage *ebiten.Image) {
	text.Draw(panelImage, "HISTORY", r.FontSourceCodePro12, graphXOffset, graphYOffset, color.White)
	graphImage := p.graph.Render()
	graphOptions := &ebiten.DrawImageOptions{}
	scaleX := float64(graphWidth) / float64(graphImage.Bounds().Dx())
	scaleY := float64(graphHeight) / float64(graphImage.Bounds().Dy())
	graphOptions.GeoM.Scale(scaleX, scaleY)
	graphOptions.GeoM.Translate(graphXOffset, graphYOffset+10)

	panelImage.DrawImage(graphImage, graphOptions)

	// draw border around graph
	left, top, right, bottom := float64(graphXOffset), float64(graphYOffset+10), float64(graphXOffset+graphWidth), float64(graphYOffset+graphHeight+10)
	ebitenutil.DrawLine(panelImage, left, top, right, top, color.White)
	ebitenutil.DrawLine(panelImage, right, top, right, bottom, color.White)
	ebitenutil.DrawLine(panelImage, left, bottom, right, bottom, color.White)
	ebitenutil.DrawLine(panelImage, left, top, left, bottom, color.White)
}

func (p *Panel) renderSelected(panelImage *ebiten.Image) {
	id := p.simulation.GetSelected()
	info := p.simulation.GetOrganismInfoByID(id)
	traits, found := p.simulation.GetOrganismTraitsByID(id)
	decisionTree := p.simulation.GetOrganismDecisionTreeByID(id)
	if info == nil || decisionTree == nil || found == false {
		return
	}

	decisionTreeString := fmt.Sprintf("DECISION TREE:\n%s", decisionTree.Print())

	infoString := fmt.Sprintf("ORGANISM ID: %7d\n"+
		"ANCESTOR ID: %7d        SIZE:      %5.2f\n"+
		"AGE:         %7d        CHILDREN: %7d\n"+
		"MUTATE CHANCE:   %3.0f        SPAWN TIME: %5d\n",
		info.ID, info.AncestorID, info.Size, info.Age, info.Children, traits.ChanceToMutateDecisionTree*100.0, traits.MinCyclesBetweenSpawns)
	bounds := text.BoundString(r.FontSourceCodePro12, infoString)
	offsetY := selectedYOffset + bounds.Dy() + padding

	text.Draw(panelImage, infoString, r.FontSourceCodePro12, selectedXOffset, selectedYOffset, color.White)
	text.Draw(panelImage, decisionTreeString, r.FontSourceCodePro10, selectedXOffset, offsetY, color.White)
}