package main

import (
    "context"
    "github.com/it-solutions-dm/terraform-provider-ansible/internal/provider"
    "github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
    providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
        Address: "registry.terraform.io/it-solutions-dm/ansible",
    })
}