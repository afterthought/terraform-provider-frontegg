package provider

import (
	"context"
	"fmt"

	"github.com/benesch/terraform-provider-frontegg/internal/restclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const fronteggWebhookPath = "https://portal.frontegg.com/frontegg/webhook"

type fronteggWebhook struct {
	ID          string   `json:"_id,omitempty"`
	DisplayName string   `json:"displayName,omitempty"`
	Description string   `json:"description,omitempty"`
	URL         string   `json:"url,omitempty"`
	Secret      string   `json:"secret,omitempty"`
	EventKeys   []string `json:"eventKeys,omitempty"`
	IsActive    bool     `json:"isActive"`
	Type        string   `json:"type,omitempty"`
	TenantID    string   `json:"tenantId,omitempty"`
	VendorID    string   `json:"vendorId,omitempty"`
	CreatedAt   string   `json:"createdAt,omitempty"`
}

func resourceFronteggWebhook() *schema.Resource {
	return &schema.Resource{
		Description: `Configures a Frontegg webhook.`,

		CreateContext: resourceFronteggWebhookCreate,
		ReadContext:   resourceFronteggWebhookRead,
		UpdateContext: resourceFronteggWebhookUpdate,
		DeleteContext: resourceFronteggWebhookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"enabled": {
				Description: "Whether the webhook is enabled.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"name": {
				Description: "A human-readable name for the webhook.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "A human-readable description of the webhook.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"url": {
				Description: "The URL to send events to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"secret": {
				Description: "A secret to include with the event.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"events": {
				Description: "The names of the events to subscribe to.",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"frontegg.user.authenticated",
						"frontegg.user.authenticatedWithSAML",
						"frontegg.user.authenticatedWithSSO",
						"frontegg.user.failedAuthentication",
						"frontegg.user.enrolledMFA",
						"frontegg.user.disabledMFA",
						"frontegg.user.created",
						"frontegg.user.signedUp",
						"frontegg.user.deleted",
						"frontegg.user.invitedToTenant",
						"frontegg.user.changedPassword",
						"frontegg.user.forgotPassword",
						"frontegg.user.removedFromTenant",
					}, false),
				},
				Required: true,
			},
			"type": {
				Description: "The type of the webhook.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"tenant_id": {
				Description: "The ID of the tenant that owns the webhook.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"vendor_id": {
				Description: "The ID of the vendor that owns the webhook.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "The timestamp at which the webhook was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceFronteggWebhookSerialize(d *schema.ResourceData) fronteggWebhook {
	return fronteggWebhook{
		IsActive:    d.Get("enabled").(bool),
		DisplayName: d.Get("name").(string),
		Description: d.Get("description").(string),
		URL:         d.Get("url").(string),
		Secret:      d.Get("secret").(string),
		EventKeys:   stringSetToList(d.Get("events").(*schema.Set)),
	}
}

func resourceFronteggWebhookDeserialize(d *schema.ResourceData, f fronteggWebhook) error {
	d.SetId(f.ID)
	if err := d.Set("enabled", f.IsActive); err != nil {
		return err
	}
	if err := d.Set("name", f.DisplayName); err != nil {
		return err
	}
	if err := d.Set("description", f.Description); err != nil {
		return err
	}
	if err := d.Set("url", f.URL); err != nil {
		return err
	}
	if err := d.Set("secret", f.Secret); err != nil {
		return err
	}
	if err := d.Set("events", f.EventKeys); err != nil {
		return err
	}
	if err := d.Set("events", f.EventKeys); err != nil {
		return err
	}
	if err := d.Set("type", f.Type); err != nil {
		return err
	}
	if err := d.Set("tenant_id", f.TenantID); err != nil {
		return err
	}
	if err := d.Set("vendor_id", f.VendorID); err != nil {
		return err
	}
	if err := d.Set("created_at", f.CreatedAt); err != nil {
		return err
	}
	return nil
}

func resourceFronteggWebhookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*restclient.Client)
	in := resourceFronteggWebhookSerialize(d)
	var out fronteggWebhook
	if err := client.Post(ctx, fronteggWebhookPath+"/custom", in, &out); err != nil {
		return diag.FromErr(err)
	}
	if err := resourceFronteggWebhookDeserialize(d, out); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceFronteggWebhookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*restclient.Client)
	var out fronteggWebhook
	if err := client.Patch(ctx, fmt.Sprintf("%s/%s", fronteggWebhookPath, d.Id()), nil, &out); err != nil {
		return diag.FromErr(err)
	}
	if err := resourceFronteggWebhookDeserialize(d, out); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func resourceFronteggWebhookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*restclient.Client)
	in := resourceFronteggWebhookSerialize(d)
	var out fronteggWebhook
	if err := client.Patch(ctx, fmt.Sprintf("%s/%s", fronteggWebhookPath, d.Id()), in, &out); err != nil {
		return diag.FromErr(err)
	}
	if err := resourceFronteggWebhookDeserialize(d, out); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceFronteggWebhookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*restclient.Client)
	if err := client.Delete(ctx, fmt.Sprintf("%s/%s", fronteggWebhookPath, d.Id()), nil); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
