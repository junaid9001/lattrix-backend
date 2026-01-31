package models

type PlanLimits struct {
	MaxApis     int64
	MinIntervel int
}

var PlanRules = map[PlanType]PlanLimits{
	PlanFree: {
		MaxApis:     10,
		MinIntervel: 300,
	},
	PlanPro: {
		MaxApis:     50,
		MinIntervel: 60,
	},
	PlanAgency: {
		MaxApis:     200,
		MinIntervel: 30,
	},
}
