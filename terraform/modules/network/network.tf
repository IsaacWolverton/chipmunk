resource "google_compute_network" "chipmunk-network" {
	name    = "chipmunk-network"
	project = var.project
	
	auto_create_subnetworks = "false"
}

resource "google_compute_subnetwork" "chipmunk-subnet" {
  name    = "chipmunk-subnet"
  project = var.project
  region  = var.region

  network       = google_compute_network.chipmunk-network.name
  ip_cidr_range = var.ip_range
}