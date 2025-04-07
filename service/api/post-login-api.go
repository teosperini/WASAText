package api

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/teosperini/WASAText/service/api/reqcontext"
	"net/http"
	"strconv"
)

/*
This function asks the db handler to log in the user who sent the message.
Returns the user ID and the
*/

// Gestione del login e creazione dell'utente
func (rt *_router) doLoginAPI(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	// Decoding the JSON containing the username
	var req Username
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request - Invalid input", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		ctx.Logger.Error("An unknown user did not provide the parameter: username")
		http.Error(w, "Bad Request - Invalid input", http.StatusBadRequest)
		return
	}

	if checkLength(req.Username) {
		ctx.Logger.Error("An unknown user provided an unprocessable username: " + req.Username)
		http.Error(w, "Unprocessable Content - The provided username is not allowed", http.StatusUnprocessableEntity)
		return
	}

	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	mediaUrl := fmt.Sprintf("%s://%s/assets/default.png", scheme, host)

	// Getting the UId from the database
	userID, errStr, errCode := rt.db.DoLoginDB(req.Username, mediaUrl)
	if errCode != nil {
		ctx.Logger.Error(errStr)
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	ctx.Logger.Info(errStr)

	// Creating the response
	response := map[string]interface{}{
		"userId": userID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.Error(ErrWriteResponse + err.Error())
		http.Error(w, "Internal Server Error - Please try again later", http.StatusInternalServerError)
		return
	}

	ctx.Logger.Infof(PrefixUser + strconv.Itoa(userID) + " logged in\n")
}
