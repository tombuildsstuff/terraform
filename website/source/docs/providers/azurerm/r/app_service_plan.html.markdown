---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_plan"
sidebar_current: "docs-azurerm-resource-app-service_plan"
description: |-
  Create an App Service Plan (formerly Server Farm).
---

# azurerm\_app\_service\_plan

Create an App Service Plan (formerly Server Farm).

## Example Usage

```
TODO
```

## Argument Reference

The following arguments are supported:

TODO

## Attributes Reference

The following attributes are exported:

* `id` - The virtual App Service Plan.


## Import

App Service Plan's can be imported using the `resource id`, e.g.

```
terraform import azurerm_app_service_plan.plan1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/serverfarms/plan1
```
