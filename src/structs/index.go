package structs

type UserRequestObject struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	Scrobbles int32  `json:"scrobbles"`
}

type GenerateRequest struct {
	ID           string      `json:"id"`
	Theme        string      `json:"theme"`
	Story        bool        `json:"story"`
	HideUsername bool        `json:"hide_username"`
	Data         interface{} `json:"data"`
	ReturnImage  bool        `json:"return_image"`
}
