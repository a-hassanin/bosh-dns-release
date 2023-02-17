package handlers_test

import (
	. "bosh-dns/dns/server/handlers"
	"bosh-dns/dns/server/handlers/handlersfakes"
	"bosh-dns/dns/server/internal/internalfakes"
	"bosh-dns/dns/server/records/dnsresolver/dnsresolverfakes"

	"errors"
	"fmt"
	"net"
	"net/http"

	"bytes"

	. "bosh-dns/dns/internal/testhelpers/question_case_helpers"

	"github.com/cloudfoundry/bosh-utils/httpclient"
	"github.com/cloudfoundry/bosh-utils/logger/loggerfakes"
	"github.com/miekg/dns"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

type closingBuffer struct {
	*bytes.Buffer
	Closed bool
}

func (cb *closingBuffer) Close() (err error) {
	cb.Closed = true
	return
}

var _ = Describe("HttpJsonHandler", func() {
	var (
		handler       HTTPJSONHandler
		fakeLogger    *loggerfakes.FakeLogger
		fakeWriter    *internalfakes.FakeResponseWriter
		fakeTruncater *dnsresolverfakes.FakeResponseTruncater
	)

	BeforeEach(func() {
		fakeLogger = &loggerfakes.FakeLogger{}
		fakeWriter = &internalfakes.FakeResponseWriter{}
		fakeTruncater = &dnsresolverfakes.FakeResponseTruncater{}
	})

	Context("response body", func() {
		var (
			client   *handlersfakes.FakeHTTPClient
			fakeBody *closingBuffer
		)

		BeforeEach(func() {
			client = &handlersfakes.FakeHTTPClient{}
			fakeBody = &closingBuffer{bytes.NewBufferString(
				`{
					"Status": 0,
					"TC": false,
					"RD": true,
					"RA": true,
					"AD": false,
					"CD": false,
					"Question":
					[
						{
							"name": "app-id.internal-domain.",
							"type": 28
						}
					],
					"Answer": []
				}`),
				false,
			}
		})

		JustBeforeEach(func() {
			handler = NewHTTPJSONHandler("http://example.com", client, fakeLogger, fakeTruncater)
		})

		Context("when the request to the http server fails", func() {
			It("closes the response body", func() {
				fakeResponse := &http.Response{StatusCode: 500, Body: fakeBody}
				client.GetReturns(fakeResponse, nil)

				req := &dns.Msg{}
				SetQuestion(req, nil, "app-id.internal-domain.", dns.TypeA)

				handler.ServeDNS(fakeWriter, req)
				Expect(client.GetCallCount()).To(Equal(1))

				Expect(fakeBody.Closed).To(Equal(true))
			})
		})

		Context("when the request is successful", func() {
			It("closes the response body", func() {
				fakeResponse := &http.Response{StatusCode: 200, Body: fakeBody}
				client.GetReturns(fakeResponse, nil)

				req := &dns.Msg{}
				SetQuestion(req, nil, "app-id.internal-domain.", dns.TypeA)

				handler.ServeDNS(fakeWriter, req)
				Expect(client.GetCallCount()).To(Equal(1))

				Expect(fakeBody.Closed).To(Equal(true))
			})
		})
	})

	Context("when requesting to a running server", func() {
		var (
			server             *ghttp.Server
			fakeServerResponse http.HandlerFunc
		)

		JustBeforeEach(func() {
			server = ghttp.NewUnstartedServer()
			server.AppendHandlers(fakeServerResponse)
			server.HTTPTestServer.Start()
			httpClient := httpclient.NewHTTPClient(httpclient.DefaultClient, fakeLogger)
			handler = NewHTTPJSONHandler(server.URL(), httpClient, fakeLogger, fakeTruncater)
		})

		AfterEach(func() {
			server.Close()
		})

		Context("successful requests", func() {
			BeforeEach(func() {
				casedName := MixCase("app-id.internal-domain.")
				fakeServerResponse = ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/", "name="+casedName+"&type=28"),
					ghttp.RespondWith(http.StatusOK, `{
						"Status": 0,
						"TC": false,
						"RD": true,
						"RA": true,
						"AD": false,
						"CD": false,
						"Question":
						[
							{
								"name": "`+casedName+`",
								"type": 28
							}
						],
						"Answer":
						[
							{
								"name": "`+casedName+`",
								"type": 1,
								"TTL": 1526,
								"data": "192.168.0.1"
							},
							{
								"name": "`+casedName+`",
								"type": 28,
								"TTL": 224,
								"data": "::1"
							}
						],
						"Additional": [ ],
						"edns_client_subnet": "12.34.56.78/0"
					}`))
			})

			It("returns a DNS response based on answer given by backend server", func() {
				var casedName string
				req := &dns.Msg{}
				SetQuestion(req, &casedName, "app-id.internal-domain.", dns.TypeAAAA)

				handler.ServeDNS(fakeWriter, req)

				Expect(fakeWriter.WriteMsgCallCount()).To(Equal(1))
				resp := fakeWriter.WriteMsgArgsForCall(0)
				Expect(resp.Question).To(Equal(req.Question))
				Expect(resp.Rcode).To(Equal(dns.RcodeSuccess))
				Expect(resp.Authoritative).To(BeTrue())
				Expect(resp.RecursionAvailable).To(BeTrue())
				Expect(resp.Answer).To(HaveLen(2))
				Expect(resp.Answer[0]).To(Equal(&dns.A{
					Hdr: dns.RR_Header{
						Name:   casedName,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    1526,
					},
					A: net.ParseIP("192.168.0.1"),
				}))

				Expect(resp.Answer[1]).To(Equal(&dns.A{
					Hdr: dns.RR_Header{
						Name:   casedName,
						Rrtype: dns.TypeAAAA,
						Class:  dns.ClassINET,
						Ttl:    224,
					},
					A: net.ParseIP("::1"),
				}))
			})

			Context("when there are no questions", func() {
				It("returns rcode success", func() {
					msg := &dns.Msg{}
					handler.ServeDNS(fakeWriter, msg)

					message := fakeWriter.WriteMsgArgsForCall(0)
					Expect(message.Rcode).To(Equal(dns.RcodeSuccess))
					Expect(message.Authoritative).To(BeTrue())
					Expect(message.RecursionAvailable).To(BeTrue())
				})
			})
		})

		Context("when it cannot reach the http server", func() {
			JustBeforeEach(func() {
				httpClient := httpclient.NewHTTPClient(httpclient.DefaultClient, fakeLogger)
				handler = NewHTTPJSONHandler("bogus-address", httpClient, fakeLogger, fakeTruncater)
			})

			It("logs the error ", func() {
				req := &dns.Msg{}
				SetQuestion(req, nil, "app-id.internal-domain.", dns.TypeA)
				handler.ServeDNS(fakeWriter, req)

				Expect(fakeLogger.ErrorCallCount()).To(Equal(1))
				tag, template, args := fakeLogger.ErrorArgsForCall(0)
				Expect(tag).To(Equal("HTTPJSONHandler"))
				msg := fmt.Sprintf(template, args...)
				Expect(msg).To(ContainSubstring("error connecting to 'bogus-address': "))
				Expect(msg).To(ContainSubstring("Performing GET request"))
			})

			It("responds with a server fail", func() {
				req := &dns.Msg{}
				SetQuestion(req, nil, "app-id.internal-domain.", dns.TypeA)
				handler.ServeDNS(fakeWriter, req)

				Expect(fakeWriter.WriteMsgCallCount()).To(Equal(1))
				resp := fakeWriter.WriteMsgArgsForCall(0)
				Expect(resp.Question).To(Equal(req.Question))
				Expect(resp.Rcode).To(Equal(dns.RcodeServerFailure))
				Expect(resp.Authoritative).To(BeTrue())
				Expect(resp.RecursionAvailable).To(BeTrue())

				Expect(resp.Answer).To(HaveLen(0))
			})
		})

		Context("when it cannot write the response message", func() {
			BeforeEach(func() {
				fakeWriter.WriteMsgReturns(errors.New("failed to write message"))
			})

			It("logs the error", func() {
				handler.ServeDNS(fakeWriter, &dns.Msg{})

				Expect(fakeLogger.ErrorCallCount()).To(Equal(1))
				tag, msg, _ := fakeLogger.ErrorArgsForCall(0)
				Expect(tag).To(Equal("HTTPJSONHandler"))
				Expect(msg).To(Equal("failed to write message"))
			})
		})

		Context("when the http server response is malformed", func() {
			BeforeEach(func() {
				fakeServerResponse = ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/"),
					ghttp.RespondWith(http.StatusOK, `{  garbage`),
				)
			})

			It("returns a serve fail response", func() {
				req := &dns.Msg{}
				SetQuestion(req, nil, "app-id.internal-domain.", dns.TypeA)
				handler.ServeDNS(fakeWriter, req)

				Expect(fakeWriter.WriteMsgCallCount()).To(Equal(1))
				resp := fakeWriter.WriteMsgArgsForCall(0)
				Expect(resp.Question).To(Equal(req.Question))
				Expect(resp.Rcode).To(Equal(dns.RcodeServerFailure))
				Expect(resp.Authoritative).To(BeTrue())
				Expect(resp.RecursionAvailable).To(BeTrue())

				Expect(resp.Answer).To(HaveLen(0))
			})

			It("logs the error", func() {
				req := &dns.Msg{}
				SetQuestion(req, nil, "app-id.internal-domain.", dns.TypeA)
				handler.ServeDNS(fakeWriter, req)

				Expect(fakeLogger.ErrorCallCount()).To(Equal(1))
				tag, template, args := fakeLogger.ErrorArgsForCall(0)
				Expect(tag).To(Equal("HTTPJSONHandler"))
				msg := fmt.Sprintf(template, args...)
				Expect(msg).To(ContainSubstring("failed to unmarshal response message"))
			})
		})

		Context("when the http server responds with non-200", func() {
			BeforeEach(func() {
				fakeServerResponse = ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/"),
					ghttp.RespondWith(http.StatusNotFound, ""),
				)
			})

			It("returns a serve fail response", func() {
				req := &dns.Msg{}
				SetQuestion(req, nil, "app-id.internal-domain.", dns.TypeA)
				handler.ServeDNS(fakeWriter, req)

				Expect(fakeWriter.WriteMsgCallCount()).To(Equal(1))
				resp := fakeWriter.WriteMsgArgsForCall(0)
				Expect(resp.Question).To(Equal(req.Question))
				Expect(resp.Rcode).To(Equal(dns.RcodeServerFailure))
				Expect(resp.Authoritative).To(BeTrue())
				Expect(resp.RecursionAvailable).To(BeTrue())

				Expect(resp.Answer).To(HaveLen(0))
			})
		})

		Context("when the https server message is truncated", func() {
			BeforeEach(func() {
				casedName := MixCase("app-id.internal-domain.")

				fakeWriter.RemoteAddrReturns(&net.UDPAddr{})
				fakeServerResponse = ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/", "name="+casedName+"&type=1"),
					ghttp.RespondWith(http.StatusOK, `{
						"Status": 0,
						"TC": true,
						"RD": true,
						"RA": true,
						"AD": false,
						"CD": false,
						"Question":
						[
							{
								"name": "`+casedName+`",
								"type": 28
							}
						],
						"Answer":
						[
							{
								"name": "`+casedName+`",
								"type": 1,
								"TTL": 1526,
								"data": "192.168.0.1"
							}
						],
						"Additional": [ ],
						"edns_client_subnet": "12.34.56.78/0"
					}`))
			})

			It("returns a truncated dns message", func() {
				req := &dns.Msg{}
				SetQuestion(req, nil, "app-id.internal-domain.", dns.TypeA)
				handler.ServeDNS(fakeWriter, req)

				Expect(fakeWriter.WriteMsgCallCount()).To(Equal(1))
				resp := fakeWriter.WriteMsgArgsForCall(0)
				Expect(resp.Rcode).To(Equal(dns.RcodeSuccess))
				Expect(resp.Truncated).To(BeTrue())
				Expect(resp.Question).To(Equal(req.Question))
				Expect(resp.Answer).To(HaveLen(1))
			})
		})

		Context("when the non truncated http server response message is too large to fit in dns message", func() {
			BeforeEach(func() {
				casedName := MixCase("app-id.internal-domain.")

				fakeWriter.RemoteAddrReturns(&net.UDPAddr{})
				fakeServerResponse = ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/", "name="+casedName+"&type=1"),
					ghttp.RespondWith(http.StatusOK, `{
						"Status": 0,
						"TC": false,
						"RD": true,
						"RA": true,
						"AD": false,
						"CD": false,
						"Question":
						[
							{
								"name": "`+casedName+`",
								"type": 28
							}
						],
						"Answer": [
							{
								"name": "`+casedName+`",
								"type": 1,
								"TTL": 1526,
								"data": "192.168.0.1"
							},
							{
								"name": "`+casedName+`",
								"type": 1,
								"TTL": 1526,
								"data": "192.168.0.2"
							},
							{
								"name": "`+casedName+`",
								"type": 1,
								"TTL": 1526,
								"data": "192.168.0.3"
							},
							{
								"name": "`+casedName+`",
								"type": 1,
								"TTL": 224,
								"data": "192.168.0.4"
							},
							{
								"name": "`+casedName+`",
								"type": 1,
								"TTL": 224,
								"data": "192.168.0.5"
							},
							{
								"name": "`+casedName+`",
								"type": 1,
								"TTL": 224,
								"data": "192.168.0.6"
							},
							{
								"name": "`+casedName+`",
								"type": 1,
								"TTL": 224,
								"data": "192.168.0.7"
							},
							{
								"name": "`+casedName+`",
								"type": 1,
								"TTL": 224,
								"data": "192.168.0.8"
							},
							{
								"name": "`+casedName+`",
								"type": 1,
								"TTL": 224,
								"data": "192.168.0.9"
							},
							{
								"name": "`+casedName+`",
								"type": 1,
								"TTL": 224,
								"data": "192.168.0.10"
							},
							{
								"name": "`+casedName+`",
								"type": 1,
								"TTL": 224,
								"data": "192.168.0.11"
							},
							{
								"name": "`+casedName+`",
								"type": 1,
								"TTL": 224,
								"data": "192.168.0.12"
							},
							{
								"name": "`+casedName+`",
								"type": 1,
								"TTL": 224,
								"data": "192.168.0.13"
							}
						],
						"Additional": [ ],
						"edns_client_subnet": "12.34.56.78/0"
					}`))
			})

			It("truncates the answers to fit", func() {
				request := &dns.Msg{}
				SetQuestion(request, nil, "app-id.internal-domain.", dns.TypeA)
				handler.ServeDNS(fakeWriter, request)

				Expect(fakeWriter.WriteMsgCallCount()).To(Equal(1))
				response := fakeWriter.WriteMsgArgsForCall(0)
				Expect(response.Rcode).To(Equal(dns.RcodeSuccess))
				Expect(response.RecursionAvailable).To(BeTrue())

				Expect(fakeTruncater.TruncateIfNeededCallCount()).To(Equal(1))
				writer, req, resp := fakeTruncater.TruncateIfNeededArgsForCall(0)
				Expect(writer).To(Equal(fakeWriter))
				Expect(req).To(Equal(request))
				Expect(resp).To(Equal(response))
			})
		})
	})
})
