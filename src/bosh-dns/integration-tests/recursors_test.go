package integration_tests

import (
	"bosh-dns/dns/server/handlers"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"bosh-dns/acceptance_tests/helpers"
	"bosh-dns/dns/server/record"

	gomegadns "bosh-dns/gomega-dns"

	"github.com/miekg/dns"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"gopkg.in/yaml.v2"
)

type testRecursor struct {
	session              *gexec.Session
	port                 int
	address              string
	configurableResponse string
}

type Config struct {
	Port                 int    `yaml:"port"`
	ConfigurableResponse string `yaml:"configurable_response"`
}

func NewTestRecursor(port int, configurableResponse string) *testRecursor {
	return &testRecursor{
		address:              "127.0.0.1",
		port:                 port,
		configurableResponse: configurableResponse,
	}
}

func (t *testRecursor) start() error {
	originalCwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	testRecursorPath, err :=
		filepath.Abs(filepath.Join(
			originalCwd,
			"..",
			"acceptance_tests",
			"dns-acceptance-release",
			"src",
			"test-recursor",
		))
	Expect(err).NotTo(HaveOccurred())

	err = os.Chdir(testRecursorPath)
	Expect(err).NotTo(HaveOccurred())

	binaryPath, err := gexec.Build("./main")
	Expect(err).NotTo(HaveOccurred())

	err = os.Chdir(originalCwd)
	Expect(err).NotTo(HaveOccurred())

	configYAML, err := yaml.Marshal(Config{
		Port:                 t.port,
		ConfigurableResponse: t.configurableResponse,
	})
	if err != nil {
		return err
	}

	configTempfile, err := os.CreateTemp("", "test-recursor")
	if err != nil {
		return err
	}
	if _, err := configTempfile.Write(configYAML); err != nil {
		return err
	}
	if err := configTempfile.Close(); err != nil {
		return err
	}

	t.session, err = gexec.Start(exec.Command(binaryPath, configTempfile.Name()),
		GinkgoWriter, GinkgoWriter)
	if err != nil {
		return err
	}

	Eventually(t.checkConnection, 5*time.Second, 500*time.Millisecond).Should(ConsistOf(
		gomegadns.MatchResponse(gomegadns.Response{"ip": "10.10.10.10", "ttl": 0}),
	))

	return nil
}

func (t *testRecursor) stop() {
	t.session.Kill().Wait()
}

func (t *testRecursor) checkConnection() []dns.RR {
	response := helpers.DigWithOptions("example.com.", t.address, helpers.DigOpts{
		SkipErrCheck:   true,
		SkipRcodeCheck: true,
		Port:           t.port,
		Timeout:        5 * time.Millisecond,
	})

	if response == nil {
		return []dns.RR{}
	}

	return response.Answer
}

