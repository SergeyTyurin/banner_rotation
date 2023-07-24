package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/SergeyTyurin/banner_rotation/configs"
	"github.com/SergeyTyurin/banner_rotation/message_broker"
	"github.com/SergeyTyurin/banner_rotation/structures"
	"github.com/stretchr/testify/require"
)

func createBanner(url string) structures.Banner {
	jsonBody := []byte(`{"info":"New Banner"}`)
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	response, _ := http.DefaultClient.Do(request)
	responseBody := new(bytes.Buffer)
	responseBody.ReadFrom(response.Body)
	var created structures.Banner
	json.Unmarshal(responseBody.Bytes(), &created)
	return created
}

func createSlot(url string) structures.Slot {
	jsonBody := []byte(`{"info":"New Slot"}`)
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	response, _ := http.DefaultClient.Do(request)
	responseBody := new(bytes.Buffer)
	responseBody.ReadFrom(response.Body)
	var created structures.Slot
	json.Unmarshal(responseBody.Bytes(), &created)
	return created
}

func addToRotation(url string, bannerId, slotId int) {
	request, _ := http.NewRequest(http.MethodPost, url, nil)
	q := request.URL.Query()
	q.Add("banner_id", strconv.Itoa(bannerId))
	q.Add("slot_id", strconv.Itoa(slotId))
	request.URL.RawQuery = q.Encode()
	http.DefaultClient.Do(request)
}

func createGroup(url string) structures.Group {
	jsonBody := []byte(`{"info":"New Group"}`)
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	response, _ := http.DefaultClient.Do(request)
	responseBody := new(bytes.Buffer)
	responseBody.ReadFrom(response.Body)
	var created structures.Group
	json.Unmarshal(responseBody.Bytes(), &created)
	return created
}

