package inspec

import (
	"testing"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/terraform"
)

func testConfig(t *testing.T, c map[string]interface{}) *terraform.ResourceConfig {
	r, err := config.NewRawConfig(c)
	if err != nil {
		t.Fatalf("config error: %s", err)
	}
	return terraform.NewResourceConfig(r)
}

func TestResourceProvisioner_Validate_good_config(t *testing.T) {
	c := testConfig(t, map[string]interface{}{
		"target": map[string]interface{}{
			"name":       "aws",
			"access_key": "blub",
			"secret_key": "no",
			"region":     "test task",
		},

		"reporter": []interface{}{
			map[string]interface{}{
				"name": "automate",
			},
			map[string]interface{}{
				"name": "json",
			},
		},
		"profiles": []string{"profile1", "profile2"},
	})

	warn, errs := Provisioner().Validate(c)
	if len(warn) > 0 {
		t.Fatalf("Warnings: %v", warn)
	}
	if len(errs) > 0 {
		t.Fatalf("Errors: %v", errs)
	}

}

func TestNoProfileSet(t *testing.T) {
	c := testConfig(t, map[string]interface{}{})

	warn, errs := Provisioner().Validate(c)
	if len(warn) > 0 {
		t.Fatalf("Warnings: %v", warn)
	}
	if len(errs) > 1 {
		t.Fatalf("Errors: %v", errs)
	}

}
