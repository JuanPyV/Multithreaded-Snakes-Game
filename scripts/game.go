package scripts

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
)

// Game : Main object of the scene. Parent of everything
type Game struct {
	snake        *Snake
	hud          *Hud
	foods        []*Food
	enemies      []*EnemySnake
	enemiesChan  []chan int
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
		foodArray[i] = GenFood()
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

	arrayEnemies := make([]*EnemySnake, game.numEnemies)
	for i := 0; i < len(arrayEnemies); i++ {
		arrayEnemies[i] = CreateEnemySnake(&game)
		time.Sleep(20)
	}
	enemiesChan := make([]chan int, game.numEnemies)
	for i := 0; i < len(enemiesChan); i++ {
		enemiesChan[i] = make(chan int)
		arrayEnemies[i].channelMovements = enemiesChan[i]
		go arrayEnemies[i].ChannelPipe()
		time.Sleep(20)
	}
	game.enemiesChan = enemiesChan
	game.enemies = arrayEnemies

	game.hud = initHud(&game)
	fmt.Printf("Food: %d \n", nFood)
	fmt.Printf("Enemies: %d \n", nEnemies)
	//fmt.Println(foodArray)
	return game
}

// gameOver ends the game
func (game *Game) gameOver() {
	game.alive = false //boolean to keep alive
}

func (game *Game) Crashed() {
	game.crashed = true
}

// Update the main process of the game
func (game *Game) Update() error {
	if game.alive {
		if game.numFood == 0 { //when all cherries has been eating the game ends
			game.hud.game.alive = false
			largest := game.enemies[0]
			for i := 1; i < len(game.enemies); i++ {
				if game.enemies[i].score > largest.score {
					largest = game.enemies[i]
				}
			}

			if game.snake.score > largest.score {
				game.hud.bigger = true
			} else {
				game.hud.bigger = false
			}

		}
		//update the channels
		game.dotTime = (game.dotTime + 1) % 5

		if err := game.snake.Update(game.dotTime); err != nil {
			game.snakeChannel <- game.dotTime
		}
		for i := 0; i < len(game.enemiesChan); i++ {
			game.enemiesChan[i] <- game.dotTime
		}
		xPos, yPos := game.snake.getHeadPos()
		for i := 0; i < len(game.foods); i++ {
			if xPos == game.foods[i].x && yPos == game.foods[i].y { //if snake eats a cherry grows
				game.foods[i].y = -20
				game.foods[i].x = -20
				game.hud.addPoint()
				game.numFood--
				game.snake.addPoint()
				break
			}
		}
		for j := 0; j < len(game.enemies); j++ {
			xPos, yPos := game.enemies[j].GetHeadPos()
			for i := 0; i < len(game.foods); i++ {
				if xPos == game.foods[i].x && yPos == game.foods[i].y { //if snake eats a cherry grows
					game.foods[i].y = -20
					game.foods[i].x = -20
					game.numFood--
					game.enemies[j].AddPoint()
					break
				}
			}
		}
	}
	for i := 0; i < game.numFood; i++ {
		if err := game.foods[i].Update(game.dotTime); err != nil {
			return err
		}
	}

	return nil
}

// Draw the whole interface
func (game *Game) Draw(screen *ebiten.Image) error {

	drawer := &ebiten.DrawImageOptions{}
	drawer.GeoM.Translate(0, 0)
	background, _, _ := ebitenutil.NewImageFromFile("images/background.png", ebiten.FilterLinear)
	screen.DrawImage(background, drawer)

	if err := game.snake.Draw(screen, game.dotTime); err != nil {
		return err
	}

	for _, enemy := range game.enemies {
		if err := enemy.Draw(screen, game.dotTime); err != nil {
			return err
		}
	}

	if err := game.hud.Draw(screen); err != nil {
		return err
	}

	for i := 0; i < len(game.foods); i++ {
		if err := game.foods[i].Draw(screen, game.dotTime); err != nil {
			return err
		}
	}
	return nil
}
