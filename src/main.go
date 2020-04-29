package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

/**
 * This program will checkpoint the container named "application"
 * every 1 minute.
 *  Parameters: none
 *  Return: none (unreachable during normal execution)
 */
func main() {
	log.Println("Starting checkpointing")

	// Create the docker client with API version matching the latest version
	//  availabel on the debian dockerd binary
	cli, err := client.NewClientWithOpts(client.WithVersion("1.40"))
	if err != nil {
		panic(err)
	}

	// Checkpoint version
	version := 0

	for {
		select {
		case <-time.After(time.Second * 10):
			log.Printf("Attempting checkpoint: %d\n", version)

			// TODO: dump check point to shared fs with version and not version in name
			err := cli.CheckpointCreate(context.Background(), "application", types.CheckpointCreateOptions{
				Exit:          false,
				CheckpointID:  fmt.Sprintf("cp-%d", version),
				CheckpointDir: "shared",
			})
			if err != nil {
				panic(err)
			}

			log.Println(" -> Success!")
			version++
			break
		}
	}
}
