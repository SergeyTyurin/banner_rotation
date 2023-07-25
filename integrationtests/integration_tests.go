package integrationtests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/SergeyTyurin/banner-rotation/configs"
	"github.com/SergeyTyurin/banner-rotation/messagebroker"
	"github.com/SergeyTyurin/banner-rotation/structures"
	"github.com/stretchr/testify/require"
)

func createBanner(url string) structures.Banner {
	jsonBody := []byte(`{"info":"New Banner"}`)
	request, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(jsonBody))
	response, _ := http.DefaultClient.Do(request)
	responseBody := new(bytes.Buffer)
	defer func() {
		_ = response.Body.Close()
	}()
	_, _ = responseBody.ReadFrom(response.Body)
	var created structures.Banner
	_ = json.Unmarshal(responseBody.Bytes(), &created)
	return created
}

func createSlot(url string) structures.Slot {
	jsonBody := []byte(`{"info":"New Slot"}`)
	request, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(jsonBody))
	response, _ := http.DefaultClient.Do(request)
	defer func() {
		_ = response.Body.Close()
	}()

	responseBody := new(bytes.Buffer)
	_, _ = responseBody.ReadFrom(response.Body)
	var created structures.Slot
	_ = json.Unmarshal(responseBody.Bytes(), &created)
	return created
}

func addToRotation(url string, bannerID, slotID int) {
	request, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, url, nil)
	q := request.URL.Query()
	q.Add("banner_id", strconv.Itoa(bannerID))
	q.Add("slot_id", strconv.Itoa(slotID))
	request.URL.RawQuery = q.Encode()
	_, _ = http.DefaultClient.Do(request) //nolint:bodyclose
}

func createGroup(url string) structures.Group {
	jsonBody := []byte(`{"info":"New Group"}`)
	request, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(jsonBody))
	response, _ := http.DefaultClient.Do(request)
	defer func() {
		_ = response.Body.Close()
	}()
	responseBody := new(bytes.Buffer)
	_, _ = responseBody.ReadFrom(response.Body)
	var created structures.Group
	_ = json.Unmarshal(responseBody.Bytes(), &created)
	return created
}

func TestBanner(t *testing.T) {
	config, _ := configs.GetAppSettings("../config/test/test_connection_config.yaml")
	url := fmt.Sprintf("http://%s:%d/%s", config.Host(), config.Port(), "banner")
	t.Run("create", func(t *testing.T) {
		jsonBody := []byte(`{"info":"New Banner"}`)
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(jsonBody))
		response, err := http.DefaultClient.Do(request)
		defer func() {
			_ = response.Body.Close()
		}()
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, response.StatusCode)
	})

	t.Run("update", func(t *testing.T) {
		created := createBanner(url)
		updated := structures.Banner{ID: created.ID, Info: "Update"}
		jsonBody, _ := json.Marshal(updated)
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, url, bytes.NewReader(jsonBody))
		response, err := http.DefaultClient.Do(request)
		defer func() {
			_ = response.Body.Close()
		}()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("delete", func(t *testing.T) {
		created := createBanner(url)
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, url, nil)
		q := request.URL.Query()
		q.Add("id", strconv.Itoa(created.ID))
		request.URL.RawQuery = q.Encode()

		response, err := http.DefaultClient.Do(request)
		defer func() {
			_ = response.Body.Close()
		}()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("get", func(t *testing.T) {
		created := createBanner(url)
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
		q := request.URL.Query()
		q.Add("id", strconv.Itoa(created.ID))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request)
		defer func() {
			_ = response.Body.Close()
		}()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		responceBody := new(bytes.Buffer)
		_, _ = responceBody.ReadFrom(response.Body)
		var fromDB structures.Banner
		_ = json.Unmarshal(responceBody.Bytes(), &fromDB)
		require.Equal(t, created.ID, fromDB.ID)
	})
}

func TestSlot(t *testing.T) {
	config, _ := configs.GetAppSettings("../config/test/test_connection_config.yaml")
	url := fmt.Sprintf("http://%s:%d/%s", config.Host(), config.Port(), "slot")
	t.Run("create", func(t *testing.T) {
		jsonBody := []byte(`{"info":"New Banner"}`)
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(jsonBody))
		response, err := http.DefaultClient.Do(request) //nolint:bodyclose
		require.NoError(t, err)
		go func() {
			_ = response.Body.Close()
		}()
		require.Equal(t, http.StatusCreated, response.StatusCode)
	})

	t.Run("update", func(t *testing.T) {
		created := createSlot(url)
		updated := structures.Slot{ID: created.ID, Info: "Update"}
		jsonBody, _ := json.Marshal(updated)
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, url, bytes.NewReader(jsonBody))
		response, err := http.DefaultClient.Do(request) //nolint:bodyclose
		go func() {
			_ = response.Body.Close()
		}()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("delete", func(t *testing.T) {
		created := createSlot(url)
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, url, nil)
		q := request.URL.Query()
		q.Add("id", strconv.Itoa(created.ID))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request) //nolint:bodyclose
		go func() {
			_ = response.Body.Close()
		}()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("get", func(t *testing.T) {
		created := createSlot(url)
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
		q := request.URL.Query()
		q.Add("id", strconv.Itoa(created.ID))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request) //nolint:bodyclose
		go func() {
			_ = response.Body.Close()
		}()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		responceBody := new(bytes.Buffer)
		_, _ = responceBody.ReadFrom(response.Body)
		var fromDB structures.Slot
		_ = json.Unmarshal(responceBody.Bytes(), &fromDB)
		require.Equal(t, created.ID, fromDB.ID)
	})
}

