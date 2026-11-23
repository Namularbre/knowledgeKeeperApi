package domain

import "context"

// Repository persists and retrieves documents.
type Repository interface {
	Create(ctx context.Context, name string) (Document, error)
	FindByID(ctx context.Context, id uint64) (Document, error)
	SearchByName(ctx context.Context, name string) ([]Document, error)
	FindByUserID(ctx context.Context, userID uint64) ([]Document, error)
	AddUserDocument(ctx context.Context, documentID, userID uint64) ([]Document, error)
	RemoveUserDocument(ctx context.Context, documentID, userID uint64) ([]Document, error)
}
