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

/*func (g *GroupMgr) AddUser(uid string, rsps chan []byte, groupId string) {
    g.GetGroup(groupId).AddUser(uid, rsps)
}

func (g *GroupMgr) RemoveUser(uid string, groupId string) {
    g.GetGroup(groupId).RemoveUser(uid)
}

func (g *GroupMgr) PushMsgToGroup(groupId string, rsp []byte, excludeUid string) {
    g.GetGroup(groupId).PushMsgToGroup(rsp, excludeUid)
}*/

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

func (g *GroupMgr) delGroup(groupId string) {
    g.lock.Lock()
    defer g.lock.Unlock()

    delete(g.groups, groupId)
}

// group
type Group struct {
    groupId string
    users   map[string]*Session
    lock    sync.RWMutex
}

func NewGroup(groupId string) *Group {
    g := &Group{}
    g.groupId = groupId
    g.users = make(map[string]*Session)
    return g
}

func (g *Group) AddUser(uid string, session *Session) {
    g.lock.Lock()
    defer g.lock.Unlock()

    g.users[uid] = session
}

func (g *Group) RemoveUser(uid string) {
    g.lock.Lock()
    defer g.lock.Unlock()

    delete(g.users, uid)
    if len(g.users) == 0 {
        groupMgr.delGroup(g.groupId)
    }
}

func (g *Group) PushMsgToGroup(rsp []byte, excludeUid string) {
    g.lock.RLock()
    defer g.lock.RUnlock()

    for uid, session := range g.users {
        if uid != excludeUid {
            session.WriteMsg(rsp)
        }
    }
}
