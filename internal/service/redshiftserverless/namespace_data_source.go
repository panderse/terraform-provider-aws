package redshiftserverless

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKDataSource("aws_redshiftserverless_namespace")
func DataSourceNamespace() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceNamespaceRead,

		Schema: map[string]*schema.Schema{
			"admin_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_iam_role_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iam_roles": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_exports": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceNamespaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RedshiftServerlessConn()

	namespaceName := d.Get("namespace_name").(string)

	resource, err := FindNamespaceByName(ctx, conn, namespaceName)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Redshift Serverless Namespace (%s): %s", namespaceName, err)
	}

	d.SetId(namespaceName)

	d.Set("admin_username", resource.AdminUsername)
	d.Set("arn", resource.NamespaceArn)
	d.Set("db_name", resource.DbName)
	d.Set("default_iam_role_arn", resource.DefaultIamRoleArn)
	d.Set("iam_roles", resource.IamRoles)
	d.Set("kms_key_id", resource.KmsKeyId)
	d.Set("log_exports", resource.LogExports)

	d.Set("namespace_id", resource.NamespaceId)

	return diags
}
