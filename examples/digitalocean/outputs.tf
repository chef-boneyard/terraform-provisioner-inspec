output "Name" {
  value = "${digitalocean_droplet.web.name}"
}

output "Public ip" {
  value = "${digitalocean_droplet.web.ipv4_address}"
}

output "Load Balancer Public ip" {
  value = "${digitalocean_loadbalancer.public.ip}"
}


