package scripts

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Food is what is eaten
type Food struct {
	xLimit int
	yLimit int
	xPos   float64
	yPos   float64
	eaten  bool
	game   *Game
	taco   ebiten.Image
}

// Here we generate the food
func GenFood(g *Game) *Food {
	c := Food{
		game:   g,
		xLimit: 30,
		yLimit: 30,
		eaten:  false,
	}
	taco, _, _ := ebitenutil.NewImageFromFile("images/taco.png", ebiten.FilterDefault)
	c.taco = *taco

	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	c.xPos = float64(random.Intn(c.xLimit) * 20)
	c.yPos = float64(random.Intn(c.yLimit) * 20)
	return &c
}

// Update to the delicious taco deletion
func (c *Food) Update(dotTime int) error {
	if c.eaten == false {
		return nil
	}
	return nil
}

// Draw the taco
func (c *Food) Draw(screen *ebiten.Image, dotTime int) error {
	drawer := &ebiten.DrawImageOptions{}
	drawer.GeoM.Translate(c.xPos, c.yPos)
	screen.DrawImage(&c.taco, drawer)
	return nil
}
