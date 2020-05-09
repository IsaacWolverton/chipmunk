package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	applicationImage           string
	applicationImagePullPolicy string
	applicationPort            int
	bucketName                 string
	chipmunk                   *Chipmunk
	checkpointVersion          int
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
	time.Sleep(time.Second * 2)

	// exec.Command("chmod", "777", "/shared")
	// exec.Command("gcsfuse", "--implicit-dirs", "chipmunk-storage", "/shared").Output()

	out, err := exec.Command("python3", "-c", fmt.Sprintf("import os; print(max([int(d.split(\"-\")[-1]) for d in os.listdir(\"/sheck/%s/.\") if \"cp-\" in d]))", applicationImage)).Output()
	if err != nil {
		checkpointVersion = -1
	} else {
		checkpointVersion, _ = strconv.Atoi(strings.TrimSuffix(string(out), "\n"))
	}
	log.Println("version:", checkpointVersion)

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
		PathPrefix: fmt.Sprintf("/sheck/%s/", applicationImage),
	}

	if checkpointVersion != -1 {
		p.ReplayPath = fmt.Sprintf("/sheck/%s/network-%d", applicationImage, checkpointVersion)
		p.SaveFile = fmt.Sprintf("network-%d", checkpointVersion+1)
	} else {
		p.SaveFile = fmt.Sprintf("network-%d", 0)
	}
	go p.ListenAndServe()

	// Checkpoint version
	version := checkpointVersion + 1

	for {
		select {
		case <-time.After(time.Second * 2):
			p.StopProxy(version)
			chipmunk.Checkpoint(version)
			p.ResumeProxy()

			version++
			break
		}
	}
}
