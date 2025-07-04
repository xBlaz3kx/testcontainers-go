package neo4j

import (
	"context"
	"fmt"
	"net/http"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	// containerPorts {
	defaultBoltPort  = "7687/tcp"
	defaultHTTPPort  = "7474/tcp"
	defaultHTTPSPort = "7473/tcp"
	// }
)

// Neo4jContainer represents the Neo4j container type used in the module
type Neo4jContainer struct {
	testcontainers.Container
}

// BoltUrl returns the bolt url for the Neo4j container, using the bolt port, in the format of neo4j://host:port
//
//nolint:revive,staticcheck //FIXME
func (c Neo4jContainer) BoltUrl(ctx context.Context) (string, error) {
	return c.PortEndpoint(ctx, defaultBoltPort, "neo4j")
}

// Deprecated: use Run instead
// RunContainer creates an instance of the Neo4j container type
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*Neo4jContainer, error) {
	return Run(ctx, "neo4j:4.4", opts...)
}

// Run creates an instance of the Neo4j container type
func Run(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (*Neo4jContainer, error) {
	request := testcontainers.ContainerRequest{
		Image: img,
		Env: map[string]string{
			"NEO4J_AUTH": "none",
		},
		ExposedPorts: []string{
			defaultBoltPort,
			defaultHTTPPort,
			defaultHTTPSPort,
		},
		WaitingFor: &wait.MultiStrategy{
			Strategies: []wait.Strategy{
				wait.NewLogStrategy("Bolt enabled on"),
				&wait.HTTPStrategy{
					Port:              defaultHTTPPort,
					StatusCodeMatcher: isHTTPOk(),
				},
			},
		},
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	}

	if len(opts) == 0 {
		opts = append(opts, WithoutAuthentication())
	}

	for _, option := range opts {
		if err := option.Customize(&genericContainerReq); err != nil {
			return nil, err
		}
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	var c *Neo4jContainer
	if container != nil {
		c = &Neo4jContainer{Container: container}
	}

	if err != nil {
		return c, fmt.Errorf("generic container: %w", err)
	}

	return c, nil
}

func isHTTPOk() func(status int) bool {
	return func(status int) bool {
		return status == http.StatusOK
	}
}
