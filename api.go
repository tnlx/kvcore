package kvcore

import (
	"fmt"
	"io"
	"net/http"
)

const (
	API_ROOT string = "https://api.kvcore.com/v2/public/"
)

// Init the API client using the token obtained from the KvCore/ Lead Dropbox.
type API struct {
	Token string
}

type Paginator struct {
	PageSize       uint16            // Default: 1000
	OnPagedSuccess func(interface{}) // Consumer of every 'PageSize' collection fetched
	OnPagedFailure func(error)       // Action on every failed iteration
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
