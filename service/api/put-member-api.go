package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/teosperini/WASAText/service/api/reqcontext"
	"net/http"
	"strconv"
)

func (rt *_router) addMemberToGroupAPI(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authId, errCode := rt.extractId(r)
	if errCode != nil {
		ctx.Logger.Error("error in the extraction of the userId")
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	convIdStr := ps.ByName("convId")

	convId, err := strconv.Atoi(convIdStr)
	if err != nil {
		ctx.Logger.Error(InvalidId + err.Error())
		http.Error(w, "Bad Request - Invalid input", http.StatusBadRequest)
		return
	}

	var req Username
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ctx.Logger.Error(err.Error())
		http.Error(w, "Bad Request - Invalid input", http.StatusBadRequest)
		return
	}

	newUser := req.Username

	errStr, errCode := rt.db.AddMemberToGroupDB(authId, convId, newUser)
	if errCode != nil {
		ctx.Logger.Error(errStr)
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	ctx.Logger.Infof(PrefixUser + strconv.Itoa(authId) + " successfully added user " + newUser + " to group " + strconv.Itoa(convId))
}
