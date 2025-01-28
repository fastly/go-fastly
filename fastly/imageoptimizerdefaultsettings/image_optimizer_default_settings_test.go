package imageoptimizerdefaultsettings_test

import (
	"strings"
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products/imageoptimizer"
)

// TestClient_ImageOptimizerDefaultSettings tests the Image Optimizer Default Settings API
//
// NOTE: to run this, all backends on the active version of the test service must have shielding
// enabled, and the test service's account must be allowed to enable Image Optimization.
func TestClient_ImageOptimizerDefaultSettings(t *testing.T) {
	t.Parallel()

	fixtureBase := ""

	testVersion := fastly.CreateTestVersion(t, fixtureBase+"version", fastly.TestDeliveryServiceID)

	var err error

	// Enable IO
	fastly.Record(t, fixtureBase+"enable_product", func(c *fastly.Client) {
		_, err = imageoptimizer.Enable(c, fastly.TestDeliveryServiceID)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure we disable IO on the service after the test
	defer func() {
		fastly.Record(t, fixtureBase+"disable_product", func(c *fastly.Client) {
			_, err = imageoptimizer.Enable(c, fastly.TestDeliveryServiceID)

			if err != nil {
				t.Fatal(err)
			}
		})
	}()

	var defaultSettings *fastly.ImageOptimizerDefaultSettings

	// Fetch
	fastly.Record(t, fixtureBase+"original_fetch", func(c *fastly.Client) {
		defaultSettings, err = c.GetImageOptimizerDefaultSettings(&fastly.GetImageOptimizerDefaultSettingsInput{
			ServiceID:      fastly.TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	originalSettings := defaultSettings

	// Reset our settings back to the original
	defer func() {
		fastly.Record(t, fixtureBase+"final_reset", func(c *fastly.Client) {
			if originalSettings != nil {
				_, err = c.UpdateImageOptimizerDefaultSettings(&fastly.UpdateImageOptimizerDefaultSettingsInput{
					ServiceID:      fastly.TestDeliveryServiceID,
					ServiceVersion: *testVersion.Number,
					// just use default resizefilter & jpegtype since it doesn't matter much, and it's annoying
					// to parse the API output strings back into enums.
					ResizeFilter: fastly.ToPointer(fastly.ImageOptimizerNearest),
					Webp:         &originalSettings.Webp,
					WebpQuality:  &originalSettings.WebpQuality,
					JpegType:     fastly.ToPointer(fastly.ImageOptimizerAuto),
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
	fastly.Record(t, fixtureBase+"update_1_patch", func(c *fastly.Client) {
		defaultSettings, err = c.UpdateImageOptimizerDefaultSettings(&fastly.UpdateImageOptimizerDefaultSettingsInput{
			ServiceID:      fastly.TestDeliveryServiceID,
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
	fastly.Record(t, fixtureBase+"update_1_get", func(c *fastly.Client) {
		defaultSettings, err = c.GetImageOptimizerDefaultSettings(&fastly.GetImageOptimizerDefaultSettingsInput{
			ServiceID:      fastly.TestDeliveryServiceID,
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

	fastly.Record(t, fixtureBase+"update_2_patch", func(c *fastly.Client) {
		defaultSettings, err = c.UpdateImageOptimizerDefaultSettings(&fastly.UpdateImageOptimizerDefaultSettingsInput{
			ServiceID:      fastly.TestDeliveryServiceID,
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
	fastly.Record(t, fixtureBase+"update_2_get", func(c *fastly.Client) {
		defaultSettings, err = c.GetImageOptimizerDefaultSettings(&fastly.GetImageOptimizerDefaultSettingsInput{
			ServiceID:      fastly.TestDeliveryServiceID,
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

	fastly.Record(t, fixtureBase+"incorrect_fetch", func(c *fastly.Client) {
		_, err = c.UpdateImageOptimizerDefaultSettings(&fastly.UpdateImageOptimizerDefaultSettingsInput{
			ServiceID:      fastly.TestDeliveryServiceID,
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
	for _, resizeFilter := range []fastly.ImageOptimizerResizeFilter{fastly.ImageOptimizerLanczos3, fastly.ImageOptimizerLanczos2, fastly.ImageOptimizerBicubic, fastly.ImageOptimizerBilinear, fastly.ImageOptimizerNearest} {
		fastly.Record(t, fixtureBase+"set_resize_filter/"+resizeFilter.String(), func(c *fastly.Client) {
			defaultSettings, err = c.UpdateImageOptimizerDefaultSettings(&fastly.UpdateImageOptimizerDefaultSettingsInput{
				ServiceID:      fastly.TestDeliveryServiceID,
				ServiceVersion: *testVersion.Number,
				ResizeFilter:   &resizeFilter,
			})
		})
		if defaultSettings.ResizeFilter != resizeFilter.String() {
			t.Fatalf("expected ResizeFilter: %s; got: %s", resizeFilter.String(), defaultSettings.ResizeFilter)
		}
	}

	for _, jpegType := range []fastly.ImageOptimizerJpegType{fastly.ImageOptimizerAuto, fastly.ImageOptimizerBaseline, fastly.ImageOptimizerProgressive} {
		fastly.Record(t, fixtureBase+"set_jpeg_type/"+jpegType.String(), func(c *fastly.Client) {
			defaultSettings, err = c.UpdateImageOptimizerDefaultSettings(&fastly.UpdateImageOptimizerDefaultSettingsInput{
				ServiceID:      fastly.TestDeliveryServiceID,
				ServiceVersion: *testVersion.Number,
				JpegType:       &jpegType,
			})
		})
		if defaultSettings.JpegType != jpegType.String() {
			t.Fatalf("expected JpegType: %s; got: %s", jpegType.String(), defaultSettings.JpegType)
		}
	}

	// Confirm a full request is accepted - that all parameters in our library match the API's expectations
	fastly.Record(t, fixtureBase+"set_full", func(c *fastly.Client) {
		defaultSettings, err = c.UpdateImageOptimizerDefaultSettings(&fastly.UpdateImageOptimizerDefaultSettingsInput{
			ServiceID:      fastly.TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
			ResizeFilter:   fastly.ToPointer(fastly.ImageOptimizerLanczos3),
			Webp:           fastly.ToPointer(false),
			WebpQuality:    fastly.ToPointer(85),
			JpegType:       fastly.ToPointer(fastly.ImageOptimizerAuto),
			JpegQuality:    fastly.ToPointer(85),
			Upscale:        fastly.ToPointer(false),
			AllowVideo:     fastly.ToPointer(false),
		})
	})
}

func TestClient_GetImageOptimizerDefaultSettings_validation(t *testing.T) {
	var err error
	_, err = fastly.TestClient.GetImageOptimizerDefaultSettings(&fastly.GetImageOptimizerDefaultSettingsInput{
		ServiceID:      "",
		ServiceVersion: 3,
	})
	if err != fastly.ErrMissingServiceID {
		t.Errorf("expected error %v; got: %v", fastly.ErrMissingServiceID, err)
	}

	_, err = fastly.TestClient.GetImageOptimizerDefaultSettings(&fastly.GetImageOptimizerDefaultSettingsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != fastly.ErrMissingServiceVersion {
		t.Errorf("expected error %v; got: %v", fastly.ErrMissingServiceVersion, err)
	}
}

func TestClient_UpdateImageOptimizerDefaultSettings_validation(t *testing.T) {
	newUpscale := true

	var err error
	_, err = fastly.TestClient.UpdateImageOptimizerDefaultSettings(&fastly.UpdateImageOptimizerDefaultSettingsInput{
		ServiceID:      "",
		ServiceVersion: 3,
		Upscale:        &newUpscale,
	})
	if err != fastly.ErrMissingServiceID {
		t.Errorf("expected error %v; got: %v", fastly.ErrMissingServiceID, err)
	}

	_, err = fastly.TestClient.UpdateImageOptimizerDefaultSettings(&fastly.UpdateImageOptimizerDefaultSettingsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
		Upscale:        &newUpscale,
	})
	if err != fastly.ErrMissingServiceVersion {
		t.Errorf("expected error %v; got: %v", fastly.ErrMissingServiceVersion, err)
	}

	_, err = fastly.TestClient.UpdateImageOptimizerDefaultSettings(&fastly.UpdateImageOptimizerDefaultSettingsInput{
		ServiceID:      "foo",
		ServiceVersion: 3,
	})
	if err != fastly.ErrMissingImageOptimizerDefaultSetting {
		t.Errorf("expected error: %v, got: %v", fastly.ErrMissingImageOptimizerDefaultSetting, err)
	}
}
