package main

import (
	"embed"
	"fmt"
	"github.com/frabjous-studios/stackedsprite"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"

	_ "image/png"
)

func main() {
	g := Game{
		camera: &ebiten.GeoM{},
		bg:     solidColor(color.Gray{150}),
	}
	g.car = NewCar(g.camera)

	g.camera.Scale(2, 2)
	g.camera.Translate(160, 120)

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	camera *ebiten.GeoM

	car *Car
	bg  *ebiten.Image

	ticks int
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw background
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(320, 240)
	screen.DrawImage(g.bg, opts)

	// draw car
	g.car.DrawTo(screen)
}

const radsPerSecond = 2 * math.Pi / 5

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.car.ApplyGas()
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.car.ApplyBrake()
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.car.TurnRight()
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.car.TurnLeft()
	}

	g.car.Update()
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

const maxVel = 0.5

type Car struct {
	sprite *stackedsprite.StackedSprite

	a      float64 // acceleration magnitude
	dx, dy float64 // facing (unit vector)
	vx, vy float64 // velocity
}

func NewCar(camera *ebiten.GeoM) *Car {
	slices := LoadSpriteRow("img/BlueCar.png", 16)
	sprite := stackedsprite.NewStackedSprite(slices)
	sprite.GlobalM = camera
	return &Car{
		sprite: sprite,
		dx:     1, dy: 0,
	}
}

func (c *Car) DrawTo(screen *ebiten.Image) {
	c.sprite.DrawTo(screen)
}

const friction = 0.98

func (c *Car) Update() {
	tps := float64(ebiten.TPS())
	ax, ay := c.a*c.dx, c.a*c.dy

	c.vx += ax / tps
	c.vy += ay / tps

	// apply friction
	c.vx *= friction
	c.vy *= friction

	c.sprite.MoveX(c.vx)
	c.sprite.MoveY(c.vy)

	// reset accel for input next frame.
	c.a = 0
}

func (c *Car) ApplyGas() {
	c.a = 1
}

func (c *Car) ApplyBrake() {
	c.a = -1
}

const turnTheta = 0.01

func (c *Car) TurnRight() {
	geom := ebiten.GeoM{}
	geom.Rotate(-turnTheta)
	c.dx, c.dy = geom.Apply(c.dx, c.dy)
	c.sprite.Rotate(-turnTheta)
}

func (c *Car) TurnLeft() {
	geom := ebiten.GeoM{}
	geom.Rotate(turnTheta)
	c.dx, c.dy = geom.Apply(c.dx, c.dy)
	c.sprite.Rotate(turnTheta)
}

//go:embed img/*.png
var images embed.FS

func LoadSpriteRow(filename string, dims int) []*image.NRGBA {
	f, err := images.Open(filename)
	if err != nil {
		panic(fmt.Errorf("error opening: %w", err))
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		panic(fmt.Errorf("error decoding: %w", err))
	}
	return LoadTiles(img, image.Rect(0, 0, dims, dims))
}

// LoadTiles splits up a tiled horizontal image strip using the provided rectangle as the starting tile.
// All tiles are presumed to be the same size, the strip ends at the end of the image.
func LoadTiles(img image.Image, rect image.Rectangle) []*image.NRGBA {
	var result []*image.NRGBA
	destRect := image.Rect(0, 0, rect.Dx(), rect.Dx())
	for x := rect.Min.X; x < img.Bounds().Max.X; x += rect.Dx() {
		frame := image.NewNRGBA(destRect)
		draw.Draw(frame, destRect, img, image.Pt(x, rect.Min.Y), draw.Over)
		result = append(result, frame)
	}
	return result
}

func solidColor(c color.Color) *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(c)
	return img
}
