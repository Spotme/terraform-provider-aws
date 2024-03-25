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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			AudioDescriptions:   expandAudioDescriptions(rawEncoderSettings["audio_descriptions"].([]interface{})),
			CaptionDescriptions: expandCaptionDescriptions(rawEncoderSettings["caption_descriptions"].([]interface{})),
			OutputGroups:        expandOutputGroups(rawEncoderSettings["output_groups"].([]interface{})),
			TimecodeConfig:      expandTimecodeConfigs(rawEncoderSettings["timecode_config"].(*schema.Set)),
			VideoDescriptions:   expandVideoDescriptions(rawEncoderSettings["video_descriptions"].([]interface{})),
			FeatureActivations:  expandFeatureActivations(rawEncoderSettings["feature_activations"].(*schema.Set)),
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
			AudioSelectorName:   aws.String(r["audio_selector_name"].(string)),
			Name:                aws.String(r["name"].(string)),
			CodecSettings:       expandAudioCodecSettings(r["codec_settings"].(*schema.Set)),
			AudioTypeControl:    aws.String(r["audio_type_control"].(string)),
			LanguageCodeControl: aws.String(r["language_code_control"].(string)),
		})
	}
	return result
}

func expandFeatureActivations(s *schema.Set) *medialive.FeatureActivations {
	if s.Len() > 0 {
		rawConfig := s.List()[0].(map[string]interface{})
		return &medialive.FeatureActivations{
			InputPrepareScheduleActions: aws.String(rawConfig["input_prepare_schedule_actions"].(string)),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: FeatureActivations config is malformed")
		return &medialive.FeatureActivations{}
	}
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

		// we can now have either hls, rmtp group, or both settings specified,
		// to ensure a working api we need to ensure not to attach an empty
		// hls or rtmp group struct
		outputGroupSettings := medialive.OutputGroupSettings{}
		hlsGroupSettings := expandHlsGroupSettings(settings["hls_group_settings"].(*schema.Set))
		rtmpGroupSettings := expandRtmpGroupSettings(settings["rtmp_group_settings"].(*schema.Set))

		// hls and rtmp group settings dont implement comparison operators
		// so we are checking for hardcoded fields in our module to not be NIL
		if hlsGroupSettings.Mode != nil {
			outputGroupSettings.HlsGroupSettings = hlsGroupSettings
		}
		if rtmpGroupSettings.AuthenticationScheme != nil {
			outputGroupSettings.RtmpGroupSettings = rtmpGroupSettings
		}

		return &outputGroupSettings
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

		videoDescName := r["video_description_name"].(string)
		var videoDescNameAws *string
		if len(videoDescName) > 0 {
			videoDescNameAws = aws.String(videoDescName)
		}

		result = append(result, &medialive.Output{
			OutputName:              aws.String(r["output_name"].(string)),
			AudioDescriptionNames:   expandStringList(r["audio_description_names"].([]interface{})),
			CaptionDescriptionNames: expandStringList(r["caption_description_names"].([]interface{})),
			OutputSettings:          expandOutputSettings(r["output_settings"].(*schema.Set)),
			VideoDescriptionName:    videoDescNameAws,
		})
	}
	return result
}

func expandOutputSettings(s *schema.Set) *medialive.OutputSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})

		// we can now have either hls or rmtp outputs, or both settings specified,
		// to ensure a working api we need to ensure not to attach an empty
		// hls or rtmp group struct
		outputSettings := medialive.OutputSettings{}
		hlsOutputSettings := expandHlsOutputSettings(settings["hls_output_settings"].(*schema.Set))
		rtmpOutputSettings := expandRtmpOutputSettings(settings["rtmp_output_settings"].(*schema.Set))

		if (*hlsOutputSettings != medialive.HlsOutputSettings{}) {
			outputSettings.HlsOutputSettings = hlsOutputSettings
		}
		if (*rtmpOutputSettings != medialive.RtmpOutputSettings{}) {
			outputSettings.RtmpOutputSettings = rtmpOutputSettings
		}

		return &outputSettings
	} else {
		log.Printf("[ERROR] MediaLive Channel: OutputSettings can not be found")
		return &medialive.OutputSettings{}
	}
}

