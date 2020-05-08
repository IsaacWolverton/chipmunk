resource "kubernetes_daemonset" "configurator" {
	depends_on = [
		var.application_pool
	]
	
	metadata {
		name = "node-configurator"
		labels = {
	  		App = "configurator"
		}
	}
	
	spec {
		selector {
			match_labels = {
	   			App = "configurator"
			}
		}
		
		template {
			metadata {
				labels = {
	  				App = "configurator"
				}
			}
			
			spec {
				toleration {
					key      = "configured"
					operator = "Equal"
					value    = "false"
					effect   = "NoSchedule"
				}
				
				affinity {
					node_affinity {
						required_during_scheduling_ignored_during_execution {
							node_selector_term {
								match_expressions {
									key      = "cloud.google.com/gke-nodepool"
									operator = "In"
									values   = ["application-pool"]
								}
							}
						}
					}
				}

				// TODO: enforce sensitive service account name from cluster module
				automount_service_account_token = true
				service_account_name            = "default"
				
				container {
					image = "gcr.io/mit-mic/configurator:v1"
					name  = "configurator"
					
					resources {
						limits {
							cpu    = "0.5"
							memory = "512Mi"
						}
						
						requests {
			  				cpu    = "250m"
			  				memory = "50Mi"
						}
		  			}

					security_context {
						privileged = true
						allow_privilege_escalation = true
					}

					volume_mount {
						name = "host-root"
						mount_path = "/host"
					}

					image_pull_policy = "Always"
				}

				volume {
					name = "host-root"
					host_path {
						path = "/"
					}
				}

				host_network = true
				host_pid = true
			}
		}
	}
}

resource "kubernetes_pod" "chipmunk" {
	depends_on = [
		var.application_pool,
		google_storage_bucket_object.chipmunk-application
	]

	metadata {
		name = "chipmunk"
		labels = {
	  		App = "chipper"
		}
	}

	spec {
		# host_network = true # Allows for host node to be exposed, nice for testing restore occured on another node
		toleration {
			key      = "configured"
			operator = "Equal"
			value    = "true"
			effect   = "NoSchedule"
		}
		
		affinity {
			node_affinity {
				required_during_scheduling_ignored_during_execution {
					node_selector_term {
						match_expressions {
							key      = "cloud.google.com/gke-nodepool"
							operator = "In"
							values   = ["application-pool"]
						}
					}
				}
			}
		}
		
		container {
			image = "gcr.io/mit-mic/checkpointer:v1"
			name  = "checkpointer"

			image_pull_policy = "Always"

			security_context {
				privileged = true
				allow_privilege_escalation = true
			}

			volume_mount {
				name = "dockerd"
				mount_path = "/var/run/docker.sock"
			}

			volume_mount {
				name = "containerd"
				mount_path = "/containers"
			}

			volume_mount {
				name = "shared"
				mount_path = "/sheck"
			}

			env {
				name  = "APPLICATION_IMAGE"
				value = var.application_image
			}

			env {
				name  = "APPLICATION_PORT"
				value = var.application_port
			}

			env {
				name  = "APPLICATION_IMAGE_PULL_POLICY"
				value = var.application_image_pull_policy
			}

			env {
				name  = "BUCKET"
				value = google_storage_bucket.chipmunk-storage.name
			}
		}

		volume {
			name = "dockerd"
			host_path {
				path = "/var/run/docker.sock"
			}
		}

		volume {
			name = "containerd"
			host_path {
				path = "/var/lib/docker/containers"
			}
		}

		volume {
			name = "shared"
			host_path {
				path = "/sheck"
			}
		}

		automount_service_account_token = true
		service_account_name = "default"
	}
}


/**
 * Test node with priviledge access to host for testing
 * TODO: remove!
 * TODO: keep! 
 */
resource "kubernetes_daemonset" "test-root-pod" {
	depends_on = [
		var.application_pool
	]
	
	metadata {
		name = "root-pod"
		labels = {
	  		App = "tester"
		}
	}
	
	spec {
		selector {
			match_labels = {
	   			App = "tester"
			}
		}
		
		template {
			metadata {
				labels = {
	  				App = "tester"
				}
			}
			
			spec {
				toleration {
					key      = "configured"
					operator = "Equal"
					value    = "true"
					effect   = "NoSchedule"
				}
				
				affinity {
					node_affinity {
						required_during_scheduling_ignored_during_execution {
							node_selector_term {
								match_expressions {
									key      = "cloud.google.com/gke-nodepool"
									operator = "In"
									values   = ["application-pool"]
								}
							}
						}
					}
				}

				automount_service_account_token = true
				service_account_name = "default"
				
				container {
					image = "gcr.io/mit-mic/syncer:v1"
					name  = "tester"

					image_pull_policy = "Always"
					
					resources {
						limits {
							cpu    = "0.5"
							memory = "512Mi"
						}
						
						requests {
			  				cpu    = "250m"
			  				memory = "50Mi"
						}
		  			}

					security_context {
						privileged = true
						allow_privilege_escalation = true
					}

					volume_mount {
						name = "host-root"
						mount_path = "/host"
					}
				}

				volume {
					name = "host-root"
					host_path {
						path = "/"
					}
				}

				host_network = true
				host_pid = true
			}
		}
	}
}
