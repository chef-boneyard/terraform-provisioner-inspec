# Digital Ocean

The example launches an Ubuntu 16.04, and installs nginx and connects it to a load balancer.

## Prerequisites

In order to use this example, you need to create a [DigitalOcean API Token](https://cloud.digitalocean.com/account/api/tokens) and export it as an environment variable

```bash
export DIGITALOCEAN_TOKEN="Put Your Token Here" 
```

## Run this terraform setup

```bash
$ ssh-keygen -t rsa -b 4096 -C "digitalocean" -f ./id_rsa
$ terraform init
$ terraform plan
$ terraform apply
```

During apply phase, terraform executes inspec:

```
digitalocean_droplet.web (inspec): Profile: DevSec Linux Security Baseline (linux-baseline)
digitalocean_droplet.web (inspec): Version: 2.2.2
digitalocean_droplet.web (inspec): Target:  local://

digitalocean_droplet.web (inspec):   ✔  os-01: Trusted hosts login
digitalocean_droplet.web (inspec):      ✔  File /etc/hosts.equiv should not exist
digitalocean_droplet.web (inspec):   ✔  os-02: Check owner and permissions for /etc/shadow
digitalocean_droplet.web (inspec):      ✔  File /etc/shadow should exist
digitalocean_droplet.web (inspec):      ✔  File /etc/shadow should be file
digitalocean_droplet.web (inspec):      ✔  File /etc/shadow should be owned by "root"
digitalocean_droplet.web (inspec):      ✔  File /etc/shadow should not be executable
digitalocean_droplet.web (inspec):      ✔  File /etc/shadow should not be readable by other
digitalocean_droplet.web (inspec):      ✔  File /etc/shadow group should eq "shadow"
digitalocean_droplet.web (inspec):      ✔  File /etc/shadow should be writable by owner
digitalocean_droplet.web (inspec):      ✔  File /etc/shadow should be readable by owner
digitalocean_droplet.web (inspec):      ✔  File /etc/shadow should be readable by group
digitalocean_droplet.web (inspec):   ✔  os-03: Check owner and permissions for /etc/passwd
digitalocean_droplet.web (inspec):      ✔  File /etc/passwd should exist
```
