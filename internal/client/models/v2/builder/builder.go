package builder

import "github.com/JIIL07/jcloud/internal/client/models"

type InfoBuilder struct {
	id       int
	metadata models.FileMetadata
	status   string
	data     []byte
}

func (b *InfoBuilder) WithID(id int) *InfoBuilder {
	b.id = id
	return b
}

func (b *InfoBuilder) WithMetadata(metadata models.FileMetadata) *InfoBuilder {
	b.metadata = metadata
	return b
}

func (b *InfoBuilder) WithStatus(status string) *InfoBuilder {
	b.status = status
	return b
}

func (b *InfoBuilder) WithData(data []byte) *InfoBuilder {
	b.data = data
	return b
}

func (b *InfoBuilder) Build() models.File {
	return models.File{
		ID:       b.id,
		Metadata: b.metadata,
		Status:   b.status,
		Data:     b.data,
	}
}
