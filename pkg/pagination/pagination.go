package pagination

import "time"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ValidateAndNormalize(req *Request) {
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 10
	}

	if req.Order != "asc" && req.Order != "desc" {
		req.Order = "desc"
	}

	if req.SortBy == "" {
		req.SortBy = "created_at"
	}

}

// Determine navigation direction
func (s *Service) GetNavigationParams(req Request) (cursor string, effectiveOrder string) {
	cursor = req.NextPage
	effectiveOrder = req.Order

	if req.PrevPage != "" {
		cursor = req.PrevPage
		// Reverse order for backward navigation
		effectiveOrder = ReverseOrder(req.Order)
	}
	return cursor, effectiveOrder
}

func (s *Service) CalculateFetchLimit(limit int) int {
	return limit + 1
}

func (s *Service) BuildResponse(data []interface{}, req *Request, getCursorFields func(interface{}) (time.Time, int64)) Response {

	response := Response{
		Items:       data,
		Limit:       req.Limit,
		HasNext:     false,
		HasPrevious: false,
		NextPage:    "",
		PrevPage:    "",
	}
	// Handle empty results
	if len(data) == 0 {
		return response
	}

	hasData := len(data) > req.Limit
	responseData := data
	if hasData {
		// Trim the extra item used to check for more data
		responseData = data[:len(data)-1]
	}
	timestamp, id := getCursorFields(responseData[len(responseData)-1])

	// previous page
	if req.PrevPage != "" {
		firstTimestamp, firstID := getCursorFields(responseData[0])
		response.HasPrevious = hasData
		response.HasNext = true
		response.Items = ReverseSlice(responseData)
		response.NextPage = EncodeCursor(firstTimestamp, firstID)
		if hasData {
			response.PrevPage = EncodeCursor(timestamp, id)
		}
	} else {
		// next page
		response.HasNext = hasData
		response.Items = responseData
		if hasData {
			response.NextPage = EncodeCursor(timestamp, id)
		}
		if req.NextPage != "" {
			response.HasPrevious = true
			firstTimestamp, firstID := getCursorFields(responseData[0])
			response.PrevPage = EncodeCursor(firstTimestamp, firstID)
		}
	}

	return response
}
