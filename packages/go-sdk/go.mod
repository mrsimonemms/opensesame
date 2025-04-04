module github.com/mrsimonemms/opensesame/packages/go-sdk

go 1.24.1

replace github.com/mrsimonemms/opensesame/packages/authentication => ../authentication

require (
	github.com/mrsimonemms/golang-helpers v0.2.1
	github.com/mrsimonemms/opensesame/packages/authentication v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.8.1
	golang.org/x/crypto v0.36.0
	google.golang.org/grpc v1.71.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
