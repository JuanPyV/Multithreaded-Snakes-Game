package scripts

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"time"

	"github.com/hajimehoshi/ebiten"
)

// Game : Main object of the scene. Parent of everything
type Game struct {
	snake        *Snake
	hud          *Hud
	foods        []*Food
	snakeChannel chan int
	alive        bool
	crashed      bool
	numFood      int
	numEnemies   int
	score        int
	dotTime      int
}


// NewGame : Starts a new game assigning variables
func NewGame(nFood int, nEnemies int) Game {
	game := Game{
		alive:      true,
		crashed:    false,
		score:      0,
		dotTime:    0,
		numFood:    nFood,
		numEnemies: nEnemies,
	}

	foodArray := make([]*Food, game.numFood) //store all the cherries
	for i := 0; i < game.numFood; i++ {
		foodArray[i] = GenFood(&game)
		time.Sleep(20)
	}
	game.foods = foodArray

	game.snake = createSnake(&game)
	game.snakeChannel = make(chan int)
	go func() {
		err := game.snake.Behavior()
		if err != nil {

		}
	}()
	game.hud = initHud(&game)
	fmt.Printf("Food: %d \n", nFood)
	fmt.Printf("Enemies: %d \n", nEnemies)
	fmt.Println(foodArray)
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
		if g.numFood == 0 { //when all cherries has been eating the game ends
			g.alive = false
		}
		//update the channels
		g.dotTime = (g.dotTime + 1) % 10

		if err := g.snake.Update(g.dotTime); err != nil {
			g.snakeChannel <- g.dotTime
		}
		xPos, yPos := g.snake.getHeadPos()
		for i := 0; i < len(g.foods); i++ {
			if xPos == g.foods[i].x && yPos == g.foods[i].y { //if snake eats a cherry grows
				g.foods[i].y = -20
				g.foods[i].x = -20
				g.hud.addPoint()
				g.numFood--
				g.snake.addPoint()
				break
			}
		}
	}

	for i := 0; i < g.numFood; i++ {
		if err := g.foods[i].Update(g.dotTime); err != nil {
			return err
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

	for i := 0; i < len(g.foods); i++ {
		if err := g.foods[i].Draw(screen, g.dotTime); err != nil {
			return err
		}
	}
	return nil
}
