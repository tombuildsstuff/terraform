---
layout: "aws"
page_title: "AWS: aws_elb_attachment"
sidebar_current: "docs-aws-resource-elb-attachment"
description: |-
  Provides an Elastic Load Balancer Attachment resource.
---

# aws\_elb\_attachment

Provides an Elastic Load Balancer Attachment resource.

~> **NOTE on ELB Instances and ELB Attachments:** Terraform currently provides
both a standalone ELB Attachment resource (describing the instances attached to
an ELB), and an [Elastic Load Balancer resource](elb.html) with
`instances` defined in-line. At this time you cannot use an ELB with in-line
instaces in conjunction with an ELB Attachment resource. Doing so will cause a
conflict of rule settings and will overwrite rules.
## Example Usage

```
# Create a new load balancer attachment
resource "aws_elb_attachment" "baz" {
  elb      = "${aws_elb.bar.id}"
  instances = ["${aws_instance.foo.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `elb` - (Required) The name of the ELB.
* `instances` - (Required) A list of instance ids to place in the ELB pool.
