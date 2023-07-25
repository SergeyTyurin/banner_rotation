package router

import (
	"net/http"
	"strings"
)

func (router *routerImpl) handleGroupsFunc(w http.ResponseWriter, r *http.Request) {
	if strings.TrimRight(r.URL.Path, "/") != "/group" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		router.handlers.GetGroup(w, r) // Получить группу пользователей
	case http.MethodPost:
		router.handlers.CreateGroup(w, r) // Создание новой группы
	case http.MethodPut:
		router.handlers.UpdateGroup(w, r) // Обновление информации о группе
	case http.MethodDelete:
		router.handlers.DeleteGroup(w, r) // Удаление группы
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (router *routerImpl) handleSlotsFunc(w http.ResponseWriter, r *http.Request) {
	if strings.TrimRight(r.URL.Path, "/") != "/slot" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		router.handlers.GetSlot(w, r) // Получить слот
	case http.MethodPost:
		router.handlers.CreateSlot(w, r) // Создание нового слота
	case http.MethodPut:
		router.handlers.UpdateSlot(w, r) // Обновление информации о слоте
	case http.MethodDelete:
		router.handlers.DeleteSlot(w, r) // Удаление слота
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (router *routerImpl) handleRotationFunc(w http.ResponseWriter, r *http.Request) {
	if strings.TrimRight(r.URL.Path, "/") != "/rotation" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		router.handlers.SelectFromRotation(w, r) // Выбор баннера из ротации
	case http.MethodPost:
		router.handlers.HandlerAddToRotation(w, r) // Добавление баннера в ротацию
	case http.MethodPut:
		router.handlers.RegisterTransition(w, r) // Зарегистрировать переход по баннеру
	case http.MethodDelete:
		router.handlers.DeleteFromRotation(w, r) // Удаление баннера из ротации
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (router *routerImpl) handleBannersFunc(w http.ResponseWriter, r *http.Request) {
	if strings.TrimRight(r.URL.Path, "/") != "/banner" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		router.handlers.GetBanner(w, r) // Получение баннер
	case http.MethodPost:
		router.handlers.CreateBanner(w, r) // Создание нового баннера
	case http.MethodPut:
		router.handlers.UpdateBanner(w, r) // Обновление информации о баннере
	case http.MethodDelete:
		router.handlers.DeleteBanner(w, r) // Удаление баннера
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
