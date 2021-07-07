package tcitoprommetrics

import (
	"strconv"
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

//Method to format TCI Stats to Prom
func FormatToPrometheus(mList MetricList) string {
	var sb strings.Builder
	for _, metric := range mList.Metrics {
		sb.WriteString("# HELP " + metric.Name + " " + metric.Description + "\n")
		sb.WriteString("# TYPE " + metric.Name + " " + metric.Type + "\n")

		for _, sample := range metric.Samples {
			sb.WriteString(metric.Name + "{")
			first := true
			for key, value := range sample.Labels {
				if first == true {
					first = false
				} else {
					sb.WriteString(",")
				}
				sb.WriteString(key + "=" + "\"" + value + "\"")
			}

			sb.WriteString("} ")
			sb.WriteString(strconv.FormatFloat(sample.Value, 'f', -1, 64) + "\n")
		}
	}

	return sb.String()

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

	list := MetricList{}

	for _, metric := range input.Metrics {

		appID := metric.App.AppID
		appName := metric.App.AppName
		appType := metric.App.AppType

		for _, instanceMetric := range metric.AppInstanceMetrics {

			for _, flow := range instanceMetric.AppInstanceMetrics.Flows {

				pLabels := make(map[string]string)
				pLabels["flowName"] = flow.FlowName
				pLabels["appName"] = appName
				pLabels["appID"] = appID
				pLabels["appType"] = appType
				pLabels["appInstance"] = instanceMetric.AppInstance

				if appType != "" && appName != "" && appID != "" {

					pMetricCount := list.Create("flow_execution_count", "Total number of times the flow is started, completed, or failed", "counter")
					pLabelsAdditionalCompleted := make(map[string]string)
					for k, v := range pLabels {
						pLabelsAdditionalCompleted[k] = v
					}
					pLabelsAdditionalCompleted["status"] = "Completed"
					pMetricCount.Add(pLabelsAdditionalCompleted, float64(flow.Completed))
					pLabelsAdditionalFailed := make(map[string]string)
					for k, v := range pLabels {
						pLabelsAdditionalFailed[k] = v
					}
					pLabelsAdditionalFailed["status"] = "Failed"
					pMetricCount.Add(pLabelsAdditionalFailed, float64(flow.Failed))
					pLabelsAdditionalStarted := make(map[string]string)
					for k, v := range pLabels {
						pLabelsAdditionalStarted[k] = v
					}
					pLabelsAdditionalStarted["status"] = "Started"
					pMetricCount.Add(pLabelsAdditionalStarted, float64(flow.Started))

					pMetricDuration := list.Create("flow_duration_msec", "Total time (in ms) taken by the flow for successful completion or failure", "gauge")
					pMetricDuration.Add(pLabels, float64(flow.AvgExecTime))

				}

			}

		}

		pLabelsApp := make(map[string]string)

		pLabelsApp["appName"] = appName
		pLabelsApp["appType"] = appType
		pLabelsApp["appID"] = appID

		logger.Debugf("Test2 appID: %s", appID)

		for _, appMetric := range metric.AppMetrics {

			logger.Debugf("Test2 appID - 1 : %s", appMetric.InstanceId)

			pLabelsApp["appInstance"] = appMetric.InstanceId

			logger.Debugf("Test appID - 11 : %v", appMetric)

			for _, a := range appMetric.TciAppInstancesCPU {

				logger.Debugf("Test appID - 2 : %s", appMetric.InstanceId)

				println(a.Labels.Status)
				if a.Labels.Status != "" {
					pAppCPUUsage := list.Create(a.Labels.Status+"_app_cpu_usage", "CPU Usage Percentage", "gauge")
					pAppCPUUsage.Add(pLabelsApp, float64(a.Value))
				}
			}

			for _, a := range appMetric.TciAppInstancesMemory {

				println(a.Labels.Status)
				if a.Labels.Status != "" {
					pAppMemoryUsage := list.Create(a.Labels.Status+"_app_memory_used", "Memory Used Percentage", "gauge")
					pAppMemoryUsage.Add(pLabelsApp, float64(a.Value))
				}
			}
		}

	}

	data := FormatToPrometheus(list)

	output := &Output{
		Data: data,
	}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}

	return true, nil
}
