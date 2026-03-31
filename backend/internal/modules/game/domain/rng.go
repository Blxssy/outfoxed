package domain

import (
	crand "crypto/rand"
	"github.com/rs/zerolog/log"
	"math/big"
	"math/rand"
)

type RNG interface {
	Intn(n int) int
}

type StdRNG struct {
	r *rand.Rand
}

func NewStdRNG(r *rand.Rand) StdRNG {
	return StdRNG{r: r}
}
func (StdRNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}

	v, err := crand.Int(crand.Reader, big.NewInt(int64(n)))
	if err != nil {
		log.Error().Err(err).Msg("Error generating random number")
		return 0
	}
	return int(v.Int64())
}

// FixedRNG для тестов: возвращает заранее заданные значения по кругу
type FixedRNG struct {
	Values []int
	i      int
}

func (f *FixedRNG) Intn(n int) int {
	if len(f.Values) == 0 {
		return 0
	}
	v := f.Values[f.i%len(f.Values)]
	f.i++

	if n <= 0 {
		return 0
	}
	if v < 0 {
		v = -v
	}
	return v % n
}
