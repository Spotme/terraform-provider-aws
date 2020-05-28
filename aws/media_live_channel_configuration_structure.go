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
	if s.Len() > 0 {
		rawAacSettings := s.List()[0].(map[string]interface{})
		return &medialive.AacSettings{
			Bitrate:         aws.Float64(rawAacSettings["bitrate"].(float64)),
			CodingMode:      aws.String(rawAacSettings["coding_mode"].(string)),
			InputType:       aws.String(rawAacSettings["input_type"].(string)),
			Profile:         aws.String(rawAacSettings["profile"].(string)),
			RateControlMode: aws.String(rawAacSettings["rate_control_mode"].(string)),
			RawFormat:       aws.String(rawAacSettings["raw_format"].(string)),
			SampleRate:      aws.Float64(rawAacSettings["sample_rate"].(float64)),
			Spec:            aws.String(rawAacSettings["spec"].(string)),
			VbrQuality:      aws.String(rawAacSettings["vbr_quality"].(string)),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: AAC Specification can not be found")
		return &medialive.AacSettings{}
	}
}

func expandOutputGroups(outputGroups []interface{}) []*medialive.OutputGroup {
	var result []*medialive.OutputGroup
	if len(outputGroups) == 0 {
		log.Printf("[ERROR] MediaLive Channel: Output group is a required field")
		return nil
	}

	for _, v := range outputGroups {
		r := v.(map[string]interface{})

		result = append(result, &medialive.OutputGroup{
			Name:                aws.String(r["name"].(string)),
			OutputGroupSettings: expandOutputGroupSettings(r["output_group_settings"].(*schema.Set)),
			Outputs:             expandOutputs(r["outputs"].([]interface{})),
		})
	}
	return result
}

func expandOutputGroupSettings(s *schema.Set) *medialive.OutputGroupSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.OutputGroupSettings{
			HlsGroupSettings: expandHlsGroupSettings(settings["hls_group_settings"].(*schema.Set)),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: OutputGroupSettings can not be found")
		return &medialive.OutputGroupSettings{}
	}
}

func expandOutputs(outputs []interface{}) []*medialive.Output {
	var result []*medialive.Output
	if len(outputs) == 0 {
		log.Printf("[WARN] MediaLive Channel: Outputs are not specified")
		return nil
	}

	for _, v := range outputs {
		r := v.(map[string]interface{})

		result = append(result, &medialive.Output{
			OutputName:            aws.String(r["output_name"].(string)),
			AudioDescriptionNames: expandAudioDescriptionNames(r["audio_description_names"].([]string)),
			OutputSettings:        expandOutputSettings(r["output_settings"].(*schema.Set)),
			VideoDescriptionName:  aws.String(r["video_description_names"].(string)),
		})
	}
	return result
}

func expandOutputSettings(s *schema.Set) *medialive.OutputSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.OutputSettings{
			HlsOutputSettings: expandHlsOutputSettings(settings["hls_output_settings"].(*schema.Set)),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: OutputSettings can not be found")
		return &medialive.OutputSettings{}
	}
}

func expandHlsOutputSettings(s *schema.Set) *medialive.HlsOutputSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.HlsOutputSettings{
			HlsSettings:     expandHlsSettings(settings["hls_output_settings"].(*schema.Set)),
			NameModifier:    aws.String(settings["name_modifier"].(string)),
			SegmentModifier: aws.String(settings["segment_modifier"].(string)),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: HlsOutputSettings can not be found")
		return &medialive.HlsOutputSettings{}
	}
}

func expandHlsSettings(s *schema.Set) *medialive.HlsSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.HlsSettings{
			StandardHlsSettings: expandStandardHlsSettings(settings["standard_hls_settings"].(*schema.Set)),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: HlsSettings can not be found")
		return &medialive.HlsSettings{}
	}
}

func expandStandardHlsSettings(s *schema.Set) *medialive.StandardHlsSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.StandardHlsSettings{
			M3u8Settings: expandM3u8settings(settings["m3u8_settings"].(*schema.Set)),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: StandardHlsSettings can not be found")
		return &medialive.StandardHlsSettings{}
	}
}

