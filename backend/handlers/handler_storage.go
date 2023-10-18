package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/interfaces"
)

type GenerateSignedURLHandler struct {
	StorageService interfaces.StorageService
	Validator      *validator.Validate
}

func (g *GenerateSignedURLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := GetLogger(ctx)
	var reqBody struct {
		FileName string `json:"fileName" validate:"required"`
	}
	defer r.Body.Close()
	if err := JsonToStruct(r, &reqBody); err != nil {
		logger.Error().Msgf("failed to parse request body: %v", err)
		ResponsdBadRequest(w, r, err)
		return
	}

	if err := g.Validator.Struct(reqBody); err != nil {
		logger.Error().Msgf("failed to validate request body: %v", err)
		ResponsdBadRequest(w, r, err)
		return
	}

	signedUrl, destinationUrl, err := g.StorageService.GenerateThumbnailPutURL(reqBody.FileName)
	if err != nil {
		logger.Error().Msgf("failed to generate signed url: %v", err)
		ResponsdInternalServerError(w, r, err)
		return
	}

	resp := struct {
		SignedUrl string `json:"signedUrl"`
		PutedUrl  string `json:"putUrl"`
	}{
		SignedUrl: signedUrl,
		PutedUrl:  destinationUrl,
	}
	if err := RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error().Msgf("failed to respond json response: %v", err)
	}
}