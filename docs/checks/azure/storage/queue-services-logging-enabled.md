---
title: queue-services-logging-enabled
---

### Explanation


Storage Analytics logs detailed information about successful and failed requests to a storage service. 

This information can be used to monitor individual requests and to diagnose issues with a storage service. 

Requests are logged on a best-effort basis.


### Possible Impact
Logging provides valuable information about access and usage

### Suggested Resolution
Enable logging for Queue Services


### Insecure Example

The following example will fail the azure-storage-queue-services-logging-enabled check.

```terraform

resource "azurerm_storage_account" "bad_example" {
    name                     = "example"
    resource_group_name      = data.azurerm_resource_group.example.name
    location                 = data.azurerm_resource_group.example.location
    account_tier             = "Standard"
    account_replication_type = "GRS"
    queue_properties  {
  }
}

```



### Secure Example

The following example will pass the azure-storage-queue-services-logging-enabled check.

```terraform

resource "azurerm_storage_account" "good_example" {
    name                     = "example"
    resource_group_name      = data.azurerm_resource_group.example.name
    location                 = data.azurerm_resource_group.example.location
    account_tier             = "Standard"
    account_replication_type = "GRS"
    queue_properties  {
    logging {
        delete                = true
        read                  = true
        write                 = true
        version               = "1.0"
        retention_policy_days = 10
    }
  }
}

```




### Related Links


- [https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/storage_account#logging](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/storage_account#logging){:target="_blank" rel="nofollow noreferrer noopener"}

- [https://docs.microsoft.com/en-us/azure/storage/common/storage-analytics-logging?tabs=dotnet](https://docs.microsoft.com/en-us/azure/storage/common/storage-analytics-logging?tabs=dotnet){:target="_blank" rel="nofollow noreferrer noopener"}


