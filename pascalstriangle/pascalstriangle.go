package main

import (
	"bytes"
	"flag"
	"fmt"
	"strconv"
)

// PascalsTriangle represents N-level Pascal's Triangle.
type PascalsTriangle struct {
	numLevels int
	triangle  [][]int
}

// Helper/Private Methods
func (t *PascalsTriangle) init() {
	// Calculate underlying matrix's dimensions.
	numRows := t.numLevels
	numCols := 2*numRows + 1

	// Initialize underlying matrix accordingly.
	t.triangle = make([][]int, numRows)
	for r := 0; r < numRows; r++ {
		t.triangle[r] = make([]int, numCols)
	}
	// 1st row shall contain 1 in middle column.
	t.triangle[0][numCols/2] = 1
}

// Generate shall create a n-level Pascal's triangle or return
// an InvalidPTDimError if n is not positive.
func (t *PascalsTriangle) Generate(n int) error {
	if n <= 0 {
		return fmt.Errorf("[NegativeNumLevels]: %d", n)
	}

	t.numLevels = n
	t.init()

	for r := 1; r < t.numLevels; r++ {
		midpoint := len(t.triangle[r]) / 2
		start := midpoint - r
		end := midpoint + r

		// Each row shall start and end with 1.
		t.triangle[r][midpoint-r] = 1
		t.triangle[r][midpoint+r] = 1

		// Every 2 spaces between start and end inclusive shall be
		// the sum of the top left and top right values from that space.
		for c := start + 2; c < end-1; c += 2 {
			topLeft := t.triangle[r-1][c-1]
			topRight := t.triangle[r-1][c+1]
			t.triangle[r][c] = topLeft + topRight
		}
	}

	return nil
}

func (t PascalsTriangle) String() string {
	var buffer bytes.Buffer
	triangle := t.triangle
	numRows := t.numLevels

	const HEADER = "Printing Pascals Triangle with %d levels\n"
	buffer.WriteString(fmt.Sprintf(HEADER, numRows))

	for r := 0; r < numRows; r++ {
		numCols := len(triangle[r])

		for c := 0; c < numCols; c++ {
			val := triangle[r][c]
			if val == 0 {
				buffer.WriteString(" ")
				continue
			}
			buffer.WriteString(strconv.Itoa(val))
		}
		buffer.WriteString("\n")
	}

	return buffer.String()
}

// Main Function
func main() {
	var numLevels int
	flag.IntVar(&numLevels, "n", 1, "Number of levels in PascalsTriangle.")
	flag.Parse()

	pt := &PascalsTriangle{}
	err := pt.Generate(numLevels)
	switch {
	case (err != nil): // There was an error.
		fmt.Println(err)
	default: // No errors encountered
		fmt.Print(pt)
	}
}
