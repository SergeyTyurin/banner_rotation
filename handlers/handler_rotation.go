package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SergeyTyurin/banner-rotation/database"
)

func (h *Handlers) HandlerAddToRotation(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("slot_id") || !r.URL.Query().Has("banner_id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slotID, slotErr := strconv.Atoi(r.URL.Query().Get("slot_id"))
	bannerID, bannerErr := strconv.Atoi(r.URL.Query().Get("banner_id"))
	if slotErr != nil || bannerErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.db.DatabaseAddToRotation(bannerID, slotID)
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

	slotID, slotErr := strconv.Atoi(r.URL.Query().Get("slot_id"))
	bannerID, bannerErr := strconv.Atoi(r.URL.Query().Get("banner_id"))
	if slotErr != nil || bannerErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.db.DatabaseDeleteFromRotation(bannerID, slotID)
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

	groupID, groupErr := strconv.Atoi(r.URL.Query().Get("group_id"))
	slotID, slotErr := strconv.Atoi(r.URL.Query().Get("slot_id"))
	bannerID, bannerErr := strconv.Atoi(r.URL.Query().Get("banner_id"))
	if groupErr != nil || slotErr != nil || bannerErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.db.DatabaseRegisterTransition(slotID, bannerID, groupID)
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
		msg := fmt.Sprintf("slot_id=%d, group_id=%d, banner_id=%d", slotID, groupID, bannerID)
		_ = h.broker.SendRegisterTransitionEvent(msg)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) SelectFromRotation(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("slot_id") || !r.URL.Query().Has("group_id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	groupID, groupErr := strconv.Atoi(r.URL.Query().Get("group_id"))
	slotID, slotErr := strconv.Atoi(r.URL.Query().Get("slot_id"))
	if groupErr != nil || slotErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bannerID, err := h.db.DatabaseSelectFromRotation(slotID, groupID)
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
		msg := fmt.Sprintf("slot_id=%d, group_id=%d, banner_id=%d", slotID, groupID, bannerID)
		_ = h.broker.SendSelectFromRotationEvent(msg)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(strconv.Itoa(bannerID)))
}
