package bitmovin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/bitmovintypes"
	"github.com/bitmovin/bitmovin-go/models"
)

func TestFactoryIsRegistered(t *testing.T) {
	_, err := provider.GetProviderFactory(Name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBitmovinFactory(t *testing.T) {
	cfg := config.Config{
		Bitmovin: &config.Bitmovin{
			APIKey:           "apikey",
			Endpoint:         "this.is.my.endpoint",
			Timeout:          uint(5),
			AccessKeyID:      "accesskey",
			SecretAccessKey:  "secretaccesskey",
			Destination:      "s3://some-output-bucket/",
			EncodingRegion:   "AWS_US_EAST_1",
			AWSStorageRegion: "US_EAST_1",
		},
	}
	provider, err := bitmovinFactory(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	bitmovinProvider, ok := provider.(*bitmovinProvider)
	if !ok {
		t.Fatalf("Wrong provider returned. Want elementalConductorProvider instance. Got %#v.", provider)
	}
	expected := &bitmovin.Bitmovin{
		APIKey:     stringToPtr("apikey"),
		APIBaseURL: stringToPtr("this.is.my.endpoint"),
	}
	if *bitmovinProvider.client.APIKey != *expected.APIKey {
		t.Errorf("Factory: wrong APIKey returned. Want %#v. Got %#v.", expected.APIKey, *bitmovinProvider.client.APIKey)
	}
	if *bitmovinProvider.client.APIBaseURL != *expected.APIBaseURL {
		t.Errorf("Factory: wrong APIKey returned. Want %#v. Got %#v.", expected.APIBaseURL, *bitmovinProvider.client.APIBaseURL)
	}
}

func TestCreatePreset(t *testing.T) {
	testPresetName := "this_is_an_audio_config_uuid"
	preset := getPreset()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/audio/aac":
			resp := models.AACCodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.AACCodecConfigurationData{
					Result: models.AACCodecConfiguration{
						ID: stringToPtr("this_is_an_audio_config_uuid"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264":
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						ID: stringToPtr(testPresetName),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	presetName, err := prov.CreatePreset(preset)
	if err != nil {
		t.Fatal(err)
	}
	if presetName != testPresetName {
		t.Error("expected ", testPresetName, "got ", presetName)
	}
}

func TestCreatePresetFailsOnAPIError(t *testing.T) {
	preset := getPreset()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/audio/aac":
			resp := models.AACCodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusError,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	_, err := prov.CreatePreset(preset)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestCreatePresetFailsOnGenericError(t *testing.T) {
	preset := getPreset()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/audio/aac":
			fmt.Fprintln(w, "Not proper json")
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	_, err := prov.CreatePreset(preset)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestDeletePreset(t *testing.T) {
	testPresetID := "i_want_to_delete_this"
	audioPresetID := "embedded_audio_id"
	customData := make(map[string]interface{})
	customData["audio"] = audioPresetID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/video/h264/" + testPresetID + "/customData":
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/audio/aac/" + audioPresetID:
			resp := models.AACCodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/" + testPresetID:
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.DeletePreset(testPresetID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeletePresetFailsOnAPIError(t *testing.T) {
	testPresetID := "i_want_to_delete_this"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/video/h264/" + testPresetID + "/customData":
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusError,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.DeletePreset(testPresetID)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestDeletePresetFailsOnGenericError(t *testing.T) {
	testPresetID := "i_want_to_delete_this"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/video/h264/" + testPresetID + "/customData":
			fmt.Fprintln(w, "Not proper json")
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.DeletePreset(testPresetID)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestGetPreset(t *testing.T) {
	testPresetID := "this_is_a_video_preset_id"
	audioPresetID := "this_is_a_audio_preset_id"
	customData := make(map[string]interface{})
	customData["audio"] = audioPresetID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/video/h264/" + testPresetID + "/customData":
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/audio/aac/" + audioPresetID:
			resp := models.AACCodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/" + testPresetID:
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	i, err := prov.GetPreset(testPresetID)
	if err != nil {
		t.Fatal(err)
	}
	expected := bitmovinPreset{
		Video: models.H264CodecConfiguration{CustomData: customData},
		Audio: models.AACCodecConfiguration{},
	}
	if !reflect.DeepEqual(i, expected) {
		t.Errorf("GetPreset: want %#v. Got %#v", expected, i)
	}
}

func TestGetPresetFailsOnAPIError(t *testing.T) {
	testPresetID := "this_is_a_video_preset_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/video/h264/" + testPresetID:
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusError,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	i, err := prov.GetPreset(testPresetID)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if i != nil {
		t.Errorf("GetPreset: got unexpected non-nil result: %#v", i)
	}
}

func TestGetPresetFailsOnGenericError(t *testing.T) {
	testPresetID := "this_is_a_video_preset_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/video/h264/" + testPresetID:
			fmt.Fprintln(w, "Not proper json")
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	i, err := prov.GetPreset(testPresetID)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if i != nil {
		t.Errorf("GetPreset: got unexpected non-nil result: %#v", i)
	}
}

func TestTranscodeWithS3Input(t *testing.T) {
	s3InputID := "this_is_the_s3_input_id"
	s3OutputID := "this_is_the_s3_output_id"
	encodingID := "this_is_the_master_encoding_id"
	manifestID := "this_is_the_master_manifest_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/inputs/s3":
			resp := models.S3InputResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.S3InputData{
					Result: models.S3InputItem{
						ID: stringToPtr(s3InputID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/outputs/s3":
			resp := models.S3OutputResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.S3OutputData{
					Result: models.S3OutputItem{
						ID: stringToPtr(s3OutputID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID1/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID1"
			customData["container"] = "mp4"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID2/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID2"
			customData["container"] = "m3u8"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID3/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID3"
			customData["container"] = "m3u8"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls":
			resp := models.HLSManifestResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.HLSManifestData{
					Result: models.HLSManifest{
						ID: stringToPtr(manifestID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings":
			resp := models.EncodingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.EncodingData{
					Result: models.Encoding{
						ID: stringToPtr(encodingID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID1",
			"/encoding/configurations/video/h264/videoID2",
			"/encoding/configurations/video/h264/videoID3":
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/streams":
			resp := models.StreamResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StreamData{
					Result: models.Stream{
						ID: stringToPtr("this_is_a_stream_id"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/mp4":
			resp := models.MP4MuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/ts":
			resp := models.TSMuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.TSMuxingData{
					Result: models.TSMuxing{
						ID: stringToPtr("this_is_a_ts_muxing_id"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/media":
			resp := models.MediaInfoResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/streams":
			resp := models.StreamInfoResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/start":
			resp := models.StartStopResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit " + r.URL.Path))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.Transcode(getJob("s3://bucket/folder/filename.mp4"))
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: encodingID,
		Status:        provider.StatusQueued,
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status: want %#v. Got %#v", expectedJobStatus, jobStatus)
	}
}

func TestTranscodeWithHTTPInput(t *testing.T) {
	httpInputID := "this_is_the_s3_input_id"
	s3OutputID := "this_is_the_s3_output_id"
	encodingID := "this_is_the_master_encoding_id"
	manifestID := "this_is_the_master_manifest_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/inputs/http":
			resp := models.HTTPInputResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.HTTPInputData{
					Result: models.HTTPInputItem{
						ID: stringToPtr(httpInputID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/outputs/s3":
			resp := models.S3OutputResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.S3OutputData{
					Result: models.S3OutputItem{
						ID: stringToPtr(s3OutputID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID1/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID1"
			customData["container"] = "mp4"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID2/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID2"
			customData["container"] = "m3u8"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID3/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID3"
			customData["container"] = "m3u8"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls":
			resp := models.HLSManifestResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.HLSManifestData{
					Result: models.HLSManifest{
						ID: stringToPtr(manifestID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings":
			resp := models.EncodingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.EncodingData{
					Result: models.Encoding{
						ID: stringToPtr(encodingID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID1",
			"/encoding/configurations/video/h264/videoID2",
			"/encoding/configurations/video/h264/videoID3":
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/streams":
			resp := models.StreamResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StreamData{
					Result: models.Stream{
						ID: stringToPtr("this_is_a_stream_id"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/mp4":
			resp := models.MP4MuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/ts":
			resp := models.TSMuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.TSMuxingData{
					Result: models.TSMuxing{
						ID: stringToPtr("this_is_a_ts_muxing_id"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/media":
			resp := models.MediaInfoResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/streams":
			resp := models.StreamInfoResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/start":
			resp := models.StartStopResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit " + r.URL.Path))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.Transcode(getJob("http://bucket.com/folder/filename.mp4"))
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: encodingID,
		Status:        provider.StatusQueued,
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status: want %#v. Got %#v", expectedJobStatus, jobStatus)
	}
}

func TestTranscodeWithHTTPSInput(t *testing.T) {
	httpsInputID := "this_is_the_s3_input_id"
	s3OutputID := "this_is_the_s3_output_id"
	encodingID := "this_is_the_master_encoding_id"
	manifestID := "this_is_the_master_manifest_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/inputs/https":
			resp := models.HTTPSInputResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.HTTPSInputData{
					Result: models.HTTPSInputItem{
						ID: stringToPtr(httpsInputID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/outputs/s3":
			resp := models.S3OutputResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.S3OutputData{
					Result: models.S3OutputItem{
						ID: stringToPtr(s3OutputID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID1/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID1"
			customData["container"] = "mp4"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID2/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID2"
			customData["container"] = "m3u8"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID3/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID3"
			customData["container"] = "m3u8"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls":
			resp := models.HLSManifestResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.HLSManifestData{
					Result: models.HLSManifest{
						ID: stringToPtr(manifestID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings":
			resp := models.EncodingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.EncodingData{
					Result: models.Encoding{
						ID: stringToPtr(encodingID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID1",
			"/encoding/configurations/video/h264/videoID2",
			"/encoding/configurations/video/h264/videoID3":
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/streams":
			resp := models.StreamResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StreamData{
					Result: models.Stream{
						ID: stringToPtr("this_is_a_stream_id"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/mp4":
			resp := models.MP4MuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/ts":
			resp := models.TSMuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.TSMuxingData{
					Result: models.TSMuxing{
						ID: stringToPtr("this_is_a_ts_muxing_id"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/media":
			resp := models.MediaInfoResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/streams":
			resp := models.StreamInfoResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/start":
			resp := models.StartStopResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit " + r.URL.Path))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.Transcode(getJob("https://bucket.com/folder/filename.mp4"))
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: encodingID,
		Status:        provider.StatusQueued,
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status: want %#v. Got %#v", expectedJobStatus, jobStatus)
	}
}

func TestTranscodeFailsOnAPIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/outputs/s3":
			resp := models.S3OutputResponse{
				Status: bitmovintypes.ResponseStatusError,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit " + r.URL.Path))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.Transcode(getJob("s3://bucket/folder/filename.mp4"))
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if jobStatus != nil {
		t.Errorf("Transcode: got unexpected non-nil result: %#v", jobStatus)
	}
}

func TestTranscodeFailsOnGenericError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/outputs/s3":
			fmt.Fprintln(w, "Not proper json")
		default:
			t.Fatal(errors.New("unexpected path hit " + r.URL.Path))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.Transcode(getJob("s3://bucket/folder/filename.mp4"))
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if jobStatus != nil {
		t.Errorf("Transcode: got unexpected non-nil result: %#v", jobStatus)
	}
}

func TestJobStatusReturnsFinishedIfEncodeAndManifestAreFinished(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("FINISHED"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/customData":
			resp := models.CustomDataResponse{
				Data: models.Data{
					Result: models.Result{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("FINISHED"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: testJobID,
		Status:        provider.StatusFinished,
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status: want %#v. Got %#v", expectedJobStatus, jobStatus)
	}
}

func TestJobStatusReturnsFinishedIfEncodeISFinishedAndNoManifestGenerationIsNeeded(t *testing.T) {
	testJobID := "this_is_a_job_id"
	customData := make(map[string]interface{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("FINISHED"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/customData":
			resp := models.CustomDataResponse{
				Data: models.Data{
					Result: models.Result{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: testJobID,
		Status:        provider.StatusFinished,
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status: want %#v. Got %#v", expectedJobStatus, jobStatus)
	}
}

func TestJobStatusReturnsStartedIfEncodeIsFinishedAndManifestIsRunning(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("FINISHED"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/customData":
			resp := models.CustomDataResponse{
				Data: models.Data{
					Result: models.Result{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("RUNNING"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: testJobID,
		Status:        provider.StatusStarted,
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status: want %#v. Got %#v", expectedJobStatus, jobStatus)
	}
}

func TestJobStatusReturnsStartedIfEncodeIsFinishedAndManifestIsCreatedAndStartReturnsNoError(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("FINISHED"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/customData":
			resp := models.CustomDataResponse{
				Data: models.Data{
					Result: models.Result{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("CREATED"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/start":
			resp := models.StartStopResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: testJobID,
		Status:        provider.StatusStarted,
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status: want %#v. Got %#v", expectedJobStatus, jobStatus)
	}
}

func TestJobStatusReturnsQueuedIfEncodeIsCreated(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("CREATED"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: testJobID,
		Status:        provider.StatusQueued,
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status: want %#v. Got %#v", expectedJobStatus, jobStatus)
	}
}

func TestJobStatusReturnsStartedIfEncodeIsRunning(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("RUNNING"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: testJobID,
		Status:        provider.StatusStarted,
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status: want %#v. Got %#v", expectedJobStatus, jobStatus)
	}
}

func TestJobStatusReturnsFailedIfEncodeFailed(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("ERROR"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: testJobID,
		Status:        provider.StatusFailed,
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status: want %#v. Got %#v", expectedJobStatus, jobStatus)
	}
}

func TestJobStatusReturnsFailedIfEncodeIsFinishedAndManifestFailed(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("FINISHED"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/customData":
			resp := models.CustomDataResponse{
				Data: models.Data{
					Result: models.Result{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("ERROR"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: testJobID,
		Status:        provider.StatusFailed,
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status: want %#v. Got %#v", expectedJobStatus, jobStatus)
	}
}

func TestJobStatusReturnsFailureOnAPIError(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusError,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: testJobID,
		Status:        provider.StatusFailed,
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status: want %#v. Got %#v", expectedJobStatus, jobStatus)
	}
}

func TestJobStatusReturnsErrorOnGenericError(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			fmt.Fprintln(w, "Not proper json")
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if jobStatus != nil {
		t.Errorf("Got unexpected non-nil JobStatus: %#v", jobStatus)
	}
}

func TestCancelJob(t *testing.T) {
	testJobID := "this_is_a_job_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/stop":
			resp := models.StartStopResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.CancelJob(testJobID)
	if err != nil {
		t.Fatal(err)
	}

}

func TestCancelJobFailsOnAPIError(t *testing.T) {
	testJobID := "this_is_a_job_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/stop":
			resp := models.StartStopResponse{
				Status: bitmovintypes.ResponseStatusError,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.CancelJob(testJobID)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestCancelJobFailsOnGenericError(t *testing.T) {
	testJobID := "this_is_a_job_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/stop":
			fmt.Fprintln(w, "Not proper json")
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.CancelJob(testJobID)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestHealthCheck(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings":
			resp := models.EncodingListResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.Healthcheck()
	if err != nil {
		t.Fatal(err)
	}
}

func TestHealthCheckFailsOnAPIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings":
			resp := models.EncodingListResponse{
				Status: bitmovintypes.ResponseStatusError,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.Healthcheck()
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestHealthCheckFailsOnGenericError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings":
			fmt.Fprintln(w, "Not proper json")
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.Healthcheck()
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestCapabilities(t *testing.T) {
	var prov bitmovinProvider
	expected := provider.Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "hls"},
		Destinations:  []string{"s3"},
	}
	cap := prov.Capabilities()
	if !reflect.DeepEqual(cap, expected) {
		t.Errorf("Capabilities: want %#v. Got %#v", expected, cap)
	}
}

func getBitmovinProvider(url string) bitmovinProvider {
	client := bitmovin.NewBitmovin("apikey", url+"/", int64(5))
	return bitmovinProvider{
		client: client,
		config: &config.Bitmovin{
			APIKey:          "apikey",
			Endpoint:        url + "/",
			Timeout:         uint(5),
			AccessKeyID:     "accesskey",
			SecretAccessKey: "secretaccesskey",
			Destination:     "s3://some-output-bucket/",
		},
	}
}

func getPreset() db.Preset {
	return db.Preset{
		Audio: db.AudioPreset{
			Bitrate: "128000",
			Codec:   "aac",
		},
		Container:   "mp4",
		Description: "my nice preset",
		Name:        "mp4_1080p",
		RateControl: "VBR",
		Video: db.VideoPreset{
			Profile:      "main",
			ProfileLevel: "3.1",
			Bitrate:      "3500000",
			Codec:        "h264",
			GopMode:      "fixed",
			GopSize:      "90",
			Height:       "1080",
		},
	}
}

func getJob(sourceMedia string) *db.Job {
	presets := []db.PresetMap{
		{
			Name: "mp4_1080p",
			ProviderMapping: map[string]string{
				Name: "videoID1",
			},
			OutputOpts: db.OutputOptions{Extension: "mp4"},
		},
		{
			Name: "hls_360p",
			ProviderMapping: map[string]string{
				Name: "videoID2",
			},
			OutputOpts: db.OutputOptions{Extension: "m3u8"},
		},
		{
			Name: "hls_480p",
			ProviderMapping: map[string]string{
				Name: "videoID3",
			},
			OutputOpts: db.OutputOptions{Extension: "m3u8"},
		},
	}
	outputs := make([]db.TranscodeOutput, len(presets))
	for i, preset := range presets {
		fileName := "output-" + preset.Name + "." + preset.OutputOpts.Extension
		if preset.OutputOpts.Extension == "m3u8" {
			fileName = "hls/output-" + preset.Name + ".m3u8"
		}
		outputs[i] = db.TranscodeOutput{
			Preset:   preset,
			FileName: fileName,
		}
	}
	job := &db.Job{
		ProviderName: Name,
		SourceMedia:  sourceMedia,
		StreamingParams: db.StreamingParams{
			SegmentDuration:  uint(4),
			PlaylistFileName: "hls/master_playlist.m3u8",
		},
		Outputs: outputs,
	}
	return job
}