func expandHlsOutputSettings(s *schema.Set) *medialive.HlsOutputSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.HlsOutputSettings{
			HlsSettings:       expandHlsSettings(settings["hls_settings"].(*schema.Set)),
			NameModifier:      aws.String(settings["name_modifier"].(string)),
			H265PackagingType: aws.String(settings["h_265_packaging_type"].(string)),
			SegmentModifier:   aws.String(settings["segment_modifier"].(string)),
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
			StandardHlsSettings:  expandStandardHlsSettings(settings["standard_hls_settings"].(*schema.Set)),
			AudioOnlyHlsSettings: expandAudioOnlyHlsSettings(settings["audio_only_hls_settings"].(*schema.Set)),
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
			AudioRenditionSets: aws.String(settings["audio_rendition_sets"].(string)),
			M3u8Settings:       expandM3u8settings(settings["m3u8_settings"].(*schema.Set)),
		}
	} else {
		return nil
	}
}

// MARK: Caption Selectors

func expandCaptionSelectors(captionSelectors []interface{}) []*medialive.CaptionSelector {
	var result []*medialive.CaptionSelector

	for _, descs := range captionSelectors {
		r := descs.(map[string]interface{})

		result = append(result, &medialive.CaptionSelector{
			Name:             aws.String(r["name"].(string)),
			SelectorSettings: expandCaptionSelectorSettings(r["selector_settings"].(*schema.Set)),
		})
	}
	return result
}

func expandCaptionSelectorSettings(s *schema.Set) *medialive.CaptionSelectorSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.CaptionSelectorSettings{
			EmbeddedSourceSettings: expandEmbeddedSourceSettings(settings["embedded_source_settings"].(*schema.Set)),
		}
	} else {
		return nil
	}
}

func expandEmbeddedSourceSettings(s *schema.Set) *medialive.EmbeddedSourceSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.EmbeddedSourceSettings{
			Convert608To708:        aws.String(settings["convert608_to708"].(string)),
			Scte20Detection:        aws.String(settings["scte20_detection"].(string)),
			Source608ChannelNumber: aws.Int64(int64(settings["source608_channel_number"].(int))),
			Source608TrackNumber:   aws.Int64(int64(settings["source608_track_number"].(int))),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: EmbeddedSourceSettings can not be found")
		return &medialive.EmbeddedSourceSettings{}
	}
}

// MARK: Caption Descriptions

func expandCaptionDescriptions(captionDescriptions []interface{}) []*medialive.CaptionDescription {
	var result []*medialive.CaptionDescription

	for _, descs := range captionDescriptions {
		r := descs.(map[string]interface{})

		result = append(result, &medialive.CaptionDescription{
			CaptionSelectorName: aws.String(r["caption_selector_name"].(string)),
			Name:                aws.String(r["name"].(string)),
			DestinationSettings: expandCaptionDestinationSettings(r["destination_settings"].(*schema.Set)),
		})
	}
	return result
}

func expandCaptionDestinationSettings(s *schema.Set) *medialive.CaptionDestinationSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.CaptionDestinationSettings{
			WebvttDestinationSettings:   expandWebvttDestinationSettings(settings["webvtt_destination_settings"].(*schema.Set)),
			EmbeddedDestinationSettings: expandEmbeddedDestinationSettings(settings["embedded_destination_settings"].(*schema.Set)),
		}
	} else {
		return nil
	}
}

func expandWebvttDestinationSettings(s *schema.Set) *medialive.WebvttDestinationSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.WebvttDestinationSettings{
			StyleControl: aws.String(settings["style_control"].(string)),
		}
	} else {
		return nil
	}
}

func expandEmbeddedDestinationSettings(s *schema.Set) *medialive.EmbeddedDestinationSettings {
	if s.Len() > 0 {
		return &medialive.EmbeddedDestinationSettings{}
	} else {
		return nil
	}
}

func expandAudioOnlyHlsSettings(s *schema.Set) *medialive.AudioOnlyHlsSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.AudioOnlyHlsSettings{
			AudioGroupId:   aws.String(settings["audio_group_id"].(string)),
			AudioTrackType: aws.String(settings["audio_track_type"].(string)),
			SegmentType:    aws.String(settings["segment_type"].(string)),
		}
	} else {
		return nil
	}
}

