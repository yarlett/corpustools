package corpustools

import (
	"math"
	"sort"
)

type Cooc struct {
	seq []int
	dat map[int]float64
}

func (cooc *Cooc) Set(key int, val float64) {
	cooc.dat[key] = val
}

func (cooc *Cooc) Inc(key int) {
	cooc.dat[key] += 1.0
}

func (cooc *Cooc) Keys() (keys []int) {
	keys = make([]int, len(cooc.dat))
	i := 0
	for k, _ := range cooc.dat {
		keys[i] = k
		i++
	}
	sort.IntSlice(keys).Sort()
	return
}

func (cooc *Cooc) Mag() (mag float64) {
	for _, val := range cooc.dat {
		mag += val * val
	}
	mag = math.Pow(mag, 0.5)
	return
}

func (cooc1 *Cooc) Prod(cooc2 *Cooc) (prod float64) {
	if len(cooc1.dat) < len(cooc2.dat) {
		for k1, v1 := range cooc1.dat {
			v2, found := cooc2.dat[k1]
			if found {
				prod += v1 * v2
			}
		}
	} else {
		for k2, v2 := range cooc2.dat {
			v1, found := cooc1.dat[k2]
			if found {
				prod += v1 * v2
			}
		}
	}
	return
}
