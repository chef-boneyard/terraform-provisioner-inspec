package main

import (
	"github.com/chris-rock/terraform-provisioner-inspec/inspec"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: func() terraform.ResourceProvisioner {
			return inspec.Provisioner()
		},
	})
}
