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

func GenerateAnchorID() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func LogAnchor(anchorID string, file models.File, message string, timestamp time.Time) error {
	hash := sha256.New()
	hash.Write(file.Data)
	hashSum := hex.EncodeToString(hash.Sum(nil))

	log := fmt.Sprintf("Anchor ID: %s\nFile ID: %d\nFilename: %s\nExtension: %s\nFilesize: %d\nStatus: %s\nHash: %s\nMessage: %s\nTimestamp: %s\n\n",
		anchorID, file.ID, file.Metadata.Name, file.Metadata.Extension, file.Metadata.Size, file.Status, hashSum, message, timestamp.Format(time.RFC3339))

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
