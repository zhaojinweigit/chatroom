// demo for websocket chatroom
//base on  gorilla websocket,use redis to save history

package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

type clientInfo struct {
	Name  string
	Photo string
	Conn  *websocket.Conn
}

type clientMsg struct {
	Name  string `json:"name"`
	Photo string `json:"photo"`
	Msg   string `json:"msg"`
	Time  string `json:"time"`
}

var maxChanNum int = 100
var addr string = "0.0.0.0:10328"
var historyKey = "history_" + addr
var RedisPool *redis.Pool

//use redis to save history
var isHistory bool = false

var broadcastChan chan clientMsg = make(chan clientMsg, maxChanNum)
var addClientchan chan clientInfo = make(chan clientInfo, maxChanNum)
var clientMap map[*websocket.Conn]clientInfo = make(map[*websocket.Conn]clientInfo)
var clientPhoto map[string]string = make(map[string]string)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	clientPhoto["zjw"] = "/image/zjw.jpg"

	if isHistory {
		initRedis()
	}

}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("start server!")
	go broadcast()

	http.HandleFunc("/login", login)
	http.HandleFunc("/history", history)
	http.Handle("/", http.FileServer(http.Dir("./html")))
	http.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("image"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.ListenAndServe(addr, nil)
	log.Println("stop server!")
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form["name"]) <= 0 {
		log.Println("dont have login name")
		return
	}
	addClient(w, r)
}

func addClient(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	var ci clientInfo
	ci.Name = r.Form["name"][0]
	ci.Conn = c
	ci.Photo = clientPhoto[strings.ToLower(ci.Name)]
	if ci.Photo == "" {
		ci.Photo = "/image/default.jpg"
	}

	log.Println("login name="+r.Form["name"][0], ",ip="+r.RemoteAddr)

	addClientchan <- ci

	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println(clientMap[c].Name, "log out,read err:", err, ",mt=", mt)
			break
		}

		var cm clientMsg
		cm.Name = clientMap[c].Name
		cm.Photo = clientMap[c].Photo
		cm.Msg = string(message)
		cm.Time = time.Now().Format("2006-01-02 15:04:05")

		broadcastChan <- cm
	}
}

func broadcast() {

	var conn redis.Conn
	if isHistory {
		conn = RedisPool.Get()
		defer conn.Close()
	}

	for {
		select {
		case cm := <-broadcastChan:
			log.Println(cm.Name, "say:", cm.Msg)
			sendmsg, err := json.Marshal(cm)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			if isHistory {
				_, err = conn.Do("lpush", historyKey, sendmsg)
				if err != nil {
					log.Println(err.Error())
				}
			}

			for k, v := range clientMap {

				err = k.WriteMessage(websocket.TextMessage, []byte(sendmsg))
				if err != nil {
					log.Println("write to ", v.Name, "err:", err)
					delete(clientMap, k)
				}
			}
		case ci := <-addClientchan:
			clientMap[ci.Conn] = ci
		}
	}
}

func initRedis() {
	var redisAddr string = "127.0.0.1:6379"
	RedisPool = &redis.Pool{
		MaxIdle:     5,
		MaxActive:   20,
		IdleTimeout: 600 * time.Second,
		Wait:        false,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", redisAddr, 500*time.Millisecond, 1*time.Second, 1*time.Second)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	log.Println("init redis pool!")
}

func history(w http.ResponseWriter, r *http.Request) {
	conn := RedisPool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("lrange", historyKey, 0, 20))
	if err != nil && err != redis.ErrNil {
		log.Println(err.Error())
		return
	}

	var historyList string = ""
	for len(values) > 0 {
		var temp string
		values, err = redis.Scan(values, &temp)
		if err != nil {
			log.Println(err.Error())
			return
		}

		historyList = "," + temp + historyList
	}
	if historyList != "" {
		historyList = historyList[1:len(historyList)]
		historyList = "[" + historyList + "]"
	} else {
		historyList = "[]"
	}
	w.Write([]byte(historyList))
}
