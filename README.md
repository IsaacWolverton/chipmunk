# chipmunk

MIT 6.824 Spring 2020 Final Project, Ishani, Matt, Isaac

## Project structure

## Requirements
- Docker (tested with v19.03.08)

## Todos

- [x] Create a github repo and a GCP project, add everyone
  - Repo: <https://github.com/IsaacWolverton/chipmunk>
  - GCP: <https://console.cloud.google.com/mit-mic>
- [ ] Setup Terraform config to create a small Kubernetes Cluster on GCP
- [ ] Set terraform backend to a cheap GCP bucket to maintain consistent state <https://www.terraform.io/docs/backends/index.html>
- [x] Create a docker container with a simple program, like the counter
  - `tests/applications/simplecounter/` with matching Dockerfile.
- [ ] Create a Kubernetes config that auto schedules this container (in the form of a deployment) 
- [ ] Run the deployment (the counter should work at this point)
- [ ] Write the CRIU checkpointing scripts
- [ ] Update the Kubernetes config to call startup scripts that attempt to restore from the latest checkpoint (if one exists)
- [ ] Run the deployment and then kill the counter, it should be auto restarted and rescheduled by Kubernetes and our setup script should restore from the latest checkpoint 
- [ ] Update the checkpointing script to log all network traffic
- [ ] Update the startup script to replay all network traffic 
- [ ] Create a simple webserver container that can be used to test the network traffic logging 
- [ ] Run the deployment and test the network traffic logging functionality 
- [ ] Update the checkpointing script to get the current version of all files when checkpointing, write the information to a metadata file associated with the checkpoint files
- [ ] Update the startup script to get the version number associated with that checkpoint and restore the files to that version
- [ ] Create a container that writes to files and closes them to test file versioning 
- [ ] Run the deployment and test the file versioning functionality 
- [ ] More thoroughly test? Maybe even make testing automated? 
- [ ] Build a fault tolerant distributed file system lol 
