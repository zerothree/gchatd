package main

import (
    "sync"
)

type UserChanMgr struct {
    users map[string]chan []byte
    lock  sync.RWMutex
}

func NewUserChanMgr() *UserChanMgr {
    m := &UserChanMgr{}
    m.users = make(map[string]chan []byte)
    return m
}

func (m *UserChanMgr) AddUser(uid string, rspChan chan []byte) {
    m.lock.Lock()
    defer m.lock.Unlock()

    oldUserRspChan, ok := m.users[uid]
    if ok {
        // kickout old user which uid is the same
        oldUserRspChan <- nil
    }
    m.users[uid] = rspChan
}

func (m *UserChanMgr) RemoveUser(uid string) {
    m.lock.Lock()
    defer m.lock.Unlock()

    delete(m.users, uid)
}

func (m *UserChanMgr) PushMsgToUser(uid string, rsp []byte) {
    m.lock.RLock()
    defer m.lock.RUnlock()

    userChan, ok := m.users[uid]
    if ok {
        userChan <- rsp
    }
}
