module github.com/infobloxopen/atlas-claims

go 1.14

require (
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/spiffe/go-spiffe/v2 v2.1.4
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/oauth2 v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20230320184635-7606e756e683 // indirect
	google.golang.org/grpc v1.54.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

// Prevent spire from increasing grpc version
replace google.golang.org/grpc => google.golang.org/grpc v1.34.0
