package azurerm

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/web"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmAppService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceCreateUpdate,
		Read:   resourceArmAppServiceRead,
		Update: resourceArmAppServiceCreateUpdate,
		Delete: resourceArmAppServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"kind": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"server_farm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"site_config": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"always_on": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"default_documents": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"number_of_workers": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"use_32_bit_worker": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"web_sockets_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"dotnet_version": {
							Type:     schema.TypeString,
							Required: true,
							// TODO: validation
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
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmAppServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sitesClient

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	kind := d.Get("kind").(string)
	location := d.Get("location").(string)
	serverFarmId := d.Get("server_farm_id").(string)
	siteConfig := expandAzureRmAppServiceSiteConfig(d)
	tags := d.Get("tags").(map[string]interface{})

	siteEnvelope := web.Site{
		Name:     &name,
		Kind:     &kind,
		Location: &location,
		SiteProperties: &web.SiteProperties{
			ServerFarmID: &serverFarmId,
			SiteConfig: &web.SiteConfig{
				SiteConfigProperties: siteConfig,
			},
		},
		Tags: expandTags(tags),
	}

	_, err := client.CreateOrUpdateSite(resourceGroup, name, siteEnvelope, "", "", "", "", make(chan struct{}))
	if err != nil {
		return err
	}

	read, err := client.GetSite(resourceGroup, name, "")
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read App Service %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServiceRead(d, meta)
}

func resourceArmAppServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sitesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["sites"]

	resp, err := client.GetSite(resourceGroup, name, "")
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on App Service Plan %s: %s", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))
	d.Set("resource_group_name", resp.ResourceGroup)
	d.Set("kind", resp.Kind)
	d.Set("server_farm_id", resp.ServerFarmID)

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmAppServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sitesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["sites"]

	resp, err := client.DeleteSite(resourceGroup, name, "", "", "", "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error issuing Azure ARM delete request of App Service Plan '%s': %s", name, err)
	}

	return nil
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
