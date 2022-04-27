package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trainer struct {
	Name string
	Age int
	City string
}

func main() {
	fmt.Println("Hello Go!")
	fmt.Println("Start to test MongoDB...")
	testMongoDB()
	log.Printf("End to test MongoDB.")
}

func testMongoDB()  {
	// set options
	clientOption := options.Client().ApplyURI("mongodb://localhost:27017")
	// connect to mongo
	client, err := mongo.Connect(context.TODO(), clientOption)
	// check err
	if err != nil {
		log.Fatal(err)
	}
	// use ping to check the connection
	err = client.Ping(context.TODO(), nil)
	// check err
	if err != nil {
		log.Fatal(err)
	}
	// info
	fmt.Println("Connected to MongoDB successfully.")
	// find collections
	collection := client.Database("test").Collection("trainers")
	// CRUD - C
	zhangsan := Trainer{"zhansan", 10, "Shanghai"}
	lisi := Trainer{"lisi", 10, "Beijing"}
	wangwu := Trainer{"wangwu", 15, "Hangzhou"}
	insertOneResult, err := collection.InsertOne(context.TODO(), zhangsan)
	// check err
	if err != nil {
		log.Fatal(err)
	}
	// insert one result
	fmt.Println("Inserted a single document:", insertOneResult.InsertedID)
	//
	trainers := []interface{}{lisi, wangwu}
	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	//
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a multiple document:", insertManyResult.InsertedIDs)
	// CRUD - U
	filter := bson.D{{"name", "zhansan"}}
	update := bson.D{{"$inc", bson.D{{"age", 1}}}}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	// update again
	updateResult, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// CRUD - R
	var result Trainer
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)
	// CRUD - R - Multiple
	// Pass these options to the Find method
	findOption := options.Find()
	findOption.SetLimit(2)
	// Here's an array in which you can store the decoded documents
	var manyResult []*Trainer
	// Passing bson.D{} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{}, findOption)
	if (err != nil) {
		log.Fatal(err)
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents on at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem Trainer
		err = cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		manyResult = append(manyResult, &elem)
	}
	if err = cur.Err(); err != nil {
		log.Fatal(err)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())
	fmt.Printf("Found multiple documents (array of pointers): %+v\n", manyResult)

	// CRUD - D
	deleteMany, err := collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteMany.DeletedCount)


	// close connection
	err = client.Disconnect(context.TODO())
	// check err
	if err != nil {
		log.Fatal(err)
	}
	// closed info
	fmt.Println("The connection to MongoDB has already been closed.")

}
