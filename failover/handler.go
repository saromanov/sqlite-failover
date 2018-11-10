package failover

import (
	"net/http"
	"strings"
)

type masterHandler struct {
	a *App
}

func (h *masterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		masters := h.a.GetMasters()
		w.Write([]byte(strings.Join(masters, ",")))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
