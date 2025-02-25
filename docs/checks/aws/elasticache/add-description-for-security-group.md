---
title: add-description-for-security-group
---

### Explanation

Security groups and security group rules should include a description for auditing purposes.

Simplifies auditing, debugging, and managing security groups.

### Possible Impact
Descriptions provide context for the firewall rule reasons

### Suggested Resolution
Add descriptions for all security groups and rules


### Insecure Example

The following example will fail the aws-elasticache-add-description-for-security-group check.

```terraform

resource "aws_security_group" "bar" {
  name = "security-group"
}

resource "aws_elasticache_security_group" "bad_example" {
  name                 = "elasticache-security-group"
  security_group_names = [aws_security_group.bar.name]
  description = ""
}

```



### Secure Example

The following example will pass the aws-elasticache-add-description-for-security-group check.

```terraform

resource "aws_security_group" "bar" {
  name = "security-group"
}

resource "aws_elasticache_security_group" "good_example" {
  name                 = "elasticache-security-group"
  security_group_names = [aws_security_group.bar.name]
  description = "something"
}

```




### Related Links


- [https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/elasticache_security_group#description](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/elasticache_security_group#description){:target="_blank" rel="nofollow noreferrer noopener"}


