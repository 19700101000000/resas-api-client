package structs

type Prefecture struct {
	PrefCode int8 `json:"prefCode"`
	PrefName string `json:"prefName"`
}

type Prefectures struct {
	Message *string `json:"message"`
	StatusCode string `json:"statusCode"`
	Result []Prefecture `json:"result"`
}
