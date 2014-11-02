package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "net"
    "reflect"
    "sync"
    "time"

    "./protocol"
)

type Session struct {
    conn              *net.TCPConn
    user              *protocol.UserInfoBean
    groups            map[string]*Group
    logined           bool
    rsps              chan []byte
    quit              chan bool
    kickedBySameLogin bool
    recvBuf           []byte
    bufDataLen        int
    rspsLock          sync.Mutex
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
                userMgr.RemoveUser(s.user.Uid)
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

    s.recvBuf = make([]byte, protocol.MAX_MSG_LEN, protocol.MAX_MSG_LEN)

    action, data, err := s.recvMsg()
    if action != protocol.ACTION_LOGIN {
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
            var bean protocol.LoginRspBean
            bean.Action = protocol.ACTION_LOGIN_RSP
            bean.ErrCode = protocol.ERRCODE_KICKED_SAMEUSER_LOGIN
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

    var bean protocol.ReqBaseBean
    for {
        if s.bufDataLen >= protocol.MAX_MSG_LEN { // assume include only one msg in buf
            err = errors.New("msg length greater than MAX_MSG_LEN")
            return
        }

        var n int
        n, err = s.conn.Read(s.recvBuf[s.bufDataLen:])
        if err != nil {
            return
        }
        if n == 0 {
            err = ErrConnClosedByPeer
            return
        }
        s.bufDataLen += n

        var beanLen int
        beanLen, err = protocol.UnMarshalReqBase(s.recvBuf[:s.bufDataLen], &bean)
        if err == nil {
            if s.bufDataLen > beanLen {
                remainLen := s.bufDataLen - beanLen
                copy(s.recvBuf[:remainLen], s.recvBuf[beanLen:s.bufDataLen])
                s.bufDataLen = remainLen
            }
            return bean.Action, s.recvBuf[:beanLen], nil
        } else if err == protocol.ErrDataNotEnough {
            continue
        } else {
            return
        }
    }

    return
}

func (s *Session) formatMsg(msg interface{}) []byte {
    data, err := json.Marshal(msg)
    if err != nil {
        panic(fmt.Sprintf("formatMsg logic err: %s", err))
    }

    return data
}

func (c *Session) WriteMsg(msg []byte) {
    c.rspsLock.Lock()
    defer c.rspsLock.Unlock()

    //drop rsp to avoid block when rsps is full
    if msg != nil && len(c.rsps) >= MAX_RSPCHAN_LEN {
        return
    }

    c.rsps <- msg
}

func (c *Session) Kick() {
    c.WriteMsg(nil)
}
