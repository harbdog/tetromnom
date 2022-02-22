package game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/gonum/mat"
)

type TetrominoShape int

const (
	ShapeI TetrominoShape = iota
	ShapeJ
	ShapeL
	ShapeS
	ShapeZ
	ShapeT
	ShapeO
)

// Tetromino represents the piece containing the individual tiles making up the shape and its tile colors
type Tetromino struct {
	x, y   int
	shape  TetrominoShape
	color  color.Color
	matrix *mat.Dense
	tiles  map[*Tile]struct{}
}

// NewTetromino creates a new random Tetromino object and its Tiles.
func NewRandomTetromino() *Tetromino {
	randShape := TetrominoShape(rand.Intn(math.MaxInt) % 7)
	return NewTetromino(randShape)
}

// NewTetromino creates a new Tetromino object and its Tiles.
func NewTetromino(shape TetrominoShape) *Tetromino {
	tetro := &Tetromino{
		x:     0,
		y:     0,
		shape: shape,
		tiles: map[*Tile]struct{}{},
	}
	switch shape {
	case ShapeI:
		tetro.initTetrominoI()
	case ShapeJ:
		tetro.initTetrominoJ()
	case ShapeL:
		tetro.initTetrominoL()
	case ShapeS:
		tetro.initTetrominoS()
	case ShapeZ:
		tetro.initTetrominoZ()
	case ShapeT:
		tetro.initTetrominoT()
	case ShapeO:
		tetro.initTetrominoO()
	default:
		panic(fmt.Errorf("unhandled shape: %v", shape))
	}

	return tetro
}

// Update updates the tiles in the Teromino's animation states.
func (t *Tetromino) Update() error {
	for t := range t.tiles {
		if err := t.Update(); err != nil {
			return err
		}
	}
	return nil
}

// Draw draws the tiles in this Tetromino to the given boardImage.
func (t *Tetromino) Draw(boardImage *ebiten.Image) {
	for t := range t.tiles {
		t.Draw(boardImage)
	}
}

// IsOutOfBounds returns true if the position of any tile is out of bounds of the board
func (t *Tetromino) IsOutOfBounds(boardCols, boardRows int) bool {
	for t := range t.tiles {
		x, y := t.Pos()
		if x < 0 || y < 0 {
			return true
		}
		if x >= boardCols || y >= boardRows {
			return true
		}
	}

	return false
}

func (t *Tetromino) IsMoving() bool {
	for t := range t.tiles {
		if t.IsMoving() {
			return true
		}
	}
	return false
}

// ShiftPositionBy shifts the coordinates the Tetromino and all tiles by given delta X/Y
func (t *Tetromino) ShiftPositionBy(dx, dy int) {
	t.x += dx
	t.y += dy
	t.shiftTilesBy(dx, dy)
}

// shiftTilesBy shifts the coordinates of all tiles by given delta X/Y
func (t *Tetromino) shiftTilesBy(dx, dy int) {
	for tile := range t.tiles {
		tile.x += dx
		tile.y += dy
	}
}

// RotatePositionCW rotates the piece ClockWise
func (t *Tetromino) RotatePositionCW() {
	t.matrix = RotateMatrixCW(t.matrix)
	t.updateTilesFromMatrix()

	// put tiles back in position of Tetromino
	t.shiftTilesBy(t.x, t.y)
}

// RotatePositionCCW rotates the piece CounterClockWise
func (t *Tetromino) RotatePositionCCW() {
	t.matrix = RotateMatrixCCW(t.matrix)
	t.updateTilesFromMatrix()

	// put tiles back in position of Tetromino
	t.shiftTilesBy(t.x, t.y)
}

// MatrixDims returns the raw matrix dimensions of the shape
func (t *Tetromino) MatrixDims() (r int, c int) {
	if t.matrix == nil {
		return 0, 0
	}
	return t.matrix.Dims()
}

// MatrixData returns the raw matrix data of the shape
func (t *Tetromino) MatrixData() []float64 {
	if t.matrix == nil {
		return nil
	}
	return t.matrix.RawMatrix().Data
}

