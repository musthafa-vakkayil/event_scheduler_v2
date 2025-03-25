package validator

import (
	"errors"
	"time"

	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
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
	if (req.RunAtDate == nil && req.RunAfterMins <= 0) || (req.RunAtDate != nil && req.RunAfterMins > 0) {
		return errors.New("provide either 'run_at_this_date' or 'run_after_x_mins', but not both")
	}

	// If `RunAtDate` is provided, ensure it's in the future
	if req.RunAtDate != nil {
		if req.RunAtDate.Before(time.Now()) {
			return errors.New("run_at_this_date cannot be in the past")
		}
	}

	return nil
}
