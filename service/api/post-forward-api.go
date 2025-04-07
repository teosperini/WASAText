package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/teosperini/WASAText/service/api/reqcontext"
	"net/http"
	"strconv"
)

func (rt *_router) forwardMessageAPI(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authId, errCode := rt.extractId(r)
	if errCode != nil {
		ctx.Logger.Error("error in the extraction of the userId")
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	newConvIdStr := ps.ByName("convId")

	newConvId, err := strconv.Atoi(newConvIdStr)
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

	newMessageId, errStr, errCode := rt.db.ForwardMessageDB(authId, newConvId, messId)
	if errCode != nil {
		ctx.Logger.Error(errStr)
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	response := map[string]interface{}{
		"messageId": newMessageId,
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.Error(ErrWriteResponse + err.Error())
		http.Error(w, "Internal Server Error - Please try again later", http.StatusInternalServerError)
		return
	}

	ctx.Logger.Infof(PrefixUser + strconv.Itoa(authId) + " forwarded message " + strconv.Itoa(newMessageId) + " to conversation " + strconv.Itoa(newConvId) + "\n")
}
