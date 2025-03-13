module github.com/na4ma4/go-certprovider

go 1.22.2
toolchain go1.23.7

retract [v0.3.0, v0.3.3] // invalid dynamic certificate name

require google.golang.org/grpc v1.71.0

require (
	golang.org/x/net v0.36.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250204164813-702378808489 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
