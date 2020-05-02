resource "google_container_node_pool" "service-pool" {
	name = "service-pool"
	
	project    = var.project
	location   = var.zone
	cluster    = google_container_cluster.chipmunk-cluster.name

	management {
		auto_repair = true
	}
	
	node_count = 2
	node_config {
		machine_type = "n1-standard-1"
		
		oauth_scopes = [
			"https://www.googleapis.com/auth/logging.write",
			"https://www.googleapis.com/auth/monitoring",
		]    
		
		metadata = {
			disable-legacy-endpoints = "true"
		} 

		preemptible = true # lmao
	}
	
	autoscaling {
		min_node_count = 1
		max_node_count = 2
	}
}

resource "google_container_node_pool" "application-pool" {
	name = "application-pool"
	
	project    = var.project
	location   = var.zone
	cluster    = google_container_cluster.chipmunk-cluster.name
	
	node_count = 1
	node_config {
		machine_type = "n1-standard-1"
		image_type = "UBUNTU"

		oauth_scopes = [
			// default permissions for kubernetes
			"https://www.googleapis.com/auth/logging.write",
			"https://www.googleapis.com/auth/monitoring",

			// gcs -- can this be read-write?
			"https://www.googleapis.com/auth/devstorage.full_control",
		]    
		
		metadata = {
			disable-legacy-endpoints = "true"
		}
		
		taint {
			effect = "NO_SCHEDULE"
			key    = "configured"
			value  = "false"
		}
	}

	autoscaling {
		min_node_count = 1
		max_node_count = 1
	}
}