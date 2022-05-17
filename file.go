package certprovider

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// FileProvider uses files for the source of certificates and keys.
type FileProvider struct {
	opts         options
	certDir      string
	identityCert tls.Certificate
	caPool       *x509.CertPool
}

// NewFileProvider returns a new FileProvider using certs from the specified directory
// optionally also can be used for gRPC clients by setting server to false.
func NewFileProvider(
	certDir string,
	opts ...Option,
) (*FileProvider, error) {
	var err error

	f := &FileProvider{ //nolint:varnamelen // `f` is fine for this scope.
		certDir: os.ExpandEnv(certDir),
		opts:    defaultOptions(),
	}

	AddSearchPath(f.certDir).apply(&f.opts)

	for _, opt := range opts {
		opt.apply(&f.opts)
	}

	if f.identityCert, err = tls.LoadX509KeyPair(f.opts.getCertFile(), f.opts.getKeyFile()); err != nil {
		return nil, fmt.Errorf("failed loading certificates: %w", err)
	}

	// populate certs.IdentityCert.Leaf. This has already been parsed, but
	// intentionally discarded by LoadX509KeyPair, for some reason.
	if f.identityCert.Leaf, err = x509.ParseCertificate(f.identityCert.Certificate[0]); err != nil {
		return nil, fmt.Errorf("failed loading certificates: %w", err)
	}

	if f.opts.loadSystemCA {
		if f.caPool, err = x509.SystemCertPool(); err != nil {
			return nil, fmt.Errorf("failed loading certificates: %w", err)
		}
	} else {
		f.caPool = x509.NewCertPool()
	}

	if caFile, ok := f.opts.getCAFile(); ok {
		ca, err := ioutil.ReadFile(caFile) //nolint:gosec // should be decided at compile time for most usages.
		if err != nil {
			return nil, fmt.Errorf("failed loading certificates: %w", err)
		}

		if ok := f.caPool.AppendCertsFromPEM(ca); !ok {
			return nil, ErrNoValidCertificates
		}
	}

	return f, nil
}

// IdentityCert returns the Identity Certificate used for the connection.
func (c *FileProvider) IdentityCert() tls.Certificate {
	return c.identityCert
}

// CAPool returns the CA Pool for the connection.
func (c *FileProvider) CAPool() *x509.CertPool {
	return c.caPool
}

// ServerOption returns the grpc.ServerOption for use with a new gRPC server.
func (c *FileProvider) ServerOption() grpc.ServerOption {
	creds := credentials.NewTLS(&tls.Config{ //nolint:gosec // default minimum is TLS1.3.
		ClientCAs:    c.CAPool(),
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{c.IdentityCert()},
		MinVersion:   c.opts.minTLSVersion,
	})

	return grpc.Creds(creds)
}

// DialOption returns the grpc.DialOption used with a gRPC client.
func (c *FileProvider) DialOption(serverName string) grpc.DialOption {
	creds := credentials.NewTLS(&tls.Config{ //nolint:gosec // default minimum is TLS1.3.
		ServerName:   serverName,
		RootCAs:      c.CAPool(),
		Certificates: []tls.Certificate{c.IdentityCert()},
		MinVersion:   c.opts.minTLSVersion,
	})

	return grpc.WithTransportCredentials(creds)
}
