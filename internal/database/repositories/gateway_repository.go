package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/antelman107/filestorage/pkg/domain"
)

type gatewayRepo struct {
	db *sqlx.DB
}

func NewGatewayRepo(db *sqlx.DB) domain.GatewayRepository {
	return &gatewayRepo{
		db: db,
	}
}

const lockServersStatement = `
SELECT id 
FROM servers
FOR UPDATE
;
`

func (r *gatewayRepo) LockServers(ctx context.Context) error {
	if _, err := getDBExecutorFromCtx(ctx, r.db).ExecContext(ctx, lockServersStatement); err != nil {
		return fmt.Errorf("failed to lock servers: %w", err)
	}

	return nil
}

const insertServerStatement = `
	INSERT INTO servers (id, url) 
	VALUES (:id, :url)
	;
`

func (r *gatewayRepo) StoreServer(ctx context.Context, file domain.Server) error {
	if _, err := getDBExecutorFromCtx(ctx, r.db).NamedExecContext(ctx, insertServerStatement, file); err != nil {
		return fmt.Errorf("failed to insert server: %w", err)
	}

	return nil
}

const deleteFileStatement = `
	DELETE FROM files 
	WHERE id = $1
	RETURNING id, name, size, is_uploaded
	;
`

func (r *gatewayRepo) DeleteFile(ctx context.Context, id string) (domain.File, error) {
	res := domain.File{}
	if err := getDBExecutorFromCtx(ctx, r.db).GetContext(ctx, &res, deleteFileStatement, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.File{}, domain.ErrNotFound
		}

		return domain.File{}, fmt.Errorf("failed to delete file %s: %w", id, err)
	}

	return res, nil
}

const selectServersUsageStatement = `
SELECT usage, 
       serverID, 
       serverURL 
FROM (
	SELECT
    DISTINCT ON (servers.id)

    COALESCE(SUM(chunks.size), 0) AS usage,
    servers.id AS serverID,
    servers.url AS serverURL,
	servers.created_at
	FROM servers
	LEFT JOIN chunks
	ON servers.id = chunks.server_id
	GROUP BY servers.id
	ORDER BY servers.id ASC    
) AS usages_unordered
ORDER BY usage ASC, created_at ASC
;
`

func (r *gatewayRepo) GetServersUsages(ctx context.Context) (domain.ServersSummary, error) {
	usages := make(domain.ServersUsages, 0)

	if err := getDBExecutorFromCtx(ctx, r.db).SelectContext(ctx, &usages, selectServersUsageStatement); err != nil {
		return domain.ServersSummary{}, fmt.Errorf("failed to select server usages: %w", err)
	}

	total := int64(0)
	for _, usage := range usages {
		total += usage.Usage
	}

	return domain.ServersSummary{
		Usages:     usages,
		TotalUsage: total,
	}, nil
}

const selectChunksByFileIDStatement = `
SELECT
    chunks.id, 
    file_id, 
    server_id, 
    index, 
    size, 
    is_uploaded, 
    servers.url AS server_url,
    hash
FROM chunks
LEFT JOIN servers
ON servers.id = chunks.server_id
WHERE chunks.file_id = $1
;
`

func (r *gatewayRepo) GetChunksByFileID(ctx context.Context, fileID string) (domain.Chunks, error) {
	chunks := make(domain.Chunks, 0)

	if err := getDBExecutorFromCtx(ctx, r.db).SelectContext(ctx, &chunks, selectChunksByFileIDStatement, fileID); err != nil {
		return nil, fmt.Errorf("failed to select chunks: %w", err)
	}

	return chunks, nil
}

const insertFileStatement = `
	INSERT INTO files (id, name, size) 
	VALUES (:id, :name, :size)
	;
`

const lockFileStatement = `
SELECT id 
FROM files
WHERE id = $1
FOR UPDATE
;
`

func (r *gatewayRepo) LockFile(ctx context.Context, id string) error {
	if _, err := getDBExecutorFromCtx(ctx, r.db).ExecContext(ctx, lockFileStatement, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}

		return fmt.Errorf("failed to lock file: %w", err)
	}

	return nil
}

func (r *gatewayRepo) StoreFile(ctx context.Context, file domain.File) error {
	if _, err := getDBExecutorFromCtx(ctx, r.db).NamedExecContext(ctx, insertFileStatement, file); err != nil {
		return fmt.Errorf("failed to insert file: %w", err)
	}

	return nil
}

const updateFileIsUploadedStatement = `
	UPDATE files 
	SET is_uploaded = TRUE
	WHERE id = $1
	;
`

func (r *gatewayRepo) UpdateFileIsUploaded(ctx context.Context, id string) error {
	if _, err := getDBExecutorFromCtx(ctx, r.db).ExecContext(ctx, updateFileIsUploadedStatement, id); err != nil {
		return fmt.Errorf("failed to update is_uploaded file %s: %w", id, err)
	}

	return nil
}

const selectUploadedFileStatement = `
SELECT
    id,
    name,
    size,
    is_uploaded
FROM files
WHERE id = $1 
  AND is_uploaded 
;
`

func (r *gatewayRepo) GetUploadedFile(ctx context.Context, fileID string) (domain.File, error) {
	file := domain.File{}

	if err := getDBExecutorFromCtx(ctx, r.db).GetContext(ctx, &file, selectUploadedFileStatement, fileID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.File{}, domain.ErrNotFound
		}

		return domain.File{}, fmt.Errorf("failed to select file: %w", err)
	}

	return file, nil
}

const insertChunksStatement = `
	INSERT INTO chunks (id, file_id, server_id, index, size) 
	VALUES (:id, :file_id, :server_id, :index, :size)
`

func (r *gatewayRepo) StoreChunks(ctx context.Context, chunks domain.Chunks) error {
	if _, err := getDBExecutorFromCtx(ctx, r.db).NamedExecContext(ctx, insertChunksStatement, chunks); err != nil {
		return fmt.Errorf("failed to store chunks: %w", err)
	}

	return nil
}

const deleteChunkStatement = `
	DELETE FROM chunks 
	WHERE id = $1
	;
`

func (r *gatewayRepo) DeleteChunk(ctx context.Context, id string) error {
	if _, err := getDBExecutorFromCtx(ctx, r.db).ExecContext(ctx, deleteChunkStatement, id); err != nil {
		return fmt.Errorf("failed to delete chunk %s: %w", id, err)
	}

	return nil
}

const updateChunkIsUploadedStatement = `
	UPDATE chunks 
	SET is_uploaded = TRUE, hash = $2
	WHERE id = $1
	;
`

func (r *gatewayRepo) UpdateUploadedChunk(ctx context.Context, id, hash string) error {
	if _, err := getDBExecutorFromCtx(ctx, r.db).ExecContext(ctx, updateChunkIsUploadedStatement, id, hash); err != nil {
		return fmt.Errorf("failed to update is_uploaded chunk %s: %w", id, err)
	}

	return nil
}

func (r *gatewayRepo) InitTransaction() (*sqlx.Tx, error) {
	return r.db.Beginx()
}
