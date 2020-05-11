package aws

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/medialive"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func resourceAwsMediaLiveInput() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsMediaLiveInputCreate,
		Read:   resourceAwsMediaLiveInputRead,
		Update: resourceAwsMediaLiveInputUpdate,
		Delete: resourceAwsMediaLiveInputDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"destinations": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"stream_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"input_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"role_arn": {
				Type:     schema.TypeString,
				Required: true,
			},

			"request_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"input_security_groups": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceAwsMediaLiveInputCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).medialiveconn

	input := &medialive.CreateInputInput{
		Type:      aws.String(d.Get("input_type").(string)),
		Name:      aws.String(d.Get("name").(string)),
		RequestId: aws.String(d.Get("request_id").(string)),
		RoleArn:   aws.String(d.Get("role_arn").(string)),
	}

	if v, ok := d.GetOk("destinations"); ok && len(v.([]interface{})) > 0 {
		input.Destinations = expandDestinations(
			v.([]interface{}),
		)
	}

	if raw, ok := d.GetOk("input_security_groups"); ok {
		list := raw.([]interface{})
		inputSecurityGroups := make([]*string, len(list))
		for i, groupId := range list {
			inputSecurityGroups[i] = aws.String(groupId.(string))
		}
	}

	if v := d.Get("tags").(map[string]interface{}); len(v) > 0 {
		input.Tags = keyvaluetags.New(v).IgnoreAws().MedialiveTags()
	}

	resp, err := conn.CreateInput(input)
	if err != nil {
		return fmt.Errorf("Error creating MediaLive Input: %s", err)
	}

	d.SetId(aws.StringValue(resp.Input.Id))

	return resourceAwsMediaLiveInputRead(d, meta)
}

func resourceAwsMediaLiveInputRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).medialiveconn

	input := &medialive.DescribeInputInput{
		InputId: aws.String(d.Id()),
	}

	resp, err := conn.DescribeInput(input)
	if err != nil {
		if isAWSErr(err, medialive.ErrCodeNotFoundException, "") {
			log.Printf("[WARN] MediaLive Input %s not found, error code (404)", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error describing MediaLive Input(%s): %s", d.Id(), err)
	}

	d.Set("arn", aws.StringValue(resp.Arn))

	if err := d.Set("tags", keyvaluetags.MedialiveKeyValueTags(resp.Tags).IgnoreAws().Map()); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	return nil
}

func resourceAwsMediaLiveInputUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).medialiveconn

	if d.HasChange("input_type") {
		input := &medialive.UpdateInputInput{
			Name: aws.String(d.Get("name").(string)),
		}

		_, err := conn.UpdateInput(input)
		if err != nil {
			if isAWSErr(err, medialive.ErrCodeNotFoundException, "") {
				log.Printf("[WARN] MediaLive Input %s not found, error code (404)", d.Id())
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error updating MediaLive Input(%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		o, n := d.GetChange("tags")
		if err := keyvaluetags.MedialiveUpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating tags: %s", err)
		}
	}

	//TODO : Check if we need to wait here too
	// if err := waitForMediaLiveInputOperation(conn, d.Id()); err != nil {
	// 	return fmt.Errorf("Error waiting for operational MediaLive Input (%s): %s", d.Id(), err)
	// }

	return resourceAwsMediaLiveInputRead(d, meta)
}

func resourceAwsMediaLiveInputDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).medialiveconn
	input := &medialive.DeleteInputInput{
		InputId: aws.String(d.Id()),
	}

	_, err := conn.DeleteInput(input)
	if err != nil {
		if isAWSErr(err, medialive.ErrCodeNotFoundException, "") {
			return nil
		}
		return fmt.Errorf("Error deleting MediaLive Input(%s): %s", d.Id(), err)
	}

	if err := waitForMediaLiveInputDeletion(conn, d.Id()); err != nil {
		return fmt.Errorf("Error waiting for deleting MediaLive Input(%s): %s", d.Id(), err)
	}

	return nil
}

func mediaLiveInputRefreshFunc(conn *medialive.MediaLive, inputId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input, err := conn.DescribeInput(&medialive.DescribeInputInput{
			InputId: aws.String(inputId),
		})

		if isAWSErr(err, medialive.ErrCodeNotFoundException, "") {
			return nil, medialive.InputStateDeleted, nil
		}

		if err != nil {
			return nil, "", fmt.Errorf("error reading MediaLive Input(%s): %s", inputId, err)
		}

		if input == nil {
			return nil, medialive.InputStateDeleted, nil
		}

		return input, aws.StringValue(input.State), nil
	}
}

func waitForMediaLiveInputDeletion(conn *medialive.MediaLive, inputId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			medialive.InputStateDetached,
			medialive.InputStateAttached,
			medialive.InputStateDeleting,
			medialive.InputStateDeleted,
		},
		Target:         []string{medialive.InputStateDeleted},
		Refresh:        mediaLiveInputRefreshFunc(conn, inputId),
		Timeout:        30 * time.Minute,
		NotFoundChecks: 1,
	}

	log.Printf("[DEBUG] Waiting for Media Live Input (%s) deletion", inputId)
	_, err := stateConf.WaitForState()

	if isAWSErr(err, medialive.ErrCodeNotFoundException, "") {
		return nil
	}

	return err
}

func expandDestinations(destinations []interface{}) []*medialive.InputDestinationRequest {
	var result []*medialive.InputDestinationRequest
	if len(destinations) == 0 {
		return nil
	}

	for _, destination := range destinations {
		r := destination.(map[string]interface{})

		result = append(result, &medialive.InputDestinationRequest{
			StreamName: aws.String(r["stream_name"].(string)),
		})
	}
	return result
}
