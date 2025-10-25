package dao

import (
	"fmt"

	"github.com/ASNMortred/AI-Hackathon/server/internal/database"
)

type MinioFileDAO struct{}

func NewMinioFileDAO() *MinioFileDAO { return &MinioFileDAO{} }

func (d *MinioFileDAO) Create(uid int, fileName, fileURL string) error {
	query := "INSERT INTO minio_files (uid, file_name, file_url) VALUES (?, ?, ?)"
	_, err := database.DB.Exec(query, uid, fileName, fileURL)
	if err != nil {
		return fmt.Errorf("failed to insert minio file record: %w", err)
	}
	return nil
}
