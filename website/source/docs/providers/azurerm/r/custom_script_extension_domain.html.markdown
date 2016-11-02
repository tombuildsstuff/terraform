---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_custom_script_extension_domain"
sidebar_current: "docs-azurerm-resource-custom-script-extension-domain"
description: |-
  Create a Custom Script Extension to join an Active Directory domain on Windows.
---

# azurerm\_custom\_script\_extension\_domain

Enables you to join an Active Directory domain on your Windows machines.

## Example Usage

```
TODO: the rest ^

resource "azurerm_custom_script_extension_domain" "test" {
  domain_name = "domain"
  organisational_path = ""
  user = "domain\\user"
  restart = true
  options = "3"
}

```
## Argument Reference

The following arguments are supported:

* `domain_name` - (Required) The command that should be executed.

* `organisational_path` - (Required) The command that should be executed.

* `user` - (Required) The user to join the domain as.

* `restart` - (Required) Should the machine be restarted after joining the domain?

* `options` - (Required) Any options required.

TODO: also mention about specifying the password in the ProtectedSettings section

## Attributes Reference

The following attributes are exported:

* `id` - The Custom Script Extension ID.
