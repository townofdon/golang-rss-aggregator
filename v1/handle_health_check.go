package v1

import (
	"net/http"

	"github.com/townofdon/tutorial-go-rss-server/util"
)

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	util.RespondWithJSON(w, 200, struct {
		Msg string `json:"msg"`
	}{
		Msg: "hello",
	})
}
