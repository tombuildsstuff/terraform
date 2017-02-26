package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/powerbiembedded"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmPowerBIWorkspaceCollection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPowerBIWorkspaceCollectionCreate,
		Read:   resourceArmPowerBIWorkspaceCollectionRead,
		Update: resourceArmPowerBIWorkspaceCollectionUpdate,
		Delete: resourceArmPowerBIWorkspaceCollectionDelete,
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

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"sku": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"tier": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: resourceAzureRMPowerBIWorkspaceCollectionSkuHash,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmPowerBIWorkspaceCollectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).workspaceCollectionsClient

	log.Printf("[INFO] preparing arguments for AzureRM PowerBI Embedded Workspace Collection Creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	skuName, skuTier := expandAzureRmPowerBIEmbeddedWorkspaceCollectionTier(d)
	tags := d.Get("tags").(map[string]interface{})

	properties := powerbiembedded.CreateWorkspaceCollectionRequest{
		Location: &location,
		Sku: &powerbiembedded.AzureSku{
			Name: &skuName,
			Tier: &skuTier,
		},
		Tags: expandTags(tags),
	}

	_, err := client.Create(resourceGroup, name, properties)
	if err != nil {
		return err
	}

	read, err := client.GetByName(resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read PowerBI Embedded Workspace Collection %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmPowerBIWorkspaceCollectionRead(d, meta)
}

func resourceArmPowerBIWorkspaceCollectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).workspaceCollectionsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["workspaceCollections"]

	skuName, skuTier := expandAzureRmPowerBIEmbeddedWorkspaceCollectionTier(d)
	tags := d.Get("tags").(map[string]interface{})

	updateProperties := powerbiembedded.UpdateWorkspaceCollectionRequest{
		Sku: &powerbiembedded.AzureSku{
			Name: &skuName,
			Tier: &skuTier,
		},
		Tags: expandTags(tags),
	}

	_, err = client.Update(resourceGroup, name, updateProperties)
	if err != nil {
		return err
	}

	return resourceArmPowerBIWorkspaceCollectionRead(d, meta)
}

func resourceArmPowerBIWorkspaceCollectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).workspaceCollectionsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["workspaceCollections"]

	resp, err := client.GetByName(resourceGroup, name)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM PowerBI Embedded Workspace Collection %s: %s", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", resp.Location)

	flattenAndSetSku(d, resp.Sku)
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmPowerBIWorkspaceCollectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).workspaceCollectionsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["workspaceCollections"]

	_, err = client.Delete(resourceGroup, name, make(chan struct{}))

	return err
}

func expandAzureRmPowerBIEmbeddedWorkspaceCollectionTier(d *schema.ResourceData) (string, string) {
	skus := d.Get("sku").(*schema.Set).List()
	sku := skus[0].(map[string]interface{})

	name := sku["name"].(string)
	tier := sku["tier"].(string)

	return name, tier
}

func flattenAndSetSku(d *schema.ResourceData, sku *powerbiembedded.AzureSku) {
	skuConfigs := &schema.Set{
		F: resourceAzureRMPowerBIWorkspaceCollectionSkuHash,
	}

	skuConfig := make(map[string]interface{}, 2)

	skuConfig["name"] = *sku.Name
	skuConfig["tier"] = *sku.Tier

	skuConfigs.Add(skuConfig)

	d.Set("sku", skuConfigs)
}

func resourceAzureRMPowerBIWorkspaceCollectionSkuHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	name := m["name"].(string)
	tier := m["tier"].(string)

	buf.WriteString(fmt.Sprintf("%s-%s", name, tier))
	return hashcode.String(buf.String())
}
