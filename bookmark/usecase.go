package bookmark

import (
	"context"
	"github.com/zhashkevych/go-clean-architecture/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UseCase interface {
	CreateBookmark(ctx context.Context, user *auth.User, url, title string) error
	GetBookmarks(ctx context.Context, user *auth.User) ([]*Bookmark, error)
	DeleteBookmark(ctx context.Context, user *auth.User, id primitive.ObjectID) error
}
