package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/teosperini/WASAText/service/api/reqcontext"
)

func (rt *_router) getConversationAPI(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	var response ConversationMessages
	messages, errStr, errCode := rt.db.GetConversationDB(authId, convId)
	if errCode != nil {
		ctx.Logger.Error(errStr)
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	response = ConvertToConversationMessages(convId, messages)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.Error(ErrWriteResponse + err.Error())
		http.Error(w, "Internal Server Error - Please try again later", http.StatusInternalServerError)
		return
	}

	ctx.Logger.Infof(PrefixUser + strconv.Itoa(authId) + " downloaded conversation " + convIdStr + "\n")
}
