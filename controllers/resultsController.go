package controllers

import (
	"encoding/json"
	"io"
	"log"
	"marathon-postgresql/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResultsController struct {
	resultsService *services.ResultsService
}

func NewResultsController(resultsService *services.ResultsService) *ResultsController {
	return &ResultsController{
		resultsService: resultsService,
	}
}

func (rh ResultsController) CreateResult(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var result models.Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Error while unmarshaling create result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := rh.resultsService.CreateResult(&result)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
