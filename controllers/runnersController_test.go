package controllers

import (
	"database/sql"
	"encoding/json"
	"marathon-postgresql/models"
	"marathon-postgresql/repositories"
	"marathon-postgresql/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetRunnersResponse(t *testing.T) {
	dbHandler, mock, _ := sqlmock.New()

	defer dbHandler.Close()

	columns := []string{"id", "first_name", "last_name", "age", "is_active", "country", "personal_best", "season_best"}
	mock.ExpectQuery("SELECT * FROM runners").
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow("1", "John", "Smith", 30, true, "United States", "02:00:41", "02:13:13").
			AddRow("2", "Barry", "Allen", 30, true, "United States", "01:18:28", "01:18:27"))

	router := initTestRouter(dbHandler)
	request, _ := http.NewRequest("GET", "/runner", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

	var runners []*models.Runner
	json.Unmarshal(recorder.Body.Bytes(), &runners)

	assert.NotEmpty(t, runners)
	assert.Equal(t, 2, len(runners))
}

func initTestRouter(dbHandler *sql.DB) *gin.Engine {
	runnersRepository := repositories.NewRunnersRepository(dbHandler)
	runnersService := services.NewRunnersService(runnersRepository, nil)
	runnersController := NewRunnersController(runnersService)

	router := gin.Default()
	router.GET("/", runnersController.GetRunnersBatch)
	return router
}
