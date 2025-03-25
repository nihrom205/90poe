package domain

type Port struct {
	key         string
	name        string
	city        string
	country     string
	alias       []string
	regions     []string
	coordinates []float64
	province    string
	timezone    string
	unlocs      []string
	code        string
}

type NewPortData struct {
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

func NewPort(data NewPortData) (Port, error) {
	return Port{
		key:         data.Key,
		name:        data.Name,
		city:        data.City,
		country:     data.Country,
		alias:       data.Alias,
		regions:     data.Regions,
		coordinates: data.Coordinates,
		province:    data.Province,
		timezone:    data.Timezone,
		unlocs:      data.Unlocs,
		code:        data.Code,
	}, nil
}

func (p Port) Key() string {
	return p.key
}

func (p Port) Name() string {
	return p.name
}

func (p Port) City() string {
	return p.city
}

func (p Port) Country() string {
	return p.country
}

func (p Port) Alias() []string {
	return p.alias
}

func (p Port) Regions() []string {
	return p.regions
}

func (p Port) Coordinates() []float64 {
	return p.coordinates
}

func (p Port) Province() string {
	return p.province
}

func (p Port) Timezone() string {
	return p.timezone
}

func (p Port) Unlocs() []string {
	return p.unlocs
}

func (p Port) Code() string {
	return p.code
}
