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
	"errors"
)

type GetNthPrime200Response struct {
	Prime int64 `json:"prime"`
}

// AssertGetNthPrime200ResponseRequired checks if the required fields are not zero-ed
func AssertGetNthPrime200ResponseRequired(obj GetNthPrime200Response) error {
	elements := map[string]interface{}{
		"prime": obj.Prime,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertGetNthPrime200ResponseConstraints checks if the values respects the defined constraints
func AssertGetNthPrime200ResponseConstraints(obj GetNthPrime200Response) error {
	if obj.Prime < 0 {
		return &ParsingError{Param: "Prime", Err: errors.New(errMsgMinValueConstraint)}
	}
	return nil
}
