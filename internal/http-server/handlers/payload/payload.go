package payload

import "github.com/arvinloc/url-shortener/internal/lib/api/response"

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}
