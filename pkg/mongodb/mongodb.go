package mongodb

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	vm "aapi/shared/view_models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetClient(ctx context.Context, config vm.MongoDbConfig) (*mongo.Client, error) {
	log.Default().Printf("GetClient UserName: %v --- pass: %v ---- Url: %v \n", config.UserName, config.Password, config.Url)
	clientOpts := options.Client().ApplyURI(config.Url)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Default().Printf("Unable to connect to Mongo: %v \n", err)
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Default().Printf("Unable to Ping to Mongo: %v \n", err)
		return nil, err
	}
	return client, err
}

func GetClientSRV(ctx context.Context, config vm.MongoDbConfig) (*mongo.Client, error) {
	clientOpts := options.Client().ApplyURI("mongodb+srv://" + config.UserName + ":" + config.Password + "@" + config.Url + "/" + config.Database + "?retryWrites=true&w=majority")

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Default().Println(err)
	}
	return client, err
}

func GetOne(client *mongo.Client, config vm.MongoDbConfig, mgCollection, mongoId string, output interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	if client == nil {
		client, err = GetClient(ctx, config)
		if err != nil {
			log.Default().Printf("Failed to get new client at GetOne :%v \n", err)
			log.Fatal(err)
		}
	}

	collection := client.Database(config.Database).Collection(mgCollection)
	objId, err := primitive.ObjectIDFromHex(mongoId)
	if err != nil {
		log.Default().Printf("GetOne for mongoId(%v) , by mongoId got err: %v \n", mongoId, err)
		return err
	}
	filter := bson.D{{
		Key:   "_id",
		Value: objId,
	}}
	document := collection.FindOne(ctx, filter)
	err = document.Decode(&output)
	if err != nil {
		log.Default().Printf("Insert.err: %v \n", err)
		return err
	}
	return err
}

func GetOneByFilter(client *mongo.Client, config vm.MongoDbConfig, mgCollection string, filter interface{}, output interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	if client == nil {
		client, err = GetClient(ctx, config)
		if err != nil {
			log.Default().Printf("Failed to get new client at GetOne :%v \n", err)
			log.Fatal(err)
		}
	}

	collection := client.Database(config.Database).Collection(mgCollection)
	document := collection.FindOne(ctx, filter)
	err = document.Decode(output)
	if err != nil {
		log.Default().Printf("GetOneByFilter.err: %v \n", err)
		return err
	}

	return err
}

