package validator

import "github.com/ahhhmadtlz/series_reader_backend/internal/config"

type Validator struct {
	uploadConfig config.Upload
}



func New(uploadConfig config.Upload) Validator {
	return Validator{
		uploadConfig: uploadConfig,
	}
}
 
 