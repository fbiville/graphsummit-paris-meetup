package container

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"strings"
)

type ContainerConfiguration struct {
	Neo4jVersion string
	Username     string
	Password     string
	Scheme       string
}

func (config ContainerConfiguration) neo4jAuthEnvVar() string {
	return fmt.Sprintf("%s/%s", config.Username, config.Password)
}

func (config ContainerConfiguration) neo4jAuthToken() neo4j.AuthToken {
	return neo4j.BasicAuth(config.Username, config.Password, "")
}

func StartCluster(ctx context.Context, config ContainerConfiguration) (testcontainers.DockerCompose, neo4j.Driver, error) {
	identifier := strings.ToLower(uuid.New().String())
	compose := testcontainers.NewLocalDockerCompose([]string{"etc/docker-compose.yml"}, identifier)
	version := config.Neo4jVersion
	execution := compose.
		WithCommand([]string{"up", "-d"}).WithEnv(map[string]string{
		"NEO4J_USERNAME": config.Username,
		"NEO4J_PASSWORD": config.Password,
		"NEO4J_VERSION":  version,
	}).WaitForService("core1", boltReadyStrategy()).Invoke()
	if execution.Error != nil {
		return nil, nil, execution.Error
	}
	driver, err := newDriver(config.Scheme, 7687, config.neo4jAuthToken())
	return compose, driver, err
}

func StartSingleInstance(ctx context.Context, config ContainerConfiguration) (testcontainers.Container, neo4j.Driver, error) {
	version := config.Neo4jVersion
	request := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("neo4j:%s", version),
		ExposedPorts: []string{"7687/tcp"},
		Env: map[string]string{
			"NEO4J_AUTH":                     config.neo4jAuthEnvVar(),
			"NEO4J_ACCEPT_LICENSE_AGREEMENT": "yes",
		},
		WaitingFor: boltReadyStrategy(),
	}
	container, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: request,
			Started:          true,
		})
	if err != nil {
		return nil, nil, err
	}
	driver, err := newNeo4jDriver(ctx, config.Scheme, container, config.neo4jAuthToken())
	return container, driver, err
}

func boltReadyStrategy() *wait.LogStrategy {
	return wait.ForLog("Bolt enabled")
}

func newNeo4jDriver(ctx context.Context, scheme string, container testcontainers.Container, auth neo4j.AuthToken) (neo4j.Driver, error) {
	port, err := container.MappedPort(ctx, "7687")
	if err != nil {
		return nil, err
	}
	return newDriver(scheme, port.Int(), auth)
}

func newDriver(scheme string, port int, auth neo4j.AuthToken) (neo4j.Driver, error) {
	_ = fmt.Sprintf("%s://localhost:%d", scheme, port)
	panic("Implement me")
}

func withDebugLogging(config *neo4j.Config) {
	config.Log = neo4j.ConsoleLogger(neo4j.DEBUG)
}
