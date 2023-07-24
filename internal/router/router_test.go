package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	test_host = "127.0.0.1"
	test_port = 3333
)

func TestCorrectURL(t *testing.T) {
	mux := NewRouter(nil, nil).Mux()
	urls := []struct {
		url    string
		method string
	}{
		{"banner", http.MethodGet},
		{"banner", http.MethodPut},
		{"banner", http.MethodDelete},
		{"banner", http.MethodPost},
		{"slot", http.MethodGet},
		{"slot", http.MethodPut},
		{"slot", http.MethodDelete},
		{"slot", http.MethodPost},
		{"group", http.MethodGet},
		{"group", http.MethodPut},
		{"group", http.MethodDelete},
		{"group", http.MethodPost},
		{"rotation", http.MethodGet},
		{"rotation", http.MethodPut},
		{"rotation", http.MethodDelete},
		{"rotation", http.MethodPost},
		{"", http.MethodPost},
	}

	for _, test := range urls {
		url := fmt.Sprintf("http://%s:%d/%s", test_host, test_port, test.url)
		request, _ := http.NewRequest(http.MethodGet, url, nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		require.NotEqual(t, http.StatusNotFound, response.Result().StatusCode, url)
	}
}

func TestIncorrectURL(t *testing.T) {
	mux := NewRouter(nil, nil).Mux()
	urls := []struct {
		url    string
		method string
	}{
		{"banner/", http.MethodGet},
		{"banner/rotation", http.MethodPut},
		{"slots/", http.MethodGet},
		{"statistic", http.MethodPut},
		{"groups", http.MethodGet},
		{"groups/", http.MethodPut},
		{"rotation/get", http.MethodGet},
	}

	for _, test := range urls {
		url := fmt.Sprintf("http://%s:%d/%s", test_host, test_port, test.url)
		request, _ := http.NewRequest(http.MethodGet, url, nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		require.Equal(t, http.StatusNotFound, response.Code, url)
	}
}

func TestIncorrectMethod(t *testing.T) {
	mux := NewRouter(nil, nil).Mux()
	urls := []struct {
		url    string
		method string
	}{
		{"banner", http.MethodPatch},
		{"slot", http.MethodPatch},
		{"group", http.MethodPatch},
		{"rotation", http.MethodPatch},
	}

	for _, test := range urls {
		url := fmt.Sprintf("http://%s:%d/%s", test_host, test_port, test.url)
		request, _ := http.NewRequest(http.MethodGet, url, nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		require.NotEqual(t, http.StatusMethodNotAllowed, response.Code, url)
	}
}
