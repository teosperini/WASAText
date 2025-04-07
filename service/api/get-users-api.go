package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/teosperini/WASAText/service/api/reqcontext"
)

func (rt *_router) getUsersAPI(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	authId, errCode := rt.extractId(r)
	if errCode != nil {
		ctx.Logger.Error("error in the extraction of the userId")
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	name := r.URL.Query().Get("name")

	users, errStr, errCode := rt.db.GetUsersDB(name)
	if errCode != nil {
		ctx.Logger.Error(errStr)
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	// users è di tipo UserDB
	// creo una struct response di tipo UserAPI (l'equivalente per l'API di UserDB)
	var response []UserAPI

	// Sposto i valori dal users di tipo UserDB sulla response di tipo UserAPI
	for _, user := range users {
		response = append(response, UserAPI{
			Username: user.Username,
			Avatar:   user.Image,
		})
	}

	responseObject := struct {
		Users []UserAPI `json:"users"`
	}{
		Users: response,
	}

	ctx.Logger.Println(responseObject)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(responseObject); err != nil {
		ctx.Logger.Error(ErrWriteResponse + err.Error())
		http.Error(w, "Internal Server Error - Please try again later", http.StatusInternalServerError)
		return
	}

	ctx.Logger.Infof(PrefixUser + strconv.Itoa(authId) + " downloaded the list of users\n")
}
