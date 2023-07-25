package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SergeyTyurin/banner-rotation/configs"
	"github.com/SergeyTyurin/banner-rotation/database"
	"github.com/stretchr/testify/require"
)

func TestBadRequestRotaton(t *testing.T) {
	d := database.NewDatabase()
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()

	h := Handlers{d, nil}
	url := fmt.Sprintf("http://%s:%d/%s", config.Host(), config.Port(), "/rotation")

	t.Run("add", func(t *testing.T) {
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, url, nil)
		q := request.URL.Query()
		q.Add("banner_id", "q")
		q.Add("slot_id", "w")
		request.URL.RawQuery = q.Encode()
		response := httptest.NewRecorder()

		h.HandlerAddToRotation(response, request)
		require.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("delete", func(t *testing.T) {
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, url, nil)
		q := request.URL.Query()
		q.Add("banner_id", "q")
		q.Add("slot_id", "w")
		request.URL.RawQuery = q.Encode()
		response := httptest.NewRecorder()

		h.DeleteFromRotation(response, request)
		require.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("select", func(t *testing.T) {
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
		q := request.URL.Query()
		q.Add("group_id", "q")
		q.Add("slot_id", "w")
		request.URL.RawQuery = q.Encode()
		response := httptest.NewRecorder()

		h.SelectFromRotation(response, request)
		require.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("register", func(t *testing.T) {
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, url, nil)
		q := request.URL.Query()
		q.Add("banner_id", "q")
		q.Add("slot_id", "w")
		q.Add("banner_id", "e")
		request.URL.RawQuery = q.Encode()
		response := httptest.NewRecorder()

		h.RegisterTransition(response, request)
		require.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestNonExistsRotaton(t *testing.T) {
	d := database.NewDatabase()
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()

	h := Handlers{d, nil}
	url := fmt.Sprintf("http://%s:%d/%s", config.Host(), config.Port(), "/rotation")

	t.Run("add", func(t *testing.T) {
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, url, nil)
		q := request.URL.Query()
		q.Add("banner_id", "-1")
		q.Add("slot_id", "-1")
		request.URL.RawQuery = q.Encode()
		response := httptest.NewRecorder()

		h.HandlerAddToRotation(response, request)
		require.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("delete", func(t *testing.T) {
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, url, nil)
		q := request.URL.Query()
		q.Add("banner_id", "-1")
		q.Add("slot_id", "-1")
		request.URL.RawQuery = q.Encode()
		response := httptest.NewRecorder()

		h.DeleteFromRotation(response, request)
		require.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("select", func(t *testing.T) {
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
		q := request.URL.Query()
		q.Add("group_id", "-1")
		q.Add("slot_id", "-1")
		request.URL.RawQuery = q.Encode()
		response := httptest.NewRecorder()

		h.SelectFromRotation(response, request)
		require.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("register", func(t *testing.T) {
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, url, nil)
		q := request.URL.Query()
		q.Add("group_id", "-1")
		q.Add("slot_id", "-1")
		q.Add("banner_id", "-1")
		request.URL.RawQuery = q.Encode()
		response := httptest.NewRecorder()

		h.RegisterTransition(response, request)
		require.Equal(t, http.StatusNotFound, response.Code)
	})
}
