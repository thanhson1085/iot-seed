package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/redis.v3"
	"io"
	"time"
)

// OplogQuery describes a query you'd like to perform on the oplog.
type OplogQuery struct {
	Session *mgo.Session
	Query   bson.M
	Timeout time.Duration
}

// Oplog is a deserialization of the fields present in an oplog entry.
type Oplog struct {
	Timestamp    bson.MongoTimestamp `bson:"ts"`
	HistoryID    int64               `bson:"h"`
	MongoVersion int                 `bson:"v"`
	Operation    string              `bson:"op"`
	Namespace    string              `bson:"ns"`
	Object       bson.M              `bson:"o"`
	//QueryObject  bson.M              `bson:"o2"`
}

// LastTime gets the timestamp of the last operation in the oplog.
// The return value can be used for construting queries on the "ts" oplog field.
func LastTime(session *mgo.Session) bson.MongoTimestamp {
	var member Oplog
	session.DB("local").C("oplog.$main").Find(nil).Sort("-$natural").One(&member)
	return member.Timestamp
}

// Tail starts the tailing of the oplog.
// Entries matching the configured query are published to channel passed to the function.
func (query *OplogQuery) Tail(logs chan Oplog, done chan bool) {
	db := query.Session.DB("local")
	collection := db.C("oplog.$main")
	iter := collection.Find(query.Query).LogReplay().Tail(query.Timeout)
	var log Oplog
	for iter.Next(&log) {
		logs <- log
	}
	iter.Close()
	done <- true
}

func printlog(buffer io.Writer, logs chan Oplog) {
	// Print logs from an oplog channel to a buffer
	for log := range logs {
		id := log.Object["_id"].(bson.ObjectId).Hex()
		fmt.Fprintf(buffer, "%s|%s|%s\n", log.Namespace, log.Operation, id)
	}
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		fmt.Printf("%s", err)
	}
	defer session.Close()

	logs := make(chan Oplog)
	done := make(chan bool)
	last := LastTime(session)

	q := OplogQuery{session, bson.M{"ts": bson.M{"$gt": last}, "ns": "TailTest.test"}, time.Second * 10}

	/*
		db := session.DB("TailTest")
		defer db.DropDatabase()
		coll := db.C("test")
		for i := 0; i < 5; i++ {
			id := bson.NewObjectId()
			if err := coll.Insert(bson.M{"name": "test_0", "_id": id}); err != nil {
				fmt.Println(err)
			}
			//fmt.Fprintf(&buffer, "TailTest.test|i|%s\n", id.Hex())
		}
	*/
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.191:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	for {
		go q.Tail(logs, done)
		for log := range logs {
			id := log.Object["_id"].(bson.ObjectId).Hex()
			client.Publish("mongdb_oblog", log.Namespace)
			fmt.Printf("%s|%s|%s\n", log.Namespace, log.Operation, id)
		}
	}
}
