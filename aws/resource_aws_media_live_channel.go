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
		Update: nil,
		Delete: nil,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"input_attachments": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"input_attachment_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"input_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"log_level": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
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
