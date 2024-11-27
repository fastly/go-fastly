package fastly

import (
	"strings"
	"testing"
)

// TestClient_ImageOptimizerDefaultSettings tests the Image Optimizer Default Settings API
//
// NOTE: to run this, all backends on the active version of the test service must have shielding
// enabled, and the test service's account must be allowed to enable Image Optimization.
func TestClient_ImageOptimizerDefaultSettings(t *testing.T) {
	t.Parallel()

	fixtureBase := "image_optimizer_default_settings/"

	testVersion := createTestVersion(t, fixtureBase+"version", TestDeliveryServiceID)

	var err error

	// Enable IO
	Record(t, fixtureBase+"enable_product", func(c *Client) {
		_, err = c.EnableProduct(&ProductEnablementInput{
			ProductID: ProductImageOptimizer,
			ServiceID: TestDeliveryServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure we disable IO on the service after the test
	defer func() {
		Record(t, fixtureBase+"disable_product", func(c *Client) {
			err = c.DisableProduct(&ProductEnablementInput{
				ProductID: ProductImageOptimizer,
				ServiceID: TestDeliveryServiceID,
			})
		})

		if err != nil {
			t.Fatal(err)
		}
	}()

	var defaultSettings *ImageOptimizerDefaultSettings

	// Fetch
	Record(t, fixtureBase+"original_fetch", func(c *Client) {
		defaultSettings, err = c.GetImageOptimizerDefaultSettings(&GetImageOptimizerDefaultSettingsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	originalSettings := defaultSettings

	// Reset our settings back to the original
	defer func() {
		Record(t, fixtureBase+"final_reset", func(c *Client) {
			if originalSettings != nil {
				_, err = c.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
					ServiceID:      TestDeliveryServiceID,
					ServiceVersion: *testVersion.Number,
					// just use default resizefilter & jpegtype since it doesn't matter much, and it's annoying
					// to parse the API output strings back into enums.
					ResizeFilter: ToPointer(ImageOptimizerNearest),
					Webp:         &originalSettings.Webp,
					WebpQuality:  &originalSettings.WebpQuality,
					JpegType:     ToPointer(ImageOptimizerAuto),
					JpegQuality:  &originalSettings.JpegQuality,
					Upscale:      &originalSettings.Upscale,
					AllowVideo:   &originalSettings.AllowVideo,
				})
			}
		})
		if err != nil {
			t.Fatal(err)
		}
	}()

	newWebp := false
	newWebpQuality := 20
	newUpscale := false

	// Change some stuff
	Record(t, fixtureBase+"update_1_patch", func(c *Client) {
		defaultSettings, err = c.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
			Webp:           &newWebp,
			WebpQuality:    &newWebpQuality,
			Upscale:        &newUpscale,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if defaultSettings.Webp != newWebp {
		t.Fatalf("expected Webp: %v; got: %v", newWebp, defaultSettings.Webp)
	}
	if defaultSettings.WebpQuality != newWebpQuality {
		t.Fatalf("expected WebpQuality %v; got: %v", newWebp, defaultSettings.WebpQuality)
	}
	if defaultSettings.Upscale != newUpscale {
		t.Fatalf("expected Upscale: %v; got: %v", newUpscale, defaultSettings.Webp)
	}

	// Confirm our changes were applied permanently
	Record(t, fixtureBase+"update_1_get", func(c *Client) {
		defaultSettings, err = c.GetImageOptimizerDefaultSettings(&GetImageOptimizerDefaultSettingsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if defaultSettings.Webp != newWebp {
		t.Fatalf("expected Webp: %v; got: %v", newWebp, defaultSettings.Webp)
	}
	if defaultSettings.WebpQuality != newWebpQuality {
		t.Fatalf("expected WebpQuality %v; got: %v", newWebp, defaultSettings.WebpQuality)
	}
	if defaultSettings.Upscale != newUpscale {
		t.Fatalf("expected Upscale: %v; got: %v", newUpscale, defaultSettings.Webp)
	}

	// Change settings again (to ensure we don't accidentally pass a test by just having the settings already set)
	newWebp = true
	newWebpQuality = 42
	newUpscale = true

	Record(t, fixtureBase+"update_2_patch", func(c *Client) {
		defaultSettings, err = c.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
			Webp:           &newWebp,
			WebpQuality:    &newWebpQuality,
			Upscale:        &newUpscale,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if defaultSettings.Webp != newWebp {
		t.Fatalf("expected Webp: %v; got: %v", newWebp, defaultSettings.Webp)
	}
	if defaultSettings.WebpQuality != newWebpQuality {
		t.Fatalf("expected WebpQuality %v; got: %v", newWebp, defaultSettings.WebpQuality)
	}
	if defaultSettings.Upscale != newUpscale {
		t.Fatalf("expected Upscale: %v; got: %v", newUpscale, defaultSettings.Webp)
	}

	// Confirm our changes were applied permanently (again)
	Record(t, fixtureBase+"update_2_get", func(c *Client) {
		defaultSettings, err = c.GetImageOptimizerDefaultSettings(&GetImageOptimizerDefaultSettingsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if defaultSettings.Webp != newWebp {
		t.Fatalf("expected Webp: %v; got: %v", newWebp, defaultSettings.Webp)
	}
	if defaultSettings.WebpQuality != newWebpQuality {
		t.Fatalf("expected WebpQuality %v; got: %v", newWebp, defaultSettings.WebpQuality)
	}
	if defaultSettings.Upscale != newUpscale {
		t.Fatalf("expected Upscale: %v; got: %v", newUpscale, defaultSettings.Webp)
	}

	// Apply a setting that produces a server-side error, and confirm it's handled well.
	newWebpQuality = 105

	Record(t, fixtureBase+"incorrect_fetch", func(c *Client) {
		_, err = c.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
			WebpQuality:    &newWebpQuality,
		})
	})
	if err == nil {
		t.Fatalf("missing err")
	}
	expectedErr := "WebP quality must be less than or equal to 100"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("expected error to include '%s'; got: %v", expectedErr, err)
	}

	// Confirm all resize_filter & jpeg_type values are accepted
	for _, resizeFilter := range []ImageOptimizerResizeFilter{ImageOptimizerLanczos3, ImageOptimizerLanczos2, ImageOptimizerBicubic, ImageOptimizerBilinear, ImageOptimizerNearest} {
		Record(t, fixtureBase+"set_resize_filter/"+resizeFilter.String(), func(c *Client) {
			defaultSettings, err = c.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *testVersion.Number,
				ResizeFilter:   &resizeFilter,
			})
		})
		if defaultSettings.ResizeFilter != resizeFilter.String() {
			t.Fatalf("expected ResizeFilter: %s; got: %s", resizeFilter.String(), defaultSettings.ResizeFilter)
		}
	}

	for _, jpegType := range []ImageOptimizerJpegType{ImageOptimizerAuto, ImageOptimizerBaseline, ImageOptimizerProgressive} {
		Record(t, fixtureBase+"set_jpeg_type/"+jpegType.String(), func(c *Client) {
			defaultSettings, err = c.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *testVersion.Number,
				JpegType:       &jpegType,
			})
		})
		if defaultSettings.JpegType != jpegType.String() {
			t.Fatalf("expected JpegType: %s; got: %s", jpegType.String(), defaultSettings.JpegType)
		}
	}

	// Confirm a full request is accepted - that all parameters in our library match the API's expectations
	Record(t, fixtureBase+"set_full", func(c *Client) {
		defaultSettings, err = c.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
			ResizeFilter:   ToPointer(ImageOptimizerLanczos3),
			Webp:           ToPointer(false),
			WebpQuality:    ToPointer(85),
			JpegType:       ToPointer(ImageOptimizerAuto),
			JpegQuality:    ToPointer(85),
			Upscale:        ToPointer(false),
			AllowVideo:     ToPointer(false),
		})
	})
}

func TestClient_GetImageOptimizerDefaultSettings_validation(t *testing.T) {
	var err error
	_, err = TestClient.GetImageOptimizerDefaultSettings(&GetImageOptimizerDefaultSettingsInput{
		ServiceID:      "",
		ServiceVersion: 3,
	})
	if err != ErrMissingServiceID {
		t.Errorf("expected error %v; got: %v", ErrMissingServiceID, err)
	}

	_, err = TestClient.GetImageOptimizerDefaultSettings(&GetImageOptimizerDefaultSettingsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("expected error %v; got: %v", ErrMissingServiceVersion, err)
	}
}

func TestClient_UpdateImageOptimizerDefaultSettings_validation(t *testing.T) {
	newUpscale := true

	var err error
	_, err = TestClient.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
		ServiceID:      "",
		ServiceVersion: 3,
		Upscale:        &newUpscale,
	})
	if err != ErrMissingServiceID {
		t.Errorf("expected error %v; got: %v", ErrMissingServiceID, err)
	}

	_, err = TestClient.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
		Upscale:        &newUpscale,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("expected error %v; got: %v", ErrMissingServiceVersion, err)
	}

	_, err = TestClient.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
		ServiceID:      "foo",
		ServiceVersion: 3,
	})
	if err != ErrMissingImageOptimizerDefaultSetting {
		t.Errorf("expected error: %v, got: %v", ErrMissingImageOptimizerDefaultSetting, err)
	}
}
