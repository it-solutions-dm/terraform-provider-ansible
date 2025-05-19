---
# ansible_playbook Resource

Runs an Ansible playbook or defined roles.

## Example
```hcl
resource "ansible_playbook" "example" {
  roles = ["common"]
  become = true
  inventory_content = data.ansible_inventory.example.content
  extra_vars = jsonencode({ env = "prod" })
  tags = ["setup"]
}
```

## Argument Reference
- `path` – (Optional) Path to a playbook file.
- `roles` – (Optional) List of roles to run (generates playbook).
- `become` – (Optional) Enables privilege escalation.
- `extra_vars` – (Optional) JSON string with extra vars.
- `inventory_content` – (Optional) In-memory inventory YAML string.
- `tags` – (Optional) Tags to limit execution.
- `check_mode` – (Optional) If true, runs with `--check`.
- `limit` – (Optional) Limits to specific hosts or groups.

---
