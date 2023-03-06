package pagination

import (
	"log"
	"math"
	"sync"

	"gorm.io/gorm"
)

type PagingResult struct {
	HasNextPage     bool  `json:"hasNextPage"`
	HasPreviousPage bool  `json:"hasPreviousPage"`
	Page            int   `json:"page"`
	PageSize        int   `json:"pageSize"`
	TotalPages      int   `json:"totalPages"`
	TotalCount      int64 `json:"totalCount"`
	Items           any   `json:"items"`
}

func NewPagingResult(q *gorm.DB, data any, page, pageSize int) (*PagingResult, error) {
	var (
		totalCount int64
		offset     int = (page - 1) * pageSize
		wg         sync.WaitGroup
	)
	wg.Add(2)
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 1
	}
	// Session mode
	tx := q.Session(&gorm.Session{PrepareStmt: true})

	// Run the count query
	go func() {
		defer wg.Done()
		if err := tx.Count(&totalCount).Error; err != nil {
			log.Println("error count data:", err)
		}

	}()

	// Fetch the data
	go func() {
		defer wg.Done()
		if err := tx.Limit(pageSize).Offset(offset).Find(data).Error; err != nil {
			log.Println("error fetching data:", err)
		}
	}()

	// Wait for the routines
	wg.Wait()

	return &PagingResult{
		HasNextPage:     int(totalCount)-offset > pageSize,
		HasPreviousPage: page > 1,
		Page:            page,
		PageSize:        pageSize,
		TotalPages:      int(math.Ceil(float64(totalCount) / float64(pageSize))),
		TotalCount:      totalCount,
		Items:           data,
	}, nil
}
