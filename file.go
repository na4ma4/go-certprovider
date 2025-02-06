package certprovider

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
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

// MustFileCertProvider returns a FileProvider or panic.
func MustFileCertProvider(
	certDir string,
	opts ...Option,
) *FileProvider {
	cp, err := NewFileProvider(certDir, opts...)
	if err != nil {
		panic(err)
	}

	return cp
}

// NewFileProvider returns a new FileProvider using certs from the specified directory
// optionally also can be used for gRPC clients by setting server to false.
func NewFileProvider(
	certDir string,
	opts ...Option,
) (*FileProvider, error) {
	f := &FileProvider{
		certDir: os.ExpandEnv(certDir),
		opts:    defaultOptions(),
	}

	AddSearchPath(f.certDir).apply(&f.opts)

	for _, opt := range opts {
		opt.apply(&f.opts)
	}

	{
		var err error
		if f.identityCert, err = tls.LoadX509KeyPair(f.opts.getCertFile(), f.opts.getKeyFile()); err != nil {
			return nil, fmt.Errorf("failed loading certificates: %w", err)
		}
	}

	// populate certs.IdentityCert.Leaf. This has already been parsed, but
	// intentionally discarded by LoadX509KeyPair, for some reason.
	{
		var err error
		if f.identityCert.Leaf, err = x509.ParseCertificate(f.identityCert.Certificate[0]); err != nil {
			return nil, fmt.Errorf("failed loading certificates: %w", err)
		}
	}

	if f.opts.loadSystemCA {
		var err error
		if f.caPool, err = x509.SystemCertPool(); err != nil {
			return nil, fmt.Errorf("failed loading certificates: %w", err)
		}
	} else {
		f.caPool = x509.NewCertPool()
	}

	if caFile, ok := f.opts.getCAFile(); ok {
		var ca []byte
		{
			var err error
			ca, err = os.ReadFile(caFile)
			if err != nil {
				return nil, fmt.Errorf("failed loading certificates: %w", err)
			}
		}

		if appendOk := f.caPool.AppendCertsFromPEM(ca); !appendOk {
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

func (c *FileProvider) KeyPair() (tls.Certificate, error) {
	return tls.LoadX509KeyPair(c.opts.getCertFile(), c.opts.getKeyFile())
}
