package database

import (
	"cloud.google.com/go/firestore"
	"context"
	"sketch-bridge/arduino-compile-server/common"
)

func GetProject(ctx context.Context, firestoreClient *firestore.Client, projectId string) (*common.Project, error) {
	doc, err := firestoreClient.Collection("versions").Doc("v1").Collection("projects").Doc(projectId).Get(ctx)
	if err != nil {
		return nil, err
	}

	var project common.Project
	project.Id = doc.Ref.ID
	if err := doc.DataTo(&project); err != nil {
		return nil, err
	}

	return &project, nil
}
