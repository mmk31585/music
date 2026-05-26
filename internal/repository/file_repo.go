package repository

import "mime/multipart"

type FileRepository interface {
	Save(fileHeader *multipart.FileHeader, destPath string) error
	Delete(filePath string) error
}
