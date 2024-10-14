package fixtures

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hizagi/fullstack-challenge-jobsity/backend/api/generated"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const taskCollection = "tasks"

const PredefinedObjectID = "670c76c9c04a05f0b0c571e5"

const TaskFixtureMethod = "TaskFixture"

func (f Fixture) TaskFixture() {
	statuses := []string{string(generated.TaskStatusComplete), string(generated.TaskStatusIncomplete), string(generated.TaskStatusInProgress)}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	objectID, err := primitive.ObjectIDFromHex(PredefinedObjectID)
	if err != nil {
		log.Fatalf("error decoding objectid during taskfixture: %v", err)
	}

	for i := 0; i < 50; i++ {
		status := statuses[r.Intn(len(statuses))]

		document := bson.M{
			"title":   fmt.Sprintf("title_%d", i),
			"content": fmt.Sprintf("content_%d", i),
			"status":  status,
		}

		if i == 1 {
			document["_id"] = objectID
		}

		_, err := f.db.Collection(taskCollection).InsertOne(context.Background(), document)
		if err != nil {
			log.Fatalf("error during taskfixture: %v", err)
		}
	}
}
