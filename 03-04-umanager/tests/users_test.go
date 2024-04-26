package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/EfimVelichkin/3rd_module_GO/03-04-umanager/internal/database"
)

func (s *IntegrationTestSuite) TestUserHandlers() {
	t := s.T()
	err := CreateSchema(s.conf.UsersService.Postgres.ConnectionURL())
	assert.NoError(t, err)

	var userID uuid.UUID

	t.Run("Create User", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		var client http.Client

		reqBody := `{"username": "pavel", "password": "test"}`
		req, err := http.NewRequest(http.MethodPost, mainURL+"users", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("List Users", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		var client http.Client

		req, err := http.NewRequest(http.MethodGet, mainURL+"users", nil)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NoError(t, err)
		defer resp.Body.Close()

		resBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var result struct {
			Users []database.User `json:"users"`
		}
		err = json.Unmarshal(resBody, &result)
		assert.NoError(t, err)
		assert.Equal(t, "pavel", result.Users[0].Username)
		userID = result.Users[0].ID
	})

	t.Run("Read User", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		var client http.Client

		req, err := http.NewRequest(http.MethodGet, mainURL+"users/"+userID.String(), nil)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		defer resp.Body.Close()

		resBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var user database.User
		err = json.Unmarshal(resBody, &user)

		assert.NoError(t, err)
		assert.Equal(t, "pavel", user.Username)
	})

	t.Run("Update User", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		var client http.Client

		reqBody := fmt.Sprintf(`{"id": "%s", "username": "admin"}`, userID.String())
		req, err := http.NewRequest(http.MethodPut, mainURL+"users/"+userID.String(), strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		defer resp.Body.Close()

		resBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.Equal(t, "", string(resBody))

		req, err = http.NewRequest(http.MethodGet, mainURL+"users/"+userID.String(), nil)
		assert.NoError(t, err)

		resp, err = client.Do(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		defer resp.Body.Close()

		resBody, err = io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var user database.User
		err = json.Unmarshal(resBody, &user)

		assert.NoError(t, err)
		assert.Equal(t, "admin", user.Username)
	})

	t.Run("Delete User", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		var client http.Client

		req, err := http.NewRequest(http.MethodDelete, mainURL+"users/"+userID.String(), nil)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.NoError(t, err)

		req, err = http.NewRequest(http.MethodGet, mainURL+"users/"+userID.String(), nil)
		assert.NoError(t, err)

		resp, err = client.Do(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.NoError(t, err)
	})

	t.Run("Read User Bad", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		var client http.Client

		req, err := http.NewRequest(http.MethodGet, mainURL+"users/bad-uuid-string", nil)
		req.Header.Set("Content-Type", "application/json")

		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Create User Bad", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		var client http.Client

		reqBody := `{"name": "pavel", "password": "test"}`
		req, err := http.NewRequest(http.MethodPost, mainURL+"users", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
