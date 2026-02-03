//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/maisiq/go-words-jar/internal/config"
	dbx "github.com/maisiq/go-words-jar/internal/db"
	"github.com/maisiq/go-words-jar/internal/models"
	"github.com/maisiq/go-words-jar/internal/repository"
	"github.com/maisiq/go-words-jar/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestPagination(t *testing.T) {
	var cursor string

	ctx := context.Background()
	client := dbx.New(config.Database{
		DSN: PGDSN,
	})
	repo := repository.New(client)

	cases := []struct {
		name          string
		params        service.QueryParams
		expectedError error
		queryForRows  string
		setNextCursor bool
		setPrevCursor bool
	}{
		// sort by en field ASC
		{"first page sort by EN ASC", service.QueryParams{Limit: 5, SortBy: "en", Desc: false}, nil, "SELECT * FROM words ORDER BY en LIMIT 5", true, false},
		{"second page sort by EN ASC", service.QueryParams{
			Limit: 5, SortBy: "en", Desc: false, Pagination: &service.Pagination{Next: true},
		}, nil, "SELECT * FROM words ORDER BY en OFFSET 5 LIMIT 5", false, true},
		{"back from second page with sort by EN ASC", service.QueryParams{
			Limit: 5, SortBy: "en", Desc: false, Pagination: &service.Pagination{Next: false},
		}, nil, "SELECT * FROM words ORDER BY en LIMIT 5", false, false},
		// sort by en field DESC
		{"first page sort by EN DESC", service.QueryParams{Limit: 5, SortBy: "en", Desc: true}, nil, "SELECT * FROM words ORDER BY en DESC LIMIT 5", true, false},
		{"second page sort by EN DESC", service.QueryParams{
			Limit: 5, SortBy: "en", Desc: true, Pagination: &service.Pagination{Next: true},
		}, nil, "SELECT * FROM words ORDER BY en DESC OFFSET 5 LIMIT 5", false, true},
		{"back from second page with sort by EN DESC", service.QueryParams{
			Limit: 5, SortBy: "en", Desc: true, Pagination: &service.Pagination{Next: false},
		}, nil, "SELECT * FROM words ORDER BY en DESC LIMIT 5", false, false},
		// sort by id field ASC
		{"first page sort by id ASC", service.QueryParams{Limit: 5, SortBy: "id", Desc: false}, nil, "SELECT * FROM words ORDER BY id LIMIT 5", true, false},
		{"second page sort by id ASC", service.QueryParams{
			Limit: 5, SortBy: "id", Desc: false, Pagination: &service.Pagination{Next: true},
		}, nil, "SELECT * FROM words ORDER BY id OFFSET 5 LIMIT 5", false, true},
		{"back from second page with sort by id ASC", service.QueryParams{
			Limit: 5, SortBy: "id", Desc: false, Pagination: &service.Pagination{Next: false},
		}, nil, "SELECT * FROM words ORDER BY id LIMIT 5", false, false},
		// sort by id field DESC
		{"first page sort by id DESC", service.QueryParams{Limit: 5, SortBy: "id", Desc: true}, nil, "SELECT * FROM words ORDER BY id DESC LIMIT 5", true, false},
		{"second page sort by id DESC", service.QueryParams{
			Limit: 5, SortBy: "id", Desc: true, Pagination: &service.Pagination{Next: true},
		}, nil, "SELECT * FROM words ORDER BY id DESC OFFSET 5 LIMIT 5", false, true},
		{"back from second page with sort by id DESC", service.QueryParams{
			Limit: 5, SortBy: "id", Desc: true, Pagination: &service.Pagination{Next: false},
		}, nil, "SELECT * FROM words ORDER BY id DESC LIMIT 5", false, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.params.Pagination != nil {
				c.params.Pagination.Pointer = cursor
			}

			pw, err := service.Paginate(ctx, c.params, func(ctx context.Context, qp service.QueryParams) ([]models.Word, error) {
				return repo.WordList(ctx, &qp)
			})
			assert.ErrorIs(t, c.expectedError, err)

			rows, _ := client.DB.QueryxContext(ctx, c.queryForRows)

			wordObjs := []models.Word{}
			for rows.Next() {
				var w models.Word
				_ = rows.StructScan(&w)
				wordObjs = append(wordObjs, w)
			}

			assert.Equal(t, len(wordObjs), len(pw.Items))

			for i, word := range wordObjs {
				assert.Equal(t, word.EN, pw.Items[i].EN)
			}

			if c.setNextCursor {
				cursor = pw.NextCursor
			}
			if c.setPrevCursor {
				cursor = pw.PrevCursor
			}
		})
	}

}
