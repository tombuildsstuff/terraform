package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMAppServiceSlot_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccAzureRMAppServiceSlot_basic, ri, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists("azurerm_app_service_slot.test"),
				),
			},
		},
	})
}

var testAccAzureRMAppServiceSlot_basic = `
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "West US"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "web"
  max_workers         = 1
  sku {
    capacity = "1"
    family   = "S"
    tier     = "Standard"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "Web"
  server_farm_id      = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on = true
  }

  tags {
    acceptance = "test"
  }
}

resource "azurerm_app_service_slot" "first" {
  name                = "acctestASS-%d"
  app_service_name    = "${azurerm_app_service.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "Web"

  site_config {
    always_on = true
  }

  tags {
    acceptance = "test"
  }
}
`

func testCheckAzureRMAppServiceSlotExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		appServiceName := rs.Primary.Attributes["app_service_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Service Slot: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).sitesClient

		resp, err := conn.GetSiteSlot(resourceGroup, appServiceName, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on sitesClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: App Service Slot %q (app service %q / resource group: %q) does not exist", name, appServiceName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMAppServiceSlotDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).sitesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_slot" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		appServiceName := rs.Primary.Attributes["app_service_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.GetSiteSlot(resourceGroup, appServiceName, name, "")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("App Service Slot still exists:\n%#v", resp)
		}
	}

	return nil
}
