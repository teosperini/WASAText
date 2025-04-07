package api

import (
	"encoding/json"
	"fmt"
	"github.com/teosperini/WASAText/service/api/reqcontext"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) postMediaAPI(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	authId, errCode := rt.extractId(r)
	if errCode != nil {
		ctx.Logger.Error("error in the extraction of the userId")
		clientErr, httpError := checkErrors(errCode)
		http.Error(w, clientErr, httpError)
		return
	}

	// Limite massimo alla dimensione dei file (10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		ctx.Logger.Println("errore, file troppo grande: ", err.Error())
		http.Error(w, "Errore: "+err.Error(), 500)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		ctx.Logger.Println("errore nella ricezione del file: ", err.Error())
		http.Error(w, "Errore: "+err.Error(), 500)
		return
	}

	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		ctx.Logger.Println("errore formato immagine invalido")
		http.Error(w, "Formato immagine non supportato", http.StatusBadRequest)
		return
	}

	// Nome univoco
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join("uploads", "images", filename)

	// Assicurati che la directory esista
	if err = os.MkdirAll("uploads/images", os.ModePerm); err != nil {
		ctx.Logger.Println("errore durante la creazione della cartella: ", err.Error())
		http.Error(w, "Errore nella creazione della cartella", http.StatusInternalServerError)
		return
	}

	// Creazione file
	out, err := os.Create(savePath)
	if err != nil {
		ctx.Logger.Println("errore durante la creazione del file: ", err.Error())
		http.Error(w, "Errore: "+err.Error(), 500)
		return
	}

	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		ctx.Logger.Println("❌ ERRORE:", err)
		http.Error(w, "Errore: "+err.Error(), 500)
		return
	}

	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	mediaUrl := fmt.Sprintf("%s://%s/media/%s", scheme, host, filename)

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(map[string]string{
		"mediaUrl": mediaUrl,
	}); err != nil {
		ctx.Logger.Error(ErrWriteResponse + err.Error())
		http.Error(w, "Internal Server Error - Please try again later", http.StatusInternalServerError)
		return
	}

	ctx.Logger.Infof(PrefixUser + strconv.Itoa(authId) + " successfully uploaded media: " + mediaUrl + " to the server")
}
