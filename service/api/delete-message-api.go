package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/teosperini/WASAText/service/api/reqcontext"
	"net/http"
	"strconv"
)

func (rt *_router) deleteMessageAPI(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	messIdStr := ps.ByName("messId")

	messId, err := strconv.Atoi(messIdStr)
	if err != nil {
		ctx.Logger.Error("invalid messId: " + err.Error())
		http.Error(w, "Bad Request - Invalid messId", http.StatusBadRequest)
		return
	}

	errStr, errCode := rt.db.DeleteMessageDB(authId, convId, messId)
	if errCode != nil {
		ctx.Logger.Error(errStr)
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	ctx.Logger.Infof(PrefixUser + strconv.Itoa(authId) + " successfully deleted message " + strconv.Itoa(messId) + " from conversation " + strconv.Itoa(convId))
}
