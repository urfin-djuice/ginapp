package types

type PaginationRequest struct {
	CurrentPage uint32 `json:"current_page" form:"current_page,default=1" binding:"omitempty"`
	PerPage     uint32 `json:"per_page" form:"per_page,default=15" binding:"omitempty"`
}

type PaginationResponse struct {
	PaginationRequest
	TotalPages   uint32 `json:"total_pages" binding:"omitempty"`
	TotalRecords uint32 `json:"total_records" binding:"omitempty"`
}

type IDReq struct {
	ID uint32 `json:"id"`
}
