package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	applicationImage           string
	applicationImagePullPolicy string
	applicationPort            int
	bucketName                 string
	chipmunk                   Chipmunk
)

func init() {
	// get the application image name from the environment
	applicationImage = os.Getenv("APPLICATION_IMAGE")
	if applicationImage == "" {
		applicationImage = "application"
	}

	// get the application image name from the environment
	var err error
	applicationPort, err = strconv.Atoi(os.Getenv("APPLICATION_PORT"))
	if err != nil {
		log.Println("failed to parse port number", err)
		applicationPort = 8080
	}

	// get the storage bucket for all chipmunk storage
	bucketName = os.Getenv("BUCKET")
	if bucketName == "" {
		bucketName = "chipmunk-storage"
	}

	// get the image pull policy
	applicationImagePullPolicy = os.Getenv("APPLICATION_IMAGE_PULL_POLICY")
	if applicationImagePullPolicy == "" {
		applicationImagePullPolicy = "pull"
	}

	// wait for all containers to be up and running. TODO: change
	time.Sleep(time.Second * 10)

	chipmunk = NewChipmunk()
}

/**
 * This program will checkpoint the container named "application"
 * every 1 minute.
 *  Parameters: none
 *  Return: none (unreachable during normal execution)
 */
func main() {
	log.Println("Starting checkpointing")

	log.Println("Starting proxy")
	localAddr := ":42069"
	targetAddr := fmt.Sprintf(":%d", applicationPort)

	p := Server{
		Addr:   localAddr,
		Target: targetAddr,
	}
	go p.ListenAndServe()

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

	// Checkpoint version
	version := 0

	for {
		select {
		case <-time.After(time.Second * 10):
			chipmunk.Checkpoint(version)
			version++
			break
		}
	}
}
