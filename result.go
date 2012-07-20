package corpustools

import (
	"sort"
)

// Result is a general container structure for assigning quantities to sequences so that the sequences can be sorted by their value.
type Result struct {
	Seq []int
	Val float64
}

type Results []Result

func (r Results) Len() int {
	return len(r)
}

func (r Results) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Results) Less(i, j int) bool {
	return r[i].Val < r[j].Val
}

// Allows Results to be reverse sorted (e.g. "sort.Sort(ResultsReverseSort{results})").
type ResultsReverseSort struct {
	sort.Interface
}

func (r ResultsReverseSort) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}
