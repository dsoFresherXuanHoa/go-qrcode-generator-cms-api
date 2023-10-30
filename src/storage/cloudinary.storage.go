package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type cloudinaryStorage struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryStore(cld *cloudinary.Cloudinary) *cloudinaryStorage {
	return &cloudinaryStorage{cld: cld}
}

// TODO: Limit Image Size
func (s *cloudinaryStorage) UploadSingleImage(ctx context.Context, encode string) (*uploader.UploadResult, error) {
	storageDir := os.Getenv("CLOUDINARY_STORAGE_DIR")
	if resp, err := s.cld.Upload.Upload(ctx, encode, uploader.UploadParams{
		Folder: storageDir,
	}); err != nil {
		fmt.Println("Error while upload image to cloudinary: " + err.Error())
		return nil, err
	} else {
		return resp, nil
	}
}
