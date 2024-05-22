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
	Total       uint32    `json:"total"`
	LastPage    int       `json:"last_page"`
	LastPageUrl string    `json:"last_page_url"`
	//PerPage     string       `json:"per_page"` // it is sometimes string, sometimes int in api response
}

type Contact struct {
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Status  uint   `json:"status"`
	Private uint8  `json:"is_private"`
}

type ContactFilter struct {
	Status          *uint8  // 0 - New, 1 - Client, 2 - Closed, 3 - Sphere, 4 - Active, 5 - Contract, 7 - Prospect
	AssignedAgentId *uint64 // KvCore User ID of the assigned agent
	Hashtags        []string
	HashtagsAndOr   string // AND (default) or OR. Used if there are more than one Hashtags provided
}

// Traverses all pages of contacts that matched the filter criteria.
// Paginator defines the actions on each paginated result (or failure)
func (api API) ListContacts(filter ContactFilter, paginator Paginator) {
	// Filter criteria
	hf := []string{}
	if len(filter.Hashtags) > 1 {
		var c string
		if strings.ToUpper(filter.HashtagsAndOr) == "OR" {
			c = "OR"
		} else {
			c = "AND"
		}
		hf = append(hf, fmt.Sprintf("filter[hashtags_and_or]=%s", c))
	}
	if len(filter.Hashtags) > 0 {
		for _, v := range filter.Hashtags {
			hf = append(hf, fmt.Sprintf("filter[hashtags][]=%s", v))
		}
	}
	if filter.AssignedAgentId != nil {
		hf = append(hf, fmt.Sprintf("filter[assigned_agent_id]=%d", *filter.AssignedAgentId))
	}
	if filter.Status != nil {
		hf = append(hf, fmt.Sprintf("filter[status]=%d", *filter.Status))
	}
	hfs := strings.Join(hf, "&")

	// Pagination criteria
	psz := paginator.PageSize
	if psz == 0 {
		psz = 1000
	}

	// API Call
	url := fmt.Sprintf("%scontacts?%s&limit=%d", API_ROOT, hfs, psz)
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
