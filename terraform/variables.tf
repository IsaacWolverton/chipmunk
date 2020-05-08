/**
 * Id of GCP project where everything will be deployed to
 */
variable "project" {
	type = string
	description = "google provider project id"
}

/**
 * Region where location based resources are deployed in
 */
variable "region" {
	type = string
	description = "region in which chipmunk cluster will be hosted"
}

/**
 * Zone where zone based resources are deployed in
 */
variable "zone" {
	type = string
	description = "zone in which chipmunk cluster will be hosted"
}

/**
 * Name of the application that is being ran in chipmunk
 */
variable "application_image" {
	type = string
	description = "name of the application that is being ran in chipmunk"
}

/**
 * Whether to pull the application or load from fs
 */
variable "application_image_pull_policy" {
	type = string
	description = "whether to pull the application [pull] or load from fs [load]"
	default = "pull"
}

/**
 * Path to the saved application image
 */
variable "application_path" {
	type = string
	description = "path to the saved application image"
	default = "../tests/application.tar"
}

/**
 * Port of the application
 */
variable "application_port" {
	type = number
	description = "port of the application"
	default = 8080
}

/**
 * Custom range for the deployment's network
 */
variable "ip_range" {
	type = string
	description = "cidr definition for subnetwork"
	default = "10.10.10.0/24"
}

/**
 * Custom secondary range for the deployment's network
 */
variable "secondary_ip_range" {
	type = string
	description = "cidr definition for cluster subnetwork"
	default = "10.11.11.0/24"
}

/**
 * Custom range for the cluster's pod network
 */
variable "pod_ip_range" {
	type = string
	description = "cidr definition for cluster pods"
	default = "10.12.0.0/14"
}