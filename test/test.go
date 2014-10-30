package main

import (
    "../src/protocol"
    "log"
)

func main() {

}

func login(uid, token, nick, groups string) {
    var bean procotol.LoginBean
    bean.Action = ACTION_LOGIN
    bean.UserInfo.Uid = uid
    bean.UserInfo.Nick = nick
    bean.Token = token
    bean.Groups = groups

}
