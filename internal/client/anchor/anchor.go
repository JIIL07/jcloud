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
}

func NewAnchor(files []models.File, message string, previousSnapshots map[int]*delta.Snapshot) (Anchor, error) {
	anchorID, err := GenerateAnchorID()
	if err != nil {
		return Anchor{}, err
	}
	timestamp := time.Now()

	deltas := make(map[int]*delta.Delta)

	for _, file := range files {
		hash := sha256.New()
		hash.Write(file.Data)
		hashSum := hex.EncodeToString(hash.Sum(nil))

		fmt.Printf("File ID: %d\nFilename: %s\nExtension: %s\nFilesize: %d\nStatus: %s\nHash: %s\n",
			file.ID, file.Metadata.Name, file.Metadata.Extension, file.Metadata.Size, file.Status, hashSum)

		newSnapshot := delta.NewSnapshot(file.Data)

		if previous, ok := previousSnapshots[file.ID]; ok {
			d := newSnapshot.CreateDelta(previous)
			if d != nil {
				deltas[file.ID] = d
			}
		}

		err := LogAnchor(anchorID, file, message, timestamp)
		if err != nil {
			return Anchor{}, err
		}

		previousSnapshots[file.ID] = newSnapshot
	}

	return Anchor{
		ID:        anchorID,
		Message:   message,
		Timestamp: timestamp,
		Deltas:    deltas,
	}, nil
}
