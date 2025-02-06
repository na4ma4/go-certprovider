package certprovider

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// DynamicProvider uses files for the source of certificates and keys.
type DynamicProvider struct {
	opts         options
	publicKey    []byte
	privateKey   []byte
	identityCert tls.Certificate
	caPool       *x509.CertPool
}

// MustDynamicCertProvider returns a DynamicCertProvider or panic.
func MustDynamicCertProvider(
	opts ...Option,
) *DynamicProvider {
	cp, err := NewDynamicProvider(opts...)
	if err != nil {
		panic(err)
	}

	return cp
}

// NewDynamicProvider returns a new DynamicProvider using dynamically generated certificates.
func NewDynamicProvider(
	opts ...Option,
) (*DynamicProvider, error) {
	f := &DynamicProvider{
		opts: defaultOptions(),
	}

	{
		var err error
		f.publicKey, f.privateKey, err = f.generateKeypair()
		if err != nil {
			return nil, fmt.Errorf("failed to generate keypair: %w", err)
		}
	}

	for _, opt := range opts {
		opt.apply(&f.opts)
	}

	{
		var err error
		if f.identityCert, err = tls.X509KeyPair(f.publicKey, f.privateKey); err != nil {
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

	if ok := f.caPool.AppendCertsFromPEM(f.publicKey); !ok {
		return nil, ErrNoValidCertificates
	}

	return f, nil
}

// IdentityCert returns the Identity Certificate used for the connection.
func (c *DynamicProvider) IdentityCert() tls.Certificate {
	return c.identityCert
}

// CAPool returns the CA Pool for the connection.
func (c *DynamicProvider) CAPool() *x509.CertPool {
	return c.caPool
}

// ServerOption returns the grpc.ServerOption for use with a new gRPC server.
func (c *DynamicProvider) ServerOption() grpc.ServerOption {
	//nolint:gosec // default minimum is TLS1.3. and skip verify false.
	creds := credentials.NewTLS(&tls.Config{
		ClientCAs:          c.CAPool(),
		ClientAuth:         tls.RequireAndVerifyClientCert,
		Certificates:       []tls.Certificate{c.IdentityCert()},
		MinVersion:         c.opts.minTLSVersion,
		InsecureSkipVerify: c.opts.serverInsecureSkipVerify,
	})

	return grpc.Creds(creds)
}

// DialOption returns the grpc.DialOption used with a gRPC client.
func (c *DynamicProvider) DialOption(serverName string) grpc.DialOption {
	//nolint:gosec // default minimum is TLS1.3. and skip verify false.
	creds := credentials.NewTLS(&tls.Config{
		ServerName:         serverName,
		RootCAs:            c.CAPool(),
		Certificates:       []tls.Certificate{c.IdentityCert()},
		MinVersion:         c.opts.minTLSVersion,
		InsecureSkipVerify: c.opts.dialInsecureSkipVerify,
	})

	return grpc.WithTransportCredentials(creds)
}

func (c *DynamicProvider) generateKeypair() ([]byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, c.opts.dynamicCertKeySize)
	if err != nil {
		return nil, nil, err
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Farnsworth Dynamic Certificate"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(c.opts.dynamicCertLifetime),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	pubBlock := &bytes.Buffer{}
	privBlock := &bytes.Buffer{}
	if err = pem.Encode(pubBlock, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return nil, nil, err
	}
	if err = pem.Encode(
		privBlock,
		&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)},
	); err != nil {
		return nil, nil, err
	}
	// fmt.Println(out.String())
	// out.Reset()
	// pem.Encode(out, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return pubBlock.Bytes(), privBlock.Bytes(), nil
}

func (c *DynamicProvider) KeyPair() (tls.Certificate, error) {
	return tls.X509KeyPair(c.publicKey, c.privateKey)
}
