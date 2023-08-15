package repositories

import (
	"database/sql"
	"marathon-postgresql/models"
	"net/http"
)

type RunnersRepository struct {
	dbHandler   *sql.DB
	transaction *sql.Tx
}

func NewRunnersRepository(dbHandler *sql.DB) *RunnersRepository {
	return &RunnersRepository{
		dbHandler: dbHandler,
	}
}

func (rr RunnersRepository) CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError) {
	query := `INSERT INTO runners(first_name, last_name, age, country) VALUES ($1, $2, $3, $4) RETURNING id`
	rows, err := rr.dbHandler.Query(query, runner.FirstName, runner.LastName, runner.Age, runner.Country)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var runnerId string
	for rows.Next() {
		err := rows.Scan(&runnerId)
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

	return &models.Runner{
		ID:        runnerId,
		FirstName: runner.FirstName,
		LastName:  runner.LastName,
		Age:       runner.Age,
		IsActive:  true,
		Country:   runner.Country,
	}, nil
}

func (rr RunnersRepository) UpdateRunner(runner *models.Runner) *models.ResponseError {
	query := `UPDATE runners SET
		first_name = $1
		last_name = $2
		age = $3
		country = $4
		WHERE id = $5
	`
	res, err := rr.dbHandler.Exec(query, runner.FirstName, runner.LastName, runner.Age, runner.Country, runner.ID)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	if rowsAffected == 0 {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (rr RunnersRepository) DeleteRunner(runnerId string) *models.ResponseError {}

func (rr RunnersRepository) GetRunner(runnerId string) (*models.Runner, *models.ResponseError) {}

func (rr RunnersRepository) GetRunnersByCountry(country string) ([]*models.Runner, *models.ResponseError) {
}

func (rr RunnersRepository) GetRunnersByYear(year int) ([]*models.Runner, *models.ResponseError) {}

func (rr RunnersRepository) GetAllRunners() ([]*models.Runner, *models.ResponseError) {}
