---
title: enable-volume-encryption
---

### Explanation

By enabling encryption on EBS volumes you protect the volume, the disk I/O and any derived snapshots from compromise if intercepted.

### Possible Impact
Unencrypted sensitive data is vulnerable to compromise.

### Suggested Resolution
Enable encryption of EBS volumes


### Insecure Example

The following example will fail the aws-ebs-enable-volume-encryption check.

```terraform

resource "aws_ebs_volume" "bad_example" {
  availability_zone = "us-west-2a"
  size              = 40

  tags = {
    Name = "HelloWorld"
  }
  encrypted = false
}

```



### Secure Example

The following example will pass the aws-ebs-enable-volume-encryption check.

```terraform

resource "aws_ebs_volume" "good_example" {
  availability_zone = "us-west-2a"
  size              = 40

  tags = {
    Name = "HelloWorld"
  }
  encrypted = true
}

```




### Related Links


- [https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume#encrypted](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume#encrypted){:target="_blank" rel="nofollow noreferrer noopener"}


