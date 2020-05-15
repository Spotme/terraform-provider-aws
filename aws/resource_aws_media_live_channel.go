package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/medialive"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func resourceAwsMediaLiveChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsMediaLiveChannelCreate,
		Read:   resourceAwsMediaLiveChannelRead,
		Update: resourceAwsMediaLiveChannelUpdate,
		Delete: resourceAwsMediaLiveChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// A standard channel has two encoding pipelines and a single pipeline channel only has one.
			"channel_class": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "STANDARD",
			},

			"destinations": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},

						"settings": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"password_param": {
										Type:     schema.TypeString,
										Required: true,
									},

									"stream_name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"url": {
										Type:     schema.TypeString,
										Required: true,
									},

									"user_name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"media_package_settings": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"multiplex_settings": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			//EgressEndpoints
			// TODO

			// Encoder Settings
			"encoder_settings": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audio_descriptions": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"audio_selector_name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"audio_type": {
										Type:     schema.TypeString,
										Required: true,
									},

									"audio_type_control": {
										Type:     schema.TypeString,
										Required: true,
									},

									"language_code": {
										Type:     schema.TypeString,
										Required: true,
									},

									"language_code_control": {
										Type:     schema.TypeString,
										Required: true,
									},

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"stream_name": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"codec_settings": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"aac_settings": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"input_type": {
																Type:     schema.TypeString,
																Required: true,
															},

															"bitrate": {
																Type:     schema.TypeString,
																Required: true,
															},

															"coding_mode": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"raw_format": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"spec": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"profile": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"rate_control_mode": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"sample_rate": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"input_attachments": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"input_attachment_name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"automatic_input_failover_settings": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"input_preference": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"secondary_input_id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"input_id": {
							Type:     schema.TypeString,
							Required: true,
						},

						"input_settings": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_end_behavior": {
										Type:     schema.TypeString,
										Required: true,
									},

									"input_filter": {
										Type:     schema.TypeString,
										Required: true,
									},

									"filter_strength": {
										Type:     schema.TypeString,
										Required: true,
									},

									"deblock_filter": {
										Type:     schema.TypeString,
										Required: true,
									},

									"denoise_filter": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"input_specification": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:     schema.TypeString,
							Required: true,
						},

						"maximum_bitrate": {
							Type:     schema.TypeString,
							Required: true,
						},

						"resolution": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			// The log level the user wants for their channel.
			"log_level": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			///PipelineDetails []*PipelineDetail `locationName:"pipelineDetails" type:"list"`
			// TODO

			"pipelines_running_count": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"reserved": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"role_arn": {
				Type:     schema.TypeString,
				Required: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceAwsMediaLiveChannelCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).medialiveconn

	input := &medialive.CreateChannelInput{
		ChannelClass: aws.String(d.Get("input_type").(string)),
		Name:         aws.String(d.Get("name").(string)),
		Reserved:     aws.String(d.Get("reserved").(string)),
		RoleArn:      aws.String(d.Get("role_arn").(string)),
	}

	if v := d.Get("tags").(map[string]interface{}); len(v) > 0 {
		input.Tags = keyvaluetags.New(v).IgnoreAws().MedialiveTags()
	}

	resp, err := conn.CreateChannel(input)
	if err != nil {
		return fmt.Errorf("Error creating MediaLive Channel: %s", err)
	}

	d.SetId(aws.StringValue(resp.Channel.Id))

	return resourceAwsMediaLiveChannelRead(d, meta)
}

func resourceAwsMediaLiveChannelRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).medialiveconn

	input := &medialive.DescribeChannelInput{
		ChannelId: aws.String(d.Id()),
	}

	resp, err := conn.DescribeChannel(input)
	if err != nil {
		if isAWSErr(err, medialive.ErrCodeNotFoundException, "") {
			log.Printf("[WARN] MediaLive Input %s not found, error code (404)", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error describing MediaLive Channel(%s): %s", d.Id(), err)
	}

	d.Set("arn", aws.StringValue(resp.Arn))
	d.Set("name", aws.StringValue(resp.Name))
	d.Set("role_arn", aws.StringValue(resp.RoleArn))

	if err := d.Set("tags", keyvaluetags.MedialiveKeyValueTags(resp.Tags).IgnoreAws().Map()); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	return nil
}

func resourceAwsMediaLiveChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAwsMediaLiveChannelRead(d, meta)
}

func resourceAwsMediaLiveChannelDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
