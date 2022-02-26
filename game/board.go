package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	backgroundColor = color.RGBA{0xF2, 0xF2, 0xF2, 0xff} // Light Gray #F2F2F2
	boardColor      = color.RGBA{0x3D, 0x3C, 0x3D, 0xff} // Dark Gray #3D3C3D
)

// Board represents the game board.
type Board struct {
	xSize, ySize            int
	dropCount, maxDropCount int
	currentPiece            *Tetromino
	pieces                  map[*Tetromino]struct{}
	gameOn                  bool
}

// NewBoard generates a new Board with giving a size.
func NewBoard(xSize, ySize int) (*Board, error) {
	b := &Board{
		xSize:        xSize,
		ySize:        ySize,
		maxDropCount: 20, // TODO: implement difficulty levels to have different cell frame drop counts
		pieces:       map[*Tetromino]struct{}{},
		gameOn:       true,
	}

	b.currentPiece = NewRandomTetromino()

	return b, nil
}

// Update updates the board state (called 60 times a second)
func (b *Board) Update(input *Input) error {
	for t := range b.pieces {
		if err := t.Update(); err != nil {
			return err
		}
	}

	if b.currentPiece != nil {
		// check if the current piece is new and tiles not yet added to the board
		pieceOnBoard := false
		if _, ok := b.pieces[b.currentPiece]; ok {
			pieceOnBoard = true
		}

		if pieceOnBoard {
			if b.dropCount > 0 {
				b.dropCount--
			} else {
				// perform downward movement if able after check for collision
				canMoveDown := b.checkValidMove(b.currentPiece, 0, -1)
				if canMoveDown {
					b.currentPiece.ShiftPositionBy(0, -1)
					b.dropCount = b.maxDropCount
				} else {
					b.currentPiece = nil
				}
			}
		} else {
			// place new piece at top center to start its descent
			_, pieceCols := b.currentPiece.MatrixDims()
			pieceX := (b.xSize / 2) - (pieceCols / 2)
			pieceY := b.ySize
			b.currentPiece.ShiftPositionBy(pieceX, pieceY)

			b.pieces[b.currentPiece] = struct{}{}
			b.dropCount = b.maxDropCount
		}
	} else {
		// check for game over conditions
		for t := range b.pieces {
			if t.IsOutOfBounds(b.xSize, b.ySize) {
				b.gameOn = false
				return nil
			}
		}

		// start the next piece
		b.currentPiece = NewRandomTetromino()
	}

	if dir, ok := input.Dir(); ok {
		if err := b.Adjust(dir); err != nil {
			return err
		}
	}
	if drop := input.Drop(); drop {
		if err := b.Drop(); err != nil {
			return err
		}
	}
	return nil
}

// checkValidMove checks if any one tile in the Tetromino will collide if it were to move at an x/y offset
func (b *Board) checkValidMove(t *Tetromino, xOffset, yOffset int) bool {
	valid := true

	bMap := b.getPlacedMap()
	for tile := range t.tiles {
		x, y := tile.Pos()
		xCheck := x + xOffset
		yCheck := y + yOffset
		if yCheck >= len(bMap) && xCheck >= 0 && xCheck < len(bMap[0]) {
			// since piece starts from above the top it can be valid when above the board
			continue
		}
		if xCheck < 0 || yCheck < 0 || xCheck >= len(bMap[yCheck]) || bMap[yCheck][xCheck] != 0 {
			valid = false
			break
		}
	}

	return valid
}

