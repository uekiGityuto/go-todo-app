package handler

import "net/http"

type Logout struct {
	Service LogoutService
}

func (l Logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := l.Service.Logout(r); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		Message string `json:"message"`
	}{
		Message: "ok",
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
