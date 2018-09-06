package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/knetic/govaluate"
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

type Combo struct {
	Numbers   []int
	Operators []int
	Results   []Eval
	Seen      map[string]bool
}

// Evaluate performs the calculation(s) given by interleaving `nums` and `ops`.
// Operator precedence and any parenthesis placement are taken into account.
// Any existing `Results` are cleared before calculation.
// The results are placed into the aptly-named `Results` field.
func (c *Combo) Evaluate() {
	c.Results = nil
	current := &Eval{combo: c}

	// iterate over all paren possibilities
	for i := 0; i < len(c.Numbers)-1; i++ {
		// outer counter indicates which index paren should come before
		for j := i + 1; j < len(c.Numbers); j++ {
			// inner counter indicates which index paren should come after
			current.Parens = []int{i, j}
			str := current.String()
			// skip permutations we've already processed
			if _, ok := c.Seen[current.String()]; ok {
				continue
			}
			c.Seen[current.String()] = true

			expression, err := govaluate.NewEvaluableExpression(str)
			if err != nil {
				log.Fatalf("error parsing [%s]: %s", str, err)
			}
			result, err := expression.Evaluate(nil)
			if err != nil {
				log.Fatalf("error evaluating [%s]: %s", str, err)
			}
			current.combo = c
			switch val := result.(type) {
			case int:
				current.Float = float64(val)
				current.Total = val
			case float64:
				current.Float = val
				current.Total = int(math.Round(val))
			}
			c.Results = append(c.Results, *current)
			current.Str = ""
		}
	}

}

type Eval struct {
	Parens []int
	Str    string // populated lazily
	Float  float64
	Total  int

	// makes it easier to stringify
	combo *Combo
}

// String transforms an `Eval` into a human-readable format
func (e *Eval) String() string {
	if e.combo == nil {
		return ""
	} else if e.Str == "" {
		expression := strings.Builder{}
		numlen := len(e.combo.Numbers)
		for n := 0; n < numlen; n++ {
			if len(e.Parens) == 2 && e.Parens[0] == n {
				expression.WriteString("(")
			}
			expression.WriteString(strconv.Itoa(e.combo.Numbers[n]))
			if len(e.Parens) == 2 && e.Parens[1] == n {
				expression.WriteString(")")
			}
			if n < (numlen - 1) {
				expression.WriteString(Operator(e.combo.Operators[n]).String())
			}
		}
		e.Str = expression.String()
	}
	return e.Str
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
	numbers := Permutations(*integers)
	if *verbose {
		log.Printf("found %d number permutations (no repetition)\n%+v\n", len(numbers), numbers)
	}

	// operators are multiple-use, taken from constants defined above (values 0-3)
	operators := Repetitions([]int{0, 1, 2, 3}, len(*integers)-1)
	if *verbose {
		log.Printf("found %d operator permutations (with repetition)\n%+v\n", len(operators), operators)
	}

	seen := map[string]bool{}
	target_f := float64(*target)

	// for each permutation of numbers
	for _, nums := range numbers {
		// for each permutation of operators
		for _, ops := range operators {
			combo := &Combo{Numbers: nums, Operators: ops, Seen: seen}
			combo.Evaluate()
			for _, result := range combo.Results {
				// skip display of dupes

				// both of these work, use the one that is easier to read
				//if (math.Round(result.Float/offBy) * offBy) == target_f {
				if result.Float >= (target_f-offBy) && result.Float <= (target_f+offBy) {
					fmt.Printf("%s = %d\n", result.String(), result.Total)
				} else if *verbose {
					log.Printf("%s = %d (%.09f)", result.String(), result.Total, result.Float)
				}
			}
		}

	}

}
