package anchor

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
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

func LogSummary(anchorID string, message string, timestamp time.Time, fileSummaries []string) string {
	logEntry := fmt.Sprintf(
		"AnchorFile ID: %s\nMessage: %s\nTimestamp: %s\nFiles:\n%s",
		anchorID, message, timestamp.Format(time.RFC3339), strings.Join(fileSummaries, "\n"),
	)
	return logEntry
}
