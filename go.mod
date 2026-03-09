module github.com/na4ma4/go-certprovider

go 1.24.0

retract [v0.3.0, v0.3.3] // invalid dynamic certificate name

retract [v0.3.0, v0.3.3] // invalid dynamic certificate name

require google.golang.org/grpc v1.79.2

require (
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
