package main

import (
	"context"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
)

/**
 * This program will checkpoint the container named "application"
 * every 1 minute.
 *  Parameters: none
 *  Return: none (unreachable during normal execution)
 */
func main() {
	log.Println("Starting checkpointing")
	time.Sleep(time.Second * 10)
	ctx := context.Background()

	// Create the docker client with API version matching the latest version
	//  availabel on the kubernetes node
	dockerClient, err := docker.NewClientWithOpts(docker.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}

	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	var containerNetwork container.NetworkMode
	for _, ctn := range containers {
		log.Println(ctn.Names[0])
		if ctn.Image == "strm/helloworld-http" {
			err := dockerClient.ContainerRemove(ctx, ctn.ID, types.ContainerRemoveOptions{Force: true})
			if err != nil {
				panic(err)
			}
		}

		if strings.Contains(ctn.Names[0], "k8s_proxy_chipmunk") {
			log.Println(ctn.HostConfig.NetworkMode)
			containerNetwork = container.NetworkMode(ctn.HostConfig.NetworkMode)
		}
	}

	log.Println(containerNetwork)

	reader, err := dockerClient.ImagePull(ctx, "strm/helloworld-http", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: "strm/helloworld-http",
		Tty:   true,
	}, &container.HostConfig{
		NetworkMode: containerNetwork,
	}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	// // Create the gcs client
	// gcsClient, err := gcs.NewClient(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// // check if bucket exists
	// bucket := gcsClient.Bucket("test-chipmunk-bucket")
	// exists, err := bucket.Attrs(ctx)
	// if err != nil {
	// 	timeoutctx, cancel := context.WithTimeout(ctx, time.Second*10)
	// 	defer cancel()
	// 	if err := bucket.Create(timeoutctx, "mit-mic", nil); err != nil {
	// 		log.Fatalf("Failed to create bucket: %v", err)
	// 	}
	// }

	// log.Println(exists)

	time.Sleep(time.Hour * 100)

	// // Checkpoint version
	// version := 0

	// for {
	// 	select {
	// 	case <-time.After(time.Second * 10):
	// 		log.Printf("Attempting checkpoint: %d\n", version)

	// 		// TODO: dump check point to shared fs with version and not version in name
	// 		err := dockerClient.CheckpointCreate(ctx, "application", types.CheckpointCreateOptions{
	// 			Exit:          false,
	// 			CheckpointID:  fmt.Sprintf("cp-%d", version),
	// 			CheckpointDir: "shared",
	// 		})
	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		log.Println(" -> Success!")
	// 		version++
	// 		break
	// 	}
	// }
}
