package inspec

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/rs/zerolog/log"
)

type ReporterConfig struct {
}

type InSpecConfig struct {
	Reporter map[string]ReporterConfig `json:"reporter,omitempty"`
	Sudo     bool                      `json:"sudo,omitempty"`
}

func Provisioner() terraform.ResourceProvisioner {
	return &schema.Provisioner{
		Schema: map[string]*schema.Schema{

			"profiles": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},

			"target": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"sudo": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"secret_key": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"access_key": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"region": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"reporter": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// name of the reporter
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						// for automate reporter
						"url": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						// for automate reporter
						"token": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
		ApplyFunc:    applyFn,
		ValidateFunc: validateFn,
	}
}

func validateFn(c *terraform.ResourceConfig) (ws []string, es []error) {
	fmt.Println("verify inspec")
	return nil, nil
}

func applyFn(ctx context.Context) error {
	log.Info().Msg("apply inspec")
	s := ctx.Value(schema.ProvRawStateKey).(*terraform.InstanceState)
	data := ctx.Value(schema.ProvConfigDataKey).(*schema.ResourceData)
	o := ctx.Value(schema.ProvOutputKey).(terraform.UIOutput)

	cliProfiles := ""
	tprofiles := getStringList(data.Get("profiles"))
	if len(tprofiles) > 0 {
		cliProfiles = strings.Join(tprofiles, " ")
	} else {
		return errors.New("new profile defined")
	}

	o.Output(fmt.Sprintf("Run the following profiles %v", cliProfiles))

	// read the target
	target := getStringMap(data.Get("target"))
	switch target["name"] {
	case "aws":
		// create target url and use aws options
		return runRemote(ctx, s, data, o, cliProfiles)
	default:
		// run on the node itself
		return runLocal(ctx, s, data, o, cliProfiles)
	}

	return nil
}

func getStringList(v interface{}) []string {
	var result []string
	switch v := v.(type) {
	case nil:
		return result
	case []interface{}:
		for _, vv := range v {
			if vv, ok := vv.(string); ok {
				result = append(result, vv)
			}
		}
		return result
	default:
		panic(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func getStringMap(v interface{}) map[string]interface{} {
	switch v := v.(type) {
	case nil:
		return make(map[string]interface{})
	case map[string]interface{}:
		return v
	default:
		panic(fmt.Sprintf("Unsupported type: %T", v))
	}
}
