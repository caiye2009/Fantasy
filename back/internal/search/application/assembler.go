package application

import (
	"back/internal/search/domain"
)

// ToSearchQuery DTO → Domain Model
func ToSearchQuery(req *SearchRequest) *domain.SearchQuery {
	sortFields := make([]domain.SortField, len(req.Sort))
	for i, s := range req.Sort {
		sortFields[i] = domain.SortField{
			Field: s.Field,
			Order: s.Order,
		}
	}

	// 处理默认值
	from := 0
	if req.From != nil {
		from = *req.From
	}

	size := 10 // 默认每页 10 条
	if req.Size != nil {
		size = *req.Size
	}

	return &domain.SearchQuery{
		Keyword: req.Query,
		Indices: req.Indices,
		Fields:  req.Fields,
		Filters: req.Filters,
		Sort:    sortFields,
		From:    from,
		Size:    size,
	}
}

// ToSearchResponse Domain Model → DTO
func ToSearchResponse(resp *domain.SearchResponse) *SearchResponse {
	results := make([]*SearchResultDTO, len(resp.Results))
	for i, r := range resp.Results {
		results[i] = &SearchResultDTO{
			Index:     r.Index,
			Type:      r.Type,
			ID:        r.ID,
			Score:     r.Score,
			Source:    r.Source,
			Highlight: r.Highlight,
		}
	}
	
	return &SearchResponse{
		Total:    resp.Total,
		Took:     resp.Took,
		MaxScore: resp.MaxScore,
		Results:  results,
	}
}

// ToIndexListResponse 索引配置 → DTO
func ToIndexListResponse() *IndexListResponse {
	indices := domain.GetAllIndices()
	infos := make([]*IndexInfoDTO, 0, len(indices))
	
	for _, idx := range indices {
		if meta, ok := domain.GetIndexMeta(idx); ok {
			infos = append(infos, &IndexInfoDTO{
				Name:   meta.Name,
				Type:   meta.Type,
				Fields: meta.DefaultFields,
			})
		}
	}
	
	return &IndexListResponse{
		Indices: infos,
	}
}