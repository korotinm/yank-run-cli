package api

// CreateReq is the body for POST /api/snippets.
type CreateReq struct {
	Body        string   `json:"body"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Author      string   `json:"author,omitempty"`
}

// CreateResp is the response from POST /api/snippets.
// New is derived from HTTP status: 201 -> true, 200 -> false (dedup).
type CreateResp struct {
	ID  string `json:"id"`
	URL string `json:"url"`
	New bool   `json:"new"`
}

// Snippet is the response from GET /api/snippets/:id.
type Snippet struct {
	ID          string   `json:"id"`
	Body        string   `json:"body"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags"`
	Author      string   `json:"author,omitempty"`
	CreatedAt   int64    `json:"createdAt"`
}

// Hit is one element of the GET /api/snippets?q=... search response.
type Hit struct {
	ID          string   `json:"id"`
	BodyPreview string   `json:"bodyPreview"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags"`
	Author      string   `json:"author,omitempty"`
	CreatedAt   int64    `json:"createdAt"`
}
