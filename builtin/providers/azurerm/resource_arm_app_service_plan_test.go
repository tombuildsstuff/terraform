package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMAppServicePlan_basic(t *testing.T) {

	ri := acctest.RandInt()
	config := fmt.Sprintf(testCheckAzureRMAppServicePlan_basic, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists("azurerm_app_service_plan.test"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_standard(t *testing.T) {

	ri := acctest.RandInt()
	config := fmt.Sprintf(testCheckAzureRMAppServicePlan_standard, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists("azurerm_app_service_plan.test"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_standardWithTags(t *testing.T) {

	ri := acctest.RandInt()
	config := fmt.Sprintf(testCheckAzureRMAppServicePlan_standardWithTags, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists("azurerm_app_service_plan.test"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_premium(t *testing.T) {

	ri := acctest.RandInt()
	config := fmt.Sprintf(testCheckAzureRMAppServicePlan_premium, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists("azurerm_app_service_plan.test"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_premiumWithTags(t *testing.T) {

	ri := acctest.RandInt()
	config := fmt.Sprintf(testCheckAzureRMAppServicePlan_premiumWithTags, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists("azurerm_app_service_plan.test"),
				),
			},
		},
	})
}

func testCheckAzureRMAppServicePlanExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		appServiceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for availability set: %s", appServiceName)
		}

		conn := testAccProvider.Meta().(*ArmClient).serverFarmsClient

		resp, err := conn.GetServerFarm(resourceGroup, appServiceName)
		if err != nil {
			return fmt.Errorf("Bad: Get on serverFarmsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: App Service Plan %q (resource group: %q) does not exist", appServiceName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMAppServicePlanDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).serverFarmsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_plan" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.GetServerFarm(resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("App Service Plan still exists:\n%#v", resp.ServerFarmWithRichSkuProperties)
		}
	}

	return nil
}

var testCheckAzureRMAppServicePlan_basic = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestasp-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "app"
  max_workers         = 1

  sku {
    capacity = 1
    family   = "B"
    tier     = "Basic"
  }
}
`

var testCheckAzureRMAppServicePlan_standard = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestasp-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "app"
  max_workers         = 1

  sku {
    capacity = 1
    family   = "S"
    tier     = "Standard"
  }
}
`

var testCheckAzureRMAppServicePlan_standardWithTags = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestasp-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "app"
  max_workers         = 1

  sku {
    capacity = 1
    family   = "S"
    tier     = "Standard"
  }

  tags {
  	Environment = "Production"
  }
}
`

var testCheckAzureRMAppServicePlan_premium = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestasp-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "app"
  max_workers         = 1

  sku {
    capacity = 1
    family   = "P"
    tier     = "Premium"
  }
}
`

var testCheckAzureRMAppServicePlan_premiumWithTags = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestasp-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "app"
  max_workers         = 1

  sku {
    capacity = 1
    family   = "P"
    tier     = "Premium"
  }

  tags {
  	Environment = "Production"
  }
}
`
