package servicebus

import (
	"io"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mssql"
)

type options struct {
	mssqlImage     string
	mssqlOptions   []testcontainers.ContainerCustomizer
	mssqlContainer *mssql.MSSQLServerContainer
	network        *testcontainers.DockerNetwork
}

func defaultOptions() options {
	return options{
		mssqlImage:     defaultMSSQLImage,
		mssqlContainer: nil,
	}
}

// Satisfy the testcontainers.CustomizeRequestOption interface
var _ testcontainers.ContainerCustomizer = (Option)(nil)

// Option is an option for the Redpanda container.
type Option func(*options) error

// Customize is a NOOP. It's defined to satisfy the testcontainers.ContainerCustomizer interface.
func (o Option) Customize(*testcontainers.GenericContainerRequest) error {
	// NOOP to satisfy interface.
	return nil
}

// WithMSSQL sets the image and options for the MSSQL container.
// By default, the image is "mcr.microsoft.com/mssql/server:2022-CU14-ubuntu-22.04".
func WithMSSQL(img string, opts ...testcontainers.ContainerCustomizer) Option {
	return func(o *options) error {
		o.mssqlImage = img
		o.mssqlOptions = opts
		return nil
	}
}

// WithAcceptEULA sets the ACCEPT_EULA environment variable to "Y" for the eventhubs container.
func WithAcceptEULA() testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Env["ACCEPT_EULA"] = "Y"

		return nil
	}
}

// WithConfig sets the eventhubs config file for the eventhubs container,
// copying the content of the reader to the container file at
// "/ServiceBus_Emulator/ConfigFiles/Config.json".
func WithConfig(r io.Reader) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Files = append(req.Files, testcontainers.ContainerFile{
			Reader:            r,
			ContainerFilePath: containerConfigFile,
			FileMode:          0o644,
		})

		return nil
	}
}
