package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	scale = 5
	w     = 1024
	h     = 768
	halfW = w / 2
	halfH = h / 2
	halfS = scale / 2
	maxX  = w/scale - 1
	maxY  = h/scale - 1
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

func (g *Grid) CheckAlive(x, y int) bool {
	return g.Rows[y].Cells[x].Alive
}

func (g *Grid) NextState() {
	nextState := make([]*Row, 0)
	for y := 0; y < h/scale; y++ {
		row := NewRow()
		for x := 0; x < w/scale; x++ {
			cell := NewCell(x, y)
			isAlive := g.Rows[y].Cells[x].Alive
			cell.Alive = isAlive
			aliveNeighbor := g.CountAliveNeighbor(x, y)
			if isAlive && (aliveNeighbor < 2 || aliveNeighbor > 3) {
				cell.Alive = false
			} else if isAlive || aliveNeighbor == 3 {
				cell.Alive = true
			}
			row.AddCell(cell)
		}
		nextState = append(nextState, row)
	}
	g.Rows = nextState
}

func (g *Grid) CountAliveNeighbor(x, y int) int {
	count := 0
	for j := -1; j < 2; j++ {
		for i := -1; i < 2; i++ {
			if i == 0 && j == 0 {
				continue
			}
			nX := x + i
			if nX < 0 {
				nX = maxX
			} else if nX > maxX {
				nX = 0
			}
			nY := y + j
			if nY < 0 {
				nY = maxY
			} else if nY > maxY {
				nY = 0
			}
			if g.CheckAlive(nX, nY) {
				count++
			}
		}
	}
	return count
}

func NewGrid() *Grid {
	rows := make([]*Row, 0)
	for y := 0; y < h/scale; y++ {
		row := NewRow()
		for x := 0; x < w/scale; x++ {
			cell := NewCell(x, y)
			if rnd.Float64() < 0.5 {
				cell.Alive = true
			}
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
	rnd  *rand.Rand
)

func Init() {
	seed := rand.NewSource(1234123)
	rnd = rand.New(seed)
	game = Game{Grid: NewGrid()}
	img = ebiten.NewImage(scale-2, scale-2)
	img.Fill(color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
}

func (g *Game) Update() error {
	mx, my := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		rx := mx / scale
		ry := my / scale
		if rx > 0 && rx < maxX && ry > 0 && ry < maxY {
			g.Grid.Rows[ry].Cells[rx].Alive = true
		}
	} else {
		g.Grid.NextState()
	}
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
