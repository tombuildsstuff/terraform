---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_powerbi_workspace_collection"
sidebar_current: "docs-azurerm-resource-powerbi-workspace-collection"
description: |-
  Creates a PowerBI Workspace Collection
---

# azurerm\_powerbi\_workspace\_collection

Creates a PowerBI Workspace Collection

## Example Usage

```
resource "azurerm_resource_group" "test" {
    name     = "exampleResourceGroup"
    location = "West US"
}

resource "azurerm_powerbi_workspace_collection" "test" {
    name                = "example-workspace-collection"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    sku {
    	name = "S1"
    	tier = "Standard"
    }

    tags {
      environment = "Production"
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Workspace Collection. Changing this forces a
    new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the Workspace Collection.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) Defines which Azure Sku to use for this Workspace Collection. The `sku` block supports fields documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.


The `sku` block supports:

* `name` - (Required) The Sku Name to use for the Workspace Collection (e.g. S1).

* `tier` - (Required) The Sku Tier to use for the Workspace Collection (e.g. Standard).

## Attributes Reference

The following attributes are exported:

* `id` - The PowerBI Workspace Collection ID.

## Import

PowerBI Workspace Collections can be imported using the `resource id`, e.g.

```
terraform import azurerm_powerbi_workspace_collection.myworkspace 
/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.PowerBI/workspaceCollections/myworkspace
```
