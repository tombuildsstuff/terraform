package azurerm

import "github.com/hashicorp/terraform/helper/schema"

func locationSchema() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		StateFunc:        azureRMNormalizeLocation,
		DiffSuppressFunc: azureRMSuppressLocationDiff,
	}
}

func resourceGroupNameSchema() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: resourceAzurermResourceGroupNameDiffSuppress,
		ValidateFunc:     validateArmResourceGroupName,
	}
}

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeMap,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validateAzureRMTags,
	}
}
