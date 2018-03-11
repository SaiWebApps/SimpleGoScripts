package intmatrix

import (
	"bytes"
	"fmt"
	"strconv"
)

// IntMatrix2D represents a 2D matrix with integers.
type IntMatrix2D struct {
	mat          [][]int
	nrows, ncols int
}

// New2D initializes and returns a nr x nc IntMatrix2D.
func New2D(nr, nc int) (*IntMatrix2D, error) {
	if nr < 0 || nc < 0 {
		return nil, fmt.Errorf("[%dx%d Matrix]: Invalid dimensions", nr, nc)
	}
	im2d := new(IntMatrix2D)
	im2d.nrows, im2d.ncols = nr, nc
	im2d.mat = make([][]int, nr)
	for i := 0; i < nc; i++ {
		im2d.mat[i] = make([]int, nc)
	}
	return im2d, nil
}

// Eye initializes and returns a n x n identity matrix.
func Eye(n int) (*IntMatrix2D, error) {
	im2d, err := New2D(n, n)
	if err != nil {
		return nil, err
	}
	// Set the diagonals to 1.
	for i := 0; i < n; i++ {
		im2d.mat[i][i] = 1
	}
	return im2d, nil
}

// FromSlice returns a IntMatrix2D instance with the values in a given slice.
func FromSlice(values []int, nr, nc int) (*IntMatrix2D, error) {
	// If values contains more ints than we can store in a nr x nc
	// matrix, then return an error.
	expectedLen, actualLen := nr*nc, len(values)
	if actualLen > expectedLen {
		return nil, fmt.Errorf("[%dx%d Matrix]: Cannot store %d values",
			nr, nc, actualLen)
	}

	// Otherwise, load the ints in values into a nr x nc matrix.
	im2d, _ := New2D(nr, nc)
	for idx, val := range values {
		r, c := idx/nc, idx%nc
		im2d.mat[r][c] = val

	}
	return im2d, nil
}

// GetDimensions returns the dimensions of this IntMatrix2D instance.
func (m IntMatrix2D) GetDimensions() (int, int) {
	return m.nrows, m.ncols
}

// At returns the integer at the target location in this IntMatrix2D instance.
func (m IntMatrix2D) At(i, j int) (int, error) {
	if i >= m.nrows || j >= m.ncols {
		return 0, fmt.Errorf("[%dx%d Matrix]: Invalid Location (%d,%d)",
			m.nrows, m.ncols, i, j)
	}
	return m.mat[i][j], nil
}

// Set sets this IntMatrix2D instance's value at (i,j) to val and returns the
// old value at that location.
func (m *IntMatrix2D) Set(i, j int, val int) (int, error) {
	if m == nil {
		return 0, fmt.Errorf("[Nil 2D Matrix]: Cannot set (%d,%d)", i, j)
	}
	origVal, err := m.At(i, j)
	if err != nil {
		m.mat[i][j] = val
	}
	return origVal, err
}

// Transpose returns the transposed version of this IntMatrix2D.
func (m IntMatrix2D) Transpose() *IntMatrix2D {
	transposed, _ := New2D(m.ncols, m.nrows)
	for i := 0; i < m.nrows; i++ {
		for j := 0; j < m.ncols; j++ {
			currentVal, _ := m.At(i, j)
			transposed.Set(j, i, currentVal)
		}
	}
	return transposed
}

// Reshape returns this IntMatrix2D's values in a new
// newNRows x newNCols IntMatrix2D.
func (m IntMatrix2D) Reshape(newNRows, newNCols int) (*IntMatrix2D, error) {
	// Both input and output should ultimately have the same number of items.
	if m.nrows*m.ncols != newNRows*newNCols {
		return nil, fmt.Errorf("[%dx%d Matrix]: Cannot reshape into %dx%d",
			m.nrows, m.ncols, newNRows, newNCols)
	}
	// Allocate a newNRows x newNCols IntMatrix2D.
	reshaped, err := New2D(newNRows, newNCols)
	if err != nil {
		return nil, err
	}
	// row-major-ordering: INDEX = ROW * NCOLS + COL
	// Given INDEX: ROW = INDEX / NCOLS, COL = INDEX % NCOLS
	for i := 0; i < newNRows; i++ {
		for j := 0; j < newNCols; j++ {
			// Calculate row-major-ordering (linear) index. This will be
			// the same in both the original and reshaped matrix.
			linearIdx := i*newNCols + j
			// Convert linearIdx into a (r,c) pair for the original matrix.
			// A.k.a - Find (r,c) in m corresponding to (i,j) in reshaped.
			origRow, origCol := linearIdx/m.ncols, linearIdx%m.ncols
			reshaped.mat[i][j] = m.mat[origRow][origCol]
		}
	}
	return reshaped, nil
}

// Plus returns the sum of the values in two IntMatrix2D instances.
func (m IntMatrix2D) Plus(other IntMatrix2D) *IntMatrix2D {
	if m.nrows != other.nrows || m.ncols != other.ncols {
		return nil
	}
	sum, _ := New2D(m.nrows, m.ncols)
	for i := 0; i < m.nrows; i++ {
		for j := 0; j < m.ncols; j++ {
			sum.mat[i][j] = m.mat[i][j] + other.mat[i][j]
		}
	}
	return sum
}

// PlusEquals adds this IntMatrix2D's values to another's in-place.
func (m *IntMatrix2D) PlusEquals(other IntMatrix2D) error {
	if m.nrows != other.nrows || m.ncols != other.ncols {
		return fmt.Errorf("[%dx%d Matrix]: Cannot add to %dx%d Matrix",
			m.nrows, m.ncols, other.nrows, other.ncols)
	}
	for i := 0; i < m.nrows; i++ {
		for j := 0; j < m.ncols; j++ {
			m.mat[i][j] += other.mat[i][j]
		}
	}
	return nil
}

// Times returns the product of two IntMatrix2D instances.
func (m IntMatrix2D) Times(other IntMatrix2D) *IntMatrix2D {
	// When multipling 2 matrices, [m1xn1] and [m2xn2], n1 shall equal m2.
	if m.ncols != other.nrows {
		return nil
	}

	product, _ := New2D(m.nrows, other.ncols)
	otherT := other.Transpose() // [m2xn2] Matrix -> [n2xm2] Matrix
	rowLen := otherT.ncols      // m2 == n1
	// Each location (i,j) in product corresponds to the dot product of row
	// row i in m and row j in otherT.
	for i := 0; i < product.nrows; i++ {
		mRow := m.mat[i]
		for j := 0; j < product.ncols; j++ {
			otherTRow := otherT.mat[j]
			for k := 0; k < rowLen; k++ {
				product.mat[i][j] += mRow[k] * otherTRow[k]
			}
		}
	}
	return product
}

func (m IntMatrix2D) String() string {
	var buffer bytes.Buffer
	for i := 0; i < m.nrows; i++ {
		for j := 0; j < m.ncols; j++ {
			val, _ := m.At(i, j)
			buffer.WriteString(strconv.Itoa(val))
			// Separate each column value with a space.
			if j < m.ncols-1 {
				buffer.WriteString(" ")
			}
		}
		// Display each row on its own line.
		buffer.WriteString("\n")
	}
	buffer.WriteString(fmt.Sprintf("[%d rows x %d columns]", m.nrows, m.ncols))
	return buffer.String()
}
