package main

import (
    "../src/protocol"
    "log"
    "net"
)

var (
    serverAddr = "localhost:8090"
)

func main() {

    log.Println("done")
}

type Client struct {
    conn net.Conn
    user *protocol.UserInfoBean
}

func NewClient(uid, nick string) *Client {
    c := &Client{}
    var err error
    c.conn, err = net.Dial("tcp4", serverAddr)
    if err != nil {
        panic(err)
    }
    c.user = &protocol.UserInfoBean{Uid: uid, Nick: nick}
    return c
}

func (c *Client) login(token string, groups []string) {
    var bean protocol.LoginBean
    bean.Action = protocol.ACTION_LOGIN
    bean.UserInfo = *c.user
    bean.Token = token
    bean.Groups = groups

}
