package main

import (
	"context"
	"fmt"
	"github.com/fbiville/graphsummit-paris-meetup/pkg/container"
	"github.com/fbiville/graphsummit-paris-meetup/pkg/errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func main() {
	ctx := context.Background()
	// single instance
	singleInstance, _, err := container.StartSingleInstance(ctx, container.ContainerConfiguration{
		Neo4jVersion: "4.4-enterprise",
		Username:     "neo4j",
		Password:     "s3cr3t",
		Scheme:       "bolt",
	})
	errors.PanicOnErr(err)
	defer func() {
		errors.PanicOnErr(singleInstance.Terminate(ctx))
	}()
	var result neo4j.Result
	// TODO: run query
	errors.PanicOnErr(err)
	printResult(result)

	// cluster
	// run 2 basic queries twice
	cluster, _, err := container.StartCluster(ctx, container.ContainerConfiguration{
		Neo4jVersion: "4.4-enterprise",
		Username:     "neo4j",
		Password:     "s3cr3t",
		Scheme:       "neo4j",
	})
	defer func() {
		errors.PanicOnErr(cluster.Down().Error)
	}()
	// TODO: run 2 queries
}

func printResult(result neo4j.Result) {
	record, err := result.Single()
	errors.PanicOnErr(err)
	value, _ := record.Get("result")
	fmt.Println(value)
}
