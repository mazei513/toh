package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func isKeysJustPress(keys []ebiten.Key) bool {
	for _, k := range keys {
		if inpututil.IsKeyJustPressed(k) {
			return true
		}
	}
	return false
}
