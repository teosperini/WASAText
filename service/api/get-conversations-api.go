package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/teosperini/WASAText/service/api/reqcontext"
)

func (rt *_router) getConversationsAPI(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	authId, errCode := rt.extractId(r)
	if errCode != nil {
		// STAMPO L'ERRORE INTEGRALE NELLA CONSOLE DEL BACKEND
		ctx.Logger.Error("error in the extraction of the userId" + strconv.Itoa(authId))
		// INVIO AL CLIENT UN ERRORE MENO PRECISO PER MIGLIOR SICUREZZA
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	var response []ConversationPreview
	conversations, errStr, errCode := rt.db.GetConversations(authId)
	if errCode != nil {
		ctx.Logger.Error(errStr)
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	err := rt.db.MarkMessagesAsDelivered(authId, conversations)
	if err != nil {
		ctx.Logger.Warn("Errore nel marcare i messaggi come delivered: " + err.Error())
	}

	for i := 0; i < len(conversations); i++ {
		response = append(response, ConvertToConversationPreview(conversations[i]))
	}

	if len(response) == 0 {
		response = []ConversationPreview{}
	}

	responseObject := struct {
		Conversations []ConversationPreview `json:"conversations"`
	}{
		Conversations: response,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(responseObject); err != nil {
		ctx.Logger.Error(ErrWriteResponse + err.Error())
		http.Error(w, "Internal Server Error - Please try again later", http.StatusInternalServerError)
		return
	}

	ctx.Logger.Infof(PrefixUser + strconv.Itoa(authId) + " downloaded all his conversations\n")
}
