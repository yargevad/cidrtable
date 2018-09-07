package twentyfour

import (
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

	// iterate over all single paren possibilities
	for i := 0; i < len(c.Numbers)-1; i++ {
		// outer counter indicates which index paren should come before
		for j := i + 1; j < len(c.Numbers); j++ {
			// inner counter indicates which index paren should come after
			current.Str = ""
			current.Parens = []int{i, j}
			str := current.String()
			// skip permutations we've already processed
			if _, ok := c.Seen[str]; ok {
				continue
			}
			c.Seen[str] = true

			current.combo = c
			current.Evaluate()
			c.Results = append(c.Results, *current)

			// FIXME: not generic enough, only supports 2- or 4-digit combos
			// generate two-paren cases, ex: (1+1)*(1+11)
			if i == 0 && j == 1 && len(c.Numbers) >= 4 {
				current.Str = ""
				current.Parens = append(current.Parens, 2, len(c.Numbers)-1)
				current.Evaluate()
				c.Results = append(c.Results, *current)
			}
		}
	}

}

type Eval struct {
	Parens []int
	Str    string // populated lazily
	Float  float64
	Total  int

	// reference to parent necessary to stringify
	combo *Combo
}

func (e *Eval) Evaluate() {
	str := e.String()
	expression, err := govaluate.NewEvaluableExpression(str)
	if err != nil {
		log.Fatalf("error parsing [%s]: %s", str, err)
	}
	result, err := expression.Evaluate(nil)
	if err != nil {
		log.Fatalf("error evaluating [%s]: %s", str, err)
	}

	switch val := result.(type) {
	case int:
		e.Float = float64(val)
		e.Total = val
	case float64:
		e.Float = val
		e.Total = int(math.Round(val))
	}
}

// String transforms an `Eval` into a human-readable format
func (e *Eval) String() string {
	if e.combo == nil {
		return ""
	} else if e.Str == "" {
		expression := strings.Builder{}
		numlen := len(e.combo.Numbers)
		for n := 0; n < numlen; n++ {
			// TODO: make this support more than 2 sets of parens
			if (len(e.Parens) >= 2 && e.Parens[0] == n) || (len(e.Parens) == 4 && e.Parens[2] == n) {
				expression.WriteString("(")
			}
			expression.WriteString(strconv.Itoa(e.combo.Numbers[n]))
			if (len(e.Parens) >= 2 && e.Parens[1] == n) || (len(e.Parens) == 4 && e.Parens[3] == n) {
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
