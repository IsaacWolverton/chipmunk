output "network_self_link" {
	value = google_compute_network.chipmunk-network.self_link
}

output "subnet_self_link" {
	value = google_compute_subnetwork.chipmunk-subnet.self_link
}