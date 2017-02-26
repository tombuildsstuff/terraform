package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMPowerBIWorkspaceCollection_importBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "azurerm_powerbi_workspace_collection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPowerBIWorkspaceCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIWorkspaceCollection_basic(rInt),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

func TestAccAzureRMPowerBIWorkspaceCollection_importBasicWithTags(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "azurerm_powerbi_workspace_collection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPowerBIWorkspaceCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIWorkspaceCollection_basicWithTags(rInt),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}
