terraform {
	required_version = ">= 0.12"
	
	backend "gcs" {
		bucket  = "chipmunk-tf"
		prefix  = "terraform/state"
	}
}

provider "google" {}

data "google_client_config" "current" { }

provider "kubernetes" {
	load_config_file = false
	host = module.cluster.endpoint
	
	client_key = base64decode(module.cluster.client_key)
	client_certificate = base64decode(module.cluster.client_certificate)
	cluster_ca_certificate = base64decode(module.cluster.ca_certifiate)
	token = data.google_client_config.current.access_token
}