func TestBanner(t *testing.T) {
	config, _ := configs.GetAppSettings("../config/connection_config.yaml")
	url := fmt.Sprintf("http://%s:%d/%s", config.Host(), config.Port(), "banner")
	t.Run("create", func(t *testing.T) {
		jsonBody := []byte(`{"info":"New Banner"}`)
		request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, response.StatusCode)
	})

	t.Run("update", func(t *testing.T) {
		created := createBanner(url)
		updated := structures.Banner{Id: created.Id, Info: "Update"}
		jsonBody, _ := json.Marshal(updated)
		request, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(jsonBody))
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("delete", func(t *testing.T) {
		created := createBanner(url)
		request, _ := http.NewRequest(http.MethodDelete, url, nil)
		q := request.URL.Query()
		q.Add("id", strconv.Itoa(created.Id))
		request.URL.RawQuery = q.Encode()

		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("get", func(t *testing.T) {
		created := createBanner(url)
		request, _ := http.NewRequest(http.MethodGet, url, nil)
		q := request.URL.Query()
		q.Add("id", strconv.Itoa(created.Id))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		responceBody := new(bytes.Buffer)
		responceBody.ReadFrom(response.Body)
		var fromDB structures.Banner
		json.Unmarshal(responceBody.Bytes(), &fromDB)
		require.Equal(t, created.Id, fromDB.Id)
	})
}

func TestSlot(t *testing.T) {
	config, _ := configs.GetAppSettings("../config/connection_config.yaml")
	url := fmt.Sprintf("http://%s:%d/%s", config.Host(), config.Port(), "slot")
	t.Run("create", func(t *testing.T) {
		jsonBody := []byte(`{"info":"New Banner"}`)
		request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, response.StatusCode)
	})

	t.Run("update", func(t *testing.T) {
		created := createSlot(url)
		updated := structures.Slot{Id: created.Id, Info: "Update"}
		jsonBody, _ := json.Marshal(updated)
		request, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(jsonBody))
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("delete", func(t *testing.T) {
		created := createSlot(url)
		request, _ := http.NewRequest(http.MethodDelete, url, nil)
		q := request.URL.Query()
		q.Add("id", strconv.Itoa(created.Id))
		request.URL.RawQuery = q.Encode()

		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("get", func(t *testing.T) {
		created := createSlot(url)
		request, _ := http.NewRequest(http.MethodGet, url, nil)
		q := request.URL.Query()
		q.Add("id", strconv.Itoa(created.Id))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		responceBody := new(bytes.Buffer)
		responceBody.ReadFrom(response.Body)
		var fromDB structures.Slot
		json.Unmarshal(responceBody.Bytes(), &fromDB)
		require.Equal(t, created.Id, fromDB.Id)
	})
}

func TestGroup(t *testing.T) {
	config, _ := configs.GetAppSettings("../config/connection_config.yaml")
	url := fmt.Sprintf("http://%s:%d/%s", config.Host(), config.Port(), "group")
	t.Run("create", func(t *testing.T) {
		jsonBody := []byte(`{"info":"New Banner"}`)
		request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, response.StatusCode)
	})

	t.Run("update", func(t *testing.T) {
		created := createGroup(url)
		updated := structures.Group{Id: created.Id, Info: "Update"}
		jsonBody, _ := json.Marshal(updated)
		request, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(jsonBody))
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("delete", func(t *testing.T) {
		created := createGroup(url)
		request, _ := http.NewRequest(http.MethodDelete, url, nil)
		q := request.URL.Query()
		q.Add("id", strconv.Itoa(created.Id))
		request.URL.RawQuery = q.Encode()

		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("get", func(t *testing.T) {
		created := createGroup(url)
		request, _ := http.NewRequest(http.MethodGet, url, nil)
		q := request.URL.Query()
		q.Add("id", strconv.Itoa(created.Id))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		responceBody := new(bytes.Buffer)
		responceBody.ReadFrom(response.Body)
		var fromDB structures.Group
		json.Unmarshal(responceBody.Bytes(), &fromDB)
		require.Equal(t, created.Id, fromDB.Id)
	})
}

func TestRotation(t *testing.T) {
	config, _ := configs.GetAppSettings("../config/connection_config.yaml")
	mq_config, _ := configs.GetMessageBrokerConfig("../config/connection_config.yaml")
	url := fmt.Sprintf("http://%s:%d", config.Host(), config.Port())
	broker := message_broker.NewBroker()
	mq_close, _ := broker.Connect(mq_config)
	defer mq_close()
	t.Run("add to rotation", func(t *testing.T) {
		banner := createBanner(url + "/banner")
		slot := createSlot(url + "/slot")
		request, _ := http.NewRequest(http.MethodPost, url+"/rotation", nil)
		q := request.URL.Query()
		q.Add("banner_id", strconv.Itoa(banner.Id))
		q.Add("slot_id", strconv.Itoa(slot.Id))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("register transition", func(t *testing.T) {
		banner := createBanner(url + "/banner")
		slot := createSlot(url + "/slot")
		group := createGroup(url + "/group")
		addToRotation(url+"/rotation", banner.Id, slot.Id)
		request, _ := http.NewRequest(http.MethodPut, url+"/rotation", nil)
		q := request.URL.Query()
		q.Add("banner_id", strconv.Itoa(banner.Id))
		q.Add("slot_id", strconv.Itoa(slot.Id))
		q.Add("group_id", strconv.Itoa(group.Id))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
		msg, err := broker.GetRegisterTransitionEvent()
		require.NoError(t, err)
		require.NotEmpty(t, msg)
		expected := fmt.Sprintf("slot_id=%d, group_id=%d, banner_id=%d",
			slot.Id, group.Id, banner.Id)
		require.Equal(t, expected, msg)
	})

	t.Run("delete from rotation", func(t *testing.T) {
		banner := createBanner(url + "/banner")
		slot := createSlot(url + "/slot")
		addToRotation(url+"/rotation", banner.Id, slot.Id)
		request, _ := http.NewRequest(http.MethodDelete, url+"/rotation", nil)
		q := request.URL.Query()
		q.Add("banner_id", strconv.Itoa(banner.Id))
		q.Add("slot_id", strconv.Itoa(slot.Id))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
		require.NoError(t, err)
	})

	t.Run("select from rotation", func(t *testing.T) {
		banner := createBanner(url + "/banner")
		slot := createSlot(url + "/slot")
		group := createGroup(url + "/group")
		addToRotation(url+"/rotation", banner.Id, slot.Id)
		request, _ := http.NewRequest(http.MethodGet, url+"/rotation", nil)
		q := request.URL.Query()
		q.Add("slot_id", strconv.Itoa(slot.Id))
		q.Add("group_id", strconv.Itoa(group.Id))
		request.URL.RawQuery = q.Encode()
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		responseBody := new(bytes.Buffer)
		responseBody.ReadFrom(response.Body)
		selectedId, _ := strconv.Atoi(responseBody.String())
		require.Equal(t, banner.Id, selectedId)

		msg, err := broker.GetSelectFromRotationEvent()
		require.NoError(t, err)
		require.NotEmpty(t, msg)
		expected := fmt.Sprintf("slot_id=%d, group_id=%d, banner_id=%d",
			slot.Id, group.Id, banner.Id)
		require.Equal(t, expected, msg)
	})
}
