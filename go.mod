module github.com/na4ma4/go-certprovider

go 1.24.0

retract [v0.3.0, v0.3.3] // invalid dynamic certificate name

retract [v0.3.0, v0.3.3] // invalid dynamic certificate name

require google.golang.org/grpc v1.76.0

require (
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
