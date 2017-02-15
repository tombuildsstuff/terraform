package azurerm

import (
	"fmt"
	"log"

	"strings"

	"github.com/Azure/azure-sdk-for-go/arm/web"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmAppServicePlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServicePlanCreateUpdate,
		Read:   resourceArmAppServicePlanRead,
		Update: resourceArmAppServicePlanCreateUpdate,
		Delete: resourceArmAppServicePlanDelete,
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

			"kind": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"max_workers": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"sku": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"family": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAppServicePlanSkuFamily,
						},

						"tier": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAppServicePlanSkuTier,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmAppServicePlanCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serverFarmsClient

	log.Printf("[INFO] preparing arguments for AzureRM App Service Plan creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	kind := d.Get("kind").(string)

	skus := d.Get("sku").(*schema.Set).List()
	sku := skus[0].(map[string]interface{})
	capacity := int32(sku["capacity"].(int))
	family := sku["family"].(string)
	skuName := fmt.Sprintf("%s%d", family, capacity)
	tier := sku["tier"].(string)

	maxWorkers := int32(d.Get("max_workers").(int))
	tags := d.Get("tags").(map[string]interface{})

	serverFarmEnvelope := web.ServerFarmWithRichSku{
		Name:     &name,
		Location: &location,
		Kind:     &kind,
		Sku: &web.SkuDescription{
			Name:     &skuName,
			Size:     &skuName,
			Capacity: &capacity,
			Family:   &family,
			Tier:     &tier,
		},
		Tags: expandTags(tags),
		ServerFarmWithRichSkuProperties: &web.ServerFarmWithRichSkuProperties{
			Name: &name,
			MaximumNumberOfWorkers: &maxWorkers,
		},
	}

	allowPendingState := false
	_, err := client.CreateOrUpdateServerFarm(resourceGroup, name, serverFarmEnvelope, &allowPendingState, make(chan (struct{})))
	if err != nil {
		return err
	}

	read, err := client.GetServerFarm(resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read AzureRM App Service Plan %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServicePlanRead(d, meta)
}

func resourceArmAppServicePlanRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmAppServicePlanDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func validateAppServicePlanSkuFamily(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	tiers := map[string]bool{
		"b": true,
		"s": true,
		"p": true,
	}

	if !tiers[value] {
		errors = append(errors, fmt.Errorf("App Service Plan Family can only be B, S or P"))
	}
	return
}

func validateAppServicePlanSkuTier(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	tiers := map[string]bool{
		"basic":    true,
		"standard": true,
		"premium":  true,
	}

	if !tiers[value] {
		errors = append(errors, fmt.Errorf("App Service Plan Tier can only be Basic, Standard or Premium"))
	}
	return
}
