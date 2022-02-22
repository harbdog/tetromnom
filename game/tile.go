package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	maxMovingCount  = 5
	maxPoppingCount = 6
)

var (
	tileSize   = 40
	tileMargin = 3
	tileImage  = ebiten.NewImage(tileSize, tileSize)
)

func init() {
	tileImage.Fill(color.White)
}

// TilePosition represents a tile position.
type TilePosition struct {
	x int
	y int
}

// Tile represents a tile information including Tetromino shape, position and animation states.
type Tile struct {
	x, y  int
	tetro *Tetromino

	poppingCount int
}

// Pos returns the tile's current position.
func (t *Tile) Pos() (int, int) {
	return t.x, t.y
}

// IsMoving returns a boolean value indicating if the tile is animating.
func (t *Tile) IsMoving() bool {
	return 0 < t.poppingCount
}

func (t *Tile) stopAnimation() {
	t.poppingCount = 0
}

// Update updates the tile's animation states.
func (t *Tile) Update() error {
	switch {
	case 0 < t.poppingCount:
		t.poppingCount--
	}
	return nil
}

func colorToScale(clr color.Color) (float64, float64, float64, float64) {
	r, g, b, a := clr.RGBA()
	rf := float64(r) / 0xffff
	gf := float64(g) / 0xffff
	bf := float64(b) / 0xffff
	af := float64(a) / 0xffff
	// Convert to non-premultiplied alpha components.
	if 0 < af {
		rf /= af
		gf /= af
		bf /= af
	}
	return rf, gf, bf, af
}

func mean(a, b int, rate float64) int {
	return int(float64(a)*(1-rate) + float64(b)*rate)
}

func meanF(a, b float64, rate float64) float64 {
	return a*(1-rate) + b*rate
}

// Draw draws the current tile to the given boardImage.
func (t *Tile) Draw(boardImage *ebiten.Image) {
	// get board height to render origin coordinate at bottom left instead of top left
	_, bh := boardImage.Size()

	op := &ebiten.DrawImageOptions{}
	x := t.x*tileSize + (t.x+1)*tileMargin
	y := bh - ((t.y+1)*tileSize + (t.y+1)*tileMargin)
	switch {
	case 0 < t.poppingCount:
		const maxScale = 1.2
		rate := 0.0
		if maxPoppingCount*2/3 <= t.poppingCount {
			// 0 to 1
			rate = 1 - float64(t.poppingCount-2*maxPoppingCount/3)/float64(maxPoppingCount/3)
		} else {
			// 1 to 0
			rate = float64(t.poppingCount) / float64(maxPoppingCount*2/3)
		}
		scale := meanF(1.0, maxScale, rate)
		op.GeoM.Translate(float64(-tileSize/2), float64(-tileSize/2))
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(tileSize/2), float64(tileSize/2))
	}
	op.GeoM.Translate(float64(x), float64(y))
	r, g, b, a := colorToScale(t.tetro.color)
	op.ColorM.Scale(r, g, b, a)
	boardImage.DrawImage(tileImage, op)
}
