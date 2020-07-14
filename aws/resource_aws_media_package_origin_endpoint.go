package aws

import (
	"fmt"
	"log"

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

			"authorization": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cdn_identifier_secret": {
							Type:     schema.TypeString,
							Required: true,
						},

						"secrets_role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"channel_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"endpoint_id": {
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

						"playlist_type": {
							Type:     schema.TypeString,
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

									"max_video_bits_per_second": {
										Type:     schema.TypeInt,
										Optional: true,
									},

									"min_video_bits_per_second": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
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
		Id:                     aws.String(d.Get("endpoint_id").(string)),
		ChannelId:              aws.String(d.Get("channel_id").(string)),
		Description:            aws.String(d.Get("description").(string)),
		ManifestName:           aws.String(d.Get("manifest_name").(string)),
		Origination:            aws.String(d.Get("origination").(string)),
		StartoverWindowSeconds: aws.Int64(int64(d.Get("startover_window_seconds").(int))),
		TimeDelaySeconds:       aws.Int64(int64(d.Get("time_delay_seconds").(int))),
	}

	if v, ok := d.GetOk("authorization"); ok {
		input.Authorization = expandAuthorization(v.(*schema.Set))
	}

	if v, ok := d.GetOk("hls_package"); ok {
		input.HlsPackage = expandHlsPackage(v.(*schema.Set))
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
	conn := meta.(*AWSClient).mediapackageconn

	input := &mediapackage.DescribeOriginEndpointInput{
		Id: aws.String(d.Get("endpoint_id").(string)),
	}

	resp, err := conn.DescribeOriginEndpoint(input)
	if err != nil {
		if isAWSErr(err, mediapackage.ErrCodeNotFoundException, "") {
			log.Printf("[WARN] MediaPackage Origin Endpoint %s not found, error code (404)", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error describing MediaPackage Origin Endpoint(%s): %s", d.Id(), err)
	}

	d.Set("arn", aws.StringValue(resp.Arn))
	d.Set("description", aws.StringValue(resp.Description))
	d.Set("manifest_name", aws.StringValue(resp.ManifestName))
	d.Set("origination", aws.StringValue(resp.Origination))
	d.Set("url", aws.StringValue(resp.Url))

	if resp.Authorization != nil {
		d.Set("authorization", flattenAuthorization(resp.Authorization))
	}

	if err := d.Set("tags", keyvaluetags.MedialiveKeyValueTags(resp.Tags).IgnoreAws().Map()); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	return nil
}

func resourceAwsMediaPackageOriginEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).mediapackageconn

	input := &mediapackage.UpdateOriginEndpointInput{
		Id:                     aws.String(d.Get("endpoint_id").(string)),
		Description:            aws.String(d.Get("description").(string)),
		ManifestName:           aws.String(d.Get("manifest_name").(string)),
		Origination:            aws.String(d.Get("origination").(string)),
		StartoverWindowSeconds: aws.Int64(int64(d.Get("startover_window_seconds").(int))),
		TimeDelaySeconds:       aws.Int64(int64(d.Get("time_delay_seconds").(int))),
	}

	if v, ok := d.GetOk("hls_package"); ok {
		input.HlsPackage = expandHlsPackage(v.(*schema.Set))
	}

	if v, ok := d.GetOk("authorization"); ok {
		input.Authorization = expandAuthorization(v.(*schema.Set))
	}

	_, err := conn.UpdateOriginEndpoint(input)
	if err != nil {
		return fmt.Errorf("Error updating MediaPackage Origin Endpoint: %s", err)
	}
	return resourceAwsMediaPackageOriginEndpointRead(d, meta)
}

func resourceAwsMediaPackageOriginEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).mediapackageconn

	input := &mediapackage.DeleteOriginEndpointInput{
		Id: aws.String(d.Id()),
	}

	_, err := conn.DeleteOriginEndpoint(input)
	if err != nil {
		if isAWSErr(err, mediapackage.ErrCodeNotFoundException, "") {
			return nil
		}
		return fmt.Errorf("Error deleting MediaPackage Origin Endpoint(%s): %s", d.Id(), err)
	}

	return nil
}

func expandHlsPackage(s *schema.Set) *mediapackage.HlsPackage {
	if s.Len() > 0 {
		rawSettings := s.List()[0].(map[string]interface{})
		return &mediapackage.HlsPackage{
			SegmentDurationSeconds:         aws.Int64(int64(rawSettings["segment_duration_seconds"].(int))),
			PlaylistWindowSeconds:          aws.Int64(int64(rawSettings["playlist_window_seconds"].(int))),
			PlaylistType:                   aws.String(rawSettings["playlist_type"].(string)),
			AdMarkers:                      aws.String(rawSettings["ad_markers"].(string)),
			AdTriggers:                     expandStringList(rawSettings["ad_triggers"].([]interface{})),
			AdsOnDeliveryRestrictions:      aws.String(rawSettings["ads_on_delivery_restrictions"].(string)),
			ProgramDateTimeIntervalSeconds: aws.Int64(int64(rawSettings["program_date_time_interval_seconds"].(int))),
			IncludeIframeOnlyStream:        aws.Bool(rawSettings["include_iframe_only_stream"].(bool)),
			UseAudioRenditionGroup:         aws.Bool(rawSettings["use_audio_rendition_group"].(bool)),
			StreamSelection:                expandStreamSelection(rawSettings["stream_selection"].(*schema.Set)),
		}
	} else {
		log.Printf("[ERROR] MediaPackage OriginEndpoint: HlsPackage settings can not be found")
		return &mediapackage.HlsPackage{}
	}
}

func expandStreamSelection(s *schema.Set) *mediapackage.StreamSelection {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &mediapackage.StreamSelection{
			MaxVideoBitsPerSecond: aws.Int64(int64(settings["max_video_bits_per_second"].(int))),
			MinVideoBitsPerSecond: aws.Int64(int64(settings["min_video_bits_per_second"].(int))),
			StreamOrder:           aws.String(settings["stream_order"].(string)),
		}
	} else {
		log.Printf("[ERROR] MediaPackage OriginEndpoint: StreamSelection settings can not be found")
		return &mediapackage.StreamSelection{}
	}
}

func expandAuthorization(s *schema.Set) *mediapackage.Authorization {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &mediapackage.Authorization{
			CdnIdentifierSecret: aws.String(settings["cdn_identifier_secret"].(string)),
			SecretsRoleArn:      aws.String(settings["secrets_role_arn"].(string)),
		}
	} else {
		return nil
	}
}

func flattenAuthorization(auth *mediapackage.Authorization) map[string]interface{} {
	m := map[string]interface{}{
		"cdn_identifier_secret": aws.StringValue(auth.CdnIdentifierSecret),
		"secrets_role_arn":      aws.StringValue(auth.SecretsRoleArn),
	}
	return m
}
