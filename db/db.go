package db

import (
	"crypto/tls"
	"log"
	"net"
	"os"

	"gopkg.in/mgo.v2"
)

type Connection interface {
	Close()
	DB() *mgo.Database
}

type conn struct {
	session *mgo.Session
}

func NewConnection() Connection {
	var c conn
	var err error
	dialInfo := mgo.DialInfo{
		Addrs: []string{
			"cluster0-shard-00-00.1gguw.mongodb.net:27017",
			"cluster0-shard-00-01.1gguw.mongodb.net:27017",
			"cluster0-shard-00-02.1gguw.mongodb.net:27017",
		},
		Username: "fireninja",
		Password: os.Getenv("ELDEN_HUB_MONGODB_PASS"),
	}

	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig) // add TLS config
		return conn, err
	}

	c.session, err = mgo.DialWithInfo(&dialInfo)
	if err != nil {
		log.Panicln(err.Error())
	}

	return &c
}

func (c *conn) Close() {
	c.session.Close()
}

func (c *conn) DB() *mgo.Database {
	return c.session.DB("eldenhub-forum")
}
