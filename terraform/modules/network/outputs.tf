/**
 * Output of the reference for the network that cluster nodes will
 * be attached to 
 */
output "network_self_link" {
	value = google_compute_network.chipmunk-network.self_link
}

/**
 * Output of the reference for the subnet that cluster nodes will
 * be attached to 
 */
output "subnet_self_link" {
	value = google_compute_subnetwork.chipmunk-subnet.self_link
}