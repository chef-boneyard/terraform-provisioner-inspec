# InSpec Terraform Provisioner

The InSpec provisioner executes InSpec during the terraform apply run. It supports verifying:

* instances
* cloud platforms like azure, aws, digitalocean or gcp

## Installation

*One-Liner Install (Linux)*

```
mkdir -p ~/.terraform.d/plugins/
curl -L -s https://api.github.com/repos/inspec/terraform-provisioner-inspec/releases/latest \
  | grep --color=none browser_download_url \
  | grep --color=none Linux_x86_64 \
  | cut -d '"' -f 4 \
  | xargs curl -L | tar zxv -C ~/.terraform.d/plugins/
```

*One-Liner Install (Mac)*

```
mkdir -p ~/.terraform.d/plugins/
curl -L -s https://api.github.com/repos/inspec/terraform-provisioner-inspec/releases/latest \
  | grep --color=none browser_download_url \
  | grep --color=none Darwin_x86_64 \
  | cut -d '"' -f 4 \
  | xargs curl -L | tar zxv -C ~/.terraform.d/plugins/
```

If you encounter issues during installation, please also have a look at [Terraform Plugin Basics](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin)

*Linux*

```
mkdir -p ~/.terraform.d/plugins/
curl -L https://github.com/inspec/terraform-provisioner-inspec/releases/download/0.1.0/terraform-provisioner-inspec_0.1.0_Linux_x86_64.tar.gz -o terraform-provisioner-inspec.tar.gz
tar -xvzf terraform-provisioner-inspec.tar.gz -C ~/.terraform.d/plugins/
```

*Mac*

```
mkdir -p ~/.terraform.d/plugins/
curl -L https://github.com/inspec/terraform-provisioner-inspec/releases/download/0.1.0/terraform-provisioner-inspec_0.1.0_Darwin_x86_64.tar.gz -o terraform-provisioner-inspec.tar.gz
tar -xvzf terraform-provisioner-inspec.tar.gz -C ~/.terraform.d/plugins/
```

## Build the provisioner plugin

Clone repository to: `$GOPATH/src/github.com/inspec/terraform-provisioner-inspec`

```sh
$ mkdir -p $GOPATH/src/github.com/inspec; cd $GOPATH/src/github.com/inspec
$ git clone git@github.com:inspec/terraform-provisioner-inspec
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/inspec/terraform-provisioner-inspec
$ dep ensure
$ make build
```

## Targets

The provisionier can be uses with any instance. E.g for AWS the following runs InSpec and verifies the security with the [DevSec baselines](https://dev-sec.io/).

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

**Cloud Platform**

InSpec has a wide-support for cloud-platforms. This allows us to verify configuration like security groups. See InSpec [AWS](https://www.inspec.io/docs/reference/resources/#aws-resources), [Azure](https://www.inspec.io/docs/reference/resources/#azure-resources) and [GCP](https://www.inspec.io/docs/reference/resources/#gcp-resources) documentation

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

