package odin

type response struct {
	Success    bool       `json:"success"`
	Data       []string   `json:"data"`
	Pagination pagination `json:"pagination"`
}

type pagination struct {
	Start []interface{} `json:"start"`
	Last  []interface{} `json:"last"`
	Limit int           `json:"limit"`
	Total int           `json:"total"`
}
