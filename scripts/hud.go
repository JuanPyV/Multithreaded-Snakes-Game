package scripts

import (
	"golang.org/x/image/font/inconsolata"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

type Hud struct {
	game      *Game
	score     int
}

// initHud Constructor
func initHud(g *Game) *Hud {
	hud := Hud{
		game:         g,
		score:        0,
	}

	return &hud
}

func (hud *Hud) addPoint() {
	hud.score++
}


// Draw the hud
func (hud *Hud) Draw(screen *ebiten.Image) error {
	text.Draw(screen, "Score: "+strconv.Itoa(hud.score), inconsolata.Bold8x16, 20, 30, color.Black)
	if hud.game.alive == false {
		textGameOver := ""
		if hud.game.crashed {
			textGameOver = "GAME OVER !!\n"
		}
		text.Draw(screen, textGameOver, inconsolata.Bold8x16, 250, 350, color.Black)
	}
	return nil
}
