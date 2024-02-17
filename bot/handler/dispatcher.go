package handler

import "github.com/slainsama/msgr_server/models"

type UpdateDispatcher struct {
	handlers []Handler
}

func NewUpdateDispatcher() *UpdateDispatcher {
	return &UpdateDispatcher{}
}

func (u *UpdateDispatcher) AddHandler(h Handler) {
	u.handlers = append(u.handlers, h)
}

func (u *UpdateDispatcher) Dispatch(tu *models.TelegramUpdate) {
	for _, h := range u.handlers {
		if h.ShouldHandle(tu) {
			go h.HandlerFunc(tu)
		}
	}
}