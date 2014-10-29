package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "net"
    "reflect"
    "time"
)

type Session struct {
    conn              *net.TCPConn
    user              *UserInfoBean
    groups            map[string]*Group
    logined           bool
    rsps              chan []byte
    quit              chan bool
    kickedBySameLogin bool
    recvBuf           []byte
    bufDataLen        int
}

func (s *Session) init() {
    s.rsps = make(chan []byte, 100)
    s.quit = make(chan bool)
    s.groups = make(map[string]*Group)
}

func (s *Session) start() {
    s.init()
    go s.recvRoutine()
    go s.sendRoutine()
}

func (s *Session) recvRoutine() {
    defer func() {
        // remove user from groupMgr and userchanMgr
        if s.logined {
            if !s.kickedBySameLogin {
                for _, group := range s.groups {
                    group.RemoveUser(s.user.Uid)
                }
                userchanMgr.RemoveUser(s.user.Uid)
                log.Printf("%s %s closing connection", s.conn.RemoteAddr(), s.user)
            } else {
                log.Printf("%s %s closing connection kicked because of same user login", s.conn.RemoteAddr(), s.user)
            }
        } else {
            log.Printf("%s closing connection no login", s.conn.RemoteAddr())
        }

        close(s.rsps)
        s.conn.Close()
        <-s.quit
    }()

    s.recvBuf = make([]byte, MAX_MSG_LEN, MAX_MSG_LEN)

    action, data, err := s.recvMsg()
    if action != ACTION_LOGIN {
        log.Printf("%s first msg's Type must be ACTION_LOGIN. curr action: %s, msg:%s", s.conn.RemoteAddr(), action, string(data))
        return
    }
    err = s.HandleLogin(data)
    if err != nil {
        return
    }

    for {
        action, data, err = s.recvMsg()
        if err != nil {
            log.Printf("%s %s recvMsg err %s", s.conn.RemoteAddr(), s.user, err)
            return
        }
        //log.Printf("%s %s msg: %s", s.conn.RemoteAddr(), s.user, string(data))

        err = s.handleMsg(action, data)
        if err != nil {
            return
        }
    }
}

func (s *Session) handleMsg(action string, data []byte) error {
    v := reflect.ValueOf(s)
    method := v.MethodByName("Handle" + action)
    if !method.IsValid() || method.Kind() != reflect.Func {
        log.Printf("%s %s not support action:%s. msg:%s", s.conn.RemoteAddr(), s.user, action, string(data))
        return dummyError
    }

    return method.Call([]reflect.Value{reflect.ValueOf(data)})[0].Interface().(error)
}

func (s *Session) sendRoutine() {
    for {
        rsp, ok := <-s.rsps
        if !ok { // rsps is closed by recvRoutine
            break
        }
        if rsp == nil { // kickout msg is sent from userchanMgr.AddUser()
            s.kickedBySameLogin = true
            var bean LoginRspBean
            bean.Action = ACTION_LOGIN_RSP
            bean.ErrCode = ERRCODE_KICKED_SAMEUSER_LOGIN
            bean.ErrMsg = fmt.Sprintf("kicked because same user login on other device")
            s.conn.SetWriteDeadline(time.Now().Add(time.Duration(conf.SendTimeout) * time.Second))
            s.conn.Write(s.formatMsg(&bean))
            s.conn.Close()
            break
        }
        s.conn.SetWriteDeadline(time.Now().Add(time.Duration(conf.SendTimeout) * time.Second))
        _, err := s.conn.Write(rsp)
        if err != nil {
            log.Printf("%s %s conn.Write err: %s", s.conn.RemoteAddr(), s.user, err)
            break
        }
    }
    s.quit <- true
}

func (s *Session) recvMsg() (action string, data []byte, err error) {
    s.conn.SetReadDeadline(time.Now().Add(time.Duration(conf.RecvTimeout) * time.Second))

    s.bufDataLen = 0
    for {
        var n int
        n, err = s.conn.Read(s.recvBuf[s.bufDataLen:])
        if err != nil {
            return
        }
        if n == 0 {
            err = errors.New("connection is closed by peer")
            return
        }
        s.bufDataLen += n

        if s.bufDataLen == MAX_MSG_LEN { // assume include only one msg in buf
            err = errors.New("msg length greater than MAX_MSG_LEN")
            return
        }
        if s.recvBuf[s.bufDataLen-1] == '\n' {
            break
        }
    }

    var baseBean ReqBaseBean
    err = json.Unmarshal(s.recvBuf[:s.bufDataLen], &baseBean)
    if err != nil {
        return
    }
    return baseBean.Action, s.recvBuf[:s.bufDataLen], nil
    //s.handleMsg(baseBean.Action, s.recvBuf[:s.bufDataLen])

    /*	offset := 0
    	for i := 0; i < s.bufDataLen; i++ {
    		if s.recvBuf[i] == '\n' {
    			var baseBean ReqBaseBean
    			err = json.Unmarshal(s.recvBuf[:i], &baseBean)
    			if err != nil {
    				return err
    			}
    			s.handleMsg(baseBean.Action, s.recvBuf[offset:i])
    			offset = i + 1
    		}
    	}*/

    return
}

func (s *Session) formatMsg(msg interface{}) []byte {
    data, err := json.Marshal(msg)
    if err != nil {
        panic(fmt.Sprintf("formatMsg logic err: %s", err))
    }

    return data
}
