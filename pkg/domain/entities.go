package domain

import (
	"github.com/google/uuid"
)

type ServerUsage struct {
	ServerID  uuid.UUID
	ServerURL string
	Usage     int64
}

type ServersUsages []*ServerUsage

type ServersSummary struct {
	Usages     ServersUsages
	TotalUsage int64
}

type File struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Size       int64     `db:"size"`
	IsUploaded bool      `db:"is_uploaded"`
}

type Chunk struct {
	ID         uuid.UUID `db:"id"`
	FileID     uuid.UUID `db:"file_id"`
	ServerID   uuid.UUID `db:"server_id"`
	ServerURL  string    `db:"server_url"`
	Index      int       `db:"index"`
	Size       int64     `db:"size"`
	IsUploaded bool      `db:"is_uploaded"`
	Hash       string    `db:"hash"`
}

type ChunkWithData struct {
	Chunk
	Data []byte
}

func (c Chunk) GetWithData(data []byte) ChunkWithData {
	return ChunkWithData{
		Chunk: c,
		Data:  data,
	}
}

type Chunks []Chunk

type Server struct {
	ID  uuid.UUID `db:"id"`
	URL string    `db:"url"`
}
