module "network" {
	source = "./modules/network"

	project = var.project
	region  = var.region

	ip_range           = var.ip_range
	secondary_ip_range = var.secondary_ip_range
}

module "cluster" {
	source = "./modules/cluster"

	project = var.project
	region  = var.region
	zone    = var.zone

	subnet  = module.network.subnet_self_link
	network = module.network.network_self_link

	pod_ip_range = var.pod_ip_range
}

module "k8s" {
	source = "./modules/k8s"

	project = var.project

	application_pool = module.cluster.application_pool
	service_pool     = module.cluster.service_pool

	application_image = var.application_image
	application_path  = var.application_path
	application_port  = var.application_port
}