func expandM3u8settings(s *schema.Set) *medialive.M3u8Settings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.M3u8Settings{
			AudioFramesPerPes:     aws.Int64(int64(settings["audio_frames_per_pes"].(int))),
			AudioPids:             aws.String(settings["audio_pids"].(string)),
			NielsenId3Behavior:    aws.String(settings["nielsen_id3_behavior"].(string)),
			PatInterval:           aws.Int64(int64(settings["pat_interval"].(int))),
			PcrControl:            aws.String(settings["pcr_control"].(string)),
			PmtPid:                aws.String(settings["pmt_pid"].(string)),
			ProgramNum:            aws.Int64(int64(settings["program_num"].(int))),
			Scte35Behavior:        aws.String(settings["scte_35_behavior"].(string)),
			Scte35Pid:             aws.String(settings["scte_35_pid"].(string)),
			TimedMetadataBehavior: aws.String(settings["timed_metadata_behavior"].(string)),
			TimedMetadataPid:      aws.String(settings["timed_metadata_pid"].(string)),
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
			CaptionLanguageSetting:     aws.String(settings["caption_language_setting"].(string)),
			CaptionLanguageMappings:    expandCaptionLanguageMapping(settings["caption_language_mapping"].([]interface{})),
			CodecSpecification:         aws.String(settings["codec_specification"].(string)),
			ClientCache:                aws.String(settings["client_cache"].(string)),
			HlsCdnSettings:             expandHlsCdnSettings(settings["hls_cdn_settings"].(*schema.Set)),
			HlsId3SegmentTagging:       aws.String(settings["hls_id3_segment_tagging"].(string)),
			IndexNSegments:             aws.Int64(int64(settings["index_n_segments"].(int))),
			InputLossAction:            aws.String(settings["input_loss_action"].(string)),
			IvInManifest:               aws.String(settings["iv_in_manifest"].(string)),
			IvSource:                   aws.String(settings["iv_source"].(string)),
			IFrameOnlyPlaylists:        aws.String(settings["iframe_only_playlists"].(string)),
			KeepSegments:               aws.Int64(int64(settings["keep_segments"].(int))),
			ManifestCompression:        aws.String(settings["manifest_compression"].(string)),
			ManifestDurationFormat:     aws.String(settings["manifest_duration_format"].(string)),
			Mode:                       aws.String(settings["mode"].(string)),
			OutputSelection:            aws.String(settings["output_selection"].(string)),
			ProgramDateTime:            aws.String(settings["program_date_time"].(string)),
			ProgramDateTimePeriod:      aws.Int64(int64(settings["program_date_time_period"].(int))),
			RedundantManifest:          aws.String(settings["redundant_manifest"].(string)),
			SegmentationMode:           aws.String(settings["segmentation_mode"].(string)),
			SegmentLength:              aws.Int64(int64(settings["segment_length"].(int))),
			Destination:                expandHlsDestinationRef(settings["destination"].(*schema.Set)),
			DirectoryStructure:         aws.String(settings["directory_structure"].(string)),
			SegmentsPerSubdirectory:    aws.Int64(int64(settings["segments_per_subdirectory"].(int))),
			StreamInfResolution:        aws.String(settings["stream_inf_resolution"].(string)),
			TimedMetadataId3Frame:      aws.String(settings["timed_metadata_id3_frame"].(string)),
			TimedMetadataId3Period:     aws.Int64(int64(settings["timed_metadata_id3_period"].(int))),
			TimestampDeltaMilliseconds: aws.Int64(int64(settings["timestamp_delta_milliseconds"].(int))),
			TsFileMode:                 aws.String(settings["ts_file_mode"].(string)),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: HlsGroupSettings can not be found")
		return &medialive.HlsGroupSettings{}
	}
}

func expandCaptionLanguageMapping(captionMappings []interface{}) []*medialive.CaptionLanguageMapping {
	var result []*medialive.CaptionLanguageMapping
	if len(captionMappings) == 0 {
		return nil
	}

	for _, v := range captionMappings {
		settings := v.(map[string]interface{})
		result = append(result, &medialive.CaptionLanguageMapping{
			CaptionChannel:      aws.Int64(int64(settings["caption_channel"].(int))),
			LanguageCode:        aws.String(settings["language_code"].(string)),
			LanguageDescription: aws.String(settings["language_description"].(string)),
		})
	}
	return result
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

func expandHlsDestinationRef(s *schema.Set) *medialive.OutputLocationRef {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.OutputLocationRef{
			DestinationRefId: aws.String(settings["destination_ref_id"].(string)),
		}
	} else {
		log.Printf("[WARN] MediaLive Channel: HLS Destination (OutputLocationRef) can not be found")
		return &medialive.OutputLocationRef{}
	}
}

