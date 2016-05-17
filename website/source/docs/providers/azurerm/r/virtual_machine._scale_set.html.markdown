---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_scale_set"
sidebar_current: "docs-azurerm-resource-virtualmachine-scalesets"
description: |-
  Create a Virtual Machine Scale Set.
-------------------------------------

# azurerm\_virtual\_machine\_scale\_set

Create a virtual machine scale set

## Example Usage

```
resource "azurerm_resource_group" "test" {
    name = "acctestrg"
    location = "West US"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn"
    address_space = ["10.0.0.0/16"]
    location = "West US"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni"
    location = "West US"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_storage_account" "test" {
    name = "accsa"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "westus"
    account_type = "Standard_LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "mytestscaleset-1"
  location = "West US"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_A0"
    tier = "Standard"
    capacity = 2
  }

  virtual_machine_os_profile {
    computer_name_prefix = "testvm"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  virtual_machine_os_profile_linux_config {
    disable_password_authentication = true
    ssh_keys {
      path = "/home/myadmin/.ssh/authorized_keys"
      key_data = "${file("~/.ssh/demo_key.pub")}"
    }
  }

  virtual_machine_network_profile {
      name = "TestNetworkProfile"
      primary = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  virtual_machine_storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  virtual_machine_storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "14.04.2-LTS"
    version   = "latest"
  }
}
```

## Argument Reference

The following arguments are supported:


