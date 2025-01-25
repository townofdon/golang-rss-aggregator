package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/townofdon/tutorial-go-rss-server/internal/database"
	"github.com/townofdon/tutorial-go-rss-server/src/api"
	"github.com/townofdon/tutorial-go-rss-server/src/models"
	"github.com/townofdon/tutorial-go-rss-server/src/util"
)

func CreateUser(w http.ResponseWriter, r *http.Request, api *api.Clients) {
	type requestParams struct {
		Name string `json:"name"`
	}

	params := requestParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		util.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := api.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		util.RespondWithError(w, 500, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	util.RespondWithJSON(w, 201, models.DatabaseUserToWire(user))
}

func GetUserByApiKey(w http.ResponseWriter, r *http.Request, api *api.Clients, user database.User) {
	util.RespondWithJSON(w, 200, models.DatabaseUserToWire(user))
}
