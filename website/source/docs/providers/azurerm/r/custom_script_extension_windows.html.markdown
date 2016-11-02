---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_custom_script_extension_windows"
sidebar_current: "docs-azurerm-resource-custom-script-extension-windows"
description: |-
  Create a Custom Script Extension to run Scripts on Windows.
---

# azurerm\_custom\_script\_extension\_windows

Enables you to run Custom Scripts on your Windows boxes

## Example Usage

```
TODO: the rest ^

resource "azurerm_custom_script_extension_windows" "test" {
  command_to_execute = "powershell.exe foo.ps1"
  file_uris {
    "you/me/there": "foo.ps1"
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
