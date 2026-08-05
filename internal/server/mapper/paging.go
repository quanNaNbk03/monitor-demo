package mapper

import (
	"github.com/quanNaNbk03/monitor-demo/internal/server/dto"
)

func ToPagingDTO(page, limit, total int) dto.Paging {
	var totalPage, from, to int
	if total%limit == 0 {
		totalPage = total / limit
	} else {
		totalPage = total/limit + 1
	}
	from = (page-1)*limit + 1
	to = page * limit
	if to > total {
		to = total
	}
	return dto.Paging{
		Page:      page,
		From:      from,
		To:        to,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
	}
}