func TestGroup(t *testing.T) {
	config, _ := configs.GetAppSettings("../config/test/test_connection_config.yaml")
	url := fmt.Sprintf("http://%s:%d/%s", config.Host(), config.Port(), "group")
	t.Run("create", func(t *testing.T) {
		jsonBody := []byte(`{"info":"New Banner"}`)
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(jsonBody))
		response, err := http.DefaultClient.Do(request) //nolint:bodyclose
		go func() {
			_ = response.Body.Close()
		}()
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, response.StatusCode)
	})

	t.Run("update", func(t *testing.T) {
		created := createGroup(url)
		updated := structures.Group{ID: created.ID, Info: "Update"}
		jsonBody, _ := json.Marshal(updated)
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, url, bytes.NewReader(jsonBody))
		response, err := http.DefaultClient.Do(request) //nolint:bodyclose
		go func() {
			_ = response.Body.Close()
		}()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("delete", func(t *testing.T) {
		created := createGroup(url)
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, url, nil)
		q := request.URL.Query()
		q.Add("id", strconv.Itoa(created.ID))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request) //nolint:bodyclose
		go func() {
			_ = response.Body.Close()
		}()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("get", func(t *testing.T) {
		created := createGroup(url)
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
		q := request.URL.Query()
		q.Add("id", strconv.Itoa(created.ID))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request) //nolint:bodyclose
		go func() {
			_ = response.Body.Close()
		}()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		responceBody := new(bytes.Buffer)
		_, _ = responceBody.ReadFrom(response.Body)
		var fromDB structures.Group
		_ = json.Unmarshal(responceBody.Bytes(), &fromDB)
		require.Equal(t, created.ID, fromDB.ID)
	})
}

func TestRotation(t *testing.T) {
	config, _ := configs.GetAppSettings("../config/test/test_connection_config.yaml")
	mqConfig, _ := configs.GetMessageBrokerConfig("../config/test/test_connection_config.yaml")
	url := fmt.Sprintf("http://%s:%d", config.Host(), config.Port())
	broker := messagebroker.NewBroker()
	mqClose, _ := broker.Connect(mqConfig)
	defer mqClose()
	t.Run("add to rotation", func(t *testing.T) {
		banner := createBanner(url + "/banner")
		slot := createSlot(url + "/slot")
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, url+"/rotation", nil)
		q := request.URL.Query()
		q.Add("banner_id", strconv.Itoa(banner.ID))
		q.Add("slot_id", strconv.Itoa(slot.ID))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request) //nolint:bodyclose
		go func() {
			_ = response.Body.Close()
		}()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("register transition", func(t *testing.T) {
		banner := createBanner(url + "/banner")
		slot := createSlot(url + "/slot")
		group := createGroup(url + "/group")
		addToRotation(url+"/rotation", banner.ID, slot.ID)
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, url+"/rotation", nil)
		q := request.URL.Query()
		q.Add("banner_id", strconv.Itoa(banner.ID))
		q.Add("slot_id", strconv.Itoa(slot.ID))
		q.Add("group_id", strconv.Itoa(group.ID))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request) //nolint:bodyclose
		go func() {
			_ = response.Body.Close()
		}()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
		msg, err := broker.GetRegisterTransitionEvent()
		require.NoError(t, err)
		require.NotEmpty(t, msg)
		expected := fmt.Sprintf("slot_id=%d, group_id=%d, banner_id=%d",
			slot.ID, group.ID, banner.ID)
		require.Equal(t, expected, msg)
	})

	t.Run("delete from rotation", func(t *testing.T) {
		banner := createBanner(url + "/banner")
		slot := createSlot(url + "/slot")
		addToRotation(url+"/rotation", banner.ID, slot.ID)
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, url+"/rotation", nil)
		q := request.URL.Query()
		q.Add("banner_id", strconv.Itoa(banner.ID))
		q.Add("slot_id", strconv.Itoa(slot.ID))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request) //nolint:bodyclose
		go func() {
			_ = response.Body.Close()
		}()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
		require.NoError(t, err)
	})

	t.Run("select from rotation", func(t *testing.T) {
		banner := createBanner(url + "/banner")
		slot := createSlot(url + "/slot")
		group := createGroup(url + "/group")
		addToRotation(url+"/rotation", banner.ID, slot.ID)
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url+"/rotation", nil)
		q := request.URL.Query()
		q.Add("slot_id", strconv.Itoa(slot.ID))
		q.Add("group_id", strconv.Itoa(group.ID))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request) //nolint:bodyclose
		go func() {
			_ = response.Body.Close()
		}()
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		responseBody := new(bytes.Buffer)
		_, _ = responseBody.ReadFrom(response.Body)
		selectedID, _ := strconv.Atoi(responseBody.String())
		require.Equal(t, banner.ID, selectedID)

		msg, err := broker.GetSelectFromRotationEvent()
		require.NoError(t, err)
		require.NotEmpty(t, msg)
		expected := fmt.Sprintf("slot_id=%d, group_id=%d, banner_id=%d",
			slot.ID, group.ID, banner.ID)
		require.Equal(t, expected, msg)
	})
}
