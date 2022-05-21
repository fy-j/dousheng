package model

import (
	"dousheng/config"
	"fmt"
	mgo "gopkg.in/mgo.v2"
	"log"
)

var (
	mongoSession *mgo.Session
)

//connect to db
func init() {
	log.Println("connecting db!")
	session, err := getMongoSession()
	if err != nil {
		log.Println("MongoDB init error!")
		log.Panic(err)
		return
	}
	mongoSession = session
	initMaxId()
	fmt.Println("id:")
	fmt.Println(userMaxId)

	log.Println("Database init done!")
}

//if connect error,it will panic
func getMongoSession() (*mgo.Session, error) {
	mgosession, err := mgo.Dial(config.Conf.Mongo.Host)
	if err != nil {
		log.Println("Mongodb dial error!")
		log.Panic(err)
		return nil, err
	}
	mgosession.SetMode(mgo.Monotonic, true)
	mgosession.SetPoolLimit(300)
	myDb := mgosession.DB(config.Conf.Mongo.Name)
	err = myDb.Login(config.Conf.Mongo.User, config.Conf.Mongo.Pwd)
	if err != nil {
		log.Println("Login wrong" + config.Conf.Mongo.User + config.Conf.Mongo.Pwd)
		log.Panic(err)
		return nil, err
	}
	return mgosession, nil
}
