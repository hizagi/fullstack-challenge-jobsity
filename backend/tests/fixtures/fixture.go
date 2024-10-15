package fixtures

import (
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
)

type Fixture struct {
	db *mongo.Database
}

func fixture(f Fixture, seedMethodName string) {
	// Get the reflect value of the method
	m := reflect.ValueOf(f).MethodByName(seedMethodName)
	// Exit if the method doesn't exist
	if !m.IsValid() {
		log.Fatal("No method called ", seedMethodName)
	}
	// Execute the method
	log.Println("Seeding", seedMethodName, "...")
	m.Call(nil)
	log.Println("Seed", seedMethodName, "succeed")
}

func Execute(db *mongo.Database, seedMethodNames ...string) {
	f := Fixture{db: db}

	seedType := reflect.TypeOf(f)

	if len(seedMethodNames) == 0 {
		log.Println("Running all seeder...")

		for i := 0; i < seedType.NumMethod(); i++ {
			method := seedType.Method(i)
			fixture(f, method.Name)
		}
	}

	for _, item := range seedMethodNames {
		fixture(f, item)
	}
}
