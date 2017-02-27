package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMPowerBIEmbeddedWorkspaceCollection_importBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "azurerm_powerbi_embedded_workspace_collection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPowerBIEmbeddedWorkspaceCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIEmbeddedWorkspaceCollection_basic(rInt),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMPowerBIEmbeddedWorkspaceCollection_importBasicWithTags(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "azurerm_powerbi_embedded_workspace_collection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPowerBIEmbeddedWorkspaceCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIEmbeddedWorkspaceCollection_basicWithTags(rInt),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
