package services

import (
	"marathon-postgresql/models"
	"net/http"
)

type RunnersService struct {
	runnersRepository *repositories.RunnersRepository
	resultsRepository *repositories.ResultsRepository
}

func NewRunnersService(
	runnersRepository *repositories.RunnersRepository,
	resultsRepository *repositories.ResultsRepository) *RunnersService {
	return &RunnersService{
		runnersRepository: runnersRepository,
		resultsRepository: resultsRepository,
	}
}

func (rs RunnersService) CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError) {
	responseErr := validateRunner(runner)
	if responseErr != nil {
		return nil, responseErr
	}

	return rs.runnersRepository.CreateRunner(runner)
}

func (rs RunnersService) UpdateRunner(runner *models.Runner) *models.ResponseError {
	responseErr := validateRunnerId(runner.ID)
	if responseErr != nil {
		return responseErr
	}

	responseErr = validateRunner(runner)
	if responseErr != nil {
		return responseErr
	}

	rs.runnersRepository.UpdateRunner(runner)
	return nil
}

func (rs RunnersService) DeleteRunner(runnerId string) *models.ResponseError {
	responseError := validateRunnerId(runnerId)
	if responseError != nil {
		return responseError
	}

	rs.runnersRepository.DeleteRunner(runnerId)
	return nil
}

func (rs RunnersService) GetRunner(runnerId string) (*models.Runner, *models.ResponseError) {
	responseErr := validateRunnerId(runnerId)
	if responseErr != nil {
		return nil, responseErr
	}

	runner, responseErr := rs.runnersRepository.GetRunner(runnerId)
	if responseErr != nil {
		return nil, responseErr
	}

	results, responseErr := rs.resultsRepository.getAllRunnerResults(runnerId)
	if responseErr != nil {
		return nil, responseErr
	}

	runner.Result = results
	return runner, nil
}

func (rs RunnersService) GetRunnersBatch(country string, year string) ([]*models.Runner, *models.ResponseError) {
}

func validateRunner(runner *models.Runner) *models.ResponseError {
	if runner.FirstName == "" {
		return &models.ResponseError{
			Message: "Invalid first name",
			Status:  http.StatusBadRequest,
		}
	}

	if runner.LastName == "" {
		return &models.ResponseError{
			Message: "Invalid last name",
			Status:  http.StatusBadRequest,
		}
	}

	if runner.Age <= 16 || runner.Age > 125 {
		return &models.ResponseError{
			Message: "Invalid age",
			Status:  http.StatusBadRequest,
		}
	}

	if runner.Country == "" {
		return &models.ResponseError{
			Message: "Invalid Country",
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}

func validateRunnerId(runnerId string) *models.ResponseError {
	if runnerId == "" {
		return &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}
