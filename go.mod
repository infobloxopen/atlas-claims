module github.com/infobloxopen/atlas-claims

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	google.golang.org/genproto v0.0.0-20200806141610-86f49bd18e98 // indirect
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0 // indirect
)

replace (
	// Avoid this PR that is importing bleeding edge grpc https://github.com/grpc-ecosystem/go-grpc-middleware/commit/5be27de402455e9e74d35ae771c2f4ee983dc726#diff-e72cb4d89d89634404b11bf1d01ac46aafdf8f38463374922ee81142fa5d9606
	github.com/grpc-ecosystem/go-grpc-middleware => github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	// grpc v1.30+ broke many api contracts, the world is not ready for it
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)
