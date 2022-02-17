// Package certprovider contains a certprovider for finding client and server certificates easily.
package certprovider

import "errors"

// ErrNoValidCertificates returned when no valid certificates are found in ca.pem.
var ErrNoValidCertificates = errors.New("no valid certificates present")
