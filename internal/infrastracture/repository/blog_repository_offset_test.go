package repository_test

import (
	"context"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/infrastracture/repository"
	"github.com/shoet/blog/internal/options"
	"github.com/shoet/blog/internal/testutil"
)

func Test_BlogRepositoryOffset_List(t *testing.T) {
	ctx := context.Background()
	clocker := &clocker.FiexedClocker{}
	db, err := testutil.NewDBPostgreSQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	testutil.RepositoryTestPrepare(t, ctx, db)

	sut := repository.NewBlogRepositoryOffset(clocker)

	testdata := []*models.Blog{
		{Id: 1, AuthorId: 1, Title: "title1", Content: "content1", Description: "description1", ThumbnailImageFileName: "thumbnail1", IsPublic: false},
		{Id: 2, AuthorId: 1, Title: "title2", Content: "content2", Description: "description2", ThumbnailImageFileName: "thumbnail2", IsPublic: true},
		{Id: 3, AuthorId: 1, Title: "title3", Content: "content3", Description: "description3", ThumbnailImageFileName: "thumbnail3", IsPublic: true},
		{Id: 4, AuthorId: 1, Title: "title4", Content: "content4", Description: "description4", ThumbnailImageFileName: "thumbnail4", IsPublic: true},
		{Id: 5, AuthorId: 1, Title: "title5", Content: "content5", Description: "description5", ThumbnailImageFileName: "thumbnail5", IsPublic: true},
		{Id: 6, AuthorId: 1, Title: "title6", Content: "content6", Description: "description6", ThumbnailImageFileName: "thumbnail6", IsPublic: true},
		{Id: 7, AuthorId: 1, Title: "title7", Content: "content7", Description: "description7", ThumbnailImageFileName: "thumbnail7", IsPublic: true},
		{Id: 8, AuthorId: 1, Title: "title8", Content: "content8", Description: "description8", ThumbnailImageFileName: "thumbnail8", IsPublic: true},
		{Id: 9, AuthorId: 1, Title: "title9", Content: "content9", Description: "description9", ThumbnailImageFileName: "thumbnail9", IsPublic: true},
		{Id: 10, AuthorId: 1, Title: "title10", Content: "content10", Description: "description10", ThumbnailImageFileName: "thumbnail10", IsPublic: true},
		{Id: 11, AuthorId: 1, Title: "title11", Content: "content11", Description: "description11", ThumbnailImageFileName: "thumbnail11", IsPublic: true},
		{Id: 12, AuthorId: 1, Title: "title12", Content: "content12", Description: "description12", ThumbnailImageFileName: "thumbnail12", IsPublic: true},
		{Id: 13, AuthorId: 1, Title: "title13", Content: "content13", Description: "description13", ThumbnailImageFileName: "thumbnail13", IsPublic: true},
		{Id: 14, AuthorId: 1, Title: "title14", Content: "content14", Description: "description14", ThumbnailImageFileName: "thumbnail14", IsPublic: true},
		{Id: 15, AuthorId: 1, Title: "title15", Content: "content15", Description: "description15", ThumbnailImageFileName: "thumbnail15", IsPublic: true},
	}

	ptrInt64 := func(v int64) *int64 {
		return &v
	}

	type args struct {
		option *options.ListBlogOptions
	}
	type wants struct {
		blogs []*models.Blog
		err   error
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "単純なLimit",
			args: args{
				option: &options.ListBlogOptions{
					IsPublic: true,
					Limit:    5,
					Page:     ptrInt64(1),
				},
			},
			wants: wants{
				blogs: []*models.Blog{
					{Id: 15, AuthorId: 1, Title: "title15", Content: "content15", Description: "description15", ThumbnailImageFileName: "thumbnail15", IsPublic: true},
					{Id: 14, AuthorId: 1, Title: "title14", Content: "content14", Description: "description14", ThumbnailImageFileName: "thumbnail14", IsPublic: true},
					{Id: 13, AuthorId: 1, Title: "title13", Content: "content13", Description: "description13", ThumbnailImageFileName: "thumbnail13", IsPublic: true},
					{Id: 12, AuthorId: 1, Title: "title12", Content: "content12", Description: "description12", ThumbnailImageFileName: "thumbnail12", IsPublic: true},
					{Id: 11, AuthorId: 1, Title: "title11", Content: "content11", Description: "description11", ThumbnailImageFileName: "thumbnail11", IsPublic: true},
				},
				err: nil,
			},
		},
		{
			name: "ページの指定",
			args: args{
				option: &options.ListBlogOptions{
					IsPublic: true,
					Limit:    5,
					Page:     ptrInt64(2),
				},
			},
			wants: wants{
				blogs: []*models.Blog{
					{Id: 10, AuthorId: 1, Title: "title10", Content: "content10", Description: "description10", ThumbnailImageFileName: "thumbnail10", IsPublic: true},
					{Id: 9, AuthorId: 1, Title: "title9", Content: "content9", Description: "description9", ThumbnailImageFileName: "thumbnail9", IsPublic: true},
					{Id: 8, AuthorId: 1, Title: "title8", Content: "content8", Description: "description8", ThumbnailImageFileName: "thumbnail8", IsPublic: true},
					{Id: 7, AuthorId: 1, Title: "title7", Content: "content7", Description: "description7", ThumbnailImageFileName: "thumbnail7", IsPublic: true},
					{Id: 6, AuthorId: 1, Title: "title6", Content: "content6", Description: "description6", ThumbnailImageFileName: "thumbnail6", IsPublic: true},
				},
				err: nil,
			},
		},
		{
			name: "末尾",
			args: args{
				option: &options.ListBlogOptions{
					IsPublic: true,
					Limit:    6,
					Page:     ptrInt64(3),
				},
			},
			wants: wants{
				blogs: []*models.Blog{
					{Id: 3, AuthorId: 1, Title: "title3", Content: "content3", Description: "description3", ThumbnailImageFileName: "thumbnail3", IsPublic: true},
					{Id: 2, AuthorId: 1, Title: "title2", Content: "content2", Description: "description2", ThumbnailImageFileName: "thumbnail2", IsPublic: true},
				},
				err: nil,
			},
		},
		{
			name: "範囲外の指定",
			args: args{
				option: &options.ListBlogOptions{
					IsPublic: true,
					Limit:    5,
					Page:     ptrInt64(4),
				},
			},
			wants: wants{
				blogs: []*models.Blog{},
				err:   nil,
			},
		},
		{
			name: "is_public=false",
			args: args{
				option: &options.ListBlogOptions{
					Limit: 5,
					Page:  ptrInt64(3),
				},
			},
			wants: wants{
				blogs: []*models.Blog{
					{Id: 5, AuthorId: 1, Title: "title5", Content: "content5", Description: "description5", ThumbnailImageFileName: "thumbnail5", IsPublic: true},
					{Id: 4, AuthorId: 1, Title: "title4", Content: "content4", Description: "description4", ThumbnailImageFileName: "thumbnail4", IsPublic: true},
					{Id: 3, AuthorId: 1, Title: "title3", Content: "content3", Description: "description3", ThumbnailImageFileName: "thumbnail3", IsPublic: true},
					{Id: 2, AuthorId: 1, Title: "title2", Content: "content2", Description: "description2", ThumbnailImageFileName: "thumbnail2", IsPublic: true},
					{Id: 1, AuthorId: 1, Title: "title1", Content: "content1", Description: "description1", ThumbnailImageFileName: "thumbnail1", IsPublic: false},
				},
				err: nil,
			},
		},
	}

	testdataVals := []goqu.Record{}
	for _, d := range testdata {
		goquVal := goqu.Record{
			"id":                        d.Id,
			"author_id":                 d.AuthorId,
			"title":                     d.Title,
			"content":                   d.Content,
			"description":               d.Description,
			"thumbnail_image_file_name": d.ThumbnailImageFileName,
			"is_public":                 d.IsPublic,
		}
		testdataVals = append(testdataVals, goquVal)
	}

	sql, params, err := goqu.
		Insert("blogs").
		Rows(testdataVals).
		ToSQL()
	if err != nil {
		t.Fatalf("failed to build sql: %v", err)
	}
	_, err = db.ExecContext(ctx, sql, params...)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := sut.List(ctx, db, tt.args.option)
			if diff := cmp.Diff(tt.wants.err, err); diff != "" {
				t.Errorf("unexpected error: %v", diff)
			}

			options := cmp.Options{
				cmpopts.IgnoreFields(models.Blog{}, "Created", "Modified", "Tags", "Content"),
			}

			if diff := cmp.Diff(tt.wants.blogs, got, options); diff != "" {
				t.Errorf("unexpected blogs: %v", diff)
			}

		})
	}

	sql, params, err = goqu.
		Delete("blogs").
		Where(goqu.I("id").Gte(1)).
		Where(goqu.I("id").Lte(15)).
		ToSQL()
	if err != nil {
		t.Fatalf("failed to build sql: %v", err)
	}
	_, err = db.ExecContext(ctx, sql, params...)
	if err != nil {
		t.Fatalf("failed to delete test data: %v", err)
	}
}
