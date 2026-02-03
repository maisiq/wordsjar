package service

import (
	"context"
	"slices"

	"github.com/maisiq/go-words-jar/internal/models"
)

type Paginated[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor"`
	PrevCursor string `json:"prev_cursor,omitempty"`
	HasNext    bool   `json:"has_next"`
	HasPrev    bool   `json:"has_prev"`
}

type Pagination struct {
	Next    bool
	Pointer interface{}
}

func Paginate(
	ctx context.Context,
	params QueryParams,
	fn func(context.Context, QueryParams) ([]models.Word, error),
) (Paginated[models.Word], error) {
	var backwardDirection bool

	limit := params.Limit
	params.Limit += 1

	if params.Pagination != nil {
		if !params.Pagination.Next {
			backwardDirection = true
		}
	}

	items, err := fn(ctx, params)

	if err != nil {
		return Paginated[models.Word]{}, err
	}

	if backwardDirection {
		slices.SortFunc(items, func(a models.Word, b models.Word) int {
			if params.SortBy == FieldID {
				if params.Desc {
					if a.ID >= b.ID {
						return -1
					} else {
						return 1
					}
				} else {
					if a.ID >= b.ID {
						return 1
					} else {
						return -1
					}
				}
			} else { // if params.SortBy == service.FieldEN
				//TODO: add sort with another fields
				if params.Desc {
					if a.EN >= b.EN {
						return -1
					} else {
						return 1
					}
				} else {
					if a.EN >= b.EN {
						return 1
					} else {
						return -1
					}
				}
			}
		})
	}

	var hasNext bool
	var hasPrev bool

	extra := len(items) > int(limit)

	if extra {
		hasNext = true // ?
	}

	if params.Pagination != nil {
		if params.Pagination.Next {
			hasPrev = true
		} else {
			hasNext = true // ^
			if extra {
				hasPrev = true
			}
		}
	}

	pi := Paginated[models.Word]{
		Items: []models.Word{},
	}

	if hasNext {
		pi.HasNext = true
		if backwardDirection && len(items) <= int(limit) {
			if len(items) > 0 {
				pi.Items = items
				pi.NextCursor = pi.Items[len(pi.Items)-1].ID
			} else {
				pi.NextCursor = ""
			}
		} else {
			pi.Items = items[:limit]
			pi.NextCursor = pi.Items[len(pi.Items)-1].ID
		}

	} else {
		pi.Items = items
	}

	if hasPrev {
		pi.PrevCursor = pi.Items[0].ID
		pi.HasPrev = true
	}
	return pi, nil
}
