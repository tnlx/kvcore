package kvcore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type API struct {
	Token string
}

type Paginator struct {
	PageSize       uint16
	OnPagedSuccess func(interface{})
	OnPagedFailure func(error)
}

// ContactsByTags lists all contacts associated with all the hashtags provided
func (api API) ContactsByTags(hashtags []string, paginator Paginator) {
	hfs := []string{}
	for _, v := range hashtags {
		hfs = append(hfs, fmt.Sprintf("filter[hashtags][]=%s", v))
	}
	hf := ""
	if len(hfs) > 0 {
		hf = hf + "filter[hashtags_and_or]=AND&"
		hf = hf + strings.Join(hfs, "&")
	}
	psz := paginator.PageSize
	if psz == 0 {
		psz = 1000
	}
	url := fmt.Sprintf("https://api.kvcore.com/v2/public/contacts?%s&limit=%d", hf, psz)

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

func (api API) get(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.Token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
