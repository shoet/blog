package services

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/infrastracture/repository"
	"github.com/shoet/blog/internal/options"
	"golang.org/x/exp/slices"
)

func NewBlogService(db *sqlx.DB, blog BlogRepository) *BlogService {
	return &BlogService{
		db:   db,
		blog: blog,
	}
}

type BlogService struct {
	db   *sqlx.DB
	blog BlogRepository
}

func (b *BlogService) DeleteBlog(ctx context.Context, id models.BlogId) error {
	tx, err := b.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// delete blogs_tags -----------------
	// select using other blog tags
	var usingTags models.BlogsTagsArray
	usingTags, err = b.blog.SelectBlogsTagsByOtherUsingBlog(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("failed to select using tags: %w", err)
	}

	//  select will delete tags
	blogsTags, err := b.blog.SelectBlogsTags(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("failed to select blogs_tags: %w", err)
	}
	var willDeleteTags []models.TagId
	for _, t := range blogsTags {
		if !slices.Contains(usingTags.TagIds(), t.TagId) {
			willDeleteTags = append(willDeleteTags, t.TagId)
		}
	}

	for _, tagId := range willDeleteTags {
		// delete tags
		if err := b.blog.DeleteTag(ctx, tx, tagId); err != nil {
			return fmt.Errorf("failed to delete tags: %w", err)
		}
		// delete blogs_tags
		if err := b.blog.DeleteBlogsTags(ctx, tx, id, tagId); err != nil {
			return fmt.Errorf("failed to delete blogs_tags: %w", err)
		}
	}

	// delete blogs ----------------------
	err = b.blog.Delete(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("failed to delete blog: %w", err)
	}

	// commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (b *BlogService) SelectTag(
	ctx context.Context, db repository.Execer, tag string,
) (*models.Tag, error) {
	tags, err := b.blog.SelectTags(ctx, db, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to select tag: %w", err)
	}
	if len(tags) == 0 {
		return nil, nil
	}
	return tags[0], nil
}

func (b *BlogService) PutBlog(ctx context.Context, blog *models.Blog) (*models.Blog, error) {
	tx, err := b.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var usingTagsByOtherBlog models.BlogsTagsArray
	if usingTagsByOtherBlog, err = b.blog.SelectBlogsTagsByOtherUsingBlog(ctx, tx, blog.Id); err != nil {
		return nil, fmt.Errorf("failed to select using tags: %w", err)
	}
	isUsing := func(tag string) bool { return slices.Contains(usingTagsByOtherBlog.TagNames(), tag) }

	var currentTags models.BlogsTagsArray
	if currentTags, err = b.blog.SelectBlogsTags(ctx, tx, blog.Id); err != nil {
		return nil, fmt.Errorf("failed to select current tags: %w", err)
	}
	isCurrent := func(tag string) bool { return slices.Contains(currentTags.TagNames(), tag) }

	isContainsNew := func(tag string) bool { return slices.Contains(blog.Tags, tag) }

	for _, tag := range blog.Tags {
		if !isCurrent(tag) {
			tags, err := b.blog.SelectTags(ctx, tx, tag)
			if err != nil {
				return nil, fmt.Errorf("failed to select tag: %w", err)
			}
			// add tags
			var tagId models.TagId
			if len(tags) == 0 {
				tagId, err = b.blog.AddTag(ctx, tx, tag)
				if err != nil {
					return nil, fmt.Errorf("failed to add tag: %w", err)
				}
			} else {
				tagId = tags[0].Id
			}
			// add blogs_tags
			if _, err := b.blog.AddBlogTag(ctx, tx, blog.Id, tagId); err != nil {
				return nil, fmt.Errorf("failed to add blogs_tags: %w", err)
			}
		}
	}

	for _, tag := range currentTags {
		if isContainsNew(tag.Name) {
			continue
		}
		if !isUsing(tag.Name) {
			// delete tags
			if err := b.blog.DeleteTag(ctx, tx, tag.TagId); err != nil {
				return nil, fmt.Errorf("failed to delete tags: %w", err)
			}
			if err := b.blog.DeleteBlogsTags(ctx, tx, blog.Id, tag.TagId); err != nil {
				return nil, fmt.Errorf("failed to delete blogs_tags: %w", err)
			}
		}
	}

	// put blog
	id, err := b.blog.Put(ctx, tx, blog)
	if err != nil {
		return nil, fmt.Errorf("failed to put blog: %w", err)
	}

	newBlog, err := b.blog.Get(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get blog: %w", err)
	}

	// commit
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return newBlog, nil
}

func (s *BlogService) ListTags(ctx context.Context, option options.ListTagsOptions) ([]*models.Tag, error) {
	tags, err := s.blog.ListTags(ctx, s.db, option)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	return tags, nil
}

func (s *BlogService) Export(ctx context.Context) error {
	return nil
}

func (s *BlogService) Validate(ctx context.Context, userId models.UserId, blog *models.Blog) error {
	if userId != blog.AuthorId {
		return fmt.Errorf("blog.AuthorId is invalid")
	}
	return nil
}
