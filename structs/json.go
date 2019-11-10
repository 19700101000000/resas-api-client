package structs

type Prefecture struct {
	PrefCode int    `json:"prefCode"`
	PrefName string `json:"prefName"`
}

type Prefectures struct {
	Message    *string      `json:"message"`
	StatusCode string       `json:"statusCode"`
	Result     []Prefecture `json:"result"`
}

type City struct {
	PrefCode    int    `json:"prefCode"`
	CityCode    string `json:"cityCode"`
	CityName    string `json:"cityName"`
	BigCityFlag string `json:"bigCityFlag"`
}

type Cities struct {
	Message    *string `json:"message"`
	StatusCode string  `json:"statusCode"`
	Result     []City  `json:"result"`
}
