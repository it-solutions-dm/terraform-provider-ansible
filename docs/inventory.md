---
# ansible_inventory Data Source

Generates an in-memory YAML inventory string used with Ansible.

## Example
```hcl
data "ansible_inventory" "example" {
  hosts = [{
    name = "db1"
    vars = {
      ansible_user = "ubuntu"
    }
  }]

  groups = {
    web = {
      hosts = ["web1", "web2"]
      vars = {
        ansible_user = "ubuntu"
      }
      children = ["app"]
    }

    app = {
      hosts = ["app1"]
    }
  }
}
```

## Attributes Reference
- `content` – Rendered YAML inventory string.

---
