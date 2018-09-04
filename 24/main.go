package main

import (
	"flag"
	"log"
	"math"
	"strconv"
	"strings"
)

// Operator represents an allowed mathematical operation
type Operator int

// String converts an Operator into a human-readable format
func (o Operator) String() string {
	switch o {
	case Add:
		// TODO: skip dupes due to commutative property
		return "+"
	case Subtract:
		return "-"
	case Multiply:
		// TODO: skip dupes due to commutative property
		return "*"
	case Divide:
		return "/"
	default:
		return "???"
	}
}

const (
	Add Operator = iota
	Subtract
	Multiply
	Divide
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

// Repetitions calculates permutations with reuse allowed, so the set
// `(a, b, c, d)` can yield `(a, a, a, a)` as a valid result.
func Repetitions(values []int, length int) [][]int {
	// https://rosettacode.org/wiki/Permutations_with_repetitions#Go
	rv := [][]int{}
	inLen := len(values)
	outLen := length

	indexes := make([]int, outLen)

	for {
		outputs := make([]int, outLen)
		// generate permutaton
		for i, x := range indexes {
			outputs[i] = values[x]
		}
		log.Printf("appending %+v\n", outputs)
		rv = append(rv, outputs)

		// increment permutation number
		for i := 0; ; {
			// increment current index
			indexes[i]++
			// run outer loop again if we're still in bounds
			if indexes[i] < inLen {
				break
			}
			// otherwise, reset current index and move on
			indexes[i] = 0
			i++
			if i == outLen {
				return rv // all permutations generated
			}
		}
	}
}

// Permutations is an implementation of Heap's algorithm
func Permutations(arr []int) [][]int {
	// https://en.wikipedia.org/wiki/Heap%27s_algorithm
	// https://stackoverflow.com/questions/30226438/generate-all-permutations-in-go
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func Evaluate(nums, ops []int) int {
	// set initial state to first number
	var tmp = float64(nums[0])
	// iterate over operators, applying next number
	// assumes length of nums is always greater than length of ops
	for n := 0; n < len(ops); n++ {
		switch Operator(ops[n]) {
		case Add:
			tmp = tmp + float64(nums[n+1])
		case Subtract:
			tmp = tmp - float64(nums[n+1])
		case Multiply:
			tmp = tmp * float64(nums[n+1])
		case Divide:
			tmp = tmp / float64(nums[n+1])
		}
	}
	// whole numbers only
	if math.Floor(tmp) != tmp {
		return -1
	}
	return int(tmp)
}

func main() {
	integers := &Integers{}
	flag.Var(integers, "n", "int to include (multiple)")

	verbose := flag.Bool("verbose", false, "verbose logging")
	target := flag.Int("target", 24, "desired result")

	flag.Parse()

	count := len(*integers)
	if count < 2 {
		log.Fatalf("must specify at least 2 numbers with --n (got %d)", count)
	}

	if *verbose {
		log.Printf("combining integers %s with target %d...\n", integers, *target)
	}

	// numbers are single-use, taken from command line
	numbers := Permutations(*integers)
	if *verbose {
		log.Printf("found %d number permutations (no repetition)\n", len(numbers))
	}

	// operators are multiple-use, taken from constants defined above (values 0-3)
	operators := Repetitions([]int{0, 1, 2, 3}, len(*integers)-1)
	if *verbose {
		log.Printf("found %d operator permutations (with repetition)\n%+v\n", len(operators), operators)
	}

	return

	for _, nums := range numbers {
		for _, ops := range operators {
			result := Evaluate(nums, ops)
			if *verbose || result == *target {
				expression := strings.Builder{}
				for n := 0; n < count; n++ {
					expression.WriteString(strconv.Itoa(nums[n]))
					if n < (count - 1) {
						expression.WriteString(Operator(ops[n]).String())
					}
				}
				if result == *target {
					log.Printf("%s = %d !!!", expression.String(), result)
				} else if *verbose {
					log.Printf("%s = %d", expression.String(), result)
				}
			}
		}

	}

}
