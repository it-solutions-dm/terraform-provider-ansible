package provider

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "gopkg.in/yaml.v2"
)

type inventoryDataSource struct{}

func NewInventoryDataSource() datasource.DataSource {
    return &inventoryDataSource{}
}

func (d *inventoryDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = "ansible_inventory"
}

func (d *inventoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "content": schema.StringAttribute{Computed: true},
            "hosts": schema.ListNestedAttribute{
                Optional: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "name": schema.StringAttribute{Required: true},
                        "vars": schema.MapAttribute{Optional: true, ElementType: types.StringType},
                    },
                },
            },
            "groups": schema.MapNestedAttribute{
                Optional: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "hosts": schema.ListAttribute{Optional: true, ElementType: types.StringType},
                        "vars": schema.MapAttribute{Optional: true, ElementType: types.StringType},
                        "children": schema.ListAttribute{Optional: true, ElementType: types.StringType},
                    },
                },
            },
        },
    }
}

type inventoryModel struct {
    Content types.String `tfsdk:"content"`
    Hosts   []struct {
        Name types.String            `tfsdk:"name"`
        Vars map[string]types.String `tfsdk:"vars"`
    } `tfsdk:"hosts"`
    Groups map[string]struct {
        Hosts    []types.String          `tfsdk:"hosts"`
        Vars     map[string]types.String `tfsdk:"vars"`
        Children []types.String          `tfsdk:"children"`
    } `tfsdk:"groups"`
}

func (d *inventoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data inventoryModel
    diags := req.Config.Get(ctx, &data)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    inventory := map[string]interface{}{}

    if len(data.Hosts) > 0 {
        hosts := map[string]interface{}{}
        for _, h := range data.Hosts {
            if len(h.Vars) > 0 {
                vars := map[string]string{}
                for k, v := range h.Vars {
                    vars[k] = v.ValueString()
                }
                hosts[h.Name.ValueString()] = vars
            } else {
                hosts[h.Name.ValueString()] = nil
            }
        }
        inventory["ungrouped"] = map[string]interface{}{
            "hosts": hosts,
        }
    }

    for groupName, group := range data.Groups {
        g := map[string]interface{}{}

        if len(group.Hosts) > 0 {
            h := map[string]interface{}{}
            for _, host := range group.Hosts {
                h[host.ValueString()] = nil
            }
            g["hosts"] = h
        }

        if len(group.Vars) > 0 {
            v := map[string]string{}
            for k, val := range group.Vars {
                v[k] = val.ValueString()
            }
            g["vars"] = v
        }

        if len(group.Children) > 0 {
            c := map[string]interface{}{}
            for _, child := range group.Children {
                c[child.ValueString()] = nil
            }
            g["children"] = c
        }

        inventory[groupName] = g
    }

    b, err := yaml.Marshal(inventory)
    if err != nil {
        resp.Diagnostics.AddError("YAML encoding error", err.Error())
        return
    }

    data.Content = types.StringValue(string(b))
    diags = resp.State.Set(ctx, &data)
    resp.Diagnostics.Append(diags...)
}
