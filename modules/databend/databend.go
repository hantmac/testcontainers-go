package databend

import (
	"context"
	"fmt"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	databendUser        = "databend"
	defaultUser         = "databend"
	defaultPassword     = "databend"
	defaultDatabaseName = "default"
)

// DatabendContainer represents the Databend container type used in the module
type DatabendContainer struct {
	testcontainers.Container
	username string
	password string
	database string
}

func WithDefaultCredentials() testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		username := req.Env["QUERY_DEFAULT_USER"]
		password := req.Env["QUERY_DEFAULT_PASSWORD"]
		if len(password) == 0 {
			req.Env["QUERY_DEFAULT_PASSWORD"] = defaultPassword
		}
		if len(username) == 0 {
			req.Env["QUERY_DEFAULT_USER"] = defaultUser
		}
		return nil
	}
}

// Run creates an instance of the Databend container type
func Run(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (*DatabendContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        img,
		ExposedPorts: []string{"8000/tcp"},
		Env: map[string]string{
			"QUERY_DEFAULT_USER":     defaultUser,
			"QUERY_DEFAULT_PASSWORD": defaultPassword,
		},
		WaitingFor: wait.ForLog("port: 8000  Databend Query Server"),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	opts = append(opts, WithDefaultCredentials())

	for _, opt := range opts {
		if err := opt.Customize(&genericContainerReq); err != nil {
			return nil, err
		}
	}

	username, ok := req.Env["QUERY_DEFAULT_USER"]
	if !ok {
		username = defaultUser
	}
	password := req.Env["QUERY_DEFAULT_PASSWORD"]

	if len(password) == 0 && password == "" && !strings.EqualFold(defaultUser, username) {
		return nil, fmt.Errorf("empty password and user")
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	var c *DatabendContainer
	if container != nil {
		c = &DatabendContainer{
			Container: container,
			password:  password,
			username:  username,
			database:  defaultDatabaseName,
		}
	}

	if err != nil {
		return c, fmt.Errorf("generic container: %w", err)
	}

	return c, nil
}

// MustConnectionString panics if the address cannot be determined.
func (c *DatabendContainer) MustConnectionString(ctx context.Context, args ...string) string {
	addr, err := c.ConnectionString(ctx, args...)
	if err != nil {
		panic(err)
	}
	return addr
}

func (c *DatabendContainer) ConnectionString(ctx context.Context, args ...string) (string, error) {
	containerPort, err := c.MappedPort(ctx, "8000/tcp")
	if err != nil {
		return "", err
	}

	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}

	extraArgs := ""
	if len(args) > 0 {
		extraArgs = strings.Join(args, "&")
	}
	if extraArgs != "" {
		extraArgs = "?" + extraArgs
	}

	// databend://databend:databend@localhost:8000/default?sslmode=disable
	connectionString := fmt.Sprintf("databend://%s:%s@%s:%s/%s%s", c.username, c.password, host, containerPort.Port(), c.database, extraArgs)
	return connectionString, nil
}

func WithUsername(username string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Env["QUERY_DEFAULT_USER"] = username
		return nil
	}
}

func WithPassword(password string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Env["QUERY_DEFAULT_PASSWORD"] = password
		return nil
	}
}
