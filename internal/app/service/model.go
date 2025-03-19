package service

type Location struct {
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}

type keyAndLocation struct {
	Key string
	Loc Location
}

type Port struct {
	Key         string
	Name        string
	City        string
	Country     string
	Alias       []string
	Regions     []string
	Coordinates []float64
	Province    string
	Timezone    string
	Unlocs      []string
	Code        string
}
