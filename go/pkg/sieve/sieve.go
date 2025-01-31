package sieve

import "sync"

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
	// default batch size to 1 million, in testing this was found this to be the fastest
	// enforce a minimum batch size of 10 million (in testing smaller was much too slow)
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
	s.readLock.RLock()
	if len(s.readBuffer) > int(n) {
		s.readLock.RUnlock()
		return s.readBuffer[n].primeNumber
	}
	s.readLock.RUnlock()

	s.writeLock.Lock()
	defer s.writeLock.Unlock()

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
	for i := 0; i < int(s.batchSize); i++ {
		prime[i] = true
	}
	for primeCheckerIndex, primeChecker := range s.writeBuffer {
		currentPrime := s.writeBuffer[primeCheckerIndex].primeNumber
		iTracker := s.writeBuffer[primeCheckerIndex].lastMultiple
		for i := primeChecker.lastMultiple + currentPrime; i < s.batchSize+s.curNum; i += currentPrime {
			// we have to subtract the current number from the last multiple to get the actual index
			prime[i-s.curNum] = false
			iTracker = i
		}
		s.writeBuffer[primeCheckerIndex].lastMultiple = iTracker
	}
	// loop through the found primes
	for i := 0; i < len(prime); i++ {
		if prime[i] {
			// don't check the new stored prime until you reach it's square, due to how this is being checked in loops, it should be number * (number - 1)
			s.writeBuffer = append(s.writeBuffer, PrimeChecker{int64(i + int(s.curNum)), int64((i+int(s.curNum))*(i+int(s.curNum)) - 1)})
		}
	}
	// ensure the curNum is updated with the batch size
	s.curNum += s.batchSize
}

func (s *SieveNthPrime) seedPrimes() {
	s.seedOnce.Do(func() {
		s.readLock.Lock()
		defer s.readLock.Unlock()
		// Run seedPrimes() only once safely

		if len(s.writeBuffer) == 0 {
			prime := make([]bool, s.batchSize+1)
			for i := 0; i < len(prime); i++ {
				prime[i] = true
			}
			p := int64(2)
			prime[0] = false
			prime[1] = false
			for p*p <= s.batchSize {
				if prime[p] {
					var iTracker int64
					for i := p * p; i < s.batchSize; i += p {
						prime[i] = false
						iTracker = i
					}
					s.writeBuffer = append(s.writeBuffer, PrimeChecker{p, iTracker})
				}
				p++
			}
			currentPrimes := len(s.writeBuffer)
			for i := int64(2); i < s.batchSize; i++ {
				if prime[i] {
					if currentPrimes != 0 {
						currentPrimes--
					} else {
						s.writeBuffer = append(s.writeBuffer, PrimeChecker{i, i * (i - 1)})
					}
				}
			}
			s.curNum += s.batchSize
		}
	})
}
