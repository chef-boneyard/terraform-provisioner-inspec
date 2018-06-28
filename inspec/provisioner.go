package inspec

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/rs/zerolog/log"
)

type ReporterConfig struct {
	Url         string `json:"url,omitempty"`
	Token       string `json:"token,omitempty"`
	NodeID      string `json:"node_uuid,omitempty"`
	NodeName    string `json:"node_name,omitempty"`
	Environment string `json:"environment,omitempty"`
	ReportUUID  string `json:"report_uuid,omitempty"`
	JobUUID     string `json:"job_uuid,omitempty"`
}

type TargetConfig struct {
	Backend  string `json:"backend,omitempty"`
	Hostname string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`

	User              string   `json:"user,omitempty"`
	Password          string   `json:"password,omitempty"`
	KeyFiles          []string `json:"key_files,omitempty"`
	SudoPassword      string   `json:"sudo_password,omitempty"`
	SudoOptions       string   `json:"sudo_options,omitempty"`
	AwsUser           string   `json:"aws_user,omitempty"`
	AwsPassword       string   `json:"aws_password,omitempty"`
	AzureClientID     string   `json:"azure_client_id,omitempty"`
	AzureClientSecret string   `json:"azure_client_secret,omitempty"`
	AzureTenantID     string   `json:"azure_tenant_id,omitempty"`

	LoginPath      string                    `json:"login_path,omitempty"`
	Sudo           bool                      `json:"sudo,omitempty"`
	Format         string                    `json:"format,omitempty"`
	Reporter       map[string]ReporterConfig `json:"reporter,omitempty"`
	Ssl            bool                      `json:"ssl,omitempty"`
	SslSelfSigned  bool                      `json:"self_signed,omitempty"`
	BackendCache   bool                      `json:"backend_cache,omitempty"`
	Region         string                    `json:"region,omitempty"`
	SubscriptionId string                    `json:"subscription_id,omitempty"`
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
						"backend": &schema.Schema{
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
				Type:     schema.TypeMap,
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

	profiles := getStringList(data.Get("profiles"))
	if len(profiles) == 0 {
		return errors.New("new profile defined")
	}

	o.Output(fmt.Sprintf("Run the following profiles %v", profiles))

	// read the target
	target := getStringMap(data.Get("target"))
	conf := parseTargetConfig(target)

	// read reporter config
	reporter := getStringMap(data.Get("reporter"))
	conf.Reporter = parseReporterConfig(reporter)

	switch conf.Backend {
	case "aws", "azure", "gcp":
		// create target url and use aws options
		return runRemote(ctx, s, data, o, profiles, conf)
	case "":
		// install inspec and run it from the instance
		return runLocal(ctx, s, data, o, profiles, conf)
	default:
		return errors.New(fmt.Sprintf("backend %s is not supported yet", conf.Backend))
	}

	return nil
}

func parseTargetConfig(target map[string]interface{}) *TargetConfig {
	conf := &TargetConfig{
		Backend:  getStringValue(target, "backend"),
		Hostname: getStringValue(target, "hostname"),

		Region:   getStringValue(target, "region"),
		User:     getStringValue(target, "user"),
		Password: getStringValue(target, "password"),

		AwsUser:     getStringValue(target, "aws_user"),
		AwsPassword: getStringValue(target, "aws_password"),

		AzureClientID:     getStringValue(target, "azure_client_id"),
		AzureClientSecret: getStringValue(target, "azure_client_secret"),
		AzureTenantID:     getStringValue(target, "azure_tenant_id"),
		SubscriptionId:    getStringValue(target, "subscription_id"),
	}
	return conf
}

func parseReporterConfig(data map[string]interface{}) map[string]ReporterConfig {
	rc := make(map[string]ReporterConfig)

	name := getStringValue(data, "name")
	if len(name) > 0 {
		rc[name] = ReporterConfig{}
	}
	return rc
}

func getStringValue(keymap map[string]interface{}, key string) string {
	v, ok := keymap[key]
	if ok {
		switch v := v.(type) {
		case string:
			return v
		}
	}

	return ""
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
