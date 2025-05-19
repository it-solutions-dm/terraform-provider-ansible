package provider

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/provider/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource"
)

type ansibleProvider struct{}

func (p *ansibleProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "ansible"
}

func (p *ansibleProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{}
}

func (p *ansibleProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {}

func (p *ansibleProvider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        NewPlaybookResource,
    }
}

func (p *ansibleProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{
        NewInventoryDataSource,
    }
}

func New() provider.Provider {
    return &ansibleProvider{}
}
