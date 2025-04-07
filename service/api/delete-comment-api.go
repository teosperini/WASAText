package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/teosperini/WASAText/service/api/reqcontext"
	"net/http"
	"strconv"
)

func (rt *_router) deleteEmojiAPI(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authId, errCode := rt.extractId(r)
	if errCode != nil {
		ctx.Logger.Error("error in the extraction of the userId")
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	convId, err := strconv.Atoi(ps.ByName("convId"))
	if err != nil {
		ctx.Logger.Error("invalid convId: " + err.Error())
		http.Error(w, "Bad Request - Invalid convId", http.StatusBadRequest)
		return
	}

	messId, err := strconv.Atoi(ps.ByName("messId"))
	if err != nil {
		ctx.Logger.Error("invalid messId: " + err.Error())
		http.Error(w, "Bad Request - Invalid messId", http.StatusBadRequest)
		return
	}

	errStr, errCode := rt.db.DeleteEmojiDB(authId, convId, messId)
	if errCode != nil {
		ctx.Logger.Error(errStr)
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Infof(PrefixUser + strconv.Itoa(authId) + " removed emoji from message " + strconv.Itoa(messId))
}
