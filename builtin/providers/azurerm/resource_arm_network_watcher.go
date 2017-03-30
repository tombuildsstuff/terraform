package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/networkwatcher"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmNetworkWatcher() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzureRMNetworkWatcherCreateUpdate,
		Read:   resourceAzureRMNetworkWatcherRead,
		Update: resourceAzureRMNetworkWatcherCreateUpdate,
		Delete: resourceAzureRMNetworkWatcherDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"tags": tagsSchema(),
		},
	}
}

func resourceAzureRMNetworkWatcherCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).networkWatcherClient

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})

	params := networkwatcher.NetworkWatcher{
		Name:     name,
		Location: azureRMNormalizeLocation(location),
		Tags:     expandTags(tags),
	}

	_, err := client.CreateOrUpdate(resourceGroup, name, params)
	if err != nil {
		return err
	}

	read, err := client.Get(resourceGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Network Watcher '%s' (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceAzureRMNetworkWatcherRead(d, meta)
}

func resourceAzureRMNetworkWatcherRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).networkWatcherClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["networkWatchers"]

	resp, err := client.Get(resourceGroup, name)
	if err != nil {
		return err
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", azureRMNormalizeLocation(resp.Location))
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceAzureRMNetworkWatcherDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).networkWatcherClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["networkWatchers"]

	_, err = client.Delete(resourceGroup, name, make(chan struct{}))
	if err != nil {
		return err
	}

	return nil
}
