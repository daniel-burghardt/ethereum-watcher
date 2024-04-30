package httphandler

import "net/http"

func (h *Handler) GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block, err := h.Repo.GetCurrentBlock()
	if err != nil {
		InternalServerError(w)
		return
	}

	err = Json(w, block)
	if err != nil {
		InternalServerError(w)
		return
	}
}
