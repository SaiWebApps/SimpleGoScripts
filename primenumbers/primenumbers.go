package main

import (
	"flag"
	"fmt"
)

// NegativeNumberError is raised when user tries to
// compute primes for a negative number.
type NegativeNumberError int

func (err NegativeNumberError) Error() string {
	const MSG = "Cannot calculate primes for a negative number"
	return fmt.Sprintf("%s (%d)!", MSG, err)
}

// GetPrimesUpTo returns a slice with all primes up to n inclusive.
// If n is negative, then it returns an empty slice and a NegativeNumberError.
func GetPrimesUpTo(n int) ([]int, error) {
	var primes []int
	if n < 0 {
		return primes, NegativeNumberError(n)
	}
	if n <= 1 {
		return primes, nil
	}

	// Each index i in the sieve corresponds to the integer (i+2).
	// If the boolean sieve[i] is true, then (i+2) is composite.
	sieve := make([]bool, n-1)
	for i := 2; i <= n; i++ {
		// Otherwise, it is prime.
		if !sieve[i-2] {
			primes = append(primes, i)
		}
		// Now, mark all multiples of (i+2) as composites.
		for j := 2 * i; j <= n; j += i {
			sieve[j-2] = true
		}
	}
	return primes, nil
}

func main() {
	n := flag.Int("n", 1, "Calculate primes up to this number inclusive.")
	flag.Parse()
	primes, err := GetPrimesUpTo(*n)
	switch err != nil {
	case true:
		fmt.Println(err)
	default:
		fmt.Printf("Prime Numbers Up To %d: %v", *n, primes)
	}
}
