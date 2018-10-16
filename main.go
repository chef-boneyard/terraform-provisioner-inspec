package main

import (
	"github.com/inspec/terraform-provisioner-inspec/inspec"
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
