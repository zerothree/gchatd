package main

import (
    "sync"
)

type GroupMgr struct {
    groups map[string]*Group
    lock   sync.RWMutex
}

func NewGroupMgr() *GroupMgr {
    g := &GroupMgr{}
    g.groups = make(map[string]*Group)
    return g
}

/*
func (g *GroupMgr) AddUser(uid string, rsps chan []byte, groupId string) {
    g.GetGroup(groupId).AddUser(uid, rsps)
}

func (g *GroupMgr) RemoveUser(uid string, groupId string) {
    g.GetGroup(groupId).RemoveUser(uid)
}

func (g *GroupMgr) PushMsgToGroup(groupId string, rsp []byte, excludeUid string) {
    g.GetGroup(groupId).PushMsgToGroup(rsp, excludeUid)
}*/

/**
* @return *Group
 */
func (g *GroupMgr) GetGroup(groupId string) *Group {
    g.lock.RLock()
    group, ok := g.groups[groupId]
    if !ok {
        g.lock.RUnlock()
        g.lock.Lock()
        group, ok = g.groups[groupId]
        if !ok { // double check
            group = NewGroup(groupId)
            g.groups[groupId] = group
        }
        g.lock.Unlock()
    } else {
        g.lock.RUnlock()
    }
    return group
}

// group
type Group struct {
    groupId string
    users   map[string]chan []byte
    lock    sync.RWMutex
}

func NewGroup(groupId string) *Group {
    g := &Group{}
    g.groupId = groupId
    g.users = make(map[string]chan []byte)
    return g
}

func (g *Group) AddUser(uid string, rsps chan []byte) {
    g.lock.Lock()
    defer g.lock.Unlock()

    g.users[uid] = rsps
}

func (g *Group) RemoveUser(uid string) {
    g.lock.Lock()
    defer g.lock.Unlock()

    delete(g.users, uid)
}

func (g *Group) PushMsgToGroup(rsp []byte, excludeUid string) {
    g.lock.RLock()
    defer g.lock.RUnlock()

    for uid, rsps := range g.users {
        if uid != excludeUid {
            if len(rsps) < MAX_RSPCHAN_LEN { // if rsps is full we drop rsp to avoid block
                rsps <- rsp
            }
        }
    }
}
