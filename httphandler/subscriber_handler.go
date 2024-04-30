package httphandler

import (
	"net/http"
)

func (h *Handler) PostSubscribe(w http.ResponseWriter, r *http.Request) {
	address := r.PathValue("address")
	if address == "" {
		BadRequest(w, "address path param must be provided")
		return
	}

	err := h.Repo.AddSubscriber(address)
	if err != nil {
		InternalServerError(w)
		return
	}

	Created(w)
}