func expandHlsBasicPutSettings(s *schema.Set) *medialive.HlsBasicPutSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.HlsBasicPutSettings{
			ConnectionRetryInterval: aws.Int64(int64(settings["connection_retry_interval"].(int))),
			FilecacheDuration:       aws.Int64(int64(settings["filecache_duration"].(int))),
			NumRetries:              aws.Int64(int64(settings["num_retries"].(int))),
			RestartDelay:            aws.Int64(int64(settings["restart_delay"].(int))),
		}
	} else {
		log.Printf("[WARN] MediaLive Channel: HlsBasicPutSettings can not be found")
		return &medialive.HlsBasicPutSettings{}
	}
}

func expandInputAttachmentSettings(s *schema.Set) *medialive.InputSettings {
	if s.Len() > 0 {
		rawInputSettings := s.List()[0].(map[string]interface{})
		return &medialive.InputSettings{
			DeblockFilter:           aws.String(rawInputSettings["deblock_filter"].(string)),
			DenoiseFilter:           aws.String(rawInputSettings["denoise_filter"].(string)),
			FilterStrength:          aws.Int64(int64(rawInputSettings["filter_strength"].(int))),
			InputFilter:             aws.String(rawInputSettings["input_filter"].(string)),
			SourceEndBehavior:       aws.String(rawInputSettings["source_end_behavior"].(string)),
			Smpte2038DataPreference: aws.String(rawInputSettings["smpte2038_data_preference"].(string)),
			CaptionSelectors:        expandCaptionSelectors(rawInputSettings["caption_selectors"].([]interface{})),
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
			SyncThreshold: aws.Int64(int64(rawTimecodeConfig["sync_threshold"].(int))),
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
			Height:          aws.Int64(int64(r["height"].(int))),
			Name:            aws.String(r["name"].(string)),
			RespondToAfd:    aws.String(r["respond_to_afd"].(string)),
			ScalingBehavior: aws.String(r["scaling_behavior"].(string)),
			Sharpness:       aws.Int64(int64(r["sharpness"].(int))),
			Width:           aws.Int64(int64(r["width"].(int))),
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
		h264Settings := &medialive.H264Settings{
			AdaptiveQuantization: aws.String(rawSettings["adaptive_quantization"].(string)),
			AfdSignaling:         aws.String(rawSettings["afd_signaling"].(string)),
			Bitrate:              aws.Int64(int64(rawSettings["bitrate"].(int))),
			BufFillPct:           aws.Int64(int64(rawSettings["buf_fill_pct"].(int))),
			BufSize:              aws.Int64(int64(rawSettings["buf_size"].(int))),
			ColorMetadata:        aws.String(rawSettings["color_metadata"].(string)),
			EntropyEncoding:      aws.String(rawSettings["entropy_encoding"].(string)),
			FlickerAq:            aws.String(rawSettings["flicker_aq"].(string)),
			ForceFieldPictures:   aws.String(rawSettings["force_field_pictures"].(string)),
			GopBReference:        aws.String(rawSettings["gop_b_reference"].(string)),
			GopClosedCadence:     aws.Int64(int64(rawSettings["gop_closed_cadence"].(int))),
			GopNumBFrames:        aws.Int64(int64(rawSettings["gop_num_b_frames"].(int))),
			GopSize:              aws.Float64(rawSettings["gop_size"].(float64)),
			GopSizeUnits:         aws.String(rawSettings["gop_size_units"].(string)),
			Level:                aws.String(rawSettings["level"].(string)),
			LookAheadRateControl: aws.String(rawSettings["look_ahead_rate_control"].(string)),
			NumRefFrames:         aws.Int64(int64(rawSettings["num_ref_frames"].(int))),
			ParControl:           aws.String(rawSettings["par_control"].(string)),
			QualityLevel:         aws.String(rawSettings["quality_level"].(string)),
			Profile:              aws.String(rawSettings["profile"].(string)),
			RateControlMode:      aws.String(rawSettings["rate_control_mode"].(string)),
			Syntax:               aws.String(rawSettings["syntax"].(string)),
			SceneChangeDetect:    aws.String(rawSettings["scene_change_detect"].(string)),
			SpatialAq:            aws.String(rawSettings["spatial_aq"].(string)),
			TemporalAq:           aws.String(rawSettings["temporal_aq"].(string)),
			TimecodeInsertion:    aws.String(rawSettings["timecode_insertion"].(string)),
		}
		framerateControl := rawSettings["framerate_control"].(string)
		h264Settings.FramerateControl = aws.String(framerateControl)

		if framerateControl == "SPECIFIED" {
			h264Settings.FramerateDenominator = aws.Int64(int64(rawSettings["framerate_denominator"].(int)))
			h264Settings.FramerateNumerator = aws.Int64(int64(rawSettings["framerate_numerator"].(int)))
		}
		return h264Settings
	} else {
		log.Printf("[ERROR] MediaLive Channel: H264Settings can not be found")
		return &medialive.H264Settings{}
	}
}

func expandRtmpGroupSettings(s *schema.Set) *medialive.RtmpGroupSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.RtmpGroupSettings{
			AuthenticationScheme: aws.String(settings["authentication_scheme"].(string)),
			CacheFullBehavior:    aws.String(settings["cache_full_behavior"].(string)),
			CacheLength:          aws.Int64(int64(settings["cache_length"].(int))),
			CaptionData:          aws.String(settings["caption_data"].(string)),
			InputLossAction:      aws.String(settings["input_loss_action"].(string)),
			RestartDelay:         aws.Int64(int64(settings["restart_delay"].(int))),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: RtmpGroupSettings can not be found")
		return &medialive.RtmpGroupSettings{}
	}
}

