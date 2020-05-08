package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var (
	applicationImage           string
	applicationImagePullPolicy string
	applicationPort            int
	bucketName                 string
	chipmunk                   *Chipmunk
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

	exec.Command("chmod", "777", "/shared")
	exec.Command("gcsfuse", "--implicit-dirs", "chipmunk-storage", "/shared").Output()

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
	p := Server{
		Addr:       ":42069",
		Target:     ":8080",
		PathPrefix: "/sheck",
	}
	go p.ListenAndServe()

	// Checkpoint version
	version := 0

	for {
		select {
		case <-time.After(time.Second * 10):
			p.StopProxy()
			chipmunk.Checkpoint(version)
			p.ResumeProxy()
			version++
			break
		}
	}
}