// GetMany return
func GetMany(client *mongo.Client,
	config vm.MongoDbConfig,
	mgCollection string,
	take int64,
	skip int64,
	filter interface{},
	output interface{}) error {
	if output == nil {
		return errors.New("model is empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	if client == nil {
		client, err = GetClient(ctx, config)
		if err != nil {
			log.Default().Printf("Failed to get new client at GetOne :%v \n", err)
			return errors.New("unable to create new mongodb client")
		}
	}

	if skip < 0 {
		skip = 0
	}

	if take == 0 {
		take = 20
	}
	opts := options.Find().SetSkip(skip).SetLimit(take)
	collection := client.Database(config.Database).Collection(mgCollection)
	cursor, err := collection.Find(ctx, filter, opts)

	if err != nil {
		log.Default().Printf("Failed to Find collection:%v \n", err)
		return errors.New("unable to Find collection")
	}

	if err = cursor.All(ctx, output); err != nil {
		log.Default().Printf("Failed to fetch cursor :%v \n", err)
		return errors.New("Failed to fetch cursor of collection " + mgCollection)

	}

	return nil
}

func Insert(client *mongo.Client, config vm.MongoDbConfig, mgCollection string, v interface{}) (interface{}, error) {
	if v == nil {
		return "", errors.New("model is empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	if client == nil {
		client, err = GetClient(ctx, config)
		if err != nil {
			log.Default().Printf("Failed to get new client at GetOne :%v \n", err)
			log.Fatal(err)
		}
	}

	log.Default().Printf("Insert.MongoDbConfig: %#v", config)
	collection := client.Database(config.Database).Collection(mgCollection)
	document, err := StructToMongoDocument(v)
	if err != nil {
		log.Default().Printf("Insert.StructToMongoDocument.err: %v", err)
		return "0", err
	}
	log.Default().Printf("Inserting %v \n", v)
	res, err := collection.InsertOne(ctx, document)
	if err != nil {
		log.Default().Printf("Insert.InsertOne.err: %v \n", err)
		return "0", err
	}
	if res == nil {
		log.Default().Printf("Insert.InsertOne.res is nil: %v \n", err)
		return "0", err
	}
	id := res.InsertedID
	if id == nil {
		log.Default().Printf("failed to extract Id: %#v \n", v)
		return "0", errors.New("insert failed")
	}
	log.Default().Printf("Inserted new id %v  data model:  %v \n", id.(primitive.ObjectID).String(), v)
	return id.(primitive.ObjectID).String(), nil
}

func InsertJson(client *mongo.Client, config vm.MongoDbConfig, mgCollection, jsonString string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	if client == nil {
		client, err = GetClient(ctx, config)
		if err != nil {
			log.Default().Printf("Failed to get new client at GetOne :%v \n", err)
			log.Fatal(err)
		}
	}

	collection := client.Database(config.Database).Collection(mgCollection)

	var m map[string]interface{}
	err = json.Unmarshal([]byte(jsonString), &m)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, m)
	return err
}

func UpdateOne(client *mongo.Client, config vm.MongoDbConfig, mgCollection string, filter interface{}, model interface{}) (int64, error) {
	if model == nil {
		return 0, errors.New("model is empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	if client == nil {
		client, err = GetClient(ctx, config)
		if err != nil {
			log.Default().Printf("Failed to get new client at GetOne :%v \n", err)
			log.Fatal(err)
		}
	}

	collection := client.Database(config.Database).Collection(mgCollection)
	res, err := collection.UpdateOne(ctx, filter, model)

	if err != nil {
		log.Default().Printf("UpdateOne.err: %v \n", err)
		return 0, err
	}
	if res == nil {
		log.Default().Printf("UpdateOne.res is nil: %v \n", err)
		return 0, err
	}

	if res.ModifiedCount == 0 {
		log.Default().Printf("Nothing to update for model: : %#v \n", model)
	} else {
		log.Default().Printf("Updated model: %#v \n with filter %#v\n", model, filter)
	}
	return res.ModifiedCount, err
}

func CloudUpdateOne(config vm.MongoDbConfig, filter interface{}, update interface{}) (interface{}, error) {
	if update == nil {
		return "", errors.New("model is empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := GetClientSRV(ctx, config)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database(config.Database).Collection(config.Collection)

	res, err := collection.UpdateOne(ctx, filter, update)
	if res == nil {
		log.Default().Printf("Err: %#v", err)
		return "", errors.New("Update failed.")
	}
	log.Default().Printf("MatchedCount: %#v ModifiedCount: %v UpsertedCount: %v", res.MatchedCount, res.ModifiedCount, res.UpsertedCount)
	return 1, nil
}

func CloudCollectionCount(config vm.MongoDbConfig, filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := GetClientSRV(ctx, config)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database(config.Database).Collection(config.Collection)

	res, err := collection.CountDocuments(ctx, filter)
	return res, err
}

func Update(client *mongo.Client, config vm.MongoDbConfig, mgCollection, id string, v interface{}) (int64, error) {
	if v == nil {
		return 0, errors.New("model is empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	if client == nil {
		client, err = GetClient(ctx, config)
		if err != nil {
			log.Default().Printf("Failed to get new client at GetOne :%v \n", err)
			log.Fatal(err)
		}
	}

	collection := client.Database(config.Database).Collection(mgCollection)
	document, err := StructToMongoDocument(v)
	if err != nil {
		log.Default().Printf("Update.err: %v \n", err)
		return 0, err
	}
	res, err := collection.UpdateByID(ctx, id, document)
	if err != nil {
		log.Default().Printf("Update.UpdateByID.err: %v \n", err)
		return 0, err
	}
	if res == nil {
		log.Default().Println("Insert.UpdateByID.res is nil")
		return 0, err
	}
	return res.ModifiedCount, err
}

func Delete(client *mongo.Client, config vm.MongoDbConfig, mgCollection, id string) (int64, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	if client == nil {
		client, err = GetClient(ctx, config)
		if err != nil {
			log.Default().Printf("Failed to get new client at GetOne :%v \n", err)
			log.Fatal(err)
		}
	}

	collection := client.Database(config.Database).Collection(mgCollection)
	filter := bson.M{"_id": bson.M{"$eq": id}}
	res, err := collection.DeleteOne(ctx, filter)
	return res.DeletedCount, err
}

func DeleteMany(client *mongo.Client, config vm.MongoDbConfig, mgCollection string, filter bson.M) (int64, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	if client == nil {
		client, err = GetClient(ctx, config)
		if err != nil {
			log.Default().Printf("Failed to get new client at DeleteMany :%v \n", err)
			log.Fatal(err)
		}
	}

	collection := client.Database(config.Database).Collection(mgCollection)
	res, err := collection.DeleteMany(ctx, filter)
	return res.DeletedCount, err
}

func StructToMongoDocument(v interface{}) (*bson.D, error) {
	document := bson.D{}

	if v == nil {
		return nil, errors.New("model is empty")
	}

	data, err := bson.Marshal(v)
	if err != nil {
		log.Default().Printf("StructToMongoDocument. Unable to marshal object: %v \n", v)
		return nil, err
	}

	err = bson.Unmarshal(data, &document)
	if err != nil {
		return nil, err
	}
	return &document, err
}

func MongoDocumentToStruct(doc *bson.D, out interface{}) error {
	if doc == nil {
		return errors.New("mongodb Document is nil")
	}

	var bsonAsByte []byte
	bsonAsByte, err := bson.Marshal(&doc)
	if err != nil {
		panic(err)
	}
	err = bson.Unmarshal(bsonAsByte, &out)
	if err != nil {
		panic(err)
	}
	return err
}
