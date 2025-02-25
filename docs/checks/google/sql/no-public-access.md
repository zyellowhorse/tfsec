---
title: no-public-access
---

### Explanation

Database instances should be configured so that they are not available over the public internet, but to internal compute resources which access them.

### Possible Impact
Public exposure of sensitive data

### Suggested Resolution
Remove public access from database instances


### Insecure Example

The following example will fail the google-sql-no-public-access check.

```terraform

resource "google_sql_database_instance" "postgres" {
	name             = "postgres-instance-a"
	database_version = "POSTGRES_11"
	
	settings {
		tier = "db-f1-micro"
	
		ip_configuration {
			ipv4_enabled = false
			authorized_networks {
				value           = "108.12.12.0/24"
				name            = "internal"
			}
	
			authorized_networks {
				value           = "0.0.0.0/0"
				name            = "internet"
			}
		}
	}
}
			
```



### Secure Example

The following example will pass the google-sql-no-public-access check.

```terraform

resource "google_sql_database_instance" "postgres" {
	name             = "postgres-instance-a"
	database_version = "POSTGRES_11"
	
	settings {
		tier = "db-f1-micro"
	
		ip_configuration {
			ipv4_enabled = false
			authorized_networks {
				value           = "108.12.12.0/24"
				name            = "internal"
			}
		}
	}
}
			
```




### Related Links


- [https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/sql_database_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/sql_database_instance){:target="_blank" rel="nofollow noreferrer noopener"}

- [https://www.cloudconformity.com/knowledge-base/gcp/CloudSQL/publicly-accessible-cloud-sql-instances.html](https://www.cloudconformity.com/knowledge-base/gcp/CloudSQL/publicly-accessible-cloud-sql-instances.html){:target="_blank" rel="nofollow noreferrer noopener"}


