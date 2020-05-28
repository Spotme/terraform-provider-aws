package aws

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/medialive"
	"github.com/google/uuid"
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateArn,
			},

			"input_security_groups": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"input_class": {
				Type:     schema.TypeInt,
				Computed: true,
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
		RequestId: aws.String(uuid.Must(uuid.NewRandom()).String()),
	}

	if v, ok := d.GetOk("destinations"); ok && len(v.([]interface{})) > 0 {
		input.Destinations = expandInputDestinations(
			v.([]interface{}),
		)
	}

	if raw, ok := d.GetOk("input_security_groups"); ok {
		input.InputSecurityGroups = convertInputSecurityGroups(raw)
	}

	if v := d.Get("tags").(map[string]interface{}); len(v) > 0 {
		input.Tags = keyvaluetags.New(v).IgnoreAws().MedialiveTags()
	}

	resp, err := conn.CreateInput(input)
	if err != nil {
		return fmt.Errorf("Error creating MediaLive Input: %s", err)
	}

	d.SetId(aws.StringValue(resp.Input.Id))

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"CREATING"},
		Target:  []string{"DETACHED", "ATTACHED"},
		Refresh: func() (interface{}, string, error) {
			input := &medialive.DescribeInputInput{
				InputId: aws.String(d.Id()),
			}
			resp, err := conn.DescribeInput(input)
			if err != nil {
				return 0, "", err
			}
			return resp, aws.StringValue(resp.State), nil
		},
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 5,
	}
	_, err = createStateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("Error waiting MediaLive Input (%s) to be created: %s", d.Id(), err)
	}

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

	if err := d.Set("destinations", flattenInputDestinations(resp.Destinations)); err != nil {
		return fmt.Errorf("error setting destinations: %s", err)
	}

	d.Set("arn", aws.StringValue(resp.Arn))
	d.Set("type", aws.StringValue(resp.Type))
	d.Set("name", aws.StringValue(resp.Name))
	d.Set("input_class", aws.StringValue(resp.InputClass))

	if err := d.Set("tags", keyvaluetags.MedialiveKeyValueTags(resp.Tags).IgnoreAws().Map()); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	return nil
}

func resourceAwsMediaLiveInputUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).medialiveconn

	input := &medialive.UpdateInputInput{
		InputId: aws.String(d.Id()),
	}

	if d.HasChange("name") {
		input.Name = aws.String(d.Get("name").(string))
	}

	if d.HasChange("stream_name") {
		if v, ok := d.GetOk("destinations"); ok && len(v.([]interface{})) > 0 {
			input.Destinations = expandInputDestinations(
				v.([]interface{}),
			)
		}
	}

	if d.HasChange("input_security_groups") {
		if raw, ok := d.GetOk("input_security_groups"); ok {
			input.InputSecurityGroups = convertInputSecurityGroups(raw)
		}
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

	if d.HasChange("tags") {
		o, n := d.GetChange("tags")
		if err := keyvaluetags.MedialiveUpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating tags: %s", err)
		}
	}

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

func flattenInputDestinations(inputDestinations []*medialive.InputDestination) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(inputDestinations))
	for _, destination := range inputDestinations {
		r := map[string]interface{}{
			"url":         aws.StringValue(destination.Url),
			"port":        aws.StringValue(destination.Ip),
			"ip":          aws.StringValue(destination.Port),
			"stream_name": obtainStreamName(destination.Url),
		}
		result = append(result, r)
	}
	return result
}

func expandInputDestinations(destinations []interface{}) []*medialive.InputDestinationRequest {
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

func convertInputSecurityGroups(raw interface{}) []*string {
	list := raw.([]interface{})
	inputSecurityGroups := make([]*string, len(list))
	for i, groupId := range list {
		inputSecurityGroups[i] = aws.String(groupId.(string))
	}
	return inputSecurityGroups
}

func obtainStreamName(streamUrl *string) string {
	resp, err := url.Parse(*streamUrl)
	if err != nil {
		log.Printf("[WARN] Was not able to obtain StreamName for MediaLive Input (%s)", err)
		return ""
	}
	return resp.Path[1:len(resp.Path)]
}
