package kvcore

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ContactsPaginated struct {
	CurrentPage uint16    `json:"current_page"`
	Data        []Contact `json:"data"`
	NextPageUrl string    `json:"next_page_url"`
	PerPage     string    `json:"per_page"`
	Total       uint32    `json:"total"`
	LastPage    int       `json:"last_page"`
	LastPageUrl string    `json:"last_page_url"`
}

type Contact struct {
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Status  uint   `json:"status"`
	Private uint8  `json:"is_private"`
}

type ContactFilter struct {
	Hashtags      []string
	HashtagsAndOr string // AND (default) or OR. Used if there are more than one Hashtags provided
}

// Traverses all pages of contacts that matched the filter criteria.
// Paginator defines the actions on each paginated result (or failure)
func (api API) ListContacts(filter ContactFilter, paginator Paginator) {
	// Filter criteria
	hf := ""
	if len(filter.Hashtags) > 1 {
		var c string
		if strings.ToUpper(filter.HashtagsAndOr) == "OR" {
			c = "OR"
		} else {
			c = "AND"
		}
		hf = hf + fmt.Sprintf("filter[hashtags_and_or]=%s&", c)
	}
	if len(filter.Hashtags) > 0 {
		hfs := []string{}
		for _, v := range filter.Hashtags {
			hfs = append(hfs, fmt.Sprintf("filter[hashtags][]=%s", v))
		}
		hf = hf + strings.Join(hfs, "&")
	}

	// Pagination criteria
	psz := paginator.PageSize
	if psz == 0 {
		psz = 1000
	}

	// API Call
	url := fmt.Sprintf("%scontacts?%s&limit=%d", API_ROOT, hf, psz)
	var contactsPaginated ContactsPaginated
	page := 1
	for url != "" {
		rsp, err := api.get(url)
		if err != nil {
			paginator.OnPagedFailure(err)
		}
		fmt.Printf("Fetching page %d: %s\n", page, url)
		page++
		err = json.Unmarshal(rsp, &contactsPaginated)
		if err != nil {
			paginator.OnPagedFailure(err)
		}
		paginator.OnPagedSuccess(contactsPaginated.Data)
		if contactsPaginated.NextPageUrl == url {
			break
		}
		url = contactsPaginated.NextPageUrl
	}
}
