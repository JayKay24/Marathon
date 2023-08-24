package controllers

import (
	"encoding/json"
	"io"
	"log"
	"marathon-postgresql/metrics"
	"marathon-postgresql/models"
	"marathon-postgresql/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type RunnersController struct {
	runnersService *services.RunnersService
	usersService   *services.UsersService
}

type ControllerType struct {
	controller string
}

func NewRunnersController(runnersService *services.RunnersService, usersService *services.UsersService) *RunnersController {
	return &RunnersController{
		runnersService: runnersService,
		usersService:   usersService,
	}
}

func (rc RunnersController) CreateRunner(ctx *gin.Context) {
	metrics.HttpRequestsController.Inc()
	rc.checkAuthorization(ctx, ROLE_ADMIN, nil)
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

	response, responseErr := rc.runnersService.CreateRunner(&runner)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (rc RunnersController) UpdateRunner(ctx *gin.Context) {
	metrics.HttpRequestsController.Inc()
	rc.checkAuthorization(ctx, ROLE_ADMIN, nil)
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading update runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println("Error while unmarshaling update runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	responseErr := rc.runnersService.UpdateRunner(&runner)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (rc RunnersController) DeleteRunner(ctx *gin.Context) {
	metrics.HttpRequestsController.Inc()
	rc.checkAuthorization(ctx, ROLE_ADMIN, nil)
	runnerId := ctx.Param("id")
	responseErr := rc.runnersService.DeleteRunner(runnerId)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (rc RunnersController) GetRunner(ctx *gin.Context) {
	metrics.HttpRequestsController.Inc()
	rc.checkAuthorization(ctx, ROLE_RUNNER, nil)
	runnerId := ctx.Param("id")
	response, responseErr := rc.runnersService.GetRunner(runnerId)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (rc RunnersController) GetRunnersBatch(ctx *gin.Context) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		metrics.GetAllRunnersTimer.Observe(f)
	}))

	defer func() {
		timer.ObserveDuration()
	}()

	metrics.HttpRequestsController.Inc()
	rc.checkAuthorization(
		ctx,
		ROLE_RUNNER,
		&ControllerType{
			controller: "GetRunnersBatch",
		})
	params := ctx.Request.URL.Query()
	country := params.Get("country")
	year := params.Get("year")
	response, responseErr := rc.runnersService.GetRunnersBatch(country, year)
	if responseErr != nil {
		metrics.GetRunnerHttpResponsesCounter.
			WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	metrics.GetRunnerHttpResponsesCounter.WithLabelValues("200").Inc()
	ctx.JSON(http.StatusOK, response)
}

func (rc RunnersController) checkAuthorization(
	ctx *gin.Context,
	userRole string,
	controller *ControllerType) {
	var roles []string
	if userRole == ROLE_RUNNER {
		roles = append(roles, ROLE_ADMIN, ROLE_RUNNER)
	} else {
		roles = append(roles, ROLE_ADMIN)
	}

	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, roles)
	if responseErr != nil {
		switch controller.controller {
		case "GetRunnersBatch":
			metrics.GetRunnerHttpResponsesCounter.
				WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
			ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
			return
		default:
			ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
			return
		}
	}
	if !auth {
		switch controller.controller {
		case "GetRunnersBatch":
			metrics.GetRunnerHttpResponsesCounter.WithLabelValues("401").Inc()
			ctx.Status(http.StatusUnauthorized)
			return
		default:
			ctx.Status(http.StatusUnauthorized)
			return
		}
	}
}
