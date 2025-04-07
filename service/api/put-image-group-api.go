package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/teosperini/WASAText/service/api/reqcontext"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (rt *_router) putImageGroupAPI(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	var req GroupImage
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request - Invalid input", http.StatusBadRequest)
		return
	}

	if req.GroupImage == "" {
		ctx.Logger.Error(PrefixUser + strconv.Itoa(authId) + " did not provide the parameter: chatImageUrl")
		http.Error(w, "Bad Request - Invalid input", http.StatusBadRequest)
		return
	}

	imagePath := filepath.Join(".", "uploads", "images", filepath.Base(req.GroupImage))
	if _, statErr := os.Stat(imagePath); os.IsNotExist(statErr) {
		http.Error(w, "Immagine non trovata sul server", http.StatusBadRequest)
		return
	}

	errStr, errCode := rt.db.UpdateImageGroupDB(authId, req.GroupImage, convId)
	if errCode != nil {
		ctx.Logger.Error(errStr)
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	ctx.Logger.Infof("Image of user " + strconv.Itoa(authId) + " successfully updated to " + imagePath)
}
