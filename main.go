package main

import (
	"errors"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("ToH")

	err := ebiten.RunGame(&game{s: &menu{}})
	if err != nil {
		if errors.Is(err, errQuit) {
			return
		}
		panic(err)
	}
}

type scene interface {
	ebiten.Game

	UpdateGame(g *game) error
}

type game struct {
	s scene
}

func (g *game) Update() error {
	err := g.s.Update()
	if err != nil {
		return err
	}
	return g.s.UpdateGame(g)
}

func (g *game) Draw(screen *ebiten.Image) {
	g.s.Draw(screen)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.s.Layout(outsideHeight, outsideWidth)
}

type disc int
type peg []disc

type board struct {
	hasInit bool

	pegs []peg

	p1 int

	quit bool
}

func (b *board) init() error {
	b.pegs = []peg{
		{2, 1, 0},
		{},
		{},
	}
	b.hasInit = true
	return nil
}

func (b *board) Update() error {
	if !b.hasInit {
		err := b.init()
		if err != nil {
			return err
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		b.quit = true
	}

	i := b.gameKeyPressed()
	if i > 0 && i != b.p1 {
		if b.p1 == 0 {
			if len(b.pegs[i-1]) > 0 {
				b.p1 = i
			}
		} else {
			p1 := b.pegs[b.p1-1]
			p2 := b.pegs[i-1]
			b.pegs[i-1] = append(p2, p1[len(p1)-1])
			b.pegs[b.p1-1] = p1[0 : len(p1)-1]
			b.p1 = 0
		}
	}

	return nil
}

func (b *board) gameKeyPressed() int {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		return 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		return 2
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		return 3
	}
	return 0
}

func (b *board) UpdateGame(g *game) error {
	if b.quit {
		g.s = &menu{}
	}
	return nil
}

func (b *board) Draw(screen *ebiten.Image) {
	max := screen.Bounds().Max
	n := len(b.pegs) + 2
	for i, peg := range b.pegs {
		x := max.X / n * (i + 1)
		post := 3 - float64(len(peg))
		ebitenutil.DrawRect(screen, float64(x)+27.5, float64(max.Y)/3+20, 5, 15+(20*post), color.White)

		for j, disc := range peg {
			const height = 15
			const width = 40
			const step = 10
			const gap = 5
			y := max.Y/3*2 - (j * (height + gap))
			dw := float64(disc) * step
			ebitenutil.DrawRect(screen, float64(x)+10-dw/2, float64(y), width+dw, 15, color.White)
		}
		if i+1 == b.p1 {
			ebitenutil.DrawRect(screen, float64(x)+25, float64(max.Y)/3, 10, 10, color.White)
		}
	}
}

func (b *board) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