func expandM3u8settings(s *schema.Set) *medialive.M3u8Settings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.M3u8Settings{
			AudioFramesPerPes:     aws.Int64(settings["audio_frames_per_pes"].(int64)),
			AudioPids:             aws.String(settings["audio_pids"].(string)),
			NielsenId3Behavior:    aws.String(settings["nielsen_id3_behavior"].(string)),
			PatInterval:           aws.Int64(settings["pat_interval"].(int64)),
			PcrControl:            aws.String(settings["pcr_control"].(string)),
			PcrPeriod:             aws.Int64(settings["pcr_period"].(int64)),
			PcrPid:                aws.String(settings["pcr_pid"].(string)),
			PmtInterval:           aws.Int64(settings["pmt_interval"].(int64)),
			PmtPid:                aws.String(settings["pmt_pid"].(string)),
			ProgramNum:            aws.Int64(settings["program_num"].(int64)),
			Scte35Behavior:        aws.String(settings["scte_35_behavior"].(string)),
			Scte35Pid:             aws.String(settings["scte_35_pid"].(string)),
			TimedMetadataBehavior: aws.String(settings["timed_metadata_behavior"].(string)),
			TimedMetadataPid:      aws.String(settings["timed_metadata_pid"].(string)),
			TransportStreamId:     aws.Int64(settings["transport_stream_id"].(int64)),
			VideoPid:              aws.String(settings["video_pid"].(string)),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: M3u8Settings can not be found")
		return &medialive.M3u8Settings{}
	}
}

func expandHlsGroupSettings(s *schema.Set) *medialive.HlsGroupSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.HlsGroupSettings{
			BaseUrlContent:         aws.String(settings["base_url_content"].(string)),
			CaptionLanguageSetting: aws.String(settings["caption_language_setting"].(string)),
			CodecSpecification:     aws.String(settings["codec_specification"].(string)),
			ConstantIv:             aws.String(settings["constant_iv"].(string)),
			ClientCache:            aws.String(settings["client_cache"].(string)),
			EncryptionType:         aws.String(settings["encryption_type"].(string)),
			HlsCdnSettings:         expandHlsCdnSettings(settings["hls_cdn_settings"].(*schema.Set)),
			HlsId3SegmentTagging:   aws.String(settings["hls_id3_segment_tagging"].(string)),
			IndexNSegments:         aws.Int64(settings["index_n_segments"].(int64)),
			InputLossAction:        aws.String(settings["input_loss_action"].(string)),
			IvInManifest:           aws.String(settings["iv_in_manifest"].(string)),
			IvSource:               aws.String(settings["iv_source"].(string)),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: HlsGroupSettings can not be found")
		return &medialive.HlsGroupSettings{}
	}
}

func expandHlsCdnSettings(s *schema.Set) *medialive.HlsCdnSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.HlsCdnSettings{
			HlsBasicPutSettings: expandHlsBasicPutSettings(settings["hls_basic_put_settings"].(*schema.Set)),
			//TODO: ADD support for Akamai CDN and MediaStore origin
			//HlsMediaStoreSettings: expandHlsMediaStoreSettings(settings["h264_settings"].(*schema.Set)),
			//HlsAkamaiSettings: expandHlsAkamaiSettings(settings["h264_settings"].(*schema.Set)),
		}
	} else {
		log.Printf("[WARN] MediaLive Channel: HlsCdnSettings can not be found")
		return &medialive.HlsCdnSettings{}
	}
}

func expandHlsBasicPutSettings(s *schema.Set) *medialive.HlsBasicPutSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.HlsBasicPutSettings{
			ConnectionRetryInterval: aws.Int64(settings["connection_retry_interval"].(int64)),
			FilecacheDuration:       aws.Int64(settings["filecache_duration"].(int64)),
			NumRetries:              aws.Int64(settings["num_retries"].(int64)),
			RestartDelay:            aws.Int64(settings["restart_delay"].(int64)),
		}
	} else {
		log.Printf("[WARN] MediaLive Channel: HlsBasicPutSettings can not be found")
		return &medialive.HlsBasicPutSettings{}
	}
}

