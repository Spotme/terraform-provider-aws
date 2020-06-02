package aws

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/medialive"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func resourceAwsMediaPackageOriginEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsMediaPackageOriginEndpointCreate,
		Read:   resourceAwsMediaPackageOriginEndpointRead,
		Update: resourceAwsMediaPackageOriginEndpointUpdate,
		Delete: resourceAwsMediaPackageOriginEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.String,
				Required: true,
			},

			"channel_id": {
				Type:     schema.String,
				Required: true,
			},

			"description": {
				Type:     schema.String,
				Optional: true,
			},

			"startover_window_seconds": {
				Type:     schema.Int,
				Optional: true,
			},

			"time_delay_seconds": {
				Type:     schema.Int,
				Optional: true,
			},

			"manifest_name": {
				Type:     schema.String,
				Optional: true,
			},

			"whitelist": {
				Type:     schema.List,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"hls_package": {
				Type:    schema.Set,
				Optinal: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_duration_seconds": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"playlist_window_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"play_list_type": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"ad_markers": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"ad_triggers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"ads_on_delivery_restrictions": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"program_date_time_interval_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"include_iframe_only_stream": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"use_audio_rendition_group": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"stream_selection": {
							Type:    schema.Set,
							Optinal: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"stream_order": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"origination": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceAwsMediaPackageOriginEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).mediapackageconn

	input := &mediapackage.CreateOriginEndpointInput{
		ChannelId:    aws.String(d.Get("channel_id").(string)),
		Id:           aws.String(d.Get("id").(string)),
		Description:  aws.String(d.Get("description").(string)),
		ManifestName: aws.String(d.Get("manifest_name").(string)),
	}

	if v := d.Get("tags").(map[string]interface{}); len(v) > 0 {
		input.Tags = keyvaluetags.New(v).IgnoreAws().MedialiveTags()
	}

	resp, err := conn.CreateOriginEndpoint(input)
	if err != nil {
		return fmt.Errorf("Error creating MediaPackage Origin Endpoint: %s", err)
	}

	d.SetId(aws.StringValue(resp.Channel.Id))

	return resourceAwsMediaPackageOriginEndpointRead(d, meta)
}

func resourceAwsMediaPackageOriginEndpointRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsMediaPackageOriginEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAwsMediaPackageOriginEndpointRead(d, meta)
}

func resourceAwsMediaPackageOriginEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
