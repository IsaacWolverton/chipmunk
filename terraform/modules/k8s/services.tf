resource "kubernetes_service" "chipmunk-proxy" {
    depends_on = [
        kubernetes_pod.chipmunk
    ]

    metadata {
        name = "chipmunk-proxy"
    }
    
    spec {
        selector = {
            App = "chipper"
        }
        
        port {
            port = var.application_port
            target_port = 42069
        }
        
        type = "LoadBalancer"
    }
}

data "kubernetes_service" "chipmunk-proxy" {
  depends_on = [
    kubernetes_service.chipmunk-proxy
  ]

  metadata {
    name = "chipmunk-proxy"
  }
}