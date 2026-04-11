package domain

import "math/rand"

type RNG interface {
	Intn(n int) int
}

type StdRNG struct {
	r *rand.Rand
}

func NewStdRNG(r *rand.Rand) *StdRNG {
	return &StdRNG{r: r}
}

func (s *StdRNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	return s.r.Intn(n)
}

// FixedRNG для тестов: возвращает заранее заданные значения по кругу.
type FixedRNG struct {
	Values []int
	i      int
}

func NewFixedRNG(values ...int) *FixedRNG {
	return &FixedRNG{Values: values}
}

func (f *FixedRNG) Intn(n int) int {
	if len(f.Values) == 0 || n <= 0 {
		return 0
	}

	v := f.Values[f.i%len(f.Values)]
	f.i++

	if v < 0 {
		v = -v
	}
	return v % n
}
