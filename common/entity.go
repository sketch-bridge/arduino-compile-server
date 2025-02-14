package common

import "time"

type RequestParameters struct {
	ProjectId string
}

type Project struct {
	Id        string    `firestore:"id"`
	Name      string    `firestore:"name"`
	Code      string    `firestore:"code"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}
