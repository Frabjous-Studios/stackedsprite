package stackedsprite

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
)

// StackedSprite represents a stacked sprite in integer world coordinates. The sprite is pre-rendered in its rotated and
// stacked form whenever theta changes.
type StackedSprite struct {
	// slices stores a list of sprite slices
	slices []*ebiten.Image

	// frame stores stacked slices, pre-rendered to the GPU
	frame *ebiten.Image

	// globalM stores a reference to the global world / camera matrix for this sprite.
	GlobalM *ebiten.GeoM

	x, y, z    int
	fx, fy, fz float64

	theta        float64
	needsReframe bool
}

// NewStackedSprite returns a new stacked sprite which uses the provided slices.
func NewStackedSprite(slices []*image.NRGBA) *StackedSprite {
	if len(slices) == 0 {
		return nil
	}
	dx, dy := slices[0].Bounds().Dx(), slices[0].Bounds().Dy()
	w := int(math.Ceil(math.Sqrt2 * float64(dx)))
	h := int(math.Ceil(math.Sqrt2*float64(dy) + float64(len(slices))))
	buf := ebiten.NewImage(w, h)
	var images []*ebiten.Image
	for _, slice := range slices {
		images = append(images, ebiten.NewImageFromImage(slice))
	}
	return &StackedSprite{
		slices:       images,
		frame:        buf,
		needsReframe: true,
	}
}

func (s *StackedSprite) DrawTo(screen *ebiten.Image) {
	if s.needsReframe {
		s.reframe()
		s.needsReframe = false
	}
	// draw after reframe
	opt := &ebiten.DrawImageOptions{}

	cx, cy := s.Origin()
	opt.GeoM.Translate(-cx, -cy)
	opt.GeoM.Translate(float64(s.x), float64(s.y-s.z))
	if s.GlobalM != nil {
		opt.GeoM.Concat(*s.GlobalM)
	}
	screen.DrawImage(s.frame, opt)
}

// Rotate rotates this sprite in units of radians.
func (s *StackedSprite) Rotate(theta float64) {
	if theta == 0 {
		return
	}
	s.theta += theta
	s.theta = math.Mod(s.theta, 2*math.Pi)
	if s.theta < 0 {
		s.theta += 2 * math.Pi
	}

	s.needsReframe = true
}

// MoveX moves this sprite along the X axis by the provided amount.
func (s *StackedSprite) MoveX(amt float64) {
	s.fx += amt
	px, fx := math.Modf(s.fx)
	s.fx = fx
	s.x = s.x + int(px)
}

// MoveY moves this sprite along the Y axis by the provided amount.
func (s *StackedSprite) MoveY(amt float64) {
	s.fy += amt
	py, fy := math.Modf(s.fy)
	s.fy = fy
	s.y = s.y + int(py)
}

// MoveZ moves this sprite along the Z axis by the provided amount.
func (s *StackedSprite) MoveZ(amt float64) {
	s.fz += amt
	pz, fz := math.Modf(s.fz)
	s.fz = fz
	s.z = s.z + int(pz)
}

func (s *StackedSprite) Origin() (float64, float64) {
	return float64(s.frame.Bounds().Dx() / 2), float64(s.frame.Bounds().Dy()/2 + len(s.slices))
}

func (s *StackedSprite) reframe() {
	s.frame.Clear()
	opt := &ebiten.DrawImageOptions{}

	// ox, oy are the origin in the frame
	ox, oy := s.Origin()
	// cx, cy are the center of each slice
	dx, dy := float64(s.slices[0].Bounds().Dx()), float64(s.slices[0].Bounds().Dy())
	opt.GeoM.Translate(-dx/2, -dy/2)
	opt.GeoM.Rotate(s.theta)
	opt.GeoM.Translate(ox, oy)
	for i := 0; i < len(s.slices); i++ {
		s.frame.DrawImage(s.slices[i], opt)
		opt.GeoM.Translate(0, -1)
	}
}
