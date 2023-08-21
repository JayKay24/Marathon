package controllers

import (
	"encoding/json"
	"io"
	"log"
	"marathon-postgresql/models"
	"marathon-postgresql/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

const ROLE_ADMIN = "admin"
const ROLE_RUNNER = "runner"

type ResultsController struct {
	resultsService *services.ResultsService
	usersService   *services.UsersService
}

func NewResultsController(resultsService *services.ResultsService, usersService *services.UsersService) *ResultsController {
	return &ResultsController{
		resultsService: resultsService,
		usersService:   usersService,
	}
}

func (rc ResultsController) CreateResult(ctx *gin.Context) {
	rc.checkAuthorization(ctx, ROLE_ADMIN)
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

	response, responseErr := rc.resultsService.CreateResult(&result)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (rc ResultsController) DeleteResult(ctx *gin.Context) {
	rc.checkAuthorization(ctx, ROLE_ADMIN)
	resultId := ctx.Param("id")
	responseErr := rc.resultsService.DeleteResult(resultId)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (rc ResultsController) checkAuthorization(ctx *gin.Context, userRole string) {
	var roles []string
	if userRole == ROLE_RUNNER {
		roles = append(roles, ROLE_ADMIN, ROLE_RUNNER)
	} else {
		roles = append(roles, ROLE_ADMIN)
	}

	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, roles)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}
}
