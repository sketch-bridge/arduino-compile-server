package common

import "time"

type RequestParameters struct {
	ProjectId string
}

type Library struct {
	Name     string
	Version  string
	Sentence string
}

type Project struct {
	Id        string    `firestore:"id"`
	Name      string    `firestore:"name"`
	Code      string    `firestore:"code"`
	Fqbn      string    `firestore:"fqbn"`
	Libraries []Library `firestore:"libraries"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}
