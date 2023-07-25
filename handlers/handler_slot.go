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

func (h *Handlers) GetSlot(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	slot, err := h.db.GetSlot(id)
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
	resp, _ := json.Marshal(slot)
	_, _ = w.Write(resp)
}

func (h *Handlers) CreateSlot(w http.ResponseWriter, r *http.Request) {
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var slot structures.Slot
	err = json.Unmarshal(requestBody.Bytes(), &slot)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdSlot, err := h.db.CreateSlot(slot)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	resp, _ := json.Marshal(createdSlot)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)
}

func (h *Handlers) UpdateSlot(w http.ResponseWriter, r *http.Request) {
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var slot structures.Slot
	err = json.Unmarshal(requestBody.Bytes(), &slot)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.db.UpdateSlot(slot)
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

func (h *Handlers) DeleteSlot(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.db.DeleteSlot(id)
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
