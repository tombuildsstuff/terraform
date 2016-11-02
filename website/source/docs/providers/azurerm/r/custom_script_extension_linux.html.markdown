---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_custom_script_extension_linux"
sidebar_current: "docs-azurerm-resource-custom-script-extension-linux"
description: |-
  Create a Custom Script Extension to run Scripts on Linux.
---

# azurerm\_custom\_script\_extension\_linux

Enables you to run Custom Scripts on your Linux boxes

## Example Usage

```
TODO: the rest ^

resource "azurerm_custom_script_extension_linux" "test" {
  command_to_execute = "bash foo.sh"
  file_uris {
    "you/me/there": "foo.sh"
  }
}

```
## Argument Reference

The following arguments are supported:

* `command_to_execute` - (Required) The command that should be executed.

* `file_uris` - (Required) A list of File URI's to download - to be used in conjunction with the `command_to_execute`.

## Attributes Reference

The following attributes are exported:

* `id` - The Custom Script Extension ID.
