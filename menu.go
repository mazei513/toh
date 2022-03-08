package main

import (
	"errors"
	"image/color"
	"io/ioutil"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type menuAction int

const (
	menuActionNone = iota
	menuActionConfirm
	menuActionUp
	menuActionDown
)

type menuOption int

const (
	menuOptionStart = iota
	menuOptionQuit
	menuOptionMax
)

func menuOptionText(o menuOption) string {
	switch o {
	case menuOptionStart:
		return "Start"
	case menuOptionQuit:
		return "Quit"
	}
	return ""
}

var errQuit = errors.New("quit")

type menu struct {
	hasInit bool

	option   menuOption
	selected bool

	font font.Face
}

func (m *menu) Draw(screen *ebiten.Image) {
	if !m.hasInit {
		return
	}

	max := screen.Bounds().Max
	text.Draw(screen, menuOptionText(m.option), m.font, max.X/3, max.Y/2, color.White)
}

func (m *menu) Update() error {
	if !m.hasInit {
		m.init()
	}

	a := m.inputUpdate()
	switch a {
	case menuActionDown:
		m.option = (m.option + 1) % menuOptionMax
	case menuActionUp:
		m.option = (m.option - 1) % menuOptionMax
	case menuActionConfirm:
		m.selected = true
	}

	return nil
}

func (m *menu) UpdateGame(g *game) error {
	if !m.selected {
		return nil
	}
	switch m.option {
	case menuOptionStart:
		g.s = &board{}
	case menuOptionQuit:
		return errQuit
	}
	return nil
}

func (*menu) Layout(outsideHeight, outsideWidth int) (screenHeight, screenWidth int) {
	return 800, 600
}

func (m *menu) init() error {
	b, err := ioutil.ReadFile("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf")
	if err != nil {
		return err
	}

	f, err := opentype.Parse(b)
	if err != nil {
		return err
	}

	ff, err := opentype.NewFace(f, &opentype.FaceOptions{Size: 14, DPI: 200})
	if err != nil {
		return err
	}

	m.font = ff

	m.hasInit = true
	return nil
}

func (*menu) inputUpdate() menuAction {
	if isKeysJustPress([]ebiten.Key{ebiten.KeyArrowDown, ebiten.KeyDown}) {
		return menuActionDown
	}
	if isKeysJustPress([]ebiten.Key{ebiten.KeyArrowUp, ebiten.KeyUp}) {
		return menuActionUp
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return menuActionConfirm
	}
	return menuActionNone
}
