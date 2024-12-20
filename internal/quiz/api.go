package quiz

import (
	"github.com/courage173/quiz-api/internal/errors"

	"github.com/courage173/quiz-api/internal/models"

	"github.com/courage173/quiz-api/pkg/log"

	routing "github.com/go-ozzo/ozzo-routing/v2"
)

func RegisterHandlers(rg *routing.RouteGroup, service Service, logger log.Logger) {
	rg.Get("", getQuiz(service, logger))
	rg.Get("/submission/<username>", getUserSubmission(service, logger))
	rg.Post("/submit", submitQuiz(service, logger))

}

func getQuiz(service Service, logger log.Logger) routing.Handler {
	return func(c *routing.Context) error {
		response := service.GetQuestions()

		if len(response) == 0 {
			logger.With(c.Request.Context()).Infof("No questions found")
			return errors.NotFound("")
		}
		return c.Write(response)
	}
}

func submitQuiz(service Service, logger log.Logger) routing.Handler {
	return func(c *routing.Context) error {
		var req models.Submission

		if err := c.Read(&req); err != nil {
			logger.With(c.Request.Context()).Errorf("Invalid request: %v", err)
			return errors.BadRequest("")
		}

		if err := req.Validate(); err != nil {
			logger.With(c.Request.Context()).Errorf("Invalid request: %v", err)
			return err
		}

		response, err := service.SubmitQuiz(req)

		if err != nil {
			logger.With(c.Request.Context()).Errorf("Error submitting quiz: %v", err)
			return errors.BadRequest(err.Error())
		}

		return c.Write(response)
	}
}

func getUserSubmission(service Service, logger log.Logger) routing.Handler {
	return func(c *routing.Context) error {
		username := c.Param("username")
		response, err := service.GetUserSubmission(username)

		if err != nil {
			logger.With(c.Request.Context()).Errorf("Error getting user submission: %v", err)
			return errors.BadRequest(err.Error())
		}

		return c.Write(response)
	}
}
