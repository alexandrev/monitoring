package tcitoprommetrics

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings for the activity
type Settings struct {
}

// Input for the activity
type Input struct {
	Query string `md:"query"`
}

// Output for the activity
type Output struct {
	Filter map[string]string `md:"filter"`
}

// ToMap for Input
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"query": i.Query,
	}
}

// FromMap for input
func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	i.Query, err = coerce.ToString(values["query"])
	if err != nil {
		return err
	}
	return nil
}

// ToMap conver to object
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"filter": o.Filter,
	}
}

// FromMap convert to object
func (o *Output) FromMap(values map[string]interface{}) error {
	o.Filter, _ = values["filter"].(map[string]string)
	return nil
}
