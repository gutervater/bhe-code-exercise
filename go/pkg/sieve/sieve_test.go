package sieve

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNthPrime(t *testing.T) {
	sieve := NewSieve(SieveInitOptions{})
	assert.Equal(t, int64(0), sieve.NthPrime(-1))
	assert.Equal(t, int64(2), sieve.NthPrime(0))
	assert.Equal(t, int64(71), sieve.NthPrime(19))
	assert.Equal(t, int64(541), sieve.NthPrime(99))
	assert.Equal(t, int64(3581), sieve.NthPrime(500))
	assert.Equal(t, int64(7793), sieve.NthPrime(986))
	assert.Equal(t, int64(17393), sieve.NthPrime(2000))
	assert.Equal(t, int64(15485867), sieve.NthPrime(1000000))
	assert.Equal(t, int64(179424691), sieve.NthPrime(10000000))
	assert.Equal(t, int64(2038074751), sieve.NthPrime(100000000))
	// repeated calls should be fast
	startTime := time.Now()
	assert.Equal(t, int64(2038074751), sieve.NthPrime(100000000))
	assert.LessOrEqual(t, time.Since(startTime).Milliseconds(), int64(50), "Cached call took too long")
}

func FuzzNthPrime(f *testing.F) {
	sieve := NewSieve(SieveInitOptions{})

	f.Fuzz(func(t *testing.T, n int64) {
		if !big.NewInt(sieve.NthPrime(n)).ProbablyPrime(0) {
			t.Errorf("the sieve produced a non-prime number at index %d", n)
		}
	})
}
