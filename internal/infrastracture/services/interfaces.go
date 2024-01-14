package services

import (
	"context"

	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/infrastracture/repository"
	"github.com/shoet/blog/internal/options"
)

//go:generate go run github.com/matryer/moq -out service_moq.go . JWTer KVSer

type BlogRepository interface {
	Add(ctx context.Context, tx infrastracture.TX, blog *models.Blog) (models.BlogId, error)
	List(ctx context.Context, tx infrastracture.TX, option options.ListBlogOptions) ([]*models.Blog, error)
	Get(ctx context.Context, tx infrastracture.TX, id models.BlogId) (*models.Blog, error)
	Delete(ctx context.Context, tx infrastracture.TX, id models.BlogId) error
	Put(ctx context.Context, tx infrastracture.TX, blog *models.Blog) (models.BlogId, error)

	AddBlogTag(ctx context.Context, tx infrastracture.TX, blogId models.BlogId, tagId models.TagId) (int64, error)
	SelectBlogsTagsByOtherUsingBlog(ctx context.Context, tx infrastracture.TX, blogId models.BlogId) ([]*models.BlogsTags, error)
	SelectBlogsTags(ctx context.Context, tx infrastracture.TX, blogId models.BlogId) ([]*models.BlogsTags, error)
	DeleteBlogsTags(ctx context.Context, tx infrastracture.TX, blogId models.BlogId, tagId models.TagId) error

	SelectTags(ctx context.Context, tx infrastracture.TX, tag string) ([]*models.Tag, error)
	AddTag(ctx context.Context, tx infrastracture.TX, tag string) (models.TagId, error)
	DeleteTag(ctx context.Context, tx infrastracture.TX, tagId models.TagId) error
	ListTags(ctx context.Context, db repository.Queryer, option options.ListTagsOptions) ([]*models.Tag, error)
}

type UserRepository interface {
	Add(ctx context.Context, db repository.Execer, user *models.User) (*models.User, error)
	Get(ctx context.Context, db repository.Queryer, id models.UserId) (*models.User, error)
	GetByEmail(ctx context.Context, db repository.Queryer, email string) (*models.User, error)
	// Delete(ctx context.Context, db store.Execer, id models.UserId) error
	// Put(ctx context.Context, db store.Execer, user *models.User) error
}

type JWTer interface {
	GenerateToken(ctx context.Context, u *models.User) (string, error)
	VerifyToken(ctx context.Context, token string) (models.UserId, error)
}

type KVSer interface {
	Save(ctx context.Context, key string, value string) error
	Load(ctx context.Context, key string) (string, error)
}