var _ = Describe("Integration", func() {
	Describe("Recursors Tests", func() {
		var (
			responses         []record.Record
			recursors         []string
			caching           bool
			environment       TestEnvironment
			recursorEnv       *testRecursor
			recursorSelection string
			excludedRecursors []string
		)

		BeforeEach(func() {
			responses = []record.Record{record.Record{
				ID:            "garbage",
				IP:            "255.255.255.255",
				InstanceIndex: "2",
			}}
			caching = false
			recursors = []string{}
			excludedRecursors = []string{}
			recursorSelection = "serial"
		})

		JustBeforeEach(func() {
			var err error

			environment = NewTestEnvironment(responses, []record.Host{}, recursors, caching, recursorSelection, excludedRecursors, false)
			if err := environment.Start(); err != nil {
				Fail(fmt.Sprintf("could not start test environment: %s", err))
			}

			recursorEnv = NewTestRecursor(6364, "1.1.1.1")
			err = recursorEnv.start()
			if err != nil {
				Fail(fmt.Sprintf("could not start test recursor: %s", err))
			}
		})

		JustAfterEach(func() {
			if err := environment.Stop(); err != nil {
				Fail(fmt.Sprintf("Failed to stop bosh-dns test environment: %s", err))
			}

			recursorEnv.stop()
		})

		Context("when the recursors are configured explicitly on the DNS server", func() {
			BeforeEach(func() {
				recursors = []string{"127.0.0.1:6364"}
			})

			It("forwards queries to the configured recursors", func() {
				dnsResponse := helpers.DigWithPort("example.com.", environment.ServerAddress(), environment.Port())

				Expect(dnsResponse).To(gomegadns.HaveFlags("qr", "aa", "rd", "ra"))
				Expect(dnsResponse.Answer).To(ConsistOf(
					gomegadns.MatchResponse(gomegadns.Response{"ip": "10.10.10.10", "ttl": 0}),
				))
			})
		})

		Context("when using cache", func() {
			BeforeEach(func() {
				caching = true
				recursors = []string{"127.0.0.1:6364"}
			})

			It("caches upstream dns entries for the duration of the TTL", func() {
				dnsResponse := helpers.DigWithPort("always-different-with-timeout-example.com.", environment.ServerAddress(), environment.Port())

				Expect(dnsResponse.Answer).To(HaveLen(1))
				dnsAnswer := dnsResponse.Answer[0]
				initialIP := gomegadns.FetchIP(dnsAnswer)

				Expect(dnsAnswer).To(gomegadns.MatchResponse(gomegadns.Response{
					"ttl": 5,
					"ip":  initialIP,
				}))

				Consistently(func() []dns.RR {
					dnsResponse := helpers.DigWithPort("always-different-with-timeout-example.com.", environment.ServerAddress(), environment.Port())
					return dnsResponse.Answer
				}, "2s", "1s").Should(ConsistOf(
					gomegadns.MatchResponse(gomegadns.Response{
						"ip": initialIP,
					}),
				))

				nextIP := net.ParseIP(initialIP).To4()
				nextIP[3]++

				Eventually(func() []dns.RR {
					dnsResponse := helpers.DigWithPort("always-different-with-timeout-example.com.", environment.ServerAddress(), environment.Port())
					return dnsResponse.Answer
				}, 5*time.Second, 500*time.Millisecond).Should(ConsistOf(
					gomegadns.MatchResponse(gomegadns.Response{
						"ip": nextIP.String(),
					}),
				))
			})

			Context("negative caching", func() {
				Context("recursor down", func() {
					It("does not negative cache", func() {

						recursorEnv.stop()

						dnsResponse := helpers.DigWithOptions("recursor-small.com.", environment.ServerAddress(), helpers.DigOpts{Port: environment.Port(), SkipRcodeCheck: true, SkipErrCheck: true})

						Expect(dnsResponse.Rcode).To(Equal(dns.RcodeNameError))
						Expect(dnsResponse.Answer).To(HaveLen(0))

						recursorEnv.start() //nolint:errcheck

						dnsResponse = helpers.DigWithOptions("recursor-small.com.", environment.ServerAddress(), helpers.DigOpts{Port: environment.Port(), SkipRcodeCheck: true, SkipErrCheck: true})

						Expect(dnsResponse.Rcode).To(Equal(dns.RcodeSuccess))
						Expect(dnsResponse.Answer).To(HaveLen(2))
					})
				})

				Context("slow recursor", func() {
					It("does not negative cache", func() {

						dnsResponse := helpers.DigWithOptions("alternating-slow-recursor.com.", environment.ServerAddress(), helpers.DigOpts{Id: 1, Timeout: 5 * time.Second, Port: environment.Port(), SkipRcodeCheck: true, SkipErrCheck: true})

						Expect(dnsResponse.Rcode).To(Equal(dns.RcodeNameError))
						Expect(dnsResponse.Answer).To(HaveLen(0))

						dnsResponse = helpers.DigWithOptions("alternating-slow-recursor.com.", environment.ServerAddress(), helpers.DigOpts{Id: 2, Port: environment.Port(), SkipRcodeCheck: true, SkipErrCheck: true})

						Expect(dnsResponse.Rcode).To(Equal(dns.RcodeSuccess))
						Expect(dnsResponse.Answer).To(HaveLen(1))
					})
				})

				Context("recursor returns error", func() {
					It("does not negative cache", func() {

						dnsResponse := helpers.DigWithOptions("alternating-nameerror-recursor.com.", environment.ServerAddress(), helpers.DigOpts{Id: 1, Port: environment.Port(), SkipRcodeCheck: true, SkipErrCheck: true})

						Expect(dnsResponse.Rcode).To(Equal(dns.RcodeNameError))
						Expect(dnsResponse.Answer).To(HaveLen(0))

						dnsResponse = helpers.DigWithOptions("alternating-nameerror-recursor.com.", environment.ServerAddress(), helpers.DigOpts{Id: 2, Port: environment.Port(), SkipRcodeCheck: true, SkipErrCheck: true})

						Expect(dnsResponse.Rcode).To(Equal(dns.RcodeSuccess))
						Expect(dnsResponse.Answer).To(HaveLen(1))
					})
				})

				When("the recursor returns NXDOMAIN for an SOA record", func() {
					It("is the only query type (SOA) where it'll negative cache the NXDOMAIN result", func() {

						dnsResponse := helpers.DigWithOptions("alternating-soa-nameerror-recursor.com.", environment.ServerAddress(), helpers.DigOpts{Id: 1, Port: environment.Port(), SkipRcodeCheck: true, SkipErrCheck: true, Type: dns.TypeSOA})

						Expect(dnsResponse.Rcode).To(Equal(dns.RcodeNameError))
						Expect(dnsResponse.Answer).To(HaveLen(0))

						// Note that we set `Id: 2`, which signals the upstream server to return a valid record, but since we're negative caching, we never ask
						dnsResponse = helpers.DigWithOptions("alternating-soa-nameerror-recursor.com.", environment.ServerAddress(), helpers.DigOpts{Id: 2, Port: environment.Port(), SkipRcodeCheck: true, SkipErrCheck: true, Type: dns.TypeSOA})

						Expect(dnsResponse.Rcode).To(Equal(dns.RcodeNameError))
						Expect(dnsResponse.Answer).To(HaveLen(0))
					})
				})
			})
		})

		Context("handling upstream recursor responses", func() {
			BeforeEach(func() {
				recursors = []string{"127.0.0.1:6364"}
			})

			It("returns success when receiving a truncated responses from a recursor", func() {
				dnsResponse := helpers.DigWithPort("truncated-recursor.com.", environment.ServerAddress(), environment.Port())
				Expect(dnsResponse).To(gomegadns.HaveFlags("qr", "aa", "tc", "rd", "ra"))
				Expect(dnsResponse.Answer).To(HaveLen(1))
			})

			It("timeouts when recursor takes longer than configured recursor_timeout", func() {
				dnsResponse := helpers.DigWithOptions("slow-recursor.com.", environment.ServerAddress(), helpers.DigOpts{SkipRcodeCheck: true, Timeout: 5 * time.Second, Port: environment.Port()})
				Expect(dnsResponse.Rcode).To(Equal(dns.RcodeNameError))
			})

			It("forwards large UDP EDNS messages", func() {
				dnsResponse := helpers.DigWithOptions("udp-9k-message.com.", environment.ServerAddress(), helpers.DigOpts{BufferSize: 65535, Port: environment.Port()})
				Expect(dnsResponse).To(gomegadns.HaveFlags("qr", "aa", "rd", "ra"))
				Expect(dnsResponse.Answer).To(HaveLen(272))
			})

			It("compresses message responses that are larger than requested UDP Size", func() {
				dnsResponse := helpers.DigWithOptions("compressed-ip-truncated-recursor-large.com.", environment.ServerAddress(), helpers.DigOpts{BufferSize: 512, Port: environment.Port()})
				Expect(dnsResponse).To(gomegadns.HaveFlags("qr", "aa", "rd", "ra"))
				Expect(dnsResponse.Len()).To(BeNumerically(">", 512)) // uncompressed length
				dnsResponse.Compress = true
				Expect(dnsResponse.Len()).To(BeNumerically("<=", 512)) // compressed length
			})

			It("truncates large dns answers if udp response size is larger than 512", func() {
				dnsResponse := helpers.DigWithPort("ip-truncated-recursor-large.com.", environment.ServerAddress(), environment.Port())
				Expect(dnsResponse).To(gomegadns.HaveFlags("qr", "aa", "tc", "rd", "ra"))
				Expect(dnsResponse.Answer).To(HaveLen(20))
			})

			It("does not bother to compress messages that are smaller than 512", func() {
				dnsResponse := helpers.DigWithPort("recursor-small.com.", environment.ServerAddress(), environment.Port())
				Expect(dnsResponse).To(gomegadns.HaveFlags("qr", "aa", "rd", "ra"))
				Expect(dnsResponse.Answer).To(HaveLen(2))
			})

			It("forwards ipv4 ARPA queries to the configured recursors", func() {
				dnsResponse := helpers.ReverseDigWithOptions("8.8.4.4", environment.ServerAddress(), helpers.DigOpts{Port: environment.Port()})
				Expect(dnsResponse).To(gomegadns.HaveFlags("qr", "aa", "rd", "ra"))
				Expect(dnsResponse.Answer).To(ConsistOf(
					gomegadns.MatchResponse(gomegadns.Response{"name": "4.4.8.8.in-addr.arpa."}),
				))
			})

			It("forwards ipv6 ARPA queries to the configured recursors", func() {
				dnsResponse := helpers.IPv6ReverseDigWithOptions("2001:4860:4860::8888", environment.ServerAddress(), helpers.DigOpts{Port: environment.Port()})

				Expect(dnsResponse).To(gomegadns.HaveFlags("qr", "aa", "rd", "ra"))
				Expect(dnsResponse.Answer).To(ConsistOf(
					gomegadns.MatchResponse(gomegadns.Response{"name": "8.8.8.8.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.6.8.4.0.6.8.4.1.0.0.2.ip6.arpa."}),
				))
			})
		})

		Context("when the upstream recursors have different responses", func() {
			const (
				testQuestion = "question_with_configurable_response."
			)
			var secondTestRecursor *testRecursor

			JustBeforeEach(func() {
				secondTestRecursor = NewTestRecursor(6365, "2.2.2.2")
				err := secondTestRecursor.start()
				if err != nil {
					Fail(fmt.Sprintf("could not start test recursor: %s", err))
				}
			})

			BeforeEach(func() {
				recursors = []string{"127.0.0.1:6364", "127.0.0.1:6365"}
			})

			JustAfterEach(func() {
				secondTestRecursor.stop()
			})

			Context("recursor selection", func() {
				Context("serial", func() {
					BeforeEach(func() {
						recursorSelection = "serial"
					})

					It("always chooses recursors serially", func() {
						err := environment.Restart()
						Expect(err).NotTo(HaveOccurred())

						dnsResponse := helpers.DigWithPort(testQuestion, environment.ServerAddress(), environment.Port())
						Expect(dnsResponse.Answer).To(HaveLen(1))

						Expect(dnsResponse.Answer[0]).Should(gomegadns.MatchResponse(gomegadns.Response{"ip": "1.1.1.1"}))
					})
				})

				Context("smart", func() {
					BeforeEach(func() {
						recursorSelection = "smart"
					})

					It("shuffles recursors", func() {
						dnsResponse := helpers.DigWithPort(testQuestion, environment.ServerAddress(), environment.Port())
						Expect(dnsResponse.Answer).To(HaveLen(1))

						initialUpstreamResponse := dnsResponse.Answer[0]

						Eventually(func() dns.RR {
							err := environment.Restart()
							Expect(err).NotTo(HaveOccurred())

							dnsResponse := helpers.DigWithPort(testQuestion, environment.ServerAddress(), environment.Port())
							Expect(dnsResponse.Answer).To(HaveLen(1))

							return dnsResponse.Answer[0]
						}, 30*time.Second).ShouldNot(Equal(initialUpstreamResponse))
					})
				})
			})

			Context("failing recursor shifts recursor preference", func() {
				var (
					firstSelectedRecursor   *testRecursor
					initialUpstreamResponse dns.RR
				)

				BeforeEach(func() {
					recursorSelection = "smart"
				})

				JustBeforeEach(func() {
					dnsResponse := helpers.DigWithPort(testQuestion, environment.ServerAddress(), environment.Port())
					Expect(dnsResponse.Answer).To(HaveLen(1))
					initialUpstreamResponse = dnsResponse.Answer[0]

					if gomegadns.FetchIP(initialUpstreamResponse) == recursorEnv.configurableResponse {
						firstSelectedRecursor = recursorEnv
					} else if gomegadns.FetchIP(initialUpstreamResponse) == secondTestRecursor.configurableResponse {
						firstSelectedRecursor = secondTestRecursor
					}
				})

				It("shifts", func() {
					By("killing the first upstream recursor")
					firstSelectedRecursor.stop()

					By("forcing the preference shift to the second upstream recursor")
					for i := 0; i < handlers.FailHistoryThreshold; i++ {
						dnsResponse := helpers.DigWithOptions(testQuestion, environment.ServerAddress(),
							helpers.DigOpts{Port: environment.Port(), Timeout: 3 * time.Second},
						)
						Expect(dnsResponse.Answer[0]).ShouldNot(Equal(initialUpstreamResponse))
						fmt.Printf("Running %d times out of %d total\n", i, handlers.FailHistoryThreshold)
					}

					By("bringing back the first upstream recursor")
					err := firstSelectedRecursor.start()
					Expect(err).NotTo(HaveOccurred())

					By("validating that we still use the second recursor")
					Consistently(func() dns.RR {
						return helpers.DigWithPort(testQuestion, environment.ServerAddress(), environment.Port()).Answer[0]
					}, 5*time.Second, 500*time.Millisecond).ShouldNot(Equal(initialUpstreamResponse))
				})
			})

			Context("failover strategy", func() {
				Context("when serial", func() {
					BeforeEach(func() {
						recursorSelection = "serial"
					})

					It("always attempts to query the first configured recursor", func() {

						dnsResponse := helpers.DigWithOptions(testQuestion, environment.ServerAddress(),
							helpers.DigOpts{Port: environment.Port(), Timeout: 3 * time.Second},
						)
						Expect(dnsResponse.Answer[0]).Should(gomegadns.MatchResponse(
							gomegadns.Response{"ip": recursorEnv.configurableResponse}))

						By("stopping the first recursor")
						recursorEnv.stop()

						By("then bosh-dns fails overs to the second recursor")
						dnsResponse = helpers.DigWithOptions(testQuestion, environment.ServerAddress(),
							helpers.DigOpts{Port: environment.Port(), Timeout: 3 * time.Second},
						)
						Expect(dnsResponse.Answer[0]).Should(gomegadns.MatchResponse(
							gomegadns.Response{"ip": secondTestRecursor.configurableResponse}))

						By("bringing the first recursor back")
						err := recursorEnv.start()
						Expect(err).NotTo(HaveOccurred())

						By("then bosh-dns resumes with the first recursor's response")
						dnsResponse = helpers.DigWithOptions(testQuestion, environment.ServerAddress(),
							helpers.DigOpts{Port: environment.Port(), Timeout: 3 * time.Second},
						)
						Expect(dnsResponse.Answer[0]).Should(gomegadns.MatchResponse(
							gomegadns.Response{"ip": recursorEnv.configurableResponse}))

					})

					It("returns an error if all recursors fail", func() {
						recursorEnv.stop()
						secondTestRecursor.stop()

						helpers.DigWithOptions(testQuestion, environment.ServerAddress(),
							helpers.DigOpts{Port: environment.Port(), SkipErrCheck: true, SkipRcodeCheck: true})
						Eventually(environment.Output(), 5*time.Second, 500*time.Millisecond).Should(gbytes.Say(`no response from recursors`))
					})
				})
			})

			Context("excluding recursors", func() {
				BeforeEach(func() {
					excludedRecursors = []string{"127.0.0.1:6364"}
				})

				It("excludes the recursor specified", func() {
					dnsResponse := helpers.DigWithPort(testQuestion, environment.ServerAddress(), environment.Port())
					Expect(dnsResponse.Answer).To(HaveLen(1))
					Expect(dnsResponse.Answer).To(ConsistOf(
						gomegadns.MatchResponse(gomegadns.Response{"ip": "2.2.2.2", "ttl": 0}),
					))
				})
			})
		})
	})
})
