package model

type Bucket struct {
	Bucket          string `json:"bucket"`
	Id              string `json:"id"`
	Name            string `json:"name"`
	ContentType     string `json:"contentType"`
	ContentLanguage string `json:"contentLanguage"`
	Kind            string `json:"kind"`
}
