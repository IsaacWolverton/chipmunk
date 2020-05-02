output "endpoint" {
	value = google_container_cluster.chipmunk-cluster.endpoint
	sensitive = true
}

output "client_key" {
	value = google_container_cluster.chipmunk-cluster.master_auth.0.client_key
	sensitive = true
}

output "client_certificate" {
	value = google_container_cluster.chipmunk-cluster.master_auth.0.client_certificate
	sensitive = true
}

output "ca_certifiate" {
	value = google_container_cluster.chipmunk-cluster.master_auth.0.cluster_ca_certificate
	sensitive = true
}

output "service_pool" {
	value = google_container_node_pool.service-pool
}

output "application_pool" {
	value = google_container_node_pool.application-pool
}

