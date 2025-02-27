package sieve

import (
	"sync"
)

type Sieve interface {
	NthPrime(n int64) int64
}

type SieveInitOptions struct {
	BatchSize           int64
	StoredPrimeCheckers []PrimeChecker
	CurNum              int64
}
type SieveNthPrime struct {
	batchSize   int64
	curNum      int64
	writeBuffer []PrimeChecker
	readLock    sync.RWMutex
	writeLock   sync.Mutex
	seedOnce    sync.Once
	readBuffer  []PrimeChecker
}

type PrimeChecker struct {
	primeNumber  int64
	lastMultiple int64
}

func NewSieve(options SieveInitOptions) Sieve {
	// default batch size to 10 million, in testing this was found this to be the fastest
	// enforce a minimum batch size of 1 million (in testing smaller was much too slow)
	if options.BatchSize <= 1000000 {
		options.BatchSize = 10000000
	}
	sieve := SieveNthPrime{
		batchSize:   options.BatchSize,
		curNum:      options.CurNum,
		writeBuffer: options.StoredPrimeCheckers,
		// ensure that the read buffer is a deep copy of the write buffer
		readBuffer: append([]PrimeChecker{}, options.StoredPrimeCheckers...),
	}
	sieve.seedPrimes()
	return &sieve
}

func (s *SieveNthPrime) NthPrime(n int64) int64 {
	// no negative primes exist
	if n < 0 {
		return 0
	}
	s.readLock.RLock()
	if len(s.readBuffer) > int(n) {
		s.readLock.RUnlock()
		return s.readBuffer[n].primeNumber
	}
	s.readLock.RUnlock()

	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	// don't attempt to search if a prior locking run already found the prime
	for len(s.writeBuffer) <= int(n) {
		s.batchSearchPrimes()
	}
	s.readLock.Lock()
	defer s.readLock.Unlock()
	s.readBuffer = append(s.readBuffer, s.writeBuffer[len(s.readBuffer):]...)
	return s.writeBuffer[n].primeNumber
}

func (s *SieveNthPrime) batchSearchPrimes() {
	prime := make([]bool, s.batchSize)
	for i := range int(s.batchSize) {
		prime[i] = true
	}
	// iterate through the primes we have found to fill the sieve
	for primeCheckerIndex := range s.writeBuffer {
		currentPrime := s.writeBuffer[primeCheckerIndex].primeNumber
		iTracker := s.writeBuffer[primeCheckerIndex].lastMultiple
		// check if the last multiple is greater than the current number
		for i := iTracker + currentPrime; i < s.batchSize+s.curNum; i += currentPrime {
			// we have to subtract the current number from the last multiple to get the actual index
			prime[i-s.curNum] = false
			iTracker = i
		}
		// update the next number we will be checking that is associated with the prime.
		s.writeBuffer[primeCheckerIndex].lastMultiple = iTracker
	}
	// loop through the found primes
	for i := range prime {
		if prime[i] {
			// don't check the new stored prime until you reach it's square, due to how this is being checked in loops, it should be number * (number - 1)
			s.writeBuffer = append(s.writeBuffer, PrimeChecker{
				primeNumber:  int64(i + int(s.curNum)),
				lastMultiple: int64((i+int(s.curNum))*(i+int(s.curNum)) - 1),
			})
		}
	}
	// ensure the curNum is updated with the batch size so we know what has to be added to the next batch
	s.curNum += s.batchSize
}

func (s *SieveNthPrime) seedPrimes() {
	s.seedOnce.Do(func() {
		s.readLock.Lock()
		defer s.readLock.Unlock()

		if len(s.writeBuffer) == 0 {
			// create a list that will serve as our sieve for eliminating non-primes
			prime := make([]bool, s.batchSize+1)
			for i := range prime {
				prime[i] = true
			}
			// we are kick starting this by setting two as a the starting point and the valid prime
			p := int64(2)
			prime[0] = false
			prime[1] = false
			// all primes less than the square will be covered by smaller primes.
			for p*p <= s.batchSize {
				if prime[p] {
					var iTracker int64
					for i := p * p; i < s.batchSize; i += p {
						prime[i] = false
						iTracker = i
					}
					// we don't want to lose all that math we did to figure out what the next number to check over the batch size
					s.writeBuffer = append(s.writeBuffer, PrimeChecker{p, iTracker})
				}
				p++
				// loop until we get to the next prime to avoid waisting cycles multiplying non-primes
				for !prime[p] {
					p++
				}

			}

			currentPrimes := len(s.writeBuffer)
			// add remaining items in the list that havent been added to primeChecker
			for i := int64(2); i < s.batchSize; i++ {
				if prime[i] {
					// skip over primes that have already been added
					if currentPrimes != 0 {
						currentPrimes--
					} else {
						// once all primes that have already been added, insert the unnadded primes, with the next numb.
						s.writeBuffer = append(s.writeBuffer, PrimeChecker{
							primeNumber:  i,
							lastMultiple: i * (i - 1),
						})
					}
				}
			}
			// ensure we track the current number now that our original primes are seeded
			s.curNum += s.batchSize
		}
	})
}
