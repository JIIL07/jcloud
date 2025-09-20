package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Delta struct {
	Offset int    // Смещение, где начинается изменение
	Data   []byte // Новые данные
	Length int    // Количество байтов для удаления (если 0 — просто вставка)
}

func ApplyDelta(oldData []byte, deltas []Delta) ([]byte, error) {
	var buffer bytes.Buffer
	lastIndex := 0

	for _, delta := range deltas {
		if delta.Offset < lastIndex || delta.Offset > len(oldData) {
			return nil, fmt.Errorf("invalid delta offset: %d", delta.Offset)
		}

		buffer.Write(oldData[lastIndex:delta.Offset])
		buffer.Write(delta.Data)
		lastIndex = delta.Offset + delta.Length
	}

	buffer.Write(oldData[lastIndex:])

	return buffer.Bytes(), nil
}

func ApplyVersionToContent(fileContent []byte, version FileVersion) ([]byte, error) {
	if version.FullVersion {
		return version.Delta, nil
	}

	deltas, err := DeserializeDelta(version.Delta)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize delta for file_id %d, version %d: %w", version.FileID, version.Version, err)
	}

	updatedContent, err := ApplyDelta(fileContent, deltas)
	if err != nil {
		return nil, fmt.Errorf("failed to apply delta for file_id %d, version %d: %w", version.FileID, version.Version, err)
	}

	return updatedContent, nil
}

func DeserializeDelta(data []byte) ([]Delta, error) {
	var deltas []Delta
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&deltas); err != nil {
		return nil, fmt.Errorf("failed to deserialize delta: %w", err)
	}
	return deltas, nil
}

func SerializeDelta(deltas []Delta) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(deltas); err != nil {
		return nil, fmt.Errorf("failed to serialize delta: %w", err)
	}
	return buffer.Bytes(), nil
}
