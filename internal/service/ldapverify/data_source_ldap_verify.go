package ldapverify

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasLDAPVerifyRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bind_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"validations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"validation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasLDAPVerifyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	requestID := d.Get("request_id").(string)

	ldapResp, _, err := conn.LDAPConfigurations.GetStatus(ctx, projectID, requestID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifyRead, projectID, err))
	}

	if err := d.Set("hostname", ldapResp.Request.Hostname); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifySetting, "hostname", d.Id(), err))
	}
	if err := d.Set("port", ldapResp.Request.Port); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifySetting, "port", d.Id(), err))
	}
	if err := d.Set("bind_username", ldapResp.Request.BindUsername); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifySetting, "bind_username", d.Id(), err))
	}
	if err := d.Set("links", FlattenLinks(ldapResp.Links)); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifySetting, "links", d.Id(), err))
	}
	if err := d.Set("validations", flattenValidations(ldapResp.Validations)); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifySetting, "validations", d.Id(), err))
	}
	if err := d.Set("request_id", ldapResp.RequestID); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifySetting, "request_id", d.Id(), err))
	}
	if err := d.Set("status", ldapResp.Status); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifySetting, "status", d.Id(), err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"request_id": ldapResp.RequestID,
	}))

	return nil
}
