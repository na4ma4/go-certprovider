package certprovider

import (
	"crypto/tls"
	"crypto/x509"

	"google.golang.org/grpc"
)

// CertificateProvider is an interface to a provider for certificates used with gRPC server and clients.
type CertificateProvider interface {
	IdentityCert() tls.Certificate
	CAPool() *x509.CertPool
	ServerOption() grpc.ServerOption
	DialOption(serverName string) grpc.DialOption
	KeyPair() (tls.Certificate, error)
}
