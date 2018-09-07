package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	tf "github.com/yargevad/golang/twentyfour"
)

// Integers stores multiple values passed in on the command line
type Integers []int

// String converts our array of integers into a human-readable format
func (integers *Integers) String() string {
	buf := &strings.Builder{}
	buf.WriteString("(")
	for index, integer := range *integers {
		buf.WriteString(strconv.Itoa(integer))
		if index < len(*integers)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(")")
	return buf.String()
}

// Set takes the string value from the command line and appends it
func (integers *Integers) Set(val string) error {
	// https://stackoverflow.com/questions/28322997/how-to-get-a-list-of-values-into-a-flag-in-golang
	str, err := strconv.Atoi(val)
	if err != nil {
		return err
	}

	*integers = append(*integers, str)
	return nil
}

func main() {
	integers := &Integers{}
	flag.Var(integers, "n", "int to include (multiple)")

	variance := flag.Int("variance", 5, "allow floating-point math errors after this many decimal places")
	verbose := flag.Bool("verbose", false, "verbose logging")
	target := flag.Int("target", 24, "desired result")

	flag.Parse()

	offBy := 1 / math.Pow(10, float64(*variance))

	count := len(*integers)
	if count < 2 {
		log.Fatalf("must specify at least 2 numbers with --n (got %d)", count)
	}

	if *verbose {
		log.Printf("combining integers %s with target %d...\n", integers, *target)
	}

	// numbers are single-use, taken from command line
	numbers := tf.Permutations(*integers)
	if *verbose {
		log.Printf("found %d number permutations (no repetition)\n%+v\n", len(numbers), numbers)
	}

	// operators are multiple-use, taken from constants defined above (values 0-3)
	operators := tf.Repetitions([]int{0, 1, 2, 3}, len(*integers)-1)
	if *verbose {
		log.Printf("found %d operator permutations (with repetition)\n%+v\n", len(operators), operators)
	}

	seen := map[string]bool{}
	target_f := float64(*target)

	// for each permutation of numbers
	for _, nums := range numbers {
		// for each permutation of operators
		for _, ops := range operators {
			combo := &tf.Combo{Numbers: nums, Operators: ops, Seen: seen}
			combo.Evaluate()
			for _, result := range combo.Results {
				// can't directly compare (rounding errors), set acceptable error bars
				if result.Float >= (target_f-offBy) && result.Float <= (target_f+offBy) {
					fmt.Printf("%s = %d\n", result.String(), result.Total)
				} else if *verbose {
					log.Printf("%s = %d (%.09f)", result.String(), result.Total, result.Float)
				}
			}
		}

	}

}
