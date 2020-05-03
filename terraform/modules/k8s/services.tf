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

resource "kubernetes_service" "test-proxy" {
    depends_on = [
        kubernetes_pod.chipmunk
    ]

    metadata {
        name = "test-proxy"
    }
    
    spec {
        selector = {
            App = "chipper"
        }
        
        port {
            port = var.application_port
            target_port = 8080
        }
        
        type = "LoadBalancer"
    }
}

data "kubernetes_service" "test-proxy" {
  depends_on = [
    kubernetes_service.test-proxy
  ]

  metadata {
    name = "test-proxy"
  }
}
