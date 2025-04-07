package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/teosperini/WASAText/service/api/reqcontext"
)

func (rt *_router) putUsernameAPI(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	authId, errCode := rt.extractId(r)
	if errCode != nil {
		ctx.Logger.Error("error in the extraction of the userId")
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	var req Username
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ctx.Logger.Error(err.Error())
		http.Error(w, "Bad Request - Invalid input", http.StatusBadRequest)
		return
	}

	// Validating the username
	if req.Username == "" {
		ctx.Logger.Error("User " + strconv.Itoa(authId) + " did not provide the parameter: username")
		http.Error(w, "Bad Request - Invalid input", http.StatusBadRequest)
		return
	}

	if checkLength(req.Username) {
		ctx.Logger.Error("User " + strconv.Itoa(authId) + " provided an unprocessable username: " + req.Username)
		http.Error(w, "Unprocessable Content - The provided username is not allowed", http.StatusUnprocessableEntity)
		return
	}

	errStr, errCode := rt.db.UpdateUsername(authId, req.Username)
	if errCode != nil {
		ctx.Logger.Error(errStr)
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	ctx.Logger.Infof("Username of user " + strconv.Itoa(authId) + " successfully updated to " + req.Username)
}
