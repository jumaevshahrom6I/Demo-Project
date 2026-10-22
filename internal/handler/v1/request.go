package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"demo-project/internal/handler/response"
)

const maxBodyBytes = 1 << 20 // 1 МБ: огромное тело не должно съесть память сервера

// decodeJSON читает тело запроса в dst. При ошибке сам отвечает 400 и возвращает false.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		response.Error(w, r, http.StatusBadRequest, "invalid json body", err)
		return false
	}
	return true
}

// parseID читает id из пути запроса. При ошибке (не число или меньше 1)
// сам отвечает 400 и возвращает false.
func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err == nil && id < 1 {
		err = errors.New("id must be at least 1")
	}
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, "invalid id", err)
		return 0, false
	}
	return id, true
}
