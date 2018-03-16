package main

import (
	"flag"
	"fmt"
	"simplegoscripts/intmatrix"
)

func initPascalsTriangle(numLevels int) (*intmatrix.IntMatrix2D, error) {
	numCols := 2*(numLevels-1) + 1
	pt, err := intmatrix.New2D(numLevels, numCols)
	if err != nil {
		return pt, err
	}
	_, err = pt.Set(0, numCols/2, 1)
	return pt, err
}

// Generate shall create a Pascal's Triangle with numLevels-levels
// for numLevels >= 1.
func Generate(numLevels int) (*intmatrix.IntMatrix2D, error) {
	if numLevels < 1 {
		return nil, fmt.Errorf("[NegativeNumLevels]: %d", numLevels)
	}
	pt, err := initPascalsTriangle(numLevels)
	if err != nil {
		return pt, err
	}
	numRows, numCols := pt.GetDimensions()

	mid := numCols / 2
	for r := 1; r < numRows; r++ {
		start, end := mid-r, mid+r
		pt.Set(r, start, 1)
		pt.Set(r, end, 1)
		for c := start + 2; c < end; c += 2 {
			diagLeft, _ := pt.At(r-1, c-1)
			diagRight, _ := pt.At(r-1, c+1)
			pt.Set(r, c, diagLeft+diagRight)
		}
	}
	return pt, nil
}

// Main Function
func main() {
	var numLevels int
	flag.IntVar(&numLevels, "n", 1, "Number of levels in PascalsTriangle.")
	flag.Parse()

	pt, err := Generate(numLevels)
	switch {
	case err != nil:
		fmt.Println(err)
	default:
		fmt.Println(pt)
	}
}
