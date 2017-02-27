package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMPowerBIWorkspaceCollection_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPowerBIWorkspaceCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIWorkspaceCollection_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIWorkspaceCollectionExists("azurerm_powerbi_workspace_collection.test"),
				),
			},
		},
	})
}

func TestAccAzureRMPowerBIWorkspaceCollection_basicWithTags(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPowerBIWorkspaceCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIWorkspaceCollection_basicWithTags(rInt),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIWorkspaceCollectionExists("azurerm_powerbi_workspace_collection.test"),
				),
			},
		},
	})
}

func TestAccAzureRMPowerBIWorkspaceCollection_update(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPowerBIWorkspaceCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIWorkspaceCollection_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIWorkspaceCollectionExists("azurerm_powerbi_workspace_collection.test"),
				),
			},
			{
				Config: testAccAzureRMPowerBIWorkspaceCollection_basicWithTags(rInt),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIWorkspaceCollectionExists("azurerm_powerbi_workspace_collection.test"),
				),
			},
		},
	})
}

func testCheckAzureRMPowerBIWorkspaceCollectionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for workspace collection: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).workspaceCollectionsClient

		resp, err := conn.GetByName(resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: GetByName on workspaceCollectionsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Workspace Collections Client %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPowerBIWorkspaceCollectionDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).workspaceCollectionsClient

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "azurerm_powerbi_workspace_collection" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.GetByName(resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("PowerBI Workspace Collection still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMPowerBIWorkspaceCollection_basic(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_powerbi_workspace_collection" "test" {
    name                = "accTestPowerBIWorkspaceCollection-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    sku {
    	name = "S1"
    	tier = "Standard"
    }
}`, rInt, rInt)
}
func testAccAzureRMPowerBIWorkspaceCollection_basicWithTags(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_powerbi_workspace_collection" "test" {
    name                = "accTestPowerBIWorkspaceCollection-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    sku {
    	name = "S1"
    	tier = "Standard"
    }

    tags {
    	hello = "world"
    }
}`, rInt, rInt)
}
