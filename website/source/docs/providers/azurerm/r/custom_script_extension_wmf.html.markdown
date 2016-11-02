---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_custom_script_extension_wmf"
sidebar_current: "docs-azurerm-resource-custom-script-extension-wmf"
description: |-
  Create a Custom Script Extension to install Windows Management Framework.
---

# azurerm\_custom\_script\_extension\_wmf

Enables you to install WMF on your Windows boxes.

## Example Usage

```
TODO: the rest ^

resource "azurerm_custom_script_extension_wmf" "test" {
  wmf_version = "5.0"
}

```
## Argument Reference

The following arguments are supported:

* `wmf_version` - (Required) The version of WMF to install (e.g. "5.0").

## Attributes Reference

The following attributes are exported:

* `id` - The Custom Script Extension ID.
