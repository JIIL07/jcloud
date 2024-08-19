package anchor

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/models"
	"os"
	"time"
)

type Anchor struct {
	ID        string
	Message   string
	Timestamp time.Time
}

func NewAnchor(files []models.File, message string) (Anchor, error) {
	anchorID, err := GenerateAnchorID()
	if err != nil {
		return Anchor{}, err
	}
	timestamp := time.Now()

	for _, file := range files {
		hash := sha256.New()
		hash.Write(file.Data)
		hashSum := hex.EncodeToString(hash.Sum(nil))

		fmt.Printf("File ID: %d\nFilename: %s\nExtension: %s\nFilesize: %d\nStatus: %s\nHash: %s\n",
			file.ID, file.Metadata.Filename, file.Metadata.Extension, file.Metadata.Filesize, file.Status, hashSum)

		err := logAnchor(anchorID, file, message, timestamp)
		if err != nil {
			return Anchor{}, err
		}
	}

	return Anchor{
		ID:        anchorID,
		Message:   message,
		Timestamp: timestamp,
	}, nil
}

func GenerateAnchorID() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
func logAnchor(anchorID string, file models.File, message string, timestamp time.Time) error {
	hash := sha256.New()
	hash.Write(file.Data)
	hashSum := hex.EncodeToString(hash.Sum(nil))

	log := fmt.Sprintf("Anchor ID: %s\nFile ID: %d\nFilename: %s\nExtension: %s\nFilesize: %d\nStatus: %s\nHash: %s\nMessage: %s\nTimestamp: %s\n\n",
		anchorID, file.ID, file.Metadata.Filename, file.Metadata.Extension, file.Metadata.Filesize, file.Status, hashSum, message, timestamp.Format(time.RFC3339))

	f, err := os.OpenFile("anchor.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(log); err != nil {
		return err
	}

	return nil
}
