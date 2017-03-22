---
layout: "opsgenie"
page_title: "OpsGenie: opsgenie_contact"
sidebar_current: "docs-opsgenie-resource-contact"
description: |-
  Manages a Contact for a User within OpsGenie.
---

# opsgenie\_team

Manages a Team within OpsGenie.

## Example Usage

```
resource "opsgenie_user" "batman" {
  username  = "user@domain.com"
  full_name = "Batman"
  role      = "User"
}

resource "opsgenie_contact" "batphone" {
  username = "${opsgenie_user.batman.username}"
  method   = "voice"
  to       = "+447123456789"
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) The Username of the User to associate this Contact with.

* `method` - (Required) The Method to Contact this user by - accepted values are `email`, `sms` or `voice`. Changing this forces a new resource to be created.

* `to` - (Required) - Address of the contact. Should be a valid email address for email contacts and a valid phone number for sms and voice contacts.

* `enabled` - (Required) - Whether this means of Contact is enabled - Defaults to `true`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the OpsGenie Contact.

## Import

Contacts can be imported using the `id`, e.g.

```
$ terraform import opsgenie_contact.contact 812be1a1-32c8-4666-a7fb-03ecc385106c
```
