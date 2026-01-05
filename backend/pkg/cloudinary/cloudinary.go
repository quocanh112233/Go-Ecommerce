package cloudinary

import (
	"context"
	"mime/multipart"

	"go-ecommerce/internal/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Client wraps Cloudinary operations
type Client struct {
	cld *cloudinary.Cloudinary
}

// UploadResult contains the result of an upload
type UploadResult struct {
	URL      string
	PublicID string
}

// NewClient creates a new Cloudinary client
func NewClient(cfg *config.CloudinaryConfig) (*Client, error) {
	cld, err := cloudinary.NewFromParams(cfg.CloudName, cfg.APIKey, cfg.APISecret)
	if err != nil {
		return nil, err
	}
	return &Client{cld: cld}, nil
}

// Upload uploads a file to Cloudinary
func (c *Client) Upload(ctx context.Context, file multipart.File, folder string) (*UploadResult, error) {
	uploadParams := uploader.UploadParams{
		Folder: folder,
	}

	result, err := c.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		URL:      result.SecureURL,
		PublicID: result.PublicID,
	}, nil
}

// Delete removes a file from Cloudinary
func (c *Client) Delete(ctx context.Context, publicID string) error {
	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
