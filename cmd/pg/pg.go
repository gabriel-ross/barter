package main

import (
	"fmt"
	"time"
)

var PROJECT_ID = "personal-gabrielross"
var PORT = ":8080"

type Demo struct {
	ID       string             `firestore:"id"`
	Balances map[string]float64 `firestore:"balances"`
}

func main() {
	fmt.Println(time.Now())

	// var err error
	// ctx := context.TODO()
	// godotenv.Load(".env")

	// // Instantiate dependencies
	// fsClient, err := firestore.NewClient(ctx, PROJECT_ID)
	// if err != nil {
	// 	log.Fatalf("Failed to create client: %s", err)
	// }
	// defer fsClient.Close()

	// // One operation, key doesn't exist
	// data1 := Demo{
	// 	ID:    fsClient.Collection("demo").NewDoc().ID,
	// 	Balances: map[string]float64{},
	// }

	// _, err = fsClient.Collection("demo").Doc(data1.ID).Set(ctx, data1)
	// if err != nil {
	// 	log.Println("error storing data1 ", err)
	// }
	// _, err = fsClient.Collection("demo").Doc(data1.ID).Update(ctx, []firestore.Update{
	// 	{Path: "balances.dollars", Value: firestore.Increment(10)},
	// })
	// if err != nil {
	// 	log.Println("error incrementing data1 ", err)
	// }
	// dsnap, err := fsClient.Collection("demo").Doc(data1.ID).Get(ctx)
	// _, err = fsClient.Collection("demo").Doc(data1.ID).Set(ctx, data1)
	// if err != nil {
	// 	log.Println("error getting data1 ", err)
	// }
	// var c Demo
	// dsnap.DataTo(&c)
	// bytes, err := json.Marshal(c)
	// if err != nil {
	// 	log.Fatal("err marshaling response to JSON", err)
	// 	return
	// }
	// fmt.Println("data1: ", string(bytes))

	// // One operation, key does exist, testing if dot notation works with maps
	// data2 := Demo{
	// 	ID: fsClient.Collection("demo").NewDoc().ID,
	// 	Balances: map[string]float64{
	// 		"dollars": 10,
	// 	},
	// }
	// _, err = fsClient.Collection("demo").Doc(data2.ID).Set(ctx, data2)
	// if err != nil {
	// 	log.Println("error storing data2 ", err)
	// }
	// _, err = fsClient.Collection("demo").Doc(data2.ID).Update(ctx, []firestore.Update{
	// 	{Path: "balances.dollars", Value: firestore.Increment(10)},
	// })
	// if err != nil {
	// 	log.Println("error incrementing data2 ", err)
	// }
	// dsnap, err = fsClient.Collection("demo").Doc(data2.ID).Get(ctx)
	// _, err = fsClient.Collection("demo").Doc(data2.ID).Set(ctx, data2)
	// if err != nil {
	// 	log.Println("error getting data2 ", err)
	// }
	// var b Demo
	// dsnap.DataTo(&b)
	// bytes, err = json.Marshal(b)
	// if err != nil {
	// 	log.Fatal("err marshaling response to JSON", err)
	// 	return
	// }
	// fmt.Println("data2: ", string(bytes))

	// // Two operations, testing unmarshaling into struct, updating, and storing
	// data3 := Demo{
	// 	ID: fsClient.Collection("demo").NewDoc().ID,
	// 	Balances: map[string]float64{
	// 		"dollars": 100,
	// 		"apples":  3,
	// 	},
	// }
	// _, err = fsClient.Collection("demo").Doc(data3.ID).Set(ctx, data3)
	// if err != nil {
	// 	log.Println("error storing data3 ", err)
	// }
	// dsnap, err = fsClient.Collection("demo").Doc(data3.ID).Get(ctx)
	// _, err = fsClient.Collection("demo").Doc(data3.ID).Set(ctx, data3)
	// if err != nil {
	// 	log.Println("error getting data3 ", err)
	// }
	// var a Demo
	// dsnap.DataTo(&a)
	// bytes, err = json.Marshal(a)
	// if err != nil {
	// 	log.Fatal("err marshaling response to JSON", err)
	// 	return
	// }
	// fmt.Println("data3: ", string(bytes))
}
