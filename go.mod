module github.com/na4ma4/go-certprovider

go 1.23.0

toolchain go1.24.0

retract [v0.3.0, v0.3.3] // invalid dynamic certificate name

retract [v0.3.0, v0.3.3] // invalid dynamic certificate name

require google.golang.org/grpc v1.75.1

require (
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
