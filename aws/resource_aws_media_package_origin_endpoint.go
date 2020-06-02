package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediapackage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"channel_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"startover_window_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"time_delay_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"manifest_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"whitelist": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"hls_package": {
				Type:     schema.TypeSet,
				Optional: true,
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
							Type:     schema.TypeSet,
							Optional: true,
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

	d.SetId(aws.StringValue(resp.Id))

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
