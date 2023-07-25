package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SergeyTyurin/banner_rotation/database"
)

func (h *Handlers) AddToRotation(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("slot_id") || !r.URL.Query().Has("banner_id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slotId, slotErr := strconv.Atoi(r.URL.Query().Get("slot_id"))
	bannerId, bannerErr := strconv.Atoi(r.URL.Query().Get("banner_id"))
	if slotErr != nil || bannerErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.db.AddToRotation(bannerId, slotId)
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

func (h *Handlers) DeleteFromRotation(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("slot_id") || !r.URL.Query().Has("banner_id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slotId, slotErr := strconv.Atoi(r.URL.Query().Get("slot_id"))
	bannerId, bannerErr := strconv.Atoi(r.URL.Query().Get("banner_id"))
	if slotErr != nil || bannerErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.db.DeleteFromRotation(bannerId, slotId)
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

func (h *Handlers) RegisterTransition(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("slot_id") || !r.URL.Query().Has("group_id") || !r.URL.Query().Has("banner_id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	groupId, groupErr := strconv.Atoi(r.URL.Query().Get("group_id"))
	slotId, slotErr := strconv.Atoi(r.URL.Query().Get("slot_id"))
	bannerId, bannerErr := strconv.Atoi(r.URL.Query().Get("banner_id"))
	if groupErr != nil || slotErr != nil || bannerErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.db.RegisterTransition(slotId, bannerId, groupId)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) || errors.Is(err, database.ErrNotInRotation) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if h.broker != nil {
		msg := fmt.Sprintf("slot_id=%d, group_id=%d, banner_id=%d", slotId, groupId, bannerId)
		_ = h.broker.SendRegisterTransitionEvent(msg)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) SelectFromRotation(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("slot_id") || !r.URL.Query().Has("group_id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	groupId, groupErr := strconv.Atoi(r.URL.Query().Get("group_id"))
	slotId, slotErr := strconv.Atoi(r.URL.Query().Get("slot_id"))
	if groupErr != nil || slotErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bannerId, err := h.db.SelectFromRotation(slotId, groupId)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if h.broker != nil {
		msg := fmt.Sprintf("slot_id=%d, group_id=%d, banner_id=%d", slotId, groupId, bannerId)
		_ = h.broker.SendSelectFromRotationEvent(msg)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(strconv.Itoa(bannerId)))
}
