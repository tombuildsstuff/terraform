package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMNetworkWatcher_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists("azurerm_network_watcher.test"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkWatcher_complete(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_complete(rInt),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists("azurerm_network_watcher.test"),
				),
			},
		},
	})
}

func testCheckAzureRMNetworkWatcherExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for network watcher: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).networkWatcherClient

		resp, err := conn.Get(resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on networkWatcherClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Network Watcher '%q' (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMNetworkWatcherDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).networkWatcherClient

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "azurerm_network_watcher" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Network Watcher still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMNetworkWatcher_basic(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "West US"
}

resource "azurerm_network_watcher" "test" {
  name                = "acctest-NWwatcher-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

`, rInt, rInt)
}

func testAccAzureRMNetworkWatcher_complete(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "West US"
}

resource "azurerm_network_watcher" "test" {
  name                = "acctest-NWwatcher-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tags {
    Environment = "Production"
  }
}

`, rInt, rInt)
}