func expandRtmpOutputSettings(s *schema.Set) *medialive.RtmpOutputSettings {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.RtmpOutputSettings{
			CertificateMode:         aws.String(settings["certificate_mode"].(string)),
			ConnectionRetryInterval: aws.Int64(int64(settings["connection_retry_interval"].(int))),
			NumRetries:              aws.Int64(int64(settings["num_retries"].(int))),
			Destination:             expandRtmpOutputDestination(settings["destination"].(*schema.Set)),
		}
	} else {
		log.Printf("[ERROR] MediaLive Channel: RtmpOutputSettings can not be found")
		return &medialive.RtmpOutputSettings{}
	}
}

func expandRtmpOutputDestination(s *schema.Set) *medialive.OutputLocationRef {
	if s.Len() > 0 {
		settings := s.List()[0].(map[string]interface{})
		return &medialive.OutputLocationRef{
			DestinationRefId: aws.String(settings["destination_ref_id"].(string)),
		}
	} else {
		return nil
	}
}

func expandChannelMaintenanceCreate(tfList []interface{}) *medialive.MaintenanceCreateSettings {
	if tfList == nil {
		return nil
	}
	m := tfList[0].(map[string]interface{})

	settings := &medialive.MaintenanceCreateSettings{}
	if v, ok := m["maintenance_day"].(string); ok && v != "" {
		settings.MaintenanceDay = aws.String(v)
	}
	if v, ok := m["maintenance_start_time"].(string); ok && v != "" {
		settings.MaintenanceStartTime = aws.String(v)
	}

	return settings
}

func expandChannelMaintenanceUpdate(tfList []interface{}) *medialive.MaintenanceUpdateSettings {
	if tfList == nil {
		return nil
	}
	m := tfList[0].(map[string]interface{})

	settings := &medialive.MaintenanceUpdateSettings{}
	if v, ok := m["maintenance_day"].(string); ok && v != "" {
		settings.MaintenanceDay = aws.String(v)
	}
	if v, ok := m["maintenance_start_time"].(string); ok && v != "" {
		settings.MaintenanceStartTime = aws.String(v)
	}
	// NOTE: This field is only available in the update struct. To allow users to set a scheduled
	// date on update, it may be worth adding to the base schema.
	// if v, ok := m["maintenance_scheduled_date"].(string); ok && v != "" {
	// 	settings.MaintenanceScheduledDate = aws.String(v)
	// }

	return settings
}

func flattenChannelMaintenance(apiObject *medialive.MaintenanceStatus) []interface{} {
	if apiObject == nil {
		return nil
	}

	m := map[string]interface{}{
		"maintenance_day":        apiObject.MaintenanceDay,
		"maintenance_start_time": apiObject.MaintenanceStartTime,
	}

	return []interface{}{m}
}
