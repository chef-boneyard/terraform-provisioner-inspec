# InSpec Terraform Provisioner

The InSpec provisioner executes InSpec during the terraform apply run. It supports verifying:

* instances
* aws, azure, gcp cloud services

## Installation

https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin


## Build the provisioner plugin

Clone repository to: `$GOPATH/src/github.com/chris-rock/terraform-provisioner-inspec`

```sh
$ mkdir -p $GOPATH/src/github.com/chris-rock; cd $GOPATH/src/github.com/chris-rock
$ git clone git@github.com:chris-rock/terraform-provisioner-inspec
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/chris-rock/terraform-provisioner-inspec
$ dep ensure
$ make build
```


## Targets

**Instances**

```
resource "aws_instance" "web" {
  connection {
    user = "ubuntu"
  }

  instance_type = "t2.micro"
  ami = "${lookup(var.aws_amis, var.aws_region)}"
  key_name = "chartmann"
  vpc_security_group_ids = ["${aws_security_group.default.id}"]
  subnet_id = "${aws_subnet.default.id}"

  # installs inspec and executes the profiles
  provisioner "inspec" {
    profiles = [
      "supermarket://dev-sec/linux-baseline",
      "supermarket://dev-sec/ssh-baseline",
    ]

    # allow pass if compliance errors happen
    on_failure = "continue"
  }
}
```

**AWS**

```
resource "null_resource" "inspec_aws" {
  // runs inspec profile against aws services
  provisioner "inspec" {
    profiles = [
      "https://github.com/chris-rock/aws-baseline",
    ]

    target {
      backend      = "aws"
      access_key = "${var.aws_access_key}"
      secret_key = "${var.aws_secret_key}"
      region     = "us-east-1"
    }

    reporter {
      name = "json"
    }

    on_failure = "continue"
  }
}

```

