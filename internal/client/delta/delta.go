package delta

import (
	"bytes"
	"crypto/sha1"
	"fmt"
)

type Snapshot struct {
	Content []byte
	Hash    string
}

type Delta struct {
	OriginalHash string
	NewHash      string
	Content      []byte
}

func NewSnapshot(content []byte) *Snapshot {
	hash := fmt.Sprintf("%x", sha1.Sum(content))
	return &Snapshot{Content: content, Hash: hash}
}

func (s *Snapshot) CreateDelta(previous *Snapshot) *Delta {
	if bytes.Equal(s.Content, previous.Content) {
		return nil
	}
	return &Delta{
		OriginalHash: previous.Hash,
		NewHash:      s.Hash,
		Content:      s.Content,
	}
}
