package integration

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"zadanie-6105/tests/api"
	"zadanie-6105/tests/client"
)

func TestGetTenders_Happy(t *testing.T) {
	_, c := client.New(t)
	tenders := api.GetTenders(c, t, 0, 0)
	for _, tender := range tenders {
		require.Equal(t, "Published", tender.Status)
	}
}

func TestGetTenders_BadRequest(t *testing.T) {
	_, c := client.New(t)

	testsTable := []struct {
		name   string
		limit  int32
		offset int32
	}{
		{
			name:  "limit <0",
			limit: -1,
		},
		{
			name:   "offset <0",
			offset: -1,
		},
		{
			name:   "limit <0 and offset <0",
			limit:  -1,
			offset: -1,
		},
	}

	for _, tt := range testsTable {
		t.Run(tt.name, func(t *testing.T) {
			GetTendersURL := client.BaseURL + fmt.Sprintf("/tenders?")
			if tt.limit != 0 {
				GetTendersURL += fmt.Sprintf("limit=%d", tt.limit)
			}
			if tt.offset != 0 {
				GetTendersURL += fmt.Sprintf("offset=%d", tt.offset)
			}
			req := client.FormRequest(http.MethodGet, GetTendersURL, nil)
			resp, err := c.Client.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}
