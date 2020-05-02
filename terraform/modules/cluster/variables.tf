// See main definition
variable "project" { }

// See main definition
variable "region" { }

// See main definition
variable "zone" { }

// See main definition
variable "subnet" { }

// See main definition
variable "network" { }

// See main definition
variable "pod_ip_range" { }

/**
 * Master auth username for the chipmunk cluster
 */
variable "username" {
    type = string 
    default = ""
    description = "master auth username for the chipmunk cluster"
}

/**
 * Master auth password for the chipmunk cluster
 */
variable "password" {
    type = string
    default = ""
    description = "master auth password for the chipmunk cluster"
}