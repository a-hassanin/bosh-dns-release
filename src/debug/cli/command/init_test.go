package command_test

import (
	"fmt"
	"net"

	"code.cloudfoundry.org/tlsconfig"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"testing"
)

func TestCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cli/command")
}

func newFakeAPIServer() *ghttp.Server {
	tlsConfig, err := tlsconfig.Build(
		tlsconfig.WithIdentityFromFile("../../../bosh-dns/dns/api/assets/test_certs/test_server.pem", "../../../bosh-dns/dns/api/assets/test_certs/test_server.key"),
		tlsconfig.WithInternalServiceDefaults(),
	).Server(
		tlsconfig.WithClientAuthenticationFromFile("../../../bosh-dns/dns/api/assets/test_certs/test_ca.pem"),
	)
	Expect(err).NotTo(HaveOccurred())

	tlsConfig.BuildNameToCertificate() //nolint:staticcheck
	server := ghttp.NewUnstartedServer()
	err = server.HTTPTestServer.Listener.Close()
	Expect(err).NotTo(HaveOccurred())

	suiteConfig, _ := GinkgoConfiguration()
	port := 2345 + suiteConfig.ParallelProcess
	server.HTTPTestServer.Listener, err = net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	Expect(err).ToNot(HaveOccurred())

	server.HTTPTestServer.TLS = tlsConfig
	server.HTTPTestServer.StartTLS()

	return server
}
