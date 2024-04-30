package httphandler

import (
	"fmt"
	"net/http"
)

func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.PathValue("address")
	if address == "" {
		BadRequest(w, "address path param must be provided")
		return
	}

	isSubscribed, err := h.Repo.SubscriberExists(address)
	if err != nil {
		InternalServerError(w)
		return
	}
	if !isSubscribed {
		BadRequest(w, fmt.Sprintf("address %s is not subscribed", address))
		return
	}

	txs, err := h.Repo.GetTransactions(address)
	if err != nil {
		InternalServerError(w)
		return
	}

	err = Json(w, txs)
	if err != nil {
		InternalServerError(w)
		return
	}
}
