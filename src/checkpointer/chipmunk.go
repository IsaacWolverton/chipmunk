package main

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	gcs "cloud.google.com/go/storage"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	docker "github.com/docker/docker/client"
)

type Chipmunk struct {
	Version   int
	docker    *docker.Client
	gcs       *gcs.Client
	container string
}

func NewChipmunk() *Chipmunk {
	chipmunk := &Chipmunk{}

	//mk mount point dir
	// os.Mkdir("/sheck/mount", 0755)

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
		//get version number and untar??
		// fsFile := bucket.Object(fmt.Sprintf("%s/fs-%d.tar", applicationImage, version))
		// fr, err := fsFile.NewReader(ctx)
		// if err != nil {
		// 	panic(err)
		// }
		// defer fr.Close()
		// err := untar(fr)
		// if err != nil {
		// 	panic(err)
		// }

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
		Tty:   false,
	}, &container.HostConfig{
		NetworkMode: container.NetworkMode(networkMode),
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeBind,
				Source: "/mount",
				Target: "/mountApp",
			},
		},
	}, nil, "")
	if err != nil {
		panic(err)
	}
	chipmunk.container = resp.ID
	log.Println("container id", resp.ID)

	var options types.ContainerStartOptions
	if checkpointVersion != -1 {
		out, err := exec.Command("cp", "-r", fmt.Sprintf("/sheck/%s/cp-%d", applicationImage, checkpointVersion), fmt.Sprintf("/containers/%s/checkpoints/", resp.ID)).Output()
		if err != nil {
			log.Println(err)
		}
		log.Println(string(out))

		options = types.ContainerStartOptions{CheckpointID: fmt.Sprintf("cp-%d", checkpointVersion)}
	} else {
		options = types.ContainerStartOptions{}
	}

	if err := chipmunk.docker.ContainerStart(ctx, resp.ID, options); err != nil {
		// if err := chipmunk.docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	// _, err = exec.Command("docker", "start", resp.ID).Output()
	// if err != nil {
	// 	panic(err)
	// }

	// out, err := exec.Command("gcsfuse", "chipmunk-storage", "/shared").Output()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s\n", out)

	return chipmunk
}

func (c *Chipmunk) Checkpoint(version int) {
hello:
	log.Printf("Attempting checkpoint: %d\n", version)

	// TODO: dump check point to shared fs with version and not version in name
	ctx := context.Background()

	err := c.docker.CheckpointCreate(ctx, c.container, types.CheckpointCreateOptions{
		Exit:          false,
		CheckpointID:  fmt.Sprintf("cp-%d", version),
		CheckpointDir: fmt.Sprintf("/sheck/%s", applicationImage),
	})
	if err != nil {
		log.Println("EROORE: %s", err)
		goto hello
	}
	log.Println(" -> CRIU Checkpoint Success!")

	// tarFunc(version)

	log.Println(" -> Filesystem Snapshot Success!")

	log.Println(" -> Success!")
}

func tarFunc(version int) {
	destinationfile := fmt.Sprintf("/sheck/%s/fs-%d.tar", applicationImage, version)
	sourcedir := "/mount"

	dir, err := os.Open(sourcedir)
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	// get list of files
	files, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}

	// create tar file
	tarfile, err := os.Create(destinationfile)
	if err != nil {
		panic(err)
	}
	defer tarfile.Close()

	var fileWriter io.WriteCloser = tarfile

	tarfileWriter := tar.NewWriter(fileWriter)
	defer tarfileWriter.Close()

	for _, fileInfo := range files {

		if fileInfo.IsDir() {
			continue
		}

		file, err := os.Open(dir.Name() + "/" + fileInfo.Name())
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// prepare the tar header
		header := new(tar.Header)
		header.Name = file.Name()
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()

		err = tarfileWriter.WriteHeader(header)
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(tarfileWriter, file)
		if err != nil {
			panic(err)
		}
	}

}

func untar(fileReader io.Reader) {

	tarBallReader := tar.NewReader(fileReader)

	// Extracting tarred files

	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}

		// get the individual filename and extract to the current directory
		filename := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory
			fmt.Println("Creating directory :", filename)
			err = os.MkdirAll(filename, os.FileMode(header.Mode)) // or use 0755 if you prefer

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		case tar.TypeReg:
			// handle normal file
			fmt.Println("Untarring :", filename)
			writer, err := os.Create(filename)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			io.Copy(writer, tarBallReader)

			err = os.Chmod(filename, os.FileMode(header.Mode))

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			writer.Close()
		default:
			fmt.Printf("Unable to untar type : %c in file %s", header.Typeflag, filename)
		}
	}
}
