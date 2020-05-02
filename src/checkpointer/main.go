package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	gcs "cloud.google.com/go/storage"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
)

var (
	applicationImage string
	bucketName       string
	networkMode      string

	dockerClient *docker.Client
	gcsClient    *gcs.Client
)

func init() {
	// get the application image name from the environment
	applicationImage = os.Getenv("APPLICATION_IMAGE")
	if applicationImage == "" {
		applicationImage = "application"
	}

	// get the storage bucket for all chipmunk storage
	bucketName = os.Getenv("BUCKET")
	if bucketName == "" {
		bucketName = "chipmunk-storage"
	}

	// wait for all containers to be up and running. TODO: change
	time.Sleep(time.Second * 10)

	// Initialize the docker client
	var err error
	dockerClient, err = docker.NewClientWithOpts(docker.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}

	// Get the list of running containers on the nodes
	ctx := context.Background()
	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, ctn := range containers {
		// If the application is already running on the node, kill it because it is not being proxied or
		//  checkpointed. This assumes the fact that only one chipmunk pod can be on a node at one time
		if ctn.Image == applicationImage {
			log.Println("killing old container")
			err := dockerClient.ContainerRemove(ctx, ctn.ID, types.ContainerRemoveOptions{Force: true})
			if err != nil {
				panic(err)
			}
		}

		// Get the network of the proxy container
		if strings.Contains(ctn.Names[0], "k8s_proxy_chipmunk") {
			log.Println("container network", ctn.HostConfig.NetworkMode)
			networkMode = ctn.HostConfig.NetworkMode
		}
	}

	// Create the gcs client
	gcsClient, err = gcs.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	// get the chipmunk bucket, assume it exists
	bucket := gcsClient.Bucket(bucketName)
	_, err = bucket.Attrs(ctx)
	if err != nil {
		panic(err)
	}
	imageFile := bucket.Object(fmt.Sprintf("%s/application.tar", applicationImage))
	r, err := imageFile.NewReader(ctx)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	// // docker pull
	// reader, err := dockerClient.ImagePull(ctx, "strm/helloworld-http", types.ImagePullOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	// io.Copy(os.Stdout, reader)

	imrsp, err := dockerClient.ImageLoad(ctx, r, true)
	if err != nil {
		panic(err)
	}
	log.Println(imrsp)

	// Create the container from the image
	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: applicationImage,
		Tty:   true,
	}, &container.HostConfig{
		NetworkMode: container.NetworkMode(networkMode),
	}, nil, "")
	if err != nil {
		panic(err)
	}

	// Finally run the container
	// TODO: start from checkpoint
	if err := dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

}

/**
 * This program will checkpoint the container named "application"
 * every 1 minute.
 *  Parameters: none
 *  Return: none (unreachable during normal execution)
 */
func main() {
	log.Println("Starting checkpointing")
	// time.Sleep(time.Second * 10)
	// ctx := context.Background()

	// // Create the docker client with API version matching the latest version
	// //  availabel on the kubernetes node
	// dockerClient, err := docker.NewClientWithOpts(docker.WithVersion("1.39"))
	// if err != nil {
	// 	panic(err)
	// }

	// // Create the gcs client
	//

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
