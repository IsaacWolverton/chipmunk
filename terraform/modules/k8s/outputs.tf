/**
 * The endpoint ip of the service that connects to the proxy from within
 * the chipmunk pod
 */
output "chipmunk-ip" {
  value = data.kubernetes_service.chipmunk-proxy.load_balancer_ingress.0.ip
}