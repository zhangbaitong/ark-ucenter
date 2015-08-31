package common

import (
	"database/sql"
	"fmt"
	 "gopkg.in/mgo.v2" 
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zhangbaitong/go-uuid/uuid"
	"github.com/dlintw/goconf"
	"log"
	"os"
	"time"
)

//the global var of db connection,logger,redis pool.it will be create only once.
var (
	Conn   *sql.DB
	Logger *log.Logger
	pool   *redis.Pool
	DBpool *DbPool
	MPool *MongoPool
)

//Get db connection from mysql
func GetDB() (db *sql.DB) {
	if DBpool == nil {		
		conf, err := goconf.ReadConfigFile("auth.conf")
		if err!=nil {
			fmt.Println(err)
			return nil
		}

		host,_:=conf.GetString("mysql", "host") 
		port,_:=conf.GetInt("mysql", "port") 
		user,_:=conf.GetString("mysql", "user") 
		password,_:=conf.GetString("mysql", "password") 
		db,_:=conf.GetString("mysql", "db") 
		db_server:=fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",user,password,host,port,db)

		DBpool = CreateDbPool(20, "mysql", db_server, true)
	}

	conn, err := DBpool.GetConn()
	if err != nil {
		fmt.Println(err)
		return nil
	}


	return conn
}

func FreeDB(db *sql.DB) {
	//DBpool.Mu.Lock()
	DBpool.PutConn(db)
	//DBpool.Mu.Unlock()
}

func GetDBInfo() (info string) {
	if DBpool == nil {
		return "DB not init\r\n"
	}

	return fmt.Sprintf("PoolSize=%d;MaxPoolSize=%d", DBpool.PoolSize, DBpool.MaxPoolSize)
}

func GetSession() (Session *mgo.Session) {
	if MPool == nil {
		conf, err := goconf.ReadConfigFile("auth.conf")
		if err!=nil {
			fmt.Println(err)
			return nil
		}

		host,_:=conf.GetString("mongodb", "host") 
		port,_:=conf.GetInt("mongodb", "port") 
		mongo_server:=fmt.Sprintf("%s:%d",host,port)

		MPool = CreateMongoPool(20,mongo_server)
	}

	Session, err := MPool.GetSession()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return Session	
}

func FreeSession(Session *mgo.Session) {
	MPool.PutSession(Session)
}

//get uuid lik a227cedf-e806-11e4-8666-3c075419d855
func GetUID() string {
	return uuid.NewUUID().String()
}

//get app logger
func Log() *log.Logger {
	if Logger == nil {
		Logger = log.New(os.Stdout, "AT-Resource : ", log.Ldate|log.Ltime|log.Lshortfile)
		Logger.Print("logger init success ...")
	}
	return Logger
}

//get redis connection pool
func GetRedisPool() *redis.Pool {
	if pool == nil {
		pool = &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {				
				conf, err := goconf.ReadConfigFile("auth.conf")
				if err!=nil {
					fmt.Println(err)
					return nil, err
				}

				host,_:=conf.GetString("redis", "host") 
				port,_:=conf.GetInt("redis", "port") 
				redis_server:=fmt.Sprintf("%s:%d",host,port)

				c, err := redis.Dial("tcp", redis_server)
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
	}
	return pool
}
