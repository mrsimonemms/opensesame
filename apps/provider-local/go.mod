module github.com/mrsimonemms/opensesame/apps/provider-local

go 1.24.1

replace github.com/mrsimonemms/opensesame/packages/authentication => ../../packages/authentication

replace github.com/mrsimonemms/opensesame/packages/go-sdk => ../../packages/go-sdk

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/go-playground/validator/v10 v10.26.0
	github.com/mrsimonemms/opensesame/packages/authentication v0.0.0-20250403160226-1416cdeec040
	github.com/mrsimonemms/opensesame/packages/go-sdk v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.34.0
	go.mongodb.org/mongo-driver/v2 v2.1.0
	google.golang.org/grpc v1.71.1
	sigs.k8s.io/yaml v1.4.0
)

require (
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mrsimonemms/golang-helpers v0.2.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/cobra v1.9.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
