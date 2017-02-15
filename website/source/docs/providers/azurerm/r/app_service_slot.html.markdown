---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_slot"
sidebar_current: "docs-azurerm-resource-app-service_slot"
description: |-
  Creates a Slot within an App Service.
---

# azurerm\_app\_service\_slot

Creates a Slot within an App Service.

## Example Usage

```
TODO
```

## Argument Reference

The following arguments are supported:

TODO

## Attributes Reference

The following attributes are exported:

* `id` - The virtual App Service Slot.


## Import

App Service Slot's can be imported using the `resource id`, e.g.

```
terraform import azurerm_app_service_slot.slot1
/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1/slots/slot1
```