func expandAudioDescriptionNames(audioDescriptionNames []string) []*string {
	var result []*string
	if len(audioDescriptionNames) == 0 {
		log.Printf("[ERROR] MediaLive Channel: No AudioDescriptionNames for Output")
		return nil
	}

	for _, v := range audioDescriptionNames {
		result = append(result, aws.String(v))
	}
	return result
}

func expandInputAttachmentSettings(s *schema.Set) *medialive.InputSettings {
	if s.Len() > 0 {
		rawInputSettings := s.List()[0].(map[string]interface{})
		return &medialive.InputSettings{
			DeblockFilter:     aws.String(rawInputSettings["deblock_filter"].(string)),
			DenoiseFilter:     aws.String(rawInputSettings["denoise_filter"].(string)),
			FilterStrength:    aws.Int64(rawInputSettings["filter_strength"].(int64)),
			InputFilter:       aws.String(rawInputSettings["input_filter"].(string)),
			SourceEndBehavior: aws.String(rawInputSettings["source_end_behavior"].(string)),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: InputSettings can not be found")
		return &medialive.InputSettings{}
	}
}

func expandTimecodeConfigs(s *schema.Set) *medialive.TimecodeConfig {
	if s.Len() > 0 {
		rawTimecodeConfig := s.List()[0].(map[string]interface{})
		return &medialive.TimecodeConfig{
			Source:        aws.String(rawTimecodeConfig["source"].(string)),
			SyncThreshold: aws.Int64(rawTimecodeConfig["sync_threshold"].(int64)),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: TimecodeConfig can not be found")
		return &medialive.TimecodeConfig{}
	}
}

func expandVideoDescriptions(videoDescriptions []interface{}) []*medialive.VideoDescription {
	var result []*medialive.VideoDescription
	if len(videoDescriptions) == 0 {
		log.Printf("[ERROR] MediaLive Channel: At least one video description is required for each encoder")
		return nil
	}

	for _, descs := range videoDescriptions {
		r := descs.(map[string]interface{})

		result = append(result, &medialive.VideoDescription{
			CodecSettings:   expandVideoCodecSettings(r["codec_settings"].(*schema.Set)),
			Height:          aws.Int64(r["height"].(int64)),
			Name:            aws.String(r["name"].(string)),
			RespondToAfd:    aws.String(r["respond_to_afd"].(string)),
			ScalingBehavior: aws.String(r["scaling_behavior"].(string)),
			Sharpness:       aws.Int64(r["sharpness"].(int64)),
			Width:           aws.Int64(r["width"].(int64)),
		})
	}
	return result
}

func expandVideoCodecSettings(s *schema.Set) *medialive.VideoCodecSettings {
	if s.Len() > 0 {
		rawVideoCodecSettings := s.List()[0].(map[string]interface{})
		return &medialive.VideoCodecSettings{
			H264Settings: expandH264Settings(rawVideoCodecSettings["h264_settings"].(*schema.Set)),
		}
	} else {
		log.Printf("[WARN] MediaLive Channel: VideoCodecSettings can not be found")
		return &medialive.VideoCodecSettings{}
	}
}

func expandH264Settings(s *schema.Set) *medialive.H264Settings {
	if s.Len() > 0 {
		rawSettings := s.List()[0].(map[string]interface{})
		return &medialive.H264Settings{
			AdaptiveQuantization: aws.String(rawSettings["adaptive_quantization"].(string)),
			AfdSignaling:         aws.String(rawSettings["afd_signaling"].(string)),
			Bitrate:              aws.Int64(rawSettings["bitrate"].(int64)),
			BufFillPct:           aws.Int64(rawSettings["buf_fill_pct"].(int64)),
			BufSize:              aws.Int64(rawSettings["buf_size"].(int64)),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: H264Settings can not be found")
		return &medialive.H264Settings{}
	}
}
