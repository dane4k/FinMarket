package imgur

import "errors"

var (
	ErrUploadingPic   = errors.New("error uploading picture to imgur.com")
	ErrDownloadingPic = errors.New("error downloading picture")
)
