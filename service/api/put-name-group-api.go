package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/teosperini/WASAText/service/api/reqcontext"
	"net/http"
	"strconv"
)

func (rt *_router) putNameGroupAPI(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authId, errCode := rt.extractId(r)
	if errCode != nil {
		ctx.Logger.Error("error in the extraction of the userId")
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}
	logrus.Println(" CI SONO ARRIVATO ")
	convIdStr := ps.ByName("convId")

	convId, err := strconv.Atoi(convIdStr)
	if err != nil {
		ctx.Logger.Error(InvalidId + err.Error())
		http.Error(w, "Bad Request - Invalid input", http.StatusBadRequest)
		return
	}
	var req GroupName
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ctx.Logger.Error(err.Error())
		http.Error(w, "Bad Request - Invalid input", http.StatusBadRequest)
		return
	}

	if req.GroupName == "" {
		ctx.Logger.Error(PrefixUser + strconv.Itoa(authId) + " did not provide the parameter: groupName")
		http.Error(w, "Bad Request - Invalid input", http.StatusBadRequest)
		return
	}

	newName := req.GroupName

	if checkLength(newName) {
		ctx.Logger.Error(PrefixUser + strconv.Itoa(authId) + " provided an unprocessable username: " + newName)
		http.Error(w, "Unprocessable Content - The provided username is not allowed", http.StatusUnprocessableEntity)
		return
	}

	errStr, errCode := rt.db.UpdateNameGroupDB(authId, convId, newName)
	if errCode != nil {
		ctx.Logger.Error(errStr)
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	ctx.Logger.Infof(PrefixUser + strconv.Itoa(authId) + " successfully updated name of group " + strconv.Itoa(convId) + " to " + newName)

}
