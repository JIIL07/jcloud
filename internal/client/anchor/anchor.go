package anchor

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/JIIL07/jcloud/internal/client/delta"
	"github.com/JIIL07/jcloud/internal/client/models"
	"time"
)

type Berth struct {
	CurrentAnchor Anchor
}

type Anchor struct {
	ID        string
	Message   string
	Timestamp time.Time
	Deltas    map[int]*delta.Delta
	Log       string
}

func NewAnchor(files []models.File, message string, previousSnapshots map[int]*delta.Snapshot) (Anchor, error) {
	anchorID, err := GenerateAnchorID()
	if err != nil {
		return Anchor{}, err
	}
	timestamp := time.Now()

	deltas := make(map[int]*delta.Delta)
	fileSummaries := make([]string, 0, len(files))

	for _, file := range files {
		hash := sha256.New()
		hash.Write(file.Serialize())
		hashSum := hex.EncodeToString(hash.Sum(nil))

		newSnapshot := delta.NewSnapshot(file.Serialize())

		var deltaInfo string

		if previous, ok := previousSnapshots[file.ID]; ok {

			if d := newSnapshot.CreateDelta(previous); d != nil {
				deltas[file.ID] = d
				deltaInfo = fmt.Sprintf("Delta created: Original Hash: %s, New Hash: %s",
					d.OriginalHash, d.NewHash)
			} else {
				deltaInfo = "No changes detected, no delta created."
			}
		} else {
			deltaInfo = "No previous snapshot, no delta created."
		}

		fileSummary := fmt.Sprintf("File ID: %d, Filename: %s, Hash: %s, Delta Info: %s",
			file.ID, file.Metadata.Name, hashSum, deltaInfo)
		fileSummaries = append(fileSummaries, fileSummary)

		previousSnapshots[file.ID] = newSnapshot
	}

	log := LogSummary(anchorID, message, timestamp, fileSummaries)

	return Anchor{
		ID:        anchorID,
		Message:   message,
		Timestamp: timestamp,
		Deltas:    deltas,
		Log:       log,
	}, nil
}
