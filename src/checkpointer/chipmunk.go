package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	gcs "cloud.google.com/go/storage"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
)

type Chipmunk struct {
	Version int
	docker  *docker.Client
	gcs     *gcs.Client
}

func NewChipmunk() *Chipmunk {
	chipmunk := &Chipmunk{}

	// Initialize the docker client
	var err error
	chipmunk.docker, err = docker.NewClientWithOpts(docker.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}

	// Get the list of running containers on the nodes
	ctx := context.Background()
	containers, err := chipmunk.docker.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	var networkMode string
	for _, ctn := range containers {
		// If the application is already running on the node, kill it because it is not being proxied or
		//  checkpointed. This assumes the fact that only one chipmunk pod can be on a node at one time
		if ctn.Image == applicationImage {
			log.Println("killing old container")
			err := chipmunk.docker.ContainerRemove(ctx, ctn.ID, types.ContainerRemoveOptions{Force: true})
			if err != nil {
				panic(err)
			}
		}

		// Get the network of the checkpointer container
		if strings.Contains(ctn.Names[0], "k8s_checkpointer_chipmunk") {
			log.Println("container network", ctn.HostConfig.NetworkMode)
			networkMode = ctn.HostConfig.NetworkMode
		}
	}

	// Create the gcs client
	chipmunk.gcs, err = gcs.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	switch applicationImagePullPolicy {
	case "pull":
		reader, err := chipmunk.docker.ImagePull(ctx, applicationImage, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
		io.Copy(os.Stdout, reader)
		break
	case "load":
		// get the chipmunk bucket, assume it exists
		bucket := chipmunk.gcs.Bucket(bucketName)
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

		imrsp, err := chipmunk.docker.ImageLoad(ctx, r, true)
		if err != nil {
			panic(err)
		}
		log.Println(imrsp)
		break
	}

	// Create the container from the image
	resp, err := chipmunk.docker.ContainerCreate(ctx, &container.Config{
		Image: applicationImage,
		Tty:   true,
	}, &container.HostConfig{
		NetworkMode: container.NetworkMode(networkMode),
	}, nil, "")
	if err != nil {
		panic(err)
	}

	log.Println("container id", resp.ID)

	// Finally run the container
	// TODO: start from checkpoint
	if err := chipmunk.docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return chipmunk
}

func (c *Chipmunk) Checkpoint(version int) {
	log.Printf("Attempting checkpoint: %d\n", version)

	// TODO: dump check point to shared fs with version and not version in name
	err := chipmunk.docker.CheckpointCreate(ctx, "application", types.CheckpointCreateOptions{
		Exit:          false,
		CheckpointID:  fmt.Sprintf("cp-%d", version),
		CheckpointDir: "shared",
	})
	if err != nil {
		panic(err)
	}

	log.Println(" -> Success!")
}
