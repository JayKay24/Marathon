package controllers

import (
	"encoding/json"
	"io"
	"log"
	"marathon-postgresql/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RunnersController struct {
	runnersService *services.RunnersService
}

func NewRunnersController(runnersService *services.RunnersService) *RunnersController {
	return &RunnersController{
		runnersService: runnersService,
	}
}

func (rh RunnersController) CreateRunner(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println("Error while unmarshaling create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseError := rh.runnersService.CreateRunner(&runner)
	if responseError != nil {
		ctx.AbortWithStatusJSON(responseError.Status, responseError)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (rh RunnersController) UpdateRunner(ctx *gin.Context) {}

func (rh RunnersController) DeleteRunner(ctx *gin.Context) {}

func (rh RunnersController) GetRunnersBatch(ctx *gin.Context) {}
