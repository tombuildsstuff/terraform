package azurerm

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/web"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmAppServiceSlot() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceSlotCreateUpdate,
		Read:   resourceArmAppServiceSlotRead,
		Update: resourceArmAppServiceSlotCreateUpdate,
		Delete: resourceArmAppServiceSlotDelete,
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

			"app_service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"site_config": appServiceSiteConfigSchema(),

			"tags": tagsSchema(),
		},
	}
}

func resourceArmAppServiceSlotCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sitesClient

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	kind := d.Get("kind").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	appServiceName := d.Get("app_service_name").(string)
	siteConfig := expandAzureRmAppServiceSiteConfig(d)
	tags := d.Get("tags").(map[string]interface{})

	siteEnvelope := web.Site{
		Name:     &name,
		Location: &location,
		Kind:     &kind,
		SiteProperties: &web.SiteProperties{
			SiteConfig: &web.SiteConfig{
				SiteConfigProperties: siteConfig,
			},
		},
		Tags: expandTags(tags),
	}

	_, err := client.CreateOrUpdateSiteSlot(resourceGroup, appServiceName, siteEnvelope, name, "", "", "", "", make(chan struct{}))
	if err != nil {
		return err
	}

	read, err := client.GetSiteSlot(resourceGroup, appServiceName, name, "")
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read App Service Slot %s (app service %s, resource group %s) ID", name, appServiceName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServiceSlotRead(d, meta)
}

func resourceArmAppServiceSlotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sitesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	appServiceName := id.Path["sites"]
	name := id.Path["slots"]

	resp, err := client.GetSiteSlot(resourceGroup, appServiceName, name, "")
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on App Service Slot %s: %s", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("kind", resp.Kind)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))
	d.Set("app_service_name", appServiceName)
	d.Set("resource_group_name", resp.ResourceGroup)

	flattenAndSetAppServiceSiteConfig(d, resp.SiteProperties)
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmAppServiceSlotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sitesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	appServiceName := id.Path["sites"]
	name := id.Path["slots"]

	resp, err := client.DeleteSiteSlot(resourceGroup, appServiceName, name, "", "", "", "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error issuing Azure ARM delete request of App Service Slot '%s': %s", name, err)
	}

	return nil
}
