package repository

import (
	"context"

	"go.redsock.ru/moti/internal/models"
)

//go:generate minimock -i Repo -o ../../mocks -g -s "_mock.go"
type Repo interface {
	// ReadFile returns file's content from repository
	ReadFile(ctx context.Context, revision models.Revision, fileName string) (string, error)

	// Archive passed storage to archive and return full path to archive
	Archive(ctx context.Context, revision models.Revision, cacheDownloadPaths models.CacheDownloadPaths) error

	// ReadRevision reads commit's revision by passed version
	// or return the latest commit if version is empty
	ReadRevision(ctx context.Context, requestedVersion models.RequestedVersion) (models.Revision, error)

	// Fetch from remote repository specified version
	Fetch(ctx context.Context, revision models.Revision) error
}
