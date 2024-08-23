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
	// Generate a unique anchor ID
	anchorID, err := GenerateAnchorID()
	if err != nil {
		return Anchor{}, err
	}
	timestamp := time.Now()

	// Initialize deltas map and file summaries slice
	deltas := make(map[int]*delta.Delta)
	fileSummaries := make([]string, 0, len(files)) // Preallocate to avoid resizing

	for _, file := range files {
		// Calculate the hash for the file data
		hash := sha256.New()
		hash.Write(file.Serialize())
		hashSum := hex.EncodeToString(hash.Sum(nil))

		// Create a new snapshot from the file data
		newSnapshot := delta.NewSnapshot(file.Serialize())

		var deltaInfo string

		// Check for previous snapshots and create a delta if applicable
		if previous, ok := previousSnapshots[file.ID]; ok {

			if d := newSnapshot.CreateDelta(previous); d != nil {
				deltas[file.ID] = d
				// Log the delta details
				deltaInfo = fmt.Sprintf("Delta created: Original Hash: %s, New Hash: %s",
					d.OriginalHash, d.NewHash)
			} else {
				deltaInfo = "No changes detected, no delta created."
			}
		} else {
			deltaInfo = "No previous snapshot, no delta created."
		}

		// Collect summary information for the file
		fileSummary := fmt.Sprintf("File ID: %d, Filename: %s, Hash: %s, Delta Info: %s",
			file.ID, file.Metadata.Name, hashSum, deltaInfo)
		fileSummaries = append(fileSummaries, fileSummary)

		// Update the previous snapshot map with the new snapshot
		previousSnapshots[file.ID] = newSnapshot
	}

	log := LogSummary(anchorID, message, timestamp, fileSummaries)

	// Return the created anchor
	return Anchor{
		ID:        anchorID,
		Message:   message,
		Timestamp: timestamp,
		Deltas:    deltas,
		Log:       log,
	}, nil
}
