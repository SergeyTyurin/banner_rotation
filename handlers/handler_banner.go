package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SergeyTyurin/banner-rotation/database"
	"github.com/SergeyTyurin/banner-rotation/structures"
)

func (h *Handlers) GetBanner(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	banner, err := h.db.DatabaseGetBanner(id)
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
	resp, _ := json.Marshal(banner)
	_, _ = w.Write(resp)
}

func (h *Handlers) CreateBanner(w http.ResponseWriter, r *http.Request) {
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var banner structures.Banner
	err = json.Unmarshal(requestBody.Bytes(), &banner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdBanner, err := h.db.DatabaseCreateBanner(banner)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	resp, _ := json.Marshal(createdBanner)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)
}

func (h *Handlers) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var banner structures.Banner
	err = json.Unmarshal(requestBody.Bytes(), &banner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.db.DatabaseUpdateBanner(banner)
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

func (h *Handlers) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.db.DatabaseDeleteBanner(id)
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
