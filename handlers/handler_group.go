package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SergeyTyurin/banner_rotation/database"
	"github.com/SergeyTyurin/banner_rotation/structures"
)

func (h *Handlers) GetGroup(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	group, err := h.db.GetGroup(id)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(group)
	_, _ = w.Write(resp)
}

func (h *Handlers) CreateGroup(w http.ResponseWriter, r *http.Request) {
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var group structures.Group
	err = json.Unmarshal(requestBody.Bytes(), &group)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdGroup, err := h.db.CreateGroup(group)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	resp, _ := json.Marshal(createdGroup)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)
}

func (h *Handlers) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var group structures.Group
	err = json.Unmarshal(requestBody.Bytes(), &group)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.db.UpdateGroup(group)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.db.DeleteGroup(id)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