// Adjust makes the current piece perform horizontal movements or rotations
func (b *Board) Adjust(dir Dir) error {
	if b.currentPiece == nil {
		return nil
	}
	for t := range b.pieces {
		t.stopAnimation()
	}

	vx, vy := dir.Vector()
	if vx != 0 {
		fmt.Printf("Adjust direction %v\n", dir)
		if b.checkValidMove(b.currentPiece, vx, 0) {
			b.currentPiece.ShiftPositionBy(vx, 0)
		}
	}

	if vy != 0 {
		// check if rotation results in valid move, if not then rotate again same direction until valid
		if vy > 0 {
			fmt.Printf("Rotate piece CCW\n")
			b.currentPiece.RotatePositionCCW()
		} else {
			fmt.Printf("Rotate piece CW\n")
			b.currentPiece.RotatePositionCW()
		}

		_, cols := b.currentPiece.MatrixDims()

		if !b.checkValidMove(b.currentPiece, 0, 0) {
			// only for 90 degree moves, check to see if moving the piece away from nearby pieces/wall helps
			xCheck := 1
			for !b.checkValidMove(b.currentPiece, -xCheck, 0) {
				if xCheck >= cols-1 {
					break
				}
				xCheck++
			}
			if b.checkValidMove(b.currentPiece, -xCheck, 0) {
				b.currentPiece.ShiftPositionBy(-xCheck, 0)
			} else {
				// no valid move, rotate back the other way
				if vy > 0 {
					b.currentPiece.RotatePositionCW()
				} else {
					b.currentPiece.RotatePositionCCW()
				}
			}
		}
	}

	return nil
}

// Drop makes the current piece drop straight down without futher ability to perform adjustments
func (b *Board) Drop() error {
	if b.currentPiece == nil {
		return nil
	}
	for t := range b.pieces {
		t.stopAnimation()
	}
	fmt.Printf("Dropping piece straight down\n")

	for dy := 0; dy < b.ySize; dy++ {
		if b.checkValidMove(b.currentPiece, 0, -1) {
			b.currentPiece.ShiftPositionBy(0, -1)
		} else {
			break
		}
	}

	// reset drop counter and clear current piece
	b.dropCount = 0
	b.currentPiece = nil

	return nil
}

// printBoard prints the board for debugging purposes
func (b *Board) printBoard() {
	bMap := b.getPlacedMap()
	for y := len(bMap) - 1; y >= 0; y-- {
		rowStr := ""
		for _, cell := range bMap[y] {
			rowStr += fmt.Sprintf(" %v", cell)
		}
		fmt.Printf("%s\n", rowStr)
	}
}

// getPlacedMap returns two dimensional array of 0s and 1s, where 1s are cells taken by a placed tile
func (b *Board) getPlacedMap() (boardMap [][]int) {
	boardMap = make([][]int, b.ySize)
	for y := range boardMap {
		boardMap[y] = make([]int, b.xSize)
	}

	for t := range b.pieces {
		if t == b.currentPiece {
			// skip tiles in the current piece since not yet placed
			continue
		}
		for tile := range t.tiles {
			x, y := tile.Pos()
			boardMap[y][x] = 1
		}
	}
	return boardMap
}

// Size returns the board size.
func (b *Board) Size() (int, int) {
	x := b.xSize*tileSize + (b.xSize+1)*tileMargin
	y := b.ySize*tileSize + (b.ySize+1)*tileMargin
	return x, y
}

// Draw draws the board to the given boardImage.
func (b *Board) Draw(boardImage *ebiten.Image) {
	boardImage.Fill(boardColor)
	for row := 0; row < b.ySize; row++ {
		for col := 0; col < b.xSize; col++ {
			op := &ebiten.DrawImageOptions{}
			x := col*tileSize + (col+1)*tileMargin
			y := row*tileSize + (row+1)*tileMargin
			op.GeoM.Translate(float64(x), float64(y))
			r, g, b, a := colorToScale(boardColor)
			op.ColorM.Scale(r, g, b, a)
			boardImage.DrawImage(tileImage, op)
		}
	}
	animatingTiles := map[*Tile]struct{}{}
	nonAnimatingTiles := map[*Tile]struct{}{}
	for t := range b.pieces {
		if t.IsMoving() {
			for tile := range t.tiles {
				animatingTiles[tile] = struct{}{}
			}
		} else {
			for tile := range t.tiles {
				nonAnimatingTiles[tile] = struct{}{}
			}
		}
	}
	for t := range nonAnimatingTiles {
		t.Draw(boardImage)
	}
	for t := range animatingTiles {
		t.Draw(boardImage)
	}

	if !b.gameOn {
		w, _ := b.Size()
		ebitenutil.DebugPrintAt(boardImage, "Game Over", w/2-tileSize, tileSize)
	}
}
