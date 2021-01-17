package keyboards

// The PerformanceChoices represent the range of performance results that can be
// achieved by a user when they attempted a route.
var PerformanceChoices = []Choice{
	{Label: "Flashed", Action: "flashed"},
	{Label: "Succeeded", Action: "succeeded"},
	{Label: "Failed", Action: "failed"},
}