// MatrixTilePositions returns only TilePosition data for tiles of the shape
func (t *Tetromino) MatrixTilePositions() map[TilePosition]struct{} {
	rows, cols := t.MatrixDims()
	positions := map[TilePosition]struct{}{}

	rData := t.MatrixData()
	tIndex := 0
	for r := 0; r < rows; r++ {
		// using bottom of matrix as Y origin instead of top
		y := rows - r - 1
		rowMult := r * cols
		for x := 0; x < cols; x++ {
			if rData[rowMult+x] == 1 {
				tPos := TilePosition{x: x, y: y}
				positions[tPos] = struct{}{}
				tIndex++
			}
		}
	}
	return positions
}

func (t *Tetromino) stopAnimation() {
	for t := range t.tiles {
		t.stopAnimation()
	}
}

// updateTiles updates the tiles object map based on the matrix of tile positions
func (t *Tetromino) updateTilesFromMatrix() {
	t.tiles = map[*Tile]struct{}{}
	tPositions := t.MatrixTilePositions()
	for tPos := range tPositions {
		tile := &Tile{
			x:            tPos.x,
			y:            tPos.y,
			tetro:        t,
			poppingCount: maxPoppingCount,
		}
		t.tiles[tile] = struct{}{}
	}
}

// initalize matrix shape as I piece
func (t *Tetromino) initTetrominoI() {
	// Tealish #00A691
	t.color = color.RGBA{0x00, 0xA6, 0x91, 0xff}

	dI := []float64{
		1,
		1,
		1,
		1,
	}
	t.matrix = mat.NewDense(4, 1, dI)
	t.updateTilesFromMatrix()
}

// initalize matrix shape as J piece
func (t *Tetromino) initTetrominoJ() {
	// Purplish #AD42EB
	t.color = color.RGBA{0xAD, 0x42, 0xEB, 0xff}

	dJ := []float64{
		0, 1,
		0, 1,
		1, 1,
	}
	t.matrix = mat.NewDense(3, 2, dJ)
	t.updateTilesFromMatrix()
}

// initalize matrix shape as L piece
func (t *Tetromino) initTetrominoL() {
	// Greenish #68FA3F
	t.color = color.RGBA{0x68, 0xFA, 0x3F, 0xff}

	dL := []float64{
		1, 0,
		1, 0,
		1, 1,
	}
	t.matrix = mat.NewDense(3, 2, dL)
	t.updateTilesFromMatrix()
}

// initalize matrix shape as S piece
func (t *Tetromino) initTetrominoS() {
	// Pinkish #FF2FA8
	t.color = color.RGBA{0xFF, 0x2F, 0xA8, 0xff}

	dS := []float64{
		0, 1, 1,
		1, 1, 0,
	}
	t.matrix = mat.NewDense(2, 3, dS)
	t.updateTilesFromMatrix()
}

// initalize matrix shape as Z piece
func (t *Tetromino) initTetrominoZ() {
	// Cyanish #00FFFA
	t.color = color.RGBA{0x00, 0xFF, 0xFA, 0xff}

	dZ := []float64{
		1, 1, 0,
		0, 1, 1,
	}
	t.matrix = mat.NewDense(2, 3, dZ)
	t.updateTilesFromMatrix()
}

// initalize matrix shape as T piece
func (t *Tetromino) initTetrominoT() {
	// Orangeish #FF8500
	t.color = color.RGBA{0xFF, 0x85, 0x00, 0xff}

	dT := []float64{
		1, 1, 1,
		0, 1, 0,
	}
	t.matrix = mat.NewDense(2, 3, dT)
	t.updateTilesFromMatrix()
}

// initalize matrix shape as Square piece
func (t *Tetromino) initTetrominoO() {
	// Reddish #FF1E1E
	t.color = color.RGBA{0xFF, 0x1E, 0x1E, 0xff}

	dSq := []float64{
		1, 1,
		1, 1,
	}
	t.matrix = mat.NewDense(2, 2, dSq)
	t.updateTilesFromMatrix()
}
