package repositories

import (
	"database/sql"
	"marathon-postgresql/models"
	"net/http"
)

type ResultsRepository struct {
	dbHandler   *sql.DB
	transaction *sql.Tx
}

func NewResultsRepository(dbHandler *sql.DB) *ResultsRepository {
	return &ResultsRepository{
		dbHandler: dbHandler,
	}
}

func (rr ResultsRepository) CreateResult(result *models.Result) (*models.Result, *models.ResponseError) {
	query := `
		INSERT INTO results(runner_id, race_result, location, position, year)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	rows, err := rr.transaction.Query(query, result.RunnerID, result.RaceResult, result.Location, result.Position, result.Year)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var resultId string
	for rows.Next() {
		err := rows.Scan(&resultId)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Result{
		ID:         resultId,
		RunnerID:   result.RunnerID,
		RaceResult: result.RaceResult,
		Location:   result.Location,
		Position:   result.Position,
		Year:       result.Year,
	}, nil
}

func (rr ResultsRepository) DeleteResult(resultId string) (*models.Result, *models.ResponseError) {
	query := `
		DELETE FROM results
		WHERE id = $1
		RETURNING runner_id, race_result, year
	`

	rows, err := rr.transaction.Query(query, resultId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var runnerId, raceResult string
	var year int
	for rows.Next() {
		err := rows.Scan(&runnerId, &raceResult, &year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Result{
		ID:         resultId,
		RunnerID:   runnerId,
		RaceResult: raceResult,
		Year:       year,
	}, nil
}
