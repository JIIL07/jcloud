package serializer

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func SerializeToGOB(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, fmt.Errorf("error encoding with GOB: %w", err)
	}
	return buf.Bytes(), nil
}

func DeserializeFromGOB(data []byte, v interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(v)
	if err != nil {
		return fmt.Errorf("error decoding from GOB: %w", err)
	}
	return nil
}
