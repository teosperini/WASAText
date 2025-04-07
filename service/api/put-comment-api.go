package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/teosperini/WASAText/service/api/reqcontext"
	"net/http"
	"strconv"
)

func (rt *_router) putEmojiAPI(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authId, errCode := rt.extractId(r)
	if errCode != nil {
		ctx.Logger.Error("error in the extraction of the userId")
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	var convId int

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

	var req EmojiAPI
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request - Invalid input", http.StatusBadRequest)
		return
	}

	errStr, errCode := rt.db.PutEmojiDB(authId, convId, messId, toEmojiDB(req))
	if errCode != nil {
		ctx.Logger.Error(errStr)
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	ctx.Logger.Infof(PrefixUser + strconv.Itoa(authId) + " successfully commented message " + messIdStr)
}
