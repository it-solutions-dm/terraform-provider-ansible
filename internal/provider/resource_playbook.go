package provider

import (
    "context"
    "os/exec"
    "strings"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "gopkg.in/yaml.v2"
    "os"
    "path/filepath"
)

type playbookResource struct{}

func NewPlaybookResource() resource.Resource {
    return &playbookResource{}
}

func (r *playbookResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = "ansible_playbook"
}

func (r *playbookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "path": schema.StringAttribute{Optional: true},
            "roles": schema.ListAttribute{Optional: true, ElementType: types.StringType},
            "extra_vars": schema.StringAttribute{Optional: true},
            "inventory_content": schema.StringAttribute{Optional: true},
            "become": schema.BoolAttribute{Optional: true},
            "check_mode": schema.BoolAttribute{Optional: true},
            "limit": schema.StringAttribute{Optional: true},
            "tags": schema.ListAttribute{Optional: true, ElementType: types.StringType},
        },
    }
}

type playbookModel struct {
    Path             types.String `tfsdk:"path"`
    Roles            types.List   `tfsdk:"roles"`
    ExtraVars        types.String `tfsdk:"extra_vars"`
    InventoryContent types.String `tfsdk:"inventory_content"`
    Become           types.Bool   `tfsdk:"become"`
    CheckMode        types.Bool   `tfsdk:"check_mode"`
    Limit            types.String `tfsdk:"limit"`
    Tags             types.List   `tfsdk:"tags"`
}

func (r *playbookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data playbookModel
    diags := req.Plan.Get(ctx, &data)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    args := []string{"-i", "-"}

    if !data.CheckMode.IsNull() && data.CheckMode.ValueBool() {
        args = append(args, "--check")
    }
    if !data.Limit.IsNull() && data.Limit.ValueString() != "" {
        args = append(args, "--limit", data.Limit.ValueString())
    }
    if !data.ExtraVars.IsNull() && data.ExtraVars.ValueString() != "" {
        args = append(args, "--extra-vars", data.ExtraVars.ValueString())
    }
    if !data.Tags.IsNull() && data.Tags.Elements() != nil {
        var tags []string
        diags := data.Tags.ElementsAs(ctx, &tags, false)
        resp.Diagnostics.Append(diags...)
        if resp.Diagnostics.HasError() {
            return
        }
        if len(tags) > 0 {
            args = append(args, "--tags", strings.Join(tags, ","))
        }
    }

    var playbookPath string
    if !data.Path.IsNull() && data.Path.ValueString() != "" {
        playbookPath = data.Path.ValueString()
    } else if !data.Roles.IsNull() {
        var roles []string
        diags := data.Roles.ElementsAs(ctx, &roles, false)
        resp.Diagnostics.Append(diags...)
        if resp.Diagnostics.HasError() {
            return
        }

        play := map[string]interface{}{
            "hosts":        "all",
            "gather_facts": false,
            "roles":        roles,
        }

        if !data.Become.IsNull() && data.Become.ValueBool() {
            play["become"] = true
        }

        playbook := []interface{}{play}
        b, err := yaml.Marshal(playbook)
        if err != nil {
            resp.Diagnostics.AddError("YAML error", err.Error())
            return
        }

        tmp := filepath.Join(os.TempDir(), "generated_playbook.yml")
        if err := os.WriteFile(tmp, b, 0644); err != nil {
            resp.Diagnostics.AddError("Temp file error", err.Error())
            return
        }
        playbookPath = tmp
    }

    args = append(args, playbookPath)
    cmd := exec.CommandContext(ctx, "ansible-playbook", args...)

    if !data.InventoryContent.IsNull() {
        cmd.Stdin = strings.NewReader(data.InventoryContent.ValueString())
    }

    output, err := cmd.CombinedOutput()
    if err != nil {
        resp.Diagnostics.AddError("Ansible error", string(output))
        return
    }
}

func (r *playbookResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {}
func (r *playbookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data playbookModel
    diags := req.Plan.Get(ctx, &data)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    r.Create(ctx, resource.CreateRequest{Plan: req.Plan}, (*resource.CreateResponse)(resp))
}
func (r *playbookResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {}