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

	testVersion := createTestVersion(t, fixtureBase+"version", testServiceID)

	var err error

	// Enable IO
	record(t, fixtureBase+"enable_product", func(c *Client) {
		_, err = c.EnableProduct(&ProductEnablementInput{
			ProductID: ProductImageOptimizer,
			ServiceID: testServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure we disable IO on the service after the test
	defer func() {
		record(t, fixtureBase+"disable_product", func(c *Client) {
			err = c.DisableProduct(&ProductEnablementInput{
				ProductID: ProductImageOptimizer,
				ServiceID: testServiceID,
			})
		})

		if err != nil {
			t.Fatal(err)
		}
	}()

	var defaultSettings *ImageOptimizerDefaultSettings

	// Fetch
	record(t, fixtureBase+"original_fetch", func(c *Client) {
		defaultSettings, err = c.GetImageOptimizerDefaultSettings(&GetImageOptimizerDefaultSettingsInput{
			ServiceID:      testServiceID,
			ServiceVersion: *testVersion.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	originalSettings := defaultSettings

	// Reset our settings back to the original
	defer func() {
		record(t, fixtureBase+"final_reset", func(c *Client) {
			if originalSettings != nil {
				_, err = c.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
					ServiceID:      testServiceID,
					ServiceVersion: *testVersion.Number,
					// ResizeFilter:   &originalSettings.ResizeFilter,
					Webp:        &originalSettings.Webp,
					WebpQuality: &originalSettings.WebpQuality,
					// JpegType:       &originalSettings.JpegType,
					JpegQuality: &originalSettings.JpegQuality,
					Upscale:     &originalSettings.Upscale,
					AllowVideo:  &originalSettings.AllowVideo,
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
	record(t, fixtureBase+"update_1_patch", func(c *Client) {
		defaultSettings, err = c.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
			ServiceID:      testServiceID,
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
	record(t, fixtureBase+"update_1_get", func(c *Client) {
		defaultSettings, err = c.GetImageOptimizerDefaultSettings(&GetImageOptimizerDefaultSettingsInput{
			ServiceID:      testServiceID,
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

	record(t, fixtureBase+"update_2_patch", func(c *Client) {
		defaultSettings, err = c.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
			ServiceID:      testServiceID,
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
	record(t, fixtureBase+"update_2_get", func(c *Client) {
		defaultSettings, err = c.GetImageOptimizerDefaultSettings(&GetImageOptimizerDefaultSettingsInput{
			ServiceID:      testServiceID,
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

	record(t, fixtureBase+"incorrect_fetch", func(c *Client) {
		_, err = c.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
			ServiceID:      testServiceID,
			ServiceVersion: *testVersion.Number,
			WebpQuality:    &newWebpQuality,
		})
	})
	if err == nil {
		t.Fatalf("missing err")
	}
	expectedErr := "Webp quality must be less than or equal to 100"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("expected error to include '%s'; got: %v", expectedErr, err)
	}

	// Confirm all resize_filter & jpeg_type values are accepted
	for _, resizeFilter := range []ResizeFilter{Lanczos3, Lanczos2, Bicubic, Bilinear, Nearest} {
		record(t, fixtureBase+"set_resize_filter/"+resizeFilter.String(), func(c *Client) {
			defaultSettings, err = c.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
				ServiceID:      testServiceID,
				ServiceVersion: *testVersion.Number,
				ResizeFilter:   &resizeFilter,
			})
		})
		if defaultSettings.ResizeFilter != resizeFilter.String() {
			t.Fatalf("expected ResizeFilter: %s; got: %s", resizeFilter.String(), defaultSettings.ResizeFilter)
		}
	}

	for _, jpegType := range []JpegType{Auto, Baseline, Progressive} {
		record(t, fixtureBase+"set_jpeg_type/"+jpegType.String(), func(c *Client) {
			defaultSettings, err = c.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
				ServiceID:      testServiceID,
				ServiceVersion: *testVersion.Number,
				JpegType:       &jpegType,
			})
		})
		if defaultSettings.JpegType != jpegType.String() {
			t.Fatalf("expected JpegType: %s; got: %s", jpegType.String(), defaultSettings.JpegType)
		}
	}

	// Confirm a full request is accepted - that all parameters in our library match the API's expectations
	record(t, fixtureBase+"set_full", func(c *Client) {
		defaultSettings, err = c.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
			ServiceID:      testServiceID,
			ServiceVersion: *testVersion.Number,
			ResizeFilter:   ToPointer(Lanczos3),
			Webp:           ToPointer(false),
			WebpQuality:    ToPointer(85),
			JpegType:       ToPointer(Auto),
			JpegQuality:    ToPointer(85),
			Upscale:        ToPointer(false),
			AllowVideo:     ToPointer(false),
		})
	})
}

func TestClient_GetImageOptimizerDefaultSettings_validation(t *testing.T) {
	var err error
	_, err = testClient.GetImageOptimizerDefaultSettings(&GetImageOptimizerDefaultSettingsInput{
		ServiceID:      "",
		ServiceVersion: 3,
	})
	if err != ErrMissingServiceID {
		t.Errorf("expected error %v; got: %v", ErrMissingServiceID, err)
	}

	_, err = testClient.GetImageOptimizerDefaultSettings(&GetImageOptimizerDefaultSettingsInput{
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
	_, err = testClient.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
		ServiceID:      "",
		ServiceVersion: 3,
		Upscale:        &newUpscale,
	})
	if err != ErrMissingServiceID {
		t.Errorf("expected error %v; got: %v", ErrMissingServiceID, err)
	}

	_, err = testClient.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
		Upscale:        &newUpscale,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("expected error %v; got: %v", ErrMissingServiceVersion, err)
	}

	_, err = testClient.UpdateImageOptimizerDefaultSettings(&UpdateImageOptimizerDefaultSettingsInput{
		ServiceID:      "foo",
		ServiceVersion: 3,
	})
	if err != ErrMissingImageOptimizerDefaultSetting {
		t.Errorf("expected error: %v, got: %v", ErrMissingImageOptimizerDefaultSetting, err)
	}
}
