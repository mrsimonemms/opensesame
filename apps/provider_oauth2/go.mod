module github.com/mrsimonemms/cloud-native-auth/apps/provider_oauth2

go 1.24.0

replace github.com/mrsimonemms/cloud-native-auth/packages/authentication => ../../packages/authentication

require github.com/mrsimonemms/cloud-native-auth/packages/authentication v0.0.0-00010101000000-000000000000

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mrsimonemms/golang-helpers v0.2.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/cobra v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/grpc v1.71.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
