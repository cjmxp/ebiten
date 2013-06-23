package sprites

import (
	"github.com/hajimehoshi/go.ebiten/graphics"
	"github.com/hajimehoshi/go.ebiten/graphics/matrix"
	"image"
	"image/color"
	"math/rand"
	"os"
	"time"
)

type Sprite struct {
	width  int
	height int
	ch     chan bool
	x      int
	y      int
	vx     int
	vy     int
}

func newSprite(screenWidth, screenHeight, width, height int) *Sprite {
	maxX := screenWidth - width
	maxY := screenHeight - height
	sprite := &Sprite{
		width:  width,
		height: height,
		ch:     make(chan bool),
		x:      rand.Intn(maxX),
		y:      rand.Intn(maxY),
		vx:     rand.Intn(2)*2 - 1,
		vy:     rand.Intn(2)*2 - 1,
	}
	go sprite.update(screenWidth, screenHeight)
	return sprite
}

func (sprite *Sprite) update(screenWidth, screenHeight int) {
	maxX := screenWidth - sprite.width
	maxY := screenHeight - sprite.height
	for {
		<-sprite.ch
		sprite.x += sprite.vx
		sprite.y += sprite.vy
		if sprite.x < 0 || maxX <= sprite.x {
			sprite.vx = -sprite.vx
		}
		if sprite.y < 0 || maxY <= sprite.y {
			sprite.vy = -sprite.vy
		}
		sprite.ch <- true
	}
}

func (sprite *Sprite) Update() {
	sprite.ch <- true
	<-sprite.ch
}

type Sprites struct {
	ebitenTexture graphics.Texture
	sprites       []*Sprite
}

func New() *Sprites {
	return &Sprites{}
}

func (game *Sprites) ScreenWidth() int {
	return 256
}

func (game *Sprites) ScreenHeight() int {
	return 240
}

func (game *Sprites) Fps() int {
	return 60
}

func (game *Sprites) Init(tf graphics.TextureFactory) {
	file, err := os.Open("ebiten.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	if game.ebitenTexture, err = tf.NewTextureFromImage(img); err != nil {
		panic(err)
	}
	game.sprites = []*Sprite{}
	for i := 0; i < 100; i++ {
		sprite := newSprite(
			game.ScreenWidth(),
			game.ScreenHeight(),
			game.ebitenTexture.Width,
			game.ebitenTexture.Height)
		game.sprites = append(game.sprites, sprite)
	}
}

func (game *Sprites) Update() {
	for _, sprite := range game.sprites {
		sprite.Update()
	}
}

func (game *Sprites) Draw(g graphics.GraphicsContext, offscreen graphics.Texture) {
	g.Fill(&color.RGBA{R: 128, G: 128, B: 255, A: 255})

	// Draw the sprites
	locations := make([]graphics.TexturePart, 0, len(game.sprites))
	texture := game.ebitenTexture
	for _, sprite := range game.sprites {
		location := graphics.TexturePart{
			LocationX: sprite.x,
			LocationY: sprite.y,
			Source: graphics.Rect{
				0, 0, texture.Width, texture.Height,
			},
		}
		locations = append(locations, location)
	}
	geometryMatrix := matrix.IdentityGeometry()
	g.DrawTextureParts(texture.ID, locations,
		geometryMatrix, matrix.IdentityColor())
}

func init() {
	rand.Seed(time.Now().UnixNano())
}