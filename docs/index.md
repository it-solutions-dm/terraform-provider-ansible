---
# Ansible Provider

The Ansible provider allows you to run Ansible playbooks or roles as Terraform resources. This is useful for integrating configuration management with your infrastructure provisioning.

## Example Usage

```hcl
resource "ansible_playbook" "example" {
  roles = ["myrole"]
  inventory_content = data.ansible_inventory.example.content
  become = true
  tags   = ["install"]
}
```

## Resources
- `ansible_playbook` — runs a playbook or role.

## Data Sources
- `ansible_inventory` — builds in-memory inventory YAML.

---
