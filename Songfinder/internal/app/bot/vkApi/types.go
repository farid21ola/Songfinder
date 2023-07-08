package vkApi

type Result struct {
	Resp []Response `json:"response"`
}

type Response struct {
	Name   string `json:"first_name"`
	Status string `json:"status"`
	Audio  Audios `json:"status_audio"`
}

type Audios struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

type Song struct {
	Artist string
	Title  string
}
