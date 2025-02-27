// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Sieve of Eratosthenes API
 *
 * This api does one thing and one thing only, get you a prime by number.
 *
 * API version: 1.0.0
 */

package autogenerated_api

import (
	"context"
	"net/http"
)

// EratosthenesAPIRouter defines the required methods for binding the api requests to a responses for the EratosthenesAPI
// The EratosthenesAPIRouter implementation should parse necessary information from the http request,
// pass the data to a EratosthenesAPIServicer to perform the required actions, then write the service results to the http response.
type EratosthenesAPIRouter interface {
	GetNthPrime(http.ResponseWriter, *http.Request)
}

// EratosthenesAPIServicer defines the api actions for the EratosthenesAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type EratosthenesAPIServicer interface {
	GetNthPrime(context.Context, int64) (ImplResponse, error)
}
