package scripts

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// EnemySnake object
type EnemySnake struct {
	game             *Game
	nBodyP           int
	lastDir          string
	sHeadU           ebiten.Image
	sHeadD           ebiten.Image
	sHeadL           ebiten.Image
	sHeadR           ebiten.Image
	horizontal       ebiten.Image
	partsOfBody      [][]float64
	seed             rand.Source
	pointsWaiting    int
	score            int
	channelMovements chan int
	collision        bool
}

// Initialize
func CreateEnemySnake(g *Game) *EnemySnake {
	e := EnemySnake{
		game:          g,
		nBodyP:        0,
		lastDir:       "right",
		pointsWaiting: 0,
		collision:     false,
	}

	e.channelMovements = make(chan int)
	e.seed = rand.NewSource(time.Now().UnixNano())
	random := rand.New(e.seed)

	e.partsOfBody = append(e.partsOfBody, []float64{float64(random.Intn(30) * 20), float64(random.Intn(30) * 20)})

	sHeadU, _, _ := ebitenutil.NewImageFromFile("images/headUEne.png", ebiten.FilterDefault)
	sHeadD, _, _ := ebitenutil.NewImageFromFile("images/headDEne.png", ebiten.FilterDefault)
	sHeadL, _, _ := ebitenutil.NewImageFromFile("images/headLEne.png", ebiten.FilterDefault)
	sHeadR, _, _ := ebitenutil.NewImageFromFile("images/headREne.png", ebiten.FilterDefault)
	horizontal, _, _ := ebitenutil.NewImageFromFile("images/bodyEne.png", ebiten.FilterDefault)
	e.sHeadU = *sHeadU
	e.sHeadD = *sHeadD
	e.sHeadL = *sHeadL
	e.sHeadR = *sHeadR
	e.horizontal = *horizontal

	return &e
}

// Enemy movement
func (s *EnemySnake) ChannelPipe() error {
	for {
		dotTime := <-s.channelMovements
		s.Direction(dotTime)
	}
}

// Direction updates of enemy
func (s *EnemySnake) Direction(dotTime int) error {
	if dotTime == 1 {
		random := rand.New(s.seed)
		action := random.Intn(4)
		changingDirection := random.Intn(3)
		posX, posY := s.GetHeadPos()
		if changingDirection == 0 {
			switch action {
			case 0:
				if posX < 560 && s.lastDir != "left" {
					s.lastDir = "right"
				} else {
					s.lastDir = "left"
				}
				return nil
			case 1:
				if posY < 660 && s.lastDir != "up" {
					s.lastDir = "down"
				} else {
					s.lastDir = "up"
				}
				return nil
			case 2:
				if posY > 20 && s.lastDir != "down" {
					s.lastDir = "up"
				} else {
					s.lastDir = "down"
				}
				return nil
			case 3:
				if posX > 20 && s.lastDir != "right" {
					s.lastDir = "left"
				} else {
					s.lastDir = "right"
				}
				return nil
			}
		}
		// Bounds the collision
		if posX >= 560 {
			s.lastDir = "left"
			return nil
		}
		if posX == 20 {
			s.lastDir = "right"
			return nil
		}
		if posY == 660 {
			s.lastDir = "up"
			return nil
		}
		if posY == 20 {
			s.lastDir = "down"
			return nil
		}
	}

	if dotTime == 1 { // Checks collision with enemy
		xPos, yPos := s.game.snake.getHeadPos()
		if s.CollSnake(xPos, yPos) {
			s.game.snake.game.crashed = true
			s.game.gameOver()
		}
	}
	return nil
}

// Draws the snake
func (s *EnemySnake) Draw(screen *ebiten.Image, dotTime int) error {
	if s.game.alive {
		s.UpdatePos(dotTime)
	}
	enemyDO := &ebiten.DrawImageOptions{}
	xPos, yPos := s.GetHeadPos()
	enemyDO.GeoM.Translate(xPos, yPos)

	if s.lastDir == "up" {
		screen.DrawImage(&s.sHeadU, enemyDO)
	} else if s.lastDir == "down" {
		screen.DrawImage(&s.sHeadD, enemyDO)
	} else if s.lastDir == "right" {
		screen.DrawImage(&s.sHeadR, enemyDO)
	} else if s.lastDir == "left" {
		screen.DrawImage(&s.sHeadL, enemyDO)
	}

	for i := 0; i < s.nBodyP; i++ {
		partDO := &ebiten.DrawImageOptions{}
		xPos, yPos := s.GetBody(i)
		partDO.GeoM.Translate(xPos, yPos)
		screen.DrawImage(&s.horizontal, partDO)
	}

	return nil
}

// Updates head position score
func (s *EnemySnake) UpdatePos(dotTime int) {
	if dotTime == 1 {
		if s.pointsWaiting > 0 {
			s.nBodyP++
			s.pointsWaiting--
		}
		switch s.lastDir {
		case "up":
			s.TranslateHeadPos(0, -20)
		case "down":
			s.TranslateHeadPos(0, +20)
		case "right":
			s.TranslateHeadPos(20, 0)
		case "left":
			s.TranslateHeadPos(-20, 0)
		}

	}
}

// Evaluating if there was a collision
func (s *EnemySnake) CollSnake(xPos, yPos float64) bool {
	for i := 0; i < len(s.partsOfBody); i++ {
		if xPos == s.partsOfBody[i][0] && yPos == s.partsOfBody[i][1] {
			return true
		}
	}
	return false
}

// Head pos is retuned
func (s *EnemySnake) GetHeadPos() (float64, float64) {
	return s.partsOfBody[0][0], s.partsOfBody[0][1]
}

// Last body pos is returned
func (s *EnemySnake) GetBody(pos int) (float64, float64) {
	return s.partsOfBody[pos+1][0], s.partsOfBody[pos+1][1]
}

// AddPoint controls game's score
func (s *EnemySnake) AddPoint() {
	s.score++
	s.pointsWaiting++
}

// Game score control
func (s *EnemySnake) AddParts(newX, newY float64) {
	s.partsOfBody = append([][]float64{{newX, newY}}, s.partsOfBody...)
	s.partsOfBody = s.partsOfBody[:s.nBodyP+1]
}

// Changes pos
func (s *EnemySnake) TranslateHeadPos(newXPos, newYPos float64) {
	newX := s.partsOfBody[0][0] + newXPos
	newY := s.partsOfBody[0][1] + newYPos
	s.AddParts(newX, newY)
}
