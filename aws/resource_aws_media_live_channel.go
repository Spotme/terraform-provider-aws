package aws

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/medialive"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

			// A list of destinations of the channel. For UDP outputs, there is onedestination
			// per output. For other types (HLS, for example), there isone destination per packager.
			"destinations": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 40),
								validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z-]+$`), "must contain only alphanumeric characters and hyphens"),
								validation.StringMatch(regexp.MustCompile(`^[a-zA-Z]`), "must begin with a letter"),
								validation.StringDoesNotMatch(regexp.MustCompile(`--`), "cannot contain two consecutive hyphens"),
								validation.StringDoesNotMatch(regexp.MustCompile(`-$`), "cannot end with a hyphen"),
							),
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
										Optional: true,
									},

									"url": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"username": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"media_package_settings": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"multiplex_settings": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

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
									// Advanced audio normalization settings.
									"audio_normalization_settings": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												// Audio normalization algorithm to use. itu17701 conforms to the CALM Act specification,
												// itu17702 conforms to the EBU R-128 specification.
												"algorithm": {
													Type:     schema.TypeString,
													Optional: true,
												},

												// When set to correctAudio the output audio is corrected using the chosen algorithm.
												// If set to measureOnly, the audio will be measured but not adjusted.
												"algorithm_control": {
													Type:     schema.TypeString,
													Optional: true,
												},

												// Target LKFS(loudness) to adjust volume to. If no value is entered, a default
												// value will be used according to the chosen algorithm. The CALM Act (1770-1)
												// recommends a target of -24 LKFS. The EBU R-128 specification (1770-2) recommends
												// a target of -23 LKFS.
												"target_lkfs": {
													Type:     schema.TypeFloat,
													Optional: true,
												},
											},
										},
									},

									// The name of the AudioSelector used as the source for this AudioDescription.
									"audio_selector_name": {
										Type:     schema.TypeString,
										Required: true,
									},

									// Applies only if audioTypeControl is useConfigured. The values for audioType
									// are defined in ISO-IEC 13818-1.
									"audio_type": {
										Type:     schema.TypeString,
										Optional: true,
									},

									// Determines how audio type is determined. followInput: If the input contains
									// an ISO 639 audioType, then that value is passed through to the output. If
									// the input contains no ISO 639 audioType, the value in Audio Type is included
									// in the output. useConfigured: The value in Audio Type is included in the
									// output.Note that this field and audioType are both ignored if inputType is
									// broadcasterMixedAd.
									"audio_type_control": {
										Type:     schema.TypeString,
										Required: true,
									},

									// Audio codec settings
									"codec_settings": {
										Type:     schema.TypeSet,
										Required: true,
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
																Type:     schema.TypeFloat,
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
																Type:     schema.TypeFloat,
																Optional: true,
															},

															"vbr_quality": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},

												//TODO:
												// Ac3 Settings
												// Eac3 Settings
												// Mp2 Settings
											},
										},
									},

									// Indicates the language of the audio output track. Only used if languageControlMode
									// is useConfigured, or there is no ISO 639 language code specified in the input.
									"language_code": {
										Type:     schema.TypeString,
										Optional: true,
									},

									// Choosing followInput will cause the ISO 639 language code of the output to
									// follow the ISO 639 language code of the input. The languageCode will be used
									// when useConfigured is set, or when followInput is selected but there is no
									// ISO 639 language code specified by the input.
									"language_code_control": {
										Type:     schema.TypeString,
										Required: true,
									},

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									//TODO: RemixSettings (settings that control how input audio channels are remixed into the output audio channels)

									"stream_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"avail_blanking": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"avail_blanking_image": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"password_param": {
													Type:     schema.TypeString,
													Optional: true,
												},

												"uri": {
													Type:     schema.TypeString,
													Required: true,
												},

												"username": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},

									"state": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						// TODO: avail_configuration (event-wide configuration settings for ad avail insertion).

						"blackout_slate": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"blackout_slate_image": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"password_param": {
													Type:     schema.TypeString,
													Optional: true,
												},

												"uri": {
													Type:     schema.TypeString,
													Required: true,
												},

												"username": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},

									"network_end_blackout": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"network_end_blackout_image": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"password_param": {
													Type:     schema.TypeString,
													Optional: true,
												},

												"uri": {
													Type:     schema.TypeString,
													Required: true,
												},

												"username": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},

									"network_id": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"state": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						// TODO: CaptionDescriptions

						"global_configuration": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"initial_audio_gain": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},

									"input_end_action": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "NONE",
									},

									"input_loss_behavior": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"black_frame_msec": {
													Type:     schema.TypeInt,
													Optional: true,
												},

												"input_loss_image_color": {
													Type:     schema.TypeString,
													Optional: true,
												},

												"input_loss_image_slate": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"password_param": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"uri": {
																Type:     schema.TypeString,
																Required: true,
															},

															"username": {
																Type:     schema.TypeString,
																Required: true,
															},
														},
													},
												},

												"input_loss_image_type": {
													Type:     schema.TypeString,
													Optional: true,
												},

												"repeat_frame_msec": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},

									"output_locking_mode": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "PIPELINE_LOCKING",
									},

									"output_timing_source": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "INPUT_CLOCK",
									},

									"support_low_framerate_inputs": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "DISABLED",
									},
								},
							},
						},

						"output_groups": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"output_group_settings": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"hls_group_settings": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"base_url_content": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"base_url_manifest": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"caption_language_setting": {
																Type:     schema.TypeString,
																Optional: true,
																Default:  "OMIT",
															},

															"codec_specification": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"constant_iv": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"client_cache": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"hls_cdn_settings": {
																Type:     schema.TypeSet,
																Required: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"hls_basic_put_settings": {
																			Type:     schema.TypeSet,
																			Required: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"connection_retry_interval": {
																						Type:     schema.TypeInt,
																						Optional: true,
																					},

																					"filecache_duration": {
																						Type:     schema.TypeInt,
																						Optional: true,
																					},

																					"num_retries": {
																						Type:     schema.TypeInt,
																						Optional: true,
																					},

																					"restart_delay": {
																						Type:     schema.TypeInt,
																						Optional: true,
																					},
																				},
																			},
																		},
																	},
																},
															},

															"encryption_type": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"hls_id3_segment_tagging": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"index_n_segments": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"input_loss_action": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"iv_in_manifest": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"iv_source": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"iframe_only_playlists": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"keep_segments": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"manifest_compression": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"manifest_duration_format": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"mode": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"output_selection": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"program_date_time": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"program_date_time_period": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"segmentation_mode": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"redundant_manifest": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"segment_length": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"destination": {
																Type:     schema.TypeSet,
																Required: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"destination_ref_id": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																	},
																},
															},

															"directory_structure": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"segments_per_subdirectory": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"stream_inf_resolution": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"timed_metadata_id3_frame": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"timed_metadata_id3_period": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"timestamp_delta_milliseconds": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"ts_file_mode": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},

									"outputs": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"audio_description_names": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},

												"caption_description_names": {
													Type:     schema.TypeList,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},

												"output_name": {
													Type:     schema.TypeString,
													Optional: true,
												},

												"output_settings": {
													Type:     schema.TypeSet,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"hls_output_settings": {
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"h_265_packaging_type": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},

																		"hls_settings": {
																			Type:     schema.TypeSet,
																			Optional: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"standard_hls_settings": {
																						Type:     schema.TypeSet,
																						Optional: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"audio_rendition_sets": {
																									Type:     schema.TypeString,
																									Optional: true,
																								},

																								"m3u8_settings": {
																									Type:     schema.TypeSet,
																									Optional: true,
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"audio_frames_per_pes": {
																												Type:     schema.TypeInt,
																												Required: true,
																											},

																											"audio_pids": {
																												Type:     schema.TypeString,
																												Optional: true,
																											},

																											"nielsen_id3_behavior": {
																												Type:     schema.TypeString,
																												Optional: true,
																											},

																											"pat_interval": {
																												Type:     schema.TypeInt,
																												Optional: true,
																											},

																											"pcr_control": {
																												Type:     schema.TypeString,
																												Optional: true,
																											},

																											"pcr_period": {
																												Type:     schema.TypeInt,
																												Optional: true,
																											},

																											"pcr_pid": {
																												Type:     schema.TypeString,
																												Optional: true,
																											},

																											"pmt_interval": {
																												Type:     schema.TypeInt,
																												Optional: true,
																											},

																											"pmt_pid": {
																												Type:     schema.TypeString,
																												Optional: true,
																											},

																											"program_num": {
																												Type:     schema.TypeInt,
																												Optional: true,
																											},

																											"scte_35_behavior": {
																												Type:     schema.TypeString,
																												Optional: true,
																											},

																											"scte_35_pid": {
																												Type:     schema.TypeString,
																												Optional: true,
																											},

																											"timed_metadata_behavior": {
																												Type:     schema.TypeString,
																												Optional: true,
																											},

																											"timed_metadata_pid": {
																												Type:     schema.TypeString,
																												Optional: true,
																											},

																											"transport_stream_id": {
																												Type:     schema.TypeInt,
																												Optional: true,
																											},

																											"video_pid": {
																												Type:     schema.TypeString,
																												Optional: true,
																											},
																										},
																									},
																								},
																							},
																						},
																					},

																					"audio_only_hls_settings": {
																						Type:     schema.TypeSet,
																						Optional: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"audio_track_type": {
																									Type:     schema.TypeString,
																									Required: true,
																								},

																								"audio_group_id": {
																									Type:     schema.TypeString,
																									Required: true,
																								},

																								"segment_type": {
																									Type:     schema.TypeString,
																									Required: true,
																								},
																							},
																						},
																					},
																				},
																			},
																		},

																		"name_modifier": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},

																		"segment_modifier": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
														},
													},
												},

												"video_description_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},

						"timecode_config": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:     schema.TypeString,
										Required: true,
									},

									// Threshold in frames beyond which output timecode is resynchronized to the
									// input timecode. Discrepancies below this threshold are permitted to avoid
									// unnecessary discontinuities in the output timecode. No timecode sync when
									// this is not specified.
									"sync_threshold": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},

						"video_descriptions": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"height": {
										Type:     schema.TypeInt,
										Required: true,
									},

									// Video codec settings.
									"codec_settings": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												// The frequency at which to capture frames for inclusion in the output. May
												// be specified in either seconds or milliseconds, as specified by captureIntervalUnits.
												"frame_capture_settings": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"capture_interval": {
																Type:     schema.TypeInt,
																Required: true,
															},

															"capture_interval_units": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},

												"h264_settings": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"adaptive_quantization": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"afd_signaling": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"bitrate": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"buf_fill_pct": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"buf_size": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"color_metadata": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"entropy_encoding": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"fixed_afd": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"flicker_aq": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"force_field_pictures": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"framerate_control": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"framerate_denominator": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"framerate_numerator": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"gop_b_reference": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"gop_closed_cadence": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"gop_num_b_frames": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"gop_size": {
																Type:     schema.TypeFloat,
																Optional: true,
															},

															"gop_size_units": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"level": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"look_ahead_rate_control": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"max_bitrate": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"min_i_interval": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"num_ref_frames": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"par_control": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"par_denominator": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"par_numerator": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"profile": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"quality_level": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"qvbr_quality_level": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"rate_control_mode": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"scan_type": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"scene_change_detect": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"slices": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"softness": {
																Type:     schema.TypeInt,
																Optional: true,
															},

															"spatial_aq": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"subgop_length": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"syntax": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"temporal_aq": {
																Type:     schema.TypeString,
																Optional: true,
															},

															"timecode_insertion": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"respond_to_afd": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "NONE",
									},

									"scaling_behavior": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "DEFAULT",
									},

									"sharpness": {
										Type:     schema.TypeInt,
										Required: true,
									},

									"width": {
										Type:     schema.TypeInt,
										Required: true,
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
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_end_behavior": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "CONTINUE",
									},

									"input_filter": {
										Type:     schema.TypeString,
										Required: true,
									},

									"filter_strength": {
										Type:     schema.TypeInt,
										Required: true,
									},

									"deblock_filter": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "DISABLED",
									},

									"denoise_filter": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "DISABLED",
									},

									"smpte2038_data_preference": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "IGNORE",
									},

									"caption_selectors": {
										Type:     schema.TypeList,
										Required: false,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"selector_settings": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"embedded_source_settings": {
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"source608_channel_number": {
																			Type:     schema.TypeInt,
																			Required: true,
																		},

																		"source608_track_number": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},

																		"convert608_to708": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},

																		"scte20_detection": {
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
				Default:  "DISABLED",
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"pipelines_running_count": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"reserved": {
				Type:     schema.TypeString,
				Computed: true,
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
		ChannelClass: aws.String(d.Get("channel_class").(string)),
		Name:         aws.String(d.Get("name").(string)),
		RoleArn:      aws.String(d.Get("role_arn").(string)),
		LogLevel:     aws.String(d.Get("log_level").(string)),
	}

	if v, ok := d.GetOk("input_specification"); ok {
		input.InputSpecification = expandInputSpecification(v.(*schema.Set))
	}

	if v, ok := d.GetOk("input_attachments"); ok && len(v.([]interface{})) > 0 {
		input.InputAttachments = expandInputAttachments(
			v.([]interface{}),
		)
	}

	if v, ok := d.GetOk("destinations"); ok && len(v.([]interface{})) > 0 {
		input.Destinations = expandDestinations(
			v.([]interface{}),
		)
	}

	if v, ok := d.GetOk("encoder_settings"); ok {
		input.EncoderSettings = expandEncoderSettings(
			v.(*schema.Set),
		)
	}

	if v := d.Get("tags").(map[string]interface{}); len(v) > 0 {
		input.Tags = keyvaluetags.New(v).IgnoreAws().MedialiveTags()
	}

	resp, err := conn.CreateChannel(input)
	if err != nil {
		return fmt.Errorf("Error creating MediaLive Channel: %s", err)
	}

	d.SetId(aws.StringValue(resp.Channel.Id))

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"CREATING"},
		Target:  []string{"IDLE"},
		Refresh: func() (interface{}, string, error) {
			input := &medialive.DescribeChannelInput{
				ChannelId: aws.String(d.Id()),
			}
			resp, err := conn.DescribeChannel(input)
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
		return fmt.Errorf("Error waiting MediaLive Channel (%s) to be created: %s", d.Id(), err)
	}

	return resourceAwsMediaLiveChannelRead(d, meta)
}

func resourceAwsMediaLiveChannelRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).medialiveconn
	ignoreTagsConfig := meta.(*AWSClient).IgnoreTagsConfig

	input := &medialive.DescribeChannelInput{
		ChannelId: aws.String(d.Id()),
	}

	resp, err := conn.DescribeChannel(input)
	if err != nil {
		if isAWSErr(err, medialive.ErrCodeNotFoundException, "") {
			log.Printf("[WARN] MediaLive Channel %s not found, error code (404)", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error describing MediaLive Channel(%s): %s", d.Id(), err)
	}

	d.Set("arn", resp.Arn)
	d.Set("name", resp.Name)
	d.Set("role_arn", resp.RoleArn)

	if err := d.Set("tags", keyvaluetags.MedialiveKeyValueTags(resp.Tags).IgnoreAws().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	return nil
}

func resourceAwsMediaLiveChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).medialiveconn

	input := &medialive.UpdateChannelInput{
		ChannelId: aws.String(d.Id()),
		Name:      aws.String(d.Get("name").(string)),
		RoleArn:   aws.String(d.Get("role_arn").(string)),
		LogLevel:  aws.String(d.Get("log_level").(string)),
	}

	if v, ok := d.GetOk("input_specification"); ok {
		input.InputSpecification = expandInputSpecification(v.(*schema.Set))
	}

	if v, ok := d.GetOk("input_attachments"); ok && len(v.([]interface{})) > 0 {
		input.InputAttachments = expandInputAttachments(
			v.([]interface{}),
		)
	}

	if v, ok := d.GetOk("destinations"); ok && len(v.([]interface{})) > 0 {
		input.Destinations = expandDestinations(
			v.([]interface{}),
		)
	}

	if v, ok := d.GetOk("encoder_settings"); ok {
		input.EncoderSettings = expandEncoderSettings(
			v.(*schema.Set),
		)
	}

	_, err := conn.UpdateChannel(input)
	if err != nil {
		return fmt.Errorf("Error updating MediaLive Channel: %s", err)
	}

	updateStateConf := &resource.StateChangeConf{
		Pending: []string{"UPDATING"},
		Target:  []string{"IDLE"},
		Refresh: func() (interface{}, string, error) {
			input := &medialive.DescribeChannelInput{
				ChannelId: aws.String(d.Id()),
			}
			resp, err := conn.DescribeChannel(input)
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
	_, err = updateStateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("Error waiting for MediaLive Channel (%s) update to be completed: %s", d.Id(), err)
	}

	return resourceAwsMediaLiveChannelRead(d, meta)
}

func resourceAwsMediaLiveChannelDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).medialiveconn
	input := &medialive.DeleteChannelInput{
		ChannelId: aws.String(d.Id()),
	}

	_, err := conn.DeleteChannel(input)
	if err != nil {
		if isAWSErr(err, medialive.ErrCodeNotFoundException, "") {
			return nil
		}
		return fmt.Errorf("Error deleting MediaLive Channel(%s): %s", d.Id(), err)
	}

	if err := waitForMediaLiveChannelDeletion(conn, d.Id()); err != nil {
		return fmt.Errorf("Error waiting for deleting MediaLive Channel(%s): %s", d.Id(), err)
	}

	return nil
}

func mediaLiveChannelRefreshFunc(conn *medialive.MediaLive, channelId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		channel, err := conn.DescribeChannel(&medialive.DescribeChannelInput{
			ChannelId: aws.String(channelId),
		})

		if isAWSErr(err, medialive.ErrCodeNotFoundException, "") {
			return nil, medialive.ChannelStateDeleted, nil
		}

		if err != nil {
			return nil, "", fmt.Errorf("error reading MediaLive Input(%s): %s", channelId, err)
		}

		if channel == nil {
			return nil, medialive.ChannelStateDeleted, nil
		}

		return channel, aws.StringValue(channel.State), nil
	}
}

func waitForMediaLiveChannelDeletion(conn *medialive.MediaLive, channelId string) error {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{medialive.ChannelStateDeleting},
		Target:         []string{medialive.ChannelStateDeleted},
		Refresh:        mediaLiveChannelRefreshFunc(conn, channelId),
		Timeout:        30 * time.Minute,
		NotFoundChecks: 1,
	}

	log.Printf("[DEBUG] Waiting for Media Live Channel (%s) deletion", channelId)
	_, err := stateConf.WaitForState()

	if isAWSErr(err, medialive.ErrCodeNotFoundException, "") {
		return nil
	}

	return err
}
