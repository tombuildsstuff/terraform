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

			"site_config": appServiceSiteConfigSchema(),

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

	flattenAndSetAppServiceSiteConfig(d, resp.SiteProperties)
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
