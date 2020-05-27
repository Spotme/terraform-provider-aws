
// MediaLive Channel structure helpers.
//
// These functions assist in pulling in data from Terraform resource
// configuration for the aws_media_live_channel resource, as there are
// several sub-fields that require their own data type, and do not necessarily
// 1-1 translate to resource configuration.

package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/medialive"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func expandInputAttachments(inputAttachments []interface{}) []*medialive.InputAttachment {
	var result []*medialive.InputAttachment
	if len(inputAttachments) == 0 {
		return nil
	}

	for _, inputAtt := range inputAttachments {
		r := inputAtt.(map[string]interface{})

		result = append(result, &medialive.InputAttachment{
			InputAttachmentName: aws.String(r["input_attachment_name"].(string)),
			InputId:             aws.String(r["input_id"].(string)),
			InputSettings:       expandInputAttachmentSettings(r["input_settings"].(*schema.Set)),
		})
	}
	return result
}

func expandInputSpecification(s *schema.Set) *medialive.InputSpecification {
	if s.Len() > 0 {
		rawInputSpecification := s.List()[0].(map[string]interface{})
		return &medialive.InputSpecification{
			Codec:          aws.String(rawInputSpecification["codec"].(string)),
			MaximumBitrate: aws.String(rawInputSpecification["maximum_bitrate"].(string)),
			Resolution:     aws.String(rawInputSpecification["resolution"].(string)),
		}
	} else {
		log.Printf("[WARN] MediaLive Channel: Input Specification can not be found")
		return &medialive.InputSpecification{}
	}
}

func expandDestinations(destinations []interface{}) []*medialive.OutputDestination {
	var result []*medialive.OutputDestination
	if len(destinations) == 0 {
		log.Printf("[WARN] MediaLive Channel: At least one output destination required")
		return nil
	}

	for _, destination := range destinations {
		r := destination.(map[string]interface{})

		result = append(result, &medialive.OutputDestination{
			Id:       aws.String(r["id"].(string)),
			Settings: expandOutputDestinationSettings(r["settings"].([]interface{})),
		})
	}
	return result
}

func expandOutputDestinationSettings(destinationSettings []interface{}) []*medialive.OutputDestinationSettings {
	var result []*medialive.OutputDestinationSettings
	if len(destinationSettings) == 0 {
		log.Printf("[ERROR] MediaLive Channel: One destination setting is required for each redundant encoder")
		return nil
	}

	for _, settings := range destinationSettings {
		r := settings.(map[string]interface{})

		result = append(result, &medialive.OutputDestinationSettings{
			PasswordParam: aws.String(r["password_param"].(string)),
			StreamName:    aws.String(r["stream_name"].(string)),
			Url:           aws.String(r["url"].(string)),
			Username:      aws.String(r["username"].(string)),
		})
	}
	return result
}

func expandEncoderSettings(s *schema.Set) *medialive.EncoderSettings {
	if s.Len() > 0 {
		rawEncoderSettings := s.List()[0].(map[string]interface{})
		return &medialive.EncoderSettings{
			AudioDescriptions: expandAudioDescriptions(rawEncoderSettings["audio_descriptions"].([]interface{})),
			OutputGroups:      expandOutputGroups(rawEncoderSettings["output_groups"].([]interface{})),
			TimecodeConfig:    expandTimecodeConfigs(rawEncoderSettings["timecode_config"].(*schema.Set)),
			VideoDescriptions: expandVideoDescriptions(rawEncoderSettings["video_descriptions"].([]interface{})),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: Encoder settings required")
		return &medialive.EncoderSettings{}
	}
}

func expandAudioDescriptions(audioDescriptions []interface{}) []*medialive.AudioDescription {
	var result []*medialive.AudioDescription
	if len(audioDescriptions) == 0 {
		log.Printf("[ERROR] MediaLive Channel: At least one audio description is required for each encoder")
		return nil
	}

	for _, descs := range audioDescriptions {
		r := descs.(map[string]interface{})

		result = append(result, &medialive.AudioDescription{
			AudioSelectorName: aws.String(r["audio_selector_name"].(string)),
			Name:              aws.String(r["name"].(string)),
			StreamName:        aws.String(r["stream_name"].(string)),
			CodecSettings:     expandAudioCodecSettings(r["codec_settings"].(*schema.Set)),
		})
	}
	return result
}

func expandAudioCodecSettings(s *schema.Set) *medialive.AudioCodecSettings {
	if s.Len() > 0 {
		rawCodecSettings := s.List()[0].(map[string]interface{})
		return &medialive.AudioCodecSettings{
			AacSettings: expandAacCodecSettings(rawCodecSettings["aac_settings"].(*schema.Set)),
		}
	} else {
		log.Printf("[WARN] MediaLive Channel: Input Specification can not be found")
		return &medialive.AudioCodecSettings{}
	}
}

func expandAacCodecSettings(s *schema.Set) *medialive.AacSettings {
	return &medialive.AacSettings{}
}

func expandOutputGroups(outputGroups []interface{}) []*medialive.OutputGroup {
	return nil
}

func expandOutputGroupSettings(outputGroupSettings []interface{}) []*medialive.OutputGroupSettings {
	return nil
}

func expandHlsGroupSettings(hlsGroupSettings []interface{}) []*medialive.HlsGroupSettings {
	return nil
}

func expandOutputSettings(outputSettings []interface{}) []*medialive.OutputSettings {
	return nil
}

func expandHlsSettings(hlsSettings []interface{}) []*medialive.HlsSettings {
	return nil
}

func expandInputAttachmentSettings(s *schema.Set) *medialive.InputSettings {
	return &medialive.InputSettings{}
}

func expandTimecodeConfigs(s *schema.Set) *medialive.TimecodeConfig {
	return nil
}

func expandVideoDescriptions(videoDescriptions []interface{}) []*medialive.VideoDescription {
	return nil
}
