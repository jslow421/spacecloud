package models

type UpcomingRocketLaunchesApiResponse struct {
	ValidAuth bool                         `json:"valid_auth"`
	Count     int                          `json:"count"`
	Limit     int                          `json:"limit"`
	Total     int                          `json:"total"`
	LastPage  int                          `json:"last_page"`
	Result    []UpcomingRocketLaunchResult `json:"result"`
}

type UpcomingRocketLaunchResult struct {
	Id                  int64                         `json:"id"`
	CosparId            string                        `json:"cospar_id"`
	SortDate            string                        `json:"sort_date"`
	Name                string                        `json:"name"`
	Provider            UpcomingRocketLaunchProvider  `json:"provider"`
	Vehicle             UpcomingRocketLaunchVehicle   `json:"vehicle"`
	Pad                 UpcomingRocketLaunchPad       `json:"pad"`
	Missions            []UpcomingRocketLaunchMission `json:"missions"`
	MissionDescription  string                        `json:"mission_description"`
	LaunchDescription   string                        `json:"launch_description"`
	WindowOpen          string                        `json:"win_open"`  // date time
	T0                  string                        `json:"t0"`        // date time
	WindowClose         string                        `json:"win_close"` // date time
	DateString          string                        `json:"date_str"`
	Tags                []Tags                        `json:"tags"`
	Slug                string                        `json:"slug"`
	WeatherSummary      string                        `json:"weather_summary"`
	WeatherTemp         float32                       `json:"weather_temp"`
	WeatherCondition    string                        `json:"weather_condition"`
	WeatherWindSpeedMph float32                       `json:"weather_wind_mph"`
	WeatherIcon         string                        `json:"weather_icon"`
	WeatherUpdated      string                        `json:"weather_updated"` // date time
	QuickText           string                        `json:"quicktext"`
	IsSuborbital        bool                          `json:"suborbital"`
	Modified            string                        `json:"modified"` // date time
}

type UpcomingRocketLaunchProvider struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type UpcomingRocketLaunchVehicle struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CompanyId int    `json:"company_id"`
	Slug      string `json:"slug"`
}

type UpcomingRocketLaunchPad struct {
	Id       int                             `json:"id"`
	Name     string                          `json:"name"`
	Location UpcomingRocketLaunchPadLocation `json:"location"`
}

type UpcomingRocketLaunchPadLocation struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	State     string `json:"state"`
	StateName string `json:"statename"`
	Country   string `json:"country"`
	Slug      string `json:"slug"`
}

type UpcomingRocketLaunchMission struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Tags struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}
