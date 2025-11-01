package validator

import (
	"fmt"
	"strings"
	"time"

	"github.com/andrewhowdencom/ruf/internal/model"
	"github.com/gorhill/cronexpr"
)

// Validate validates a list of calls and returns a list of errors.
func Validate(calls []*model.Call) []error {
	var errs []error
	for _, call := range calls {
		if err := validateCall(call); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func validateCall(call *model.Call) error {
	var errs []string
	if call.Subject == "" {
		errs = append(errs, "subject is required")
	}
	if call.Content == "" {
		errs = append(errs, "content is required")
	}
	if len(call.Destinations) == 0 {
		errs = append(errs, "at least one destination is required")
	}
	if len(call.Triggers) == 0 {
		errs = append(errs, "at least one trigger is required")
	}

	for _, trigger := range call.Triggers {
		if err := validateTrigger(trigger); err != nil {
			errs = append(errs, err.Error())
		}
	}

	for _, destination := range call.Destinations {
		if err := validateDestination(destination); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation failed for call '%s': %s", call.Subject, strings.Join(errs, ", "))
	}
	return nil
}

func validateTrigger(trigger model.Trigger) error {
	var errs []string
	if trigger.ScheduledAt != (time.Time{}) {
		// No validation needed, as the YAML parser will fail on invalid date/time formats.
	}
	if trigger.Cron != "" {
		if _, err := cronexpr.Parse(trigger.Cron); err != nil {
			errs = append(errs, fmt.Sprintf("invalid cron expression: %s", err))
		}
	}
	if trigger.Delta != "" {
		if _, err := time.ParseDuration(trigger.Delta); err != nil {
			errs = append(errs, fmt.Sprintf("invalid delta: %s", err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("validation failed for trigger: %s", strings.Join(errs, ", "))
	}
	return nil
}

func validateDestination(destination model.Destination) error {
	switch destination.Type {
	case "slack", "email":
		// Valid
	default:
		return fmt.Errorf("invalid destination type: %s", destination.Type)
	}
	return nil
}
