---
subcategory: "MediaPackage"
layout: "aws"
page_title: "AWS: aws_medialive_input_security_group"
description: |-
  Provides an AWS Elemental Input Security Group.
---

# Resource: aws_medialive_input_security_group

Provides an AWS Elemental MediaLive Input Security Group.

## Example Usage


## Argument Reference

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Id of the input security group
* `arn` - The Arn of the input security group

## Import

Media Live Input Security Group can be imported via the input security group id, e.g.

```
$ terraform import aws_medialive_input_security_group.test 1234567
```
