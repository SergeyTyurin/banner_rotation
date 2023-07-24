package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SergeyTyurin/banner_rotation/internal/configs"
	"github.com/SergeyTyurin/banner_rotation/internal/database"
	"github.com/SergeyTyurin/banner_rotation/structures"
	"github.com/stretchr/testify/require"
)

func TestBadRequestBanner(t *testing.T) {
	d := database.NewDatabase()
	config, _ := configs.GetDBConnectionConfig("../../config/connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection()

	h := Handlers{d, nil}
	url := fmt.Sprintf("http://%s:%d/%s", config.Host(), config.Port(), "/banner")

	t.Run("create", func(t *testing.T) {
		jsonBody := []byte("incorrect")
		request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
		response := httptest.NewRecorder()

		h.CreateBanner(response, request)
		require.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("update", func(t *testing.T) {
		jsonBody := []byte("incorrect")
		request, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(jsonBody))
		response := httptest.NewRecorder()

		h.UpdateBanner(response, request)
		require.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("delete", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, url, nil)
		response := httptest.NewRecorder()

		q := request.URL.Query()
		q.Add("id", "bad")
		request.URL.RawQuery = q.Encode()
		h.DeleteBanner(response, request)

		require.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("get", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, url, nil)
		response := httptest.NewRecorder()

		q := request.URL.Query()
		q.Add("id", "bad")
		request.URL.RawQuery = q.Encode()
		h.GetBanner(response, request)

		require.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestNonExistsBanner(t *testing.T) {
	d := database.NewDatabase()
	config, _ := configs.GetDBConnectionConfig("../../config/connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection()

	h := Handlers{d, nil}
	url := fmt.Sprintf("http://%s:%d/%s", config.Host(), config.Port(), "/banner")

	t.Run("update", func(t *testing.T) {
		entity := structures.Banner{Id: -1, Info: "entity"}
		jsonBody, _ := json.Marshal(entity)
		request, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(jsonBody))
		response := httptest.NewRecorder()

		h.UpdateBanner(response, request)
		require.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("delete", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, url, nil)
		response := httptest.NewRecorder()

		q := request.URL.Query()
		q.Add("id", "-1")
		request.URL.RawQuery = q.Encode()
		h.DeleteBanner(response, request)

		require.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("get", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, url, nil)
		response := httptest.NewRecorder()

		q := request.URL.Query()
		q.Add("id", "-1")
		request.URL.RawQuery = q.Encode()
		h.GetBanner(response, request)

		require.Equal(t, http.StatusNotFound, response.Code)
	})
}
