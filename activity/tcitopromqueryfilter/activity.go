package tcitoprommetrics

import (
	"strings"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// New function for the activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	err := metadata.MapToStruct(ctx.Settings(), &Settings{}, true)
	if err != nil {
		return nil, err
	}

	act := &Activity{}

	return act, nil
}

// Activity is an activity that is used to invoke a REST Operation
type Activity struct {
}

// Metadata for the activity
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Create the hash
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return false, err
	}

	logger := ctx.Logger()
	if logger.DebugEnabled() {
		logger.Debugf("Input params: %s", input)
	}

	filter := make(map[string]string)

	tokens := strings.Split(input.Query, ",")
	for _, token := range tokens {
		if token != "" {
			parts := strings.Split(token, "=")
			if len(parts) == 2 {
				filter[parts[0]] = parts[1]
			}
		}
	}

	output := &Output{
		Filter: filter,
	}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}

	return true, nil
}
