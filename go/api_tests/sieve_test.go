package api_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verify Seive API correctly returns primes", func() {

	// using describe table, run through all the tests
	DescribeTable("nth prime number", func(nth int64, expected int64, timeout string) {
		Eventually(func() error {
			prime, response, err := primeClient.EratosthenesAPI.GetNthPrime(ctx, nth).Execute()
			if err != nil {
				return err
			}
			if response.StatusCode != 200 {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}
			if prime.Prime != expected {
				return fmt.Errorf("prime %d != expected %d", prime.Prime, expected)
			}
			return nil
		}, timeout, "100ms").Should(Succeed(), "Failed to get expected prime")
	},
		Entry("nth prime 0", int64(0), int64(2), "50ms"),
		Entry("nth prime 19", int64(19), int64(71), "50ms"),
		Entry("nth prime 99", int64(99), int64(541), "50ms"),
		Entry("nth prime 500", int64(500), int64(3581), "50ms"),
		Entry("nth prime 986", int64(986), int64(7793), "50ms"),
		Entry("nth prime 2000", int64(2000), int64(17393), "50ms"),
		Entry("nth prime 1000000", int64(1000000), int64(15485867), "2s"),
		Entry("nth prime 10000000", int64(10000000), int64(179424691), "5s"),
		Entry("nth prime 100000000", int64(100000000), int64(2038074751), "50s"),
		Entry("getting the largest prime again should be quick", int64(100000000), int64(2038074751), "50ms"),
	)
	It("should validate the client prevents invalid requests", func() {
		_, response, err := primeClient.EratosthenesAPI.GetNthPrime(ctx, -2).Execute()
		Expect(err).To(HaveOccurred())
		Expect(response).To(BeNil())
	})
	It("should return a bad request for a non-numeric value", func() {
		// since the client catches the error, we need to call the server directly
		resp, err := http.Get(baseUrl + "/primebynumber/abc")
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(400))
		Expect(resp.Body.Close()).To(Succeed())
	})
	It("should return a bad request for a negative value", func() {
		resp, err := http.Get(baseUrl + "/primebynumber/-1")
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(400))
		Expect(resp.Body.Close()).To(Succeed())
	})
})
