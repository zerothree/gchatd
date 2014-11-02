package main

import (
    "encoding/json"
    "io/ioutil"
    "time"
)

var (
    userMgr  = NewUserMgr()
    groupMgr = NewGroupMgr()
    conf     = &Conf{}
)

const (
    defaultConfFile    = "../config/server.conf"
    defaultPort        = 8090
    defaultRecvTimeout = 10
    defaultSendTimeout = 2
)

func main() {
    data, err := ioutil.ReadFile(defaultConfFile)
    if err != nil {
        panic(err)
    }
    err = json.Unmarshal(data, conf)
    if err != nil {
        panic(err)
    }

    s := Server{port: conf.Port}
    s.start()

    time.Sleep(time.Hour)

    println("done")
}
