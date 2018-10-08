provider "digitalocean" {
  # You need to set this in your .bashrc
  # export DIGITALOCEAN_TOKEN="Your API TOKEN"
  #
}

# Create a new SSH key
resource "digitalocean_ssh_key" "key" {
  name       = "Terraform Example"
  public_key = "${file("./id_rsa.pub")}"
}

resource "digitalocean_droplet" "web" {
  ssh_keys           = ["${digitalocean_ssh_key.key.fingerprint}"]
  image              = "ubuntu-16-04-x64"
  region             = "${var.do_region}"
  size               = "s-1vcpu-1gb"
  private_networking = true
  backups            = true
  ipv6               = true
  name               = "nginx-web-ams3"

  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",
      "sudo apt-get update",
      "sudo apt-get -y install nginx",
    ]

    connection {
      type     = "ssh"
      private_key = "${file("./id_rsa")}"
      user     = "root"
      timeout  = "2m"
    }
  }

  # installs inspec and executes the profiles
  provisioner "inspec" {
    profiles = [
      "supermarket://dev-sec/linux-baseline",
      "supermarket://dev-sec/ssh-baseline",
    ]

    reporter {
      name = "cli"
    }

    connection {
      type     = "ssh"
      private_key = "${file("./id_rsa")}"
      user     = "root"
      timeout  = "2m"
    }

    on_failure = "continue"
  }
}

resource "digitalocean_floating_ip" "web" {
  droplet_id = "${digitalocean_droplet.web.id}"
  region     = "${digitalocean_droplet.web.region}"
}

resource "digitalocean_loadbalancer" "public" {
  name = "loadbalancer-1"
  region = "${var.do_region}"

  forwarding_rule {
    entry_port = 80
    entry_protocol = "http"

    target_port = 80
    target_protocol = "http"
  }

  healthcheck {
    port = 22
    protocol = "tcp"
  }

  droplet_ids = ["${digitalocean_droplet.web.id}"]
}