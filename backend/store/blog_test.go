package store

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shoet/blog/clocker"
	"github.com/shoet/blog/models"
	"github.com/shoet/blog/options"
	"github.com/shoet/blog/testutil"
)

func Test_BlogRepository_Add(t *testing.T) {
	clocker := &clocker.FiexedClocker{}
	ctx := context.Background()
	db, err := testutil.NewDBMySQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	testutil.RepositoryTestPrepare(t, ctx, db)

	sut := NewBlogRepository(clocker)

	type args struct {
		blog *models.Blog
	}

	type want struct {
		blog *models.Blog
	}

	tests := []struct {
		id   string
		args args
		want want
	}{
		{
			id: "success",
			args: args{
				blog: &models.Blog{
					AuthorId:               1,
					Title:                  "title",
					Content:                "content",
					Description:            "description",
					ThumbnailImageFileName: "thumbnail_image_file_name",
					IsPublic:               true,
					Created:                clocker.Now(),
					Modified:               clocker.Now(),
				},
			},
			want: want{
				blog: &models.Blog{
					AuthorId:               1,
					Title:                  "title",
					Content:                "content",
					Description:            "description",
					ThumbnailImageFileName: "thumbnail_image_file_name",
					IsPublic:               true,
					Created:                clocker.Now(),
					Modified:               clocker.Now(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			tx := db.MustBegin()
			defer tx.Rollback()
			blogId, err := sut.Add(ctx, tx, tt.args.blog)
			if err != nil {
				t.Fatalf("failed to add blog: %v", err)
			}

			row := tx.QueryRowContext(ctx, "SELECT * FROM blogs WHERE id = ?", blogId)
			var got models.Blog
			if err := row.Scan(
				&got.Id, &got.AuthorId, &got.Title, &got.Content, &got.Description,
				&got.ThumbnailImageFileName, &got.IsPublic, &got.Created, &got.Modified,
			); err != nil {
				t.Fatalf("failed to scan row: %v", err)
			}

			opt := cmpopts.IgnoreFields(models.Blog{}, "Id")
			if diff := cmp.Diff(tt.want.blog, &got, opt); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func generateTestBlogs(t *testing.T, count int, now time.Time) []*models.Blog {
	t.Helper()
	blogs := make([]*models.Blog, count)
	for i := 0; i < count; i++ {
		b := &models.Blog{
			AuthorId:               1,
			Title:                  fmt.Sprintf("title%d", i),
			Content:                fmt.Sprintf("content%d", i),
			Description:            fmt.Sprintf("description%d", i),
			ThumbnailImageFileName: fmt.Sprintf("thumbnail_image_file_name%d", i),
			IsPublic:               true,
			Created:                now,
			Modified:               now,
		}
		blogs[i] = b
	}
	return blogs
}

func generateTestBlogsWithPublic(t *testing.T, count int, now time.Time) []*models.Blog {
	t.Helper()
	blogs := make([]*models.Blog, count)
	for i := 0; i < count; i++ {
		b := &models.Blog{
			AuthorId:               1,
			Title:                  fmt.Sprintf("title%d", i),
			Content:                fmt.Sprintf("content%d", i),
			Description:            fmt.Sprintf("description%d", i),
			ThumbnailImageFileName: fmt.Sprintf("thumbnail_image_file_name%d", i),
			IsPublic:               true,
			Created:                now,
			Modified:               now,
		}
		if i%2 == 0 {
			b.IsPublic = false
		}
		blogs[i] = b
	}
	return blogs
}

func Test_BlogRepository_List(t *testing.T) {
	clocker := &clocker.FiexedClocker{}
	ctx := context.Background()
	db, err := testutil.NewDBMySQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	testutil.RepositoryTestPrepare(t, ctx, db)

	sut := NewBlogRepository(clocker)

	type args struct {
		blogs    []*models.Blog
		limit    *int
		isPublic bool
	}

	type want struct {
		blogs []*models.Blog
		count int
	}

	tests := []struct {
		id   string
		args args
		want want
	}{
		{
			id: "success",
			args: args{
				blogs:    generateTestBlogs(t, 20, clocker.Now()),
				isPublic: true,
				limit:    func() *int { v := 20; return &v }(),
			},
			want: want{
				count: 20,
			},
		},
		{
			id: "limit 10",
			args: args{
				blogs: generateTestBlogs(t, 20, clocker.Now()),
				limit: func() *int { v := 10; return &v }(),
			},
			want: want{
				count: 10,
			},
		},
		{
			id: "isNotPublic",
			args: args{
				blogs:    generateTestBlogsWithPublic(t, 20, clocker.Now()),
				limit:    nil,
				isPublic: false,
			},
			want: want{
				count: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			tx := db.MustBegin()
			defer tx.Rollback()

			for _, b := range tt.args.blogs {
				prepareTask := `
				INSERT INTO blogs
					(
						author_id, title, content, description, 
						thumbnail_image_file_name, is_public, created, modified)
				VALUES
					(?, ?, ?, ?, ?, ?, ?, ?)
				`
				_, err := tx.ExecContext(
					ctx, prepareTask,
					b.AuthorId, b.Title, b.Content, b.Description,
					b.ThumbnailImageFileName, b.IsPublic, b.Created, b.Modified)
				if err != nil {
					t.Fatalf("failed to prepare task: %v", err)
				}

			}

			listOption := options.ListBlogOptions{}

			if tt.args.limit != nil {
				var l int64 = int64(*tt.args.limit)
				listOption.Limit = &l
			}
			listOption.IsPublic = tt.args.isPublic

			blogs, err := sut.List(ctx, tx, listOption)
			if err != nil {
				t.Fatalf("failed to list blogs: %v", err)
			}
			if tt.want.count != len(blogs) {
				t.Fatalf("failed to count blogs: %v", err)
			}

		})
	}
}