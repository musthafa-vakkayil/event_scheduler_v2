package validator

import (
	"errors"

	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"github.com/musthafa-vakkayil/event_scheduler_v2/utils"
)

func ValidateAPIEventRequest(req models.CreateAPIEventRequest) error {
	// Validation for API Trigger
	if req.Type == "API" && req.ApiMethod != "GET" && req.ApiMethod != "DELETE" && req.ApiPayload == nil {
		return errors.New("api payload is required for methods other than GET and delete")
	}
	return nil
}

func ValidateScheduleEventRequest(req models.CreateScheduledEventRequest) error {
	// Either `RunAtDate` or `RunAfterMins` should be provided
	if (req.RunAtDate == "" && req.RunAfterMins <= 0) || (req.RunAtDate != "" && req.RunAfterMins > 0) {
		return errors.New("provide either 'run_at_this_date' or 'run_after_x_mins', but not both")
	}

	// Validate and parse time if provided
	if req.RunAtDate != "" {
		_, err := utils.ParseAndValidateTime(req.RunAtDate)
		if err != nil {
			return err
		}
	}

	return nil
}
