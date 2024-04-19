package fastly

import (
	"encoding/json"
	"fmt"
)

// ResizeFilter is a base for the different ResizeFilter variants.
type ResizeFilter int64

func (r ResizeFilter) String() string {
	switch r {
	case Lanczos3:
		return "lanczos3"
	case Lanczos2:
		return "lanczos2"
	case Bicubic:
		return "bicubic"
	case Bilinear:
		return "bilinear"
	case Nearest:
		return "nearest"
	}
	return "lanczos3" // default
}

const (
	Lanczos3 ResizeFilter = iota
	Lanczos2
	Bicubic
	Bilinear
	Nearest
)

// ImageOptimizerDefaultSettings represents the returned default settings for Image Optimizer on a given service.
type ImageOptimizerDefaultSettings struct {
	// Type of filter to use while resizing an image.
	ResizeFilter string `json:"resize_filter"`
	// Whether or not to default to webp output when the client supports it. This is equivalent to adding "auto=webp" to all image optimizer requests.
	Webp bool `json:"webp"`
	// Default quality to use with webp output. This can be overriden with the second option in the "quality" URL parameter on specific image optimizer
	// requests.
	WebpQuality int `json:"webp_quality"`
	// Default type of jpeg output to use. This can be overriden with "format=bjpeg" and "format=pjpeg" on specific image optimizer requests.
	JpegType string `json:"jpeg_type"`
	// Default quality to use with jpeg output. This can be overridden with the "quality" parameter on specific image optimizer requests.
	JpegQuality int `json:"jpeg_quality"`
	// Whether or not we should allow output images to render at sizes larger than input.
	Upscale bool `json:"upscale"`
	// Enables GIF to MP4 transformations on this service.
	AllowVideo bool `json:"allow_video"`
}

// GetImageOptimizerDefaultSettingsInput is used as input to the
// GetImageOptimizerDefaultSettings function.
type GetImageOptimizerDefaultSettingsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// UpdateImageOptimizerDefaultSettingsInput is used as input to the
// UpdateImageOptimizerDefaultSettings function.
//
// A minimum of one optional field is required.
type UpdateImageOptimizerDefaultSettingsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string `json:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `json:"-"`
	// Type of filter to use while resizing an image.
	ResizeFilter *ResizeFilter `json:"resize_filter,omitempty"`
	// Whether or not to default to webp output when the client supports it. This is equivalent to adding "auto=webp" to all image optimizer requests.
	Webp *bool `json:"webp,omitempty"`
	// Default quality to use with webp output. This can be overriden with the second option in the "quality" URL parameter on specific image optimizer
	// requests.
	WebpQuality *int `json:"webp_quality,omitempty"`
	// Default type of jpeg output to use. This can be overriden with "format=bjpeg" and "format=pjpeg" on specific image optimizer requests.
	JpegType *string `json:"jpeg_type,omitempty"`
	// Default quality to use with jpeg output. This can be overridden with the "quality" parameter on specific image optimizer requests.
	JpegQuality *int `json:"jpeg_quality,omitempty"`
	// Whether or not we should allow output images to render at sizes larger than input.
	Upscale *bool `json:"upscale,omitempty"`
	// Enables GIF to MP4 transformations on this service.
	AllowVideo *bool `json:"allow_video,omitempty"`
}

// GetImageOptimizerDefaultSettings retrives the current Image Optimizer default settings on a given service version.
//
// Returns (nil, nil) if no default settings are set.
func (c *Client) GetImageOptimizerDefaultSettings(i *GetImageOptimizerDefaultSettingsInput) (*ImageOptimizerDefaultSettings, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/image_optimizer_default_settings", i.ServiceID, i.ServiceVersion)

	resp, err := c.Get(path, nil)
	if err != nil {
		if herr, ok := err.(*HTTPError); ok {
			if herr.StatusCode == 404 {
				// API endpoint returns 404 for services without Image Optimizer settings set.
				return nil, nil
			}
		}
		return nil, err
	}
	defer resp.Body.Close()

	var iods *ImageOptimizerDefaultSettings
	if err := json.NewDecoder(resp.Body).Decode(&iods); err != nil {
		return nil, err
	}

	return iods, nil
}

// UpdateImageOptimizerDefaultSettings changes one or more Image Optimizer default settings on a given service version.
func (c *Client) UpdateImageOptimizerDefaultSettings(i *UpdateImageOptimizerDefaultSettingsInput) (*ImageOptimizerDefaultSettings, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}
	if i.ResizeFilter == nil && i.Webp == nil && i.WebpQuality == nil && i.JpegType == nil && i.JpegQuality == nil && i.Upscale == nil && i.AllowVideo == nil {
		return nil, ErrMissingImageOptimizerDefaultSetting
	}

	path := fmt.Sprintf("/service/%s/version/%d/image_optimizer_default_settings", i.ServiceID, i.ServiceVersion)

	resp, err := c.PatchJSON(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var iods *ImageOptimizerDefaultSettings
	if err := json.NewDecoder(resp.Body).Decode(&iods); err != nil {
		return nil, err
	}

	return iods, nil
}
