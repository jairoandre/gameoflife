package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	scale = 10
	w     = 1024
	h     = 768
	halfW = w / 2
	halfH = h / 2
	halfS = scale / 2
)

type Cell struct {
	X     int
	Y     int
	Alive bool
}

func NewCell(x, y int) *Cell {
	return &Cell{
		X:     x,
		Y:     y,
		Alive: false,
	}
}

func (c *Cell) Draw(screen *ebiten.Image) {
	if !c.Alive {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.X)*scale, float64(c.Y)*scale)
	//op.GeoM.Translate(halfW-halfS, halfH-halfS)
	screen.DrawImage(img, op)
}

type Row struct {
	Cells []*Cell
}

func NewRow() *Row {
	cells := make([]*Cell, 0)
	return &Row{
		Cells: cells,
	}
}

func (r *Row) AddCell(cell *Cell) {
	r.Cells = append(r.Cells, cell)
}

type Grid struct {
	Rows []*Row
}

func NewGrid() *Grid {
	rows := make([]*Row, 0)
	for y := 0; y < h/scale; y++ {
		row := NewRow()
		for x := 0; x < w/scale; x++ {
			cell := NewCell(x, y)
			cell.Alive = true
			row.AddCell(cell)
		}
		rows = append(rows, row)
	}
	return &Grid{
		Rows: rows,
	}
}

type Game struct {
	Grid *Grid
}

var (
	img  *ebiten.Image
	game Game
)

func Init() {
	game = Game{Grid: NewGrid()}
	img = ebiten.NewImage(scale-2, scale-2)
	img.Fill(color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, row := range g.Grid.Rows {
		for _, cell := range row.Cells {
			cell.Draw(screen)
		}
	}
	msg := fmt.Sprintf(`TPS: %0.2f FPS: %0.2f`, ebiten.CurrentTPS(), ebiten.CurrentFPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return w, h
}

func main() {
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Spiral")
	Init()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
