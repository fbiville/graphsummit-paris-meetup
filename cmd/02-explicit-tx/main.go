package main

import (
	"context"
	"github.com/fbiville/graphsummit-paris-meetup/pkg/container"
	"github.com/fbiville/graphsummit-paris-meetup/pkg/errors"
)

func main() {
	ctx := context.Background()
	singleInstance, _, err := container.StartSingleInstance(ctx, container.ContainerConfiguration{
		Neo4jVersion: "4.4-enterprise",
		Username:     "neo4j",
		Password:     "s3cr3t",
		Scheme:       "neo4j",
	})
	errors.PanicOnErr(err)
	defer func() {
		errors.PanicOnErr(singleInstance.Terminate(ctx))
	}()
	// TODO: run query in transaction
	// TODO: run autocommit query in transaction
}
