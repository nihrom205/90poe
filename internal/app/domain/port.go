package domain

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

//type NewPortData struct {
//	Key         string
//	Name        string
//	City        string
//	Country     string
//	Alias       []string
//	Regions     []string
//	Coordinates []float64
//	Province    string
//	Timezone    string
//	Unlocs      []string
//	Code        string
//}

//func NewPort(data NewPortData) (Port, error) {
//	return Port{
//		key:         data.Key,
//		name:        data.Name,
//		city:        data.City,
//		country:     data.Country,
//		alias:       data.Alias,
//		regions:     data.Regions,
//		coordinates: data.Coordinates,
//		province:    data.Province,
//		timezone:    data.Timezone,
//		unlocs:      data.Unlocs,
//		code:        data.Code,
//	}, nil
//}
//
//func (p Port) GetName() string {
//	return p.name
//}
