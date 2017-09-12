package healthiness_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	"bosh-dns/dns/server/healthiness"
	"bosh-dns/dns/server/healthiness/healthinessfakes"

	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HealthChecker", func() {
	var (
		ip            string
		fakeClient    *healthinessfakes.FakeHTTPClientGetter
		healthChecker healthiness.HealthChecker

		responseBody string
		responseCode int
		response     *http.Response
	)

	BeforeEach(func() {
		fakeClient = &healthinessfakes.FakeHTTPClientGetter{}
		healthChecker = healthiness.NewHealthChecker(fakeClient, 8081)

		responseCode = 200
		responseBody = `{"state":"running"}`
	})

	JustBeforeEach(func() {
		response = &http.Response{
			StatusCode: responseCode,
			Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
		}
		fakeClient.GetReturns(response, nil)
	})

	Describe("GetStatus", func() {
		Context("when healthy", func() {
			BeforeEach(func() {
				ip = "127.0.0.1"
				responseBody = `{"state":"running"}`
			})

			It("returns true", func() {
				Expect(healthChecker.GetStatus(ip)).To(BeTrue())
				Expect(fakeClient.GetCallCount()).To(Equal(1))
				Expect(fakeClient.GetArgsForCall(0)).To(Equal(fmt.Sprintf("https://%s:8081/health", ip)))
			})
		})

		Context("when unhealthy", func() {
			BeforeEach(func() {
				ip = "127.0.0.2"
				responseBody = `{"state":"stopped"}`
			})

			It("returns false", func() {
				Expect(healthChecker.GetStatus(ip)).To(BeFalse())
				Expect(fakeClient.GetCallCount()).To(Equal(1))
				Expect(fakeClient.GetArgsForCall(0)).To(Equal(fmt.Sprintf("https://%s:8081/health", ip)))
			})
		})

		Context("when unable to fetch status", func() {
			BeforeEach(func() {
				ip = "127.0.0.3"
			})

			It("returns false", func() {
				fakeClient.GetReturns(nil, errors.New("fake connect err"))

				Expect(healthChecker.GetStatus(ip)).To(BeFalse())
				Expect(fakeClient.GetCallCount()).To(Equal(1))
				Expect(fakeClient.GetArgsForCall(0)).To(Equal(fmt.Sprintf("https://%s:8081/health", ip)))
			})
		})

		Context("when status is invalid json", func() {
			BeforeEach(func() {
				ip = "127.0.0.3"
				responseBody = `duck?`
			})

			It("returns false", func() {
				Expect(healthChecker.GetStatus(ip)).To(BeFalse())
				Expect(fakeClient.GetCallCount()).To(Equal(1))
				Expect(fakeClient.GetArgsForCall(0)).To(Equal(fmt.Sprintf("https://%s:8081/health", ip)))
			})
		})

		Context("when response is not 200 OK", func() {
			BeforeEach(func() {
				ip = "127.0.0.3"
				responseCode = 400
			})

			It("returns false", func() {
				Expect(healthChecker.GetStatus(ip)).To(BeFalse())
				Expect(fakeClient.GetCallCount()).To(Equal(1))
				Expect(fakeClient.GetArgsForCall(0)).To(Equal(fmt.Sprintf("https://%s:8081/health", ip)))
			})
		})
	})
})
