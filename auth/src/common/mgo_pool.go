package common

import (
	"fmt"
    "time"
    "sync"
	"errors"
     "gopkg.in/mgo.v2" 
)

type MongoPool struct {
    ServerName string
    MaxPoolSize int
    PoolSize int
    Mu          sync.Mutex
    Sessions       chan *mgo.Session
}

//建立Mongo接池，传入参数为：连接池最大连接数，服务名的连接数，是主库连接池还是备库连接池
func CreateMongoPool(maxPoolsize int, strServerName string) *MongoPool {

        dbPool := &MongoPool{MaxPoolSize: maxPoolsize, ServerName: strServerName}
        flag := make(chan bool, dbPool.MaxPoolSize/2)
        go func() {
        for i := 0; i < dbPool.MaxPoolSize/2; i++ {
            session, err := mgo.Dial(strServerName)
            if err != nil {
                fmt.Println(err)
            }
            dbPool.PutSession(session)
            flag <- true
        }
        }()

        for i := 0; i < dbPool.MaxPoolSize/2; i++ {
        <-flag
        }
     dbPool.PoolSize= dbPool.MaxPoolSize/2

        return dbPool
}

//从连接池中获取连接
func (this *MongoPool) GetSession() (*mgo.Session, error) {
    if this.PoolSize <this.MaxPoolSize && len(this.Sessions) == 0 {
        go func() {
            this.Mu.Lock()
            if(this.PoolSize >=this.MaxPoolSize) {
                return
            }

            for i := 0; i < this.MaxPoolSize/2; i++ {
                session, err := mgo.Dial(this.ServerName)
                if err != nil {
                    fmt.Println("连接数据库失败")
                    fmt.Println(err)
                }
                this.PutSession(session)
            }
            this.PoolSize=this.MaxPoolSize
            this.Mu.Unlock()
        }()
    }

    this.Mu.Lock()
    defer this.Mu.Unlock()
    //判断是否能在3秒内获取连接，如果不能就报错
    select {
    //读取通道里的数据库连接，如果读不到就返回报错
    case connChan, ok := <-this.Sessions:
        {
            if ok {
                return connChan, nil
            } else {
                return nil, errors.New("数据库连接获取异常，可能已经被关闭！")
            }
        }
    //如果被阻塞三秒仍没有获取到连接，则就返回错误
    case <-time.After(time.Second * 3):
        return nil, errors.New("获取MongoDB数据库连接超时！")
    }
}

//把连接放入连接池中
func (this *MongoPool) PutSession(session *mgo.Session) {
    this.Mu.Lock()
    defer this.Mu.Unlock()
    if this.Sessions == nil {
        this.Sessions = make(chan *mgo.Session, this.MaxPoolSize)
    }
    if len(this.Sessions) >= this.MaxPoolSize {
        session.Close()
        return
    }
    this.Sessions <- session
}