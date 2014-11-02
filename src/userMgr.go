package main

import (
    "sync"
)

type UserMgr struct {
    users map[string]*Session
    lock  sync.RWMutex
}

func NewUserMgr() *UserMgr {
    m := &UserMgr{}
    m.users = make(map[string]*Session)
    return m
}

func (m *UserMgr) AddUser(uid string, session *Session) {
    m.lock.Lock()
    defer m.lock.Unlock()

    oldUserSession, ok := m.users[uid]
    if ok {
        // kickout old user which uid is the same
        oldUserSession.Kick()
    }
    m.users[uid] = session
}

func (m *UserMgr) RemoveUser(uid string) {
    m.lock.Lock()
    defer m.lock.Unlock()

    delete(m.users, uid)
}

func (m *UserMgr) PushMsgToUser(uid string, rsp []byte) {
    m.lock.RLock()
    defer m.lock.RUnlock()

    userSession, ok := m.users[uid]
    if ok {
        userSession.WriteMsg(rsp)
    }
}
