/*
 * Sieve of Eratosthenes API
 *
 * This api does one thing and one thing only, get you a prime by number.
 *
 * API version: 1.0.0
 */

package server

import (
	"context"
	"net/http"
	"ssse-exercise-sieve/autogenerated_api"
	"ssse-exercise-sieve/pkg/sieve"
)

// SieveApiService is a service that implements the logic for the DefaultAPIServicer
// This service should implement the business logic for every endpoint for the DefaultAPI API.
// Include any external packages or services that will be required by this service.
type SieveApiService struct {
	sieve sieve.Sieve
}

var _ = autogenerated_api.EratosthenesAPIServicer(&SieveApiService{})

// NewDefaultAPIService creates a default api service
func NewEratosthenesAPIService() *SieveApiService {
	// TODO: update this method to allow for storing the precomputed primes in a file
	return &SieveApiService{
		sieve: sieve.NewSieve(sieve.SieveInitOptions{}),
	}
}

// GetNthPrime - Retrieve nth prime number
func (s *SieveApiService) GetNthPrime(ctx context.Context, nthprime int64) (autogenerated_api.ImplResponse, error) {
	nthPrime := s.sieve.NthPrime(int64(nthprime))
	return autogenerated_api.Response(http.StatusOK, autogenerated_api.GetNthPrime200Response{Prime: nthPrime}), nil
}
