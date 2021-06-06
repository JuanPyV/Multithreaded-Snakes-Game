package scripts

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Food is what is eaten
type Food struct {
	x       float64
	y       float64
	eaten   bool
	foodImg ebiten.Image
}

// Here we generate the foods
func GenFood() *Food {
	food := Food{
		eaten:  false,
	}
	foodImg, _, _ := ebitenutil.NewImageFromFile("images/cherry.png", ebiten.FilterDefault)
	food.foodImg = *foodImg

	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	food.x = float64(random.Intn(28)*20 + 20)
	food.y = float64(random.Intn(30)*20 + 60)
	return &food
}

// Update to the delicious foodImg deletion
func (food *Food) Update(dotTime int) error {
	return nil
}

// Draw the foodImg
func (food *Food) Draw(screen *ebiten.Image, dotTime int) error {
	drawer := &ebiten.DrawImageOptions{}
	drawer.GeoM.Translate(food.x, food.y)
	screen.DrawImage(&food.foodImg, drawer)
	return nil
}
