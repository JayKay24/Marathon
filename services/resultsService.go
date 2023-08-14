package services

import (
	"marathon-postgresql/models"
	"net/http"
	"time"
)

type ResultsService struct {
	resultsRepository *repositories.ResultsRepository
	runnersRepository *repositories.RunnersRepository
}

func NewResultsService(
	resultsRepo *repositories.ResultsRepository,
	runnersRepo *repositories.RunnersRepository) *ResultsService {
	return &ResultsService{
		resultsRepository: resultsRepo,
		runnersRepository: runnersRepo,
	}
}

func (rs ResultsService) CreateResult(result *models.Result) (*models.Result, *models.ResponseError) {
	if result.RunnerID == "" {
		return nil, &models.ResponseError{
			Message: "Invalid runner id",
			Status:  http.StatusBadRequest,
		}
	}

	if result.RaceResult == "" {
		return nil, &models.ResponseError{
			Message: "Invalid race result",
			Status:  http.StatusBadRequest,
		}
	}

	if result.Location == "" {
		return nil, &models.ResponseError{
			Message: "Invalid Location",
			Status:  http.StatusBadRequest,
		}
	}

	if result.Position < 0 {
		return nil, &models.ResponseError{
			Message: "Invalid Position",
			Status:  http.StatusBadRequest,
		}
	}

	currentYear := time.Now().Year()
	if result.Year < 0 || result.Year > currentYear {
		return nil, &models.ResponseError{
			Message: "Invalid Year",
			Status:  http.StatusBadRequest,
		}
	}

	raceResult, err := parseRaceResult(result.RaceResult)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid race result",
			Status:  http.StatusBadRequest,
		}
	}

	response, responseErr := rs.resultsRepository.CreateResult(result)
	if responseErr != nil {
		return nil, responseErr
	}

	runner, responseErr := rs.runnersRepository.GetRunner(result.RunnerID)
	if responseErr != nil {
		return nil, responseErr
	}

	if runner == nil {
		return nil, &models.ResponseError{
			Message: "Runner not found",
			Status:  http.StatusNotFound,
		}
	}

	// update runner's personal best
	if runner.PersonalBest == "" {
		runner.PersonalBest = result.RaceResult
	} else {
		personalBest, err := parseRaceResult(runner.PersonalBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "Failed to parse personal best",
				Status:  http.StatusInternalServerError,
			}
		}

		if raceResult < personalBest {
			runner.PersonalBest = result.RaceResult
		}
	}

	// update runner's season best
	if result.Year == currentYear {
		if runner.SeasonBest == "" {
			runner.SeasonBest = result.RaceResult
		} else {
			seasonBest, err := parseRaceResult(runner.SeasonBest)
			if err != nil {
				return nil, &models.ResponseError{
					Message: "Failed to parse season best",
					Status:  http.StatusInternalServerError,
				}
			}

			if raceResult < seasonBest {
				runner.SeasonBest = result.RaceResult
			}
		}
	}

	responseErr = rs.resultsRepository.UpdateRunnerResults(runner)
	if responseErr != nil {
		return nil, responseErr
	}

	return response, nil
}

func (rs ResultsService) DeleteResult(resultId string) *models.ResponseError {}

func parseRaceResult(timeString string) (time.Duration, error) {
	return time.ParseDuration(
		timeString[0:2] + "h" +
			timeString[3:5] + "m" +
			timeString[6:8] + "s")
}
