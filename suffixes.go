package corpustools

// Suffix struct helps with the constructions of a suffix array for a corpus.
type Suffix struct {
	corpus_ptr *[]int
	corpus_ix  int
}

// List of suffixes.
type Suffixes []*Suffix

// Define Len() and Swap() and Less() on elements of Suffixes so that Suffixes can be sorted into order, yielding indices of suffix array.
func (s Suffixes) Len() int      { return len(s) }
func (s Suffixes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type BySuffix struct{ Suffixes }
func (s BySuffix) Less(i, j int) bool {
	corpusi := []int(*s.Suffixes[i].corpus_ptr)
	ixi := s.Suffixes[i].corpus_ix
	corpusj := []int(*s.Suffixes[j].corpus_ptr)
	ixj := s.Suffixes[j].corpus_ix
	// Find the first point at which the suffixes differ.
	for ; (ixi < len(corpusi)) && (ixj < len(corpusj)) && (corpusi[ixi] == corpusj[ixj]); ixi, ixj = ixi+1, ixj+1 {
	}
	// If we're within bounds, the lesser is the suffix with the lower value. If we're out of bounds the shorter suffix is the lesser.
	if ixi < len(corpusi) && ixj < len(corpusj) {
		return corpusi[ixi] < corpusj[ixj]
	} else if ixj < len(corpusj) {
		return true
	}
	return false
}