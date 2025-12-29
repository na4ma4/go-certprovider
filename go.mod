module github.com/na4ma4/go-certprovider

go 1.24.0

retract [v0.3.0, v0.3.3] // invalid dynamic certificate name

retract [v0.3.0, v0.3.3] // invalid dynamic certificate name

require google.golang.org/grpc v1.78.0

require (
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
