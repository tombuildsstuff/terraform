package azurerm

import (
	"github.com/Azure/azure-sdk-for-go/arm/web"
	"github.com/hashicorp/terraform/helper/schema"
)

func appServiceSiteConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"always_on": {
					Type:     schema.TypeBool,
					Optional: true,
				},

				"default_documents": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Set:      schema.HashString,
				},

				"number_of_workers": {
					Type:     schema.TypeInt,
					Optional: true,
				},

				"use_32_bit_worker": {
					Type:     schema.TypeBool,
					Optional: true,
				},

				"web_sockets_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
				},

				"dotnet_version": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"java_version": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"node_version": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"python_version": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"php_version": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
		// TODO: enable me
		// Set: resourceAzureRMAppServiceSiteConfigurationHash,
	}
}

func expandAzureRmAppServiceSiteConfig(d *schema.ResourceData) *web.SiteConfigProperties {

	configs := d.Get("site_config").(*schema.Set).List()
	config := configs[0].(map[string]interface{})

	// required
	alwaysOn := config["always_on"].(bool)
	dotnetVersion := config["dotnet_version"].(string)
	numberOfWorkers := int32(config["number_of_workers"].(int))
	use32BitWorkerProcess := config["use_32_bit_worker"].(bool)
	webSocketsEnabled := config["web_sockets_enabled"].(bool)

	defaultDocumentStrings := config["default_documents"].(*schema.Set).List()
	defaultDocuments := make([]string, len(defaultDocumentStrings))
	for i, v := range defaultDocumentStrings {
		defaultDocuments[i] = v.(string)
	}

	// optional
	javaVersion := config["java_version"].(string)
	nodeVersion := config["node_version"].(string)
	pythonVersion := config["python_version"].(string)
	phpVersion := config["php_version"].(string)

	return &web.SiteConfigProperties{
		NumberOfWorkers:       &numberOfWorkers,
		DefaultDocuments:      &defaultDocuments,
		NetFrameworkVersion:   &dotnetVersion,
		AlwaysOn:              &alwaysOn,
		NodeVersion:           &nodeVersion,
		WebSocketsEnabled:     &webSocketsEnabled,
		Use32BitWorkerProcess: &use32BitWorkerProcess,
		PythonVersion:         &pythonVersion,
		PhpVersion:            &phpVersion,
		JavaVersion:           &javaVersion,
	}
}

func flattenAndSetAppServiceSiteConfig(d *schema.ResourceData, properties *web.SiteProperties) {

}