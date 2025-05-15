package api

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/teosperini/WASAText/service/api/reqcontext"
)

func (rt *_router) postConversationAPI(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	authId, errCode := rt.extractId(r)
	if errCode != nil {
		ctx.Logger.Error("error in the extraction of the userId")
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	var req ConversationCreateAPI
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ctx.Logger.Error("errore nella decodifica del json: " + err.Error())
		http.Error(w, "Bad Request - Invalid input", http.StatusBadRequest)
		return
	}

	if req.ChatType == "group" && req.GroupImageUrl == "" {
		host := r.Host
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		req.GroupImageUrl = fmt.Sprintf("%s://%s/assets/group_default.png", scheme, host)
	}

	logrus.Println("l'immagine è " + req.GroupImageUrl)

	var response []MessageToClient
	convId, messages, errStr, errCode := rt.db.PostConversationDB(authId, convCreateToDB(req))
	if errCode != nil {
		ctx.Logger.Error(errStr)
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}
	for i := 0; i < len(messages); i++ {
		response = append(response, ConvertToMessageToClient(messages[i]))
	}

	responseObject := ConversationMessages{
		ConversationID: convId,
		Messages:       response,
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(responseObject); err != nil {
		ctx.Logger.Error(ErrWriteResponse + err.Error())
		http.Error(w, "Internal Server Error - Please try again later", http.StatusInternalServerError)
		return
	}

	ctx.Logger.Infof(PrefixUser + strconv.Itoa(authId) + " created conversation " + strconv.Itoa(convId) + "\n")
}
