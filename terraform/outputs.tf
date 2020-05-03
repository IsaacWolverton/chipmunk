/**
 * The endpoint ip of the application running within chipmunk
 */
output "chipmunk-ip" {
  value = module.k8s.chipmunk-ip
}
output "test-ip" {
  value = module.k8s.test-ip
}
