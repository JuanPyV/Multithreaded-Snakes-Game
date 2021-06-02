package scripts

import (
	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
)

// Game : Main object of the scene. Parent of everything
type Game struct {
	snake        *Snake
	snakeChannel chan int
	hud          *Hud
	numFood      int
	numEnemies   int
	alive        bool
	crashed      bool
	score        int
	dotTime      int
}

// NewGame : Starts a new game assigning variables
func NewGame(food int, enemies int) Game {
	game := Game{
		alive:      true,
		crashed:      false,
		score:      0,
		dotTime:    0,
		numFood:    food,
		numEnemies: enemies,
	}

	game.snake = createSnake(&game)
	game.snakeChannel = make(chan int)
	go func() {
		err := game.snake.Behavior()
		if err != nil {

		}
	}()
	game.hud = initHud(&game)
	return game
}

// gameOver ends the game
func (g *Game) gameOver() {
	g.alive = false //boolean to keep alive
}

func (g *Game) Crashed() {
	g.crashed = true
}

// Update the main process of the game
func (g *Game) Update() error {
	if g.alive {
		//update the channels
		g.dotTime = (g.dotTime + 1) % 10

		if err := g.snake.Update(g.dotTime); err != nil {
			g.snakeChannel <- g.dotTime
		}
	}
	return nil
}

// Draw the whole interface
func (g *Game) Draw(screen *ebiten.Image) error {

	drawer := &ebiten.DrawImageOptions{}
	drawer.GeoM.Translate(0, 0)
	background, _, _ := ebitenutil.NewImageFromFile("images/background.png", ebiten.FilterLinear)
	screen.DrawImage(background, drawer)

	if err := g.snake.Draw(screen, g.dotTime); err != nil {
		return err
	}

	if err := g.hud.Draw(screen); err != nil {
		return err
	}

	return nil
}
