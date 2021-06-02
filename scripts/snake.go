package scripts

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Position struct {
	X float64
	Y float64
}

// Snake : Object which the player controls
type Snake struct {
	game         *Game
	len          int
	dir          int
	headImgRight ebiten.Image
	headImgLeft  ebiten.Image
	headImgDown  ebiten.Image
	headImgUp    ebiten.Image
	tailImg      ebiten.Image
	parts        []Position
	score        int
	behavior     chan int
}

// createSnake : Generates a snake
func createSnake(g *Game) *Snake {
	snake := Snake{
		game: g,
		len:  0,
		dir:  1,
	}
	snake.behavior = make(chan int)
	snake.parts = append(snake.parts, Position{300, 300})
	headimgright, _, _ := ebitenutil.NewImageFromFile("images/snakeHeadRight.png", ebiten.FilterLinear)
	headimgleft, _, _ := ebitenutil.NewImageFromFile("images/snakeHeadLeft.png", ebiten.FilterLinear)
	headimgdown, _, _ := ebitenutil.NewImageFromFile("images/snakeHeadDown.png", ebiten.FilterLinear)
	headimgup, _, _ := ebitenutil.NewImageFromFile("images/snakeHeadUp.png", ebiten.FilterLinear)
	tailimg, _, _ := ebitenutil.NewImageFromFile("images/snakeTail.png", ebiten.FilterLinear)
	snake.headImgRight = *headimgright
	snake.headImgLeft = *headimgleft
	snake.headImgDown = *headimgdown
	snake.headImgUp = *headimgup
	snake.tailImg = *tailimg
	return &snake
}

func (snake *Snake) Behavior() error {
	dotTime := <-snake.behavior
	for {
		err := snake.Update(dotTime)
		if err != nil {
			return err
		}
		dotTime = <-snake.behavior
	}
}

func (snake *Snake) Update(dotTime int) error {
	// Right : 3
	// Left: 1
	// Down: 2
	// Up: 5
	if snake.game.alive {
		if ebiten.IsKeyPressed(ebiten.KeyRight) && snake.dir != 3 && snake.dir != 1{
			snake.dir = 3
			return nil
		} else if ebiten.IsKeyPressed(ebiten.KeyLeft) && snake.dir != 1 && snake.dir != 3 {
			snake.dir = 1
			return nil
		}else if ebiten.IsKeyPressed(ebiten.KeyDown) && snake.dir != 2 && snake.dir != 5{
			snake.dir = 2
			return nil
		} else if ebiten.IsKeyPressed(ebiten.KeyUp) && snake.dir != 5 && snake.dir != 2 {
			snake.dir = 5
			return nil
		}

		if dotTime == 1 { // if collision with bounds
			xPos, yPos := snake.getHeadPos()
			if xPos < 20 || xPos > 560 || yPos < 60 || yPos > 660{
				snake.game.Crashed()
				snake.game.gameOver()
			}
		}
	}
	return nil
}

// Draw the snake
func (snake *Snake) Draw(screen *ebiten.Image, dotTime int) error {
	if snake.game.alive {
		snake.UpdatePos(dotTime)
	}

	drawer := &ebiten.DrawImageOptions{}
	x, y := snake.getHeadPos()
	drawer.GeoM.Translate(x, y)

	if snake.dir == 5 {
		screen.DrawImage(&snake.headImgUp, drawer)
	} else if snake.dir == 2 {
		screen.DrawImage(&snake.headImgDown, drawer)
	} else if snake.dir == 3 {
		screen.DrawImage(&snake.headImgRight, drawer)
	} else if snake.dir == 1 {
		screen.DrawImage(&snake.headImgLeft, drawer)
	}

	return nil
}

// UpdatePos turn to a direction
func (snake *Snake) UpdatePos(dotTime int) {
	if dotTime == 1 {
		switch snake.dir {
		case 5:
			snake.turnDir(0, -20)
		case 2:
			snake.turnDir(0, +20)
		case 3:
			snake.turnDir(20, 0)
		case 1:
			snake.turnDir(-20, 0)
		}
	}
}

func (snake *Snake) addPoint() {
	snake.score++
}

func (snake *Snake) getHeadPos() (float64, float64) { // get position of the head
	return snake.parts[0].X, snake.parts[0].Y
}

func (snake *Snake) turnDir(newXPos, newYPos float64) {
	newX := snake.parts[0].X + newXPos
	newY := snake.parts[0].Y + newYPos
	snake.updateParts(newX, newY)
}

func (snake *Snake) updateParts(newX, newY float64) {
	snake.parts = append([]Position{{newX, newY}})
}

