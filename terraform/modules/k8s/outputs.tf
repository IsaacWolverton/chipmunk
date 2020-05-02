output "chipmunk-ip" {
  value = data.kubernetes_service.chipmunk-proxy.load_balancer_ingress.0.ip
}