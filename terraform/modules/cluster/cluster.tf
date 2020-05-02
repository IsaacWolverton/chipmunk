resource "google_container_cluster" "chipmunk-cluster" {
	name     = "chipmunk-cluster"
	project  = var.project
	location = var.zone
	
	remove_default_node_pool = true
	initial_node_count       = 1
	
	network    = var.network
	subnetwork = var.subnet
	cluster_ipv4_cidr = var.pod_ip_range
	
	master_auth {
		username = var.username
		password = var.password
		
		client_certificate_config {
			issue_client_certificate = false
		}
	}
	
	addons_config {
		http_load_balancing {
			disabled = true
		}
		
		horizontal_pod_autoscaling {
			disabled = false
		}
	}
}

resource "kubernetes_cluster_role_binding" "cluster-admin" {
	metadata {
		name = "cluster-role"
	}
	
	role_ref {
		api_group = "rbac.authorization.k8s.io"
		kind = "ClusterRole"
		name = "cluster-admin"
	}
	
	subject {
		kind = "ServiceAccount"
		name = "default"
	}
}