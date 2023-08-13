package httpinout

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/juicyluv/structure-experiments/internal/pkg/logger"
)

func Json(w http.ResponseWriter, val any, status int) {
	data, err := json.Marshal(val)
	if err != nil {
		logger.Get().Fatal(err.Error())
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func ReadJSON(r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		return fmt.Errorf("decoding json: %w", err)
	}

	return nil
}
