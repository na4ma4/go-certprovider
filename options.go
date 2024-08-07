package certprovider

import (
	"crypto/tls"
	"errors"
	"os"
	"path"
	"strings"
	"time"
)

type options struct {
	certFile                 string
	keyFile                  string
	searchPath               []string
	caFile                   string
	loadSystemCA             bool
	minTLSVersion            uint16
	dialInsecureSkipVerify   bool
	serverInsecureSkipVerify bool
	dynamicCertLifetime      time.Duration
	dynamicCertKeySize       int
}

func defaultOptions() options {
	return options{
		certFile:                 "",
		keyFile:                  "",
		searchPath:               []string{},
		caFile:                   "",
		loadSystemCA:             false,
		minTLSVersion:            tls.VersionTLS13,
		dialInsecureSkipVerify:   false,
		serverInsecureSkipVerify: false,
		dynamicCertLifetime:      time.Hour,
		dynamicCertKeySize:       2048, //nolint:mnd // default key size for dynamic certificates.
	}
}

func (o options) fileExistsInPath(fileName string) (string, bool) {
	if len(fileName) == 0 {
		return fileName, false
	}

	fileName = os.ExpandEnv(fileName)

	if path.IsAbs(fileName) {
		return fileName, o.fileExists(fileName)
	}

	for _, sp := range o.searchPath {
		spFile := path.Join(sp, fileName)
		if o.fileExists(spFile) {
			return spFile, true
		}
	}

	return fileName, false
}

func (o options) fileExists(fileName string) bool {
	if len(fileName) == 0 {
		return false
	}

	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		return false
	} else if err == nil {
		return true
	}

	return false
}

func (o options) getCAFile() (string, bool) {
	return o.fileExistsInPath(o.caFile)
}

func (o options) getCertFile() string {
	fileName, _ := o.fileExistsInPath(o.certFile)

	return fileName
}

func (o options) getKeyFile() string {
	fileName, _ := o.fileExistsInPath(o.keyFile)

	return fileName
}

// A Option sets options such as file paths, if a CA should be loaded, etc.
type Option interface {
	apply(*options)
}

// funcOption wraps a function that modifies options into an
// implementation of the Option interface.
type funcOption struct {
	f func(*options)
}

func (fdo *funcOption) apply(do *options) {
	fdo.f(do)
}

func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// ProviderFromString returns a CertificateProviderType from a supplied string.
func ProviderFromString(in string, defaultProvider Option) Option {
	switch strings.ToLower(in) {
	case "serverprovider", "server":
		return ServerProvider()
	case "clientprovider", "client":
		return ClientProvider()
	case "certprovider", "cert":
		return CertProvider()
	}

	return defaultProvider
}

// ClientProvider sets the file names to the defaults for a mTLS Client.
func ClientProvider() Option {
	return newFuncOption(func(o *options) {
		o.caFile = "ca.pem" //nolint:goconst // less readable as a constant.
		o.certFile = "client.pem"
		o.keyFile = "client-key.pem"
		o.loadSystemCA = true
	})
}

// ServerProvider sets the file names to the defaults for a mTLS Server.
func ServerProvider() Option {
	return newFuncOption(func(o *options) {
		o.caFile = "ca.pem"
		o.certFile = "server.pem"
		o.keyFile = "server-key.pem"
	})
}

// CertProvider sets the file names to the defaults for a mTLS Server.
func CertProvider() Option {
	return newFuncOption(func(o *options) {
		o.caFile = "ca.pem"
		o.certFile = "cert.pem"
		o.keyFile = "key.pem"
	})
}

// CertFilename sets the certificate filename to a specific filename.
func CertFilename(filename string) Option {
	return newFuncOption(func(o *options) {
		o.certFile = filename
	})
}

// KeyFilename sets the key filename to a specific filename.
func KeyFilename(filename string) Option {
	return newFuncOption(func(o *options) {
		o.keyFile = filename
	})
}

// CAFilename sets the certificate authority filename to a specific filename.
func CAFilename(filename string) Option {
	return newFuncOption(func(o *options) {
		o.caFile = filename
	})
}

// UseSystemCAPool sets whether the provider should include the system CA pool.
func UseSystemCAPool(enable bool) Option {
	return newFuncOption(func(o *options) {
		o.loadSystemCA = enable
	})
}

// AddSearchPath adds a search path for the files.
func AddSearchPath(path string) Option {
	return newFuncOption(func(o *options) {
		if o.searchPath == nil {
			o.searchPath = []string{}
		}

		o.searchPath = append(o.searchPath, path)
	})
}

// MinTLSVersion sets a minimum TLS version.
func MinTLSVersion(tlsVer uint16) Option {
	return newFuncOption(func(opt *options) {
		switch tlsVer {
		case tls.VersionTLS10,
			tls.VersionTLS11,
			tls.VersionTLS12,
			tls.VersionTLS13:
			opt.minTLSVersion = tlsVer
		}
	})
}

// InsecureSkipVerifyOnDial sets the InsecureSkipVerify on the DialOptions.
func InsecureSkipVerifyOnDial(verify bool) Option {
	return newFuncOption(func(o *options) {
		o.dialInsecureSkipVerify = verify
	})
}

// InsecureSkipVerifyOnServer sets the InsecureSkipVerify on the ServerOptions.
func InsecureSkipVerifyOnServer(verify bool) Option {
	return newFuncOption(func(o *options) {
		o.serverInsecureSkipVerify = verify
	})
}

// DynamicCertLifetime sets the lifetime of a dynamic certificate.
func DynamicCertLifetime(certLifetime time.Duration) Option {
	return newFuncOption(func(o *options) {
		o.dynamicCertLifetime = certLifetime
	})
}

// DynamicCertKeySize sets the key size of a dynamic certificate.
func DynamicCertKeySize(certKeySize int) Option {
	return newFuncOption(func(o *options) {
		o.dynamicCertKeySize = certKeySize
	})
}
