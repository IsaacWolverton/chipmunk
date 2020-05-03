/**
 * Kubernetes cluster endpoint, used for kubectl, k8s module, etc
 */
output "endpoint" {
	value = google_container_cluster.chipmunk-cluster.endpoint
	sensitive = true
}

/**
 * Key for authenticating control of the cluster
 */
output "client_key" {
	value = google_container_cluster.chipmunk-cluster.master_auth.0.client_key
	sensitive = true
}

/**
 * Certificate for authenticating control of the cluster
 */
output "client_certificate" {
	value = google_container_cluster.chipmunk-cluster.master_auth.0.client_certificate
	sensitive = true
}

/**
 * Cluster's ca certificate
 */
output "ca_certifiate" {
	value = google_container_cluster.chipmunk-cluster.master_auth.0.cluster_ca_certificate
	sensitive = true
}

/**
 * Reference to the service pool used for enforcing hiearchy during terraform apply
 */
output "service_pool" {
	value = google_container_node_pool.service-pool
}

/**
 * Reference to the application pool used for enforcing hiearchy during terraform apply
 */
output "application_pool" {
	value = google_container_node_pool.application-pool
}
