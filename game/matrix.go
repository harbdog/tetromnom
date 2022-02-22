package game

import (
	"gonum.org/v1/gonum/mat"
)

func RotateMatrixCW(m *mat.Dense) *mat.Dense {
	if m == nil {
		return nil
	}

	d := mat.DenseCopyOf(m.T())
	reverseMatrixCols(d)

	return d
}

func RotateMatrixCCW(m *mat.Dense) *mat.Dense {
	if m == nil {
		return nil
	}

	d := mat.DenseCopyOf(m.T())
	reverseMatrixRows(d)

	return d
}

func reverseMatrixCols(m *mat.Dense) {
	rows, cols := m.Dims()

	data := m.RawMatrix().Data
	rData := make([]float64, rows*cols)
	for r := 0; r < rows; r++ {
		rowMult := r * cols
		for c := 0; c < cols; c++ {
			x := cols - c - 1
			rData[rowMult+x] = data[rowMult+c]
		}
	}

	reverse := mat.NewDense(rows, cols, rData)
	m.Copy(reverse)
}

func reverseMatrixRows(m *mat.Dense) {
	rows, cols := m.Dims()

	data := m.RawMatrix().Data
	rData := make([]float64, rows*cols)
	for r := 0; r < rows; r++ {
		y := rows - r - 1
		yMult := y * cols
		rowMult := r * cols
		for c := 0; c < cols; c++ {
			rData[yMult+c] = data[rowMult+c]
		}
	}

	reverse := mat.NewDense(rows, cols, rData)
	m.Copy(reverse)
}
