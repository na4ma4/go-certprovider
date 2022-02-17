package certprovider

// // CertificateProviderType is the enum for certificate providers.
// type CertificateProviderType int

// const (
// 	// ServerProviderType uses server.pem/server-key.pem.
// 	ServerProviderType CertificateProviderType = iota
// 	// ClientProviderType uses client.pem/client-key.pem.
// 	ClientProviderType
// 	// CertProviderType uses cert.pem/key.pem.
// 	CertProviderType
// )

// // String returns the string representation of CertificateProviderType.
// func (c CertificateProviderType) String() string {
// 	switch c {
// 	case ServerProviderType:
// 		return "ServerProvider"
// 	case ClientProviderType:
// 		return "ClientProvider"
// 	case CertProviderType:
// 		return "CertProvider"
// 	}

// 	return ""
// }

// // CertProviderFromString returns a CertificateProviderType from a supplied string.
// func CertProviderFromString(in string) CertificateProviderType {
// 	switch strings.ToLower(in) {
// 	case "serverprovider", "server":
// 		return ServerProviderType
// 	case "clientprovider", "client":
// 		return ClientProviderType
// 	case "certprovider", "cert":
// 		return CertProviderType
// 	}

// 	return ClientProviderType
// }
