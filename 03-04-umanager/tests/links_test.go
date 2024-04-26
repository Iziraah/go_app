package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/EfimVelichkin/3rd_module_GO/03-04-umanager/internal/database"
)

func (s *IntegrationTestSuite) TestLinkHandlers() {
	t := s.T()

	var linkID primitive.ObjectID

	t.Run("Create Link", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		var client http.Client

		reqBody := `{
			"title": "main page",
			"url": "https://gb.ru/",
			"tags": [
				"edu"
			]
		}`
		req, err := http.NewRequest(http.MethodPost, mainURL+"links", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("List Links", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		var client http.Client

		req, err := http.NewRequest(http.MethodGet, mainURL+"links", nil)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NoError(t, err)
		defer resp.Body.Close()

		resBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var result struct {
			Links []database.Link `json:"links"`
		}
		err = json.Unmarshal(resBody, &result)
		assert.NoError(t, err)
		assert.Equal(t, "https://gb.ru/", result.Links[0].URL)
		linkID = result.Links[0].ID
	})

	t.Run("Read Link", func(t *testing.T) {
		var client http.Client

		req, err := http.NewRequest(http.MethodGet, mainURL+"links/"+linkID.Hex(), nil)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		defer resp.Body.Close()

		resBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var link database.Link
		err = json.Unmarshal(resBody, &link)

		assert.NoError(t, err)
		assert.Equal(t, "https://gb.ru/", link.URL)
	})

	t.Run("Update Link", func(t *testing.T) {
		var client http.Client

		reqBody := fmt.Sprintf(`{"id": "%s", "url": "https://ya.ru"}`, linkID.Hex())
		req, err := http.NewRequest(http.MethodPut, mainURL+"links/"+linkID.Hex(), strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.NoError(t, err)

		defer resp.Body.Close()

		resBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.Equal(t, "", string(resBody))

		req, err = http.NewRequest(http.MethodGet, mainURL+"links/"+linkID.Hex(), nil)
		assert.NoError(t, err)

		resp, err = client.Do(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		defer resp.Body.Close()

		resBody, err = io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var link database.Link
		err = json.Unmarshal(resBody, &link)

		assert.NoError(t, err)
		assert.Equal(t, "https://ya.ru", link.URL)
	})

	t.Run("Delete Link", func(t *testing.T) {
		var client http.Client

		req, err := http.NewRequest(http.MethodDelete, mainURL+"links/"+linkID.Hex(), nil)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.NoError(t, err)

		req, err = http.NewRequest(http.MethodGet, mainURL+"links/"+linkID.Hex(), nil)
		assert.NoError(t, err)

		resp, err = client.Do(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.NoError(t, err)
	})

	t.Run("Read Link Bad", func(t *testing.T) {
		var client http.Client

		req, err := http.NewRequest(http.MethodGet, mainURL+"links/bad-id-string", nil)
		req.Header.Set("Content-Type", "application/json")

		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
