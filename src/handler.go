package main

import (
    "encoding/json"
    "log"
)

func (s *Session) HandleLogin(data []byte) error {
    if s.logined {
        log.Printf("%s %s already logined, now login again", s.conn.RemoteAddr(), s.user)
        return dummyError
    }

    var bean LoginBean
    err := json.Unmarshal(data, &bean)
    if err != nil {
        log.Printf("%s Unmarshal login msg error. %s.  msg:%s", s.conn.RemoteAddr(), err, string(data))
        return err
    }

    // check uid and token
    // to do ...

    s.user = &bean.UserInfo
    s.logined = true

    log.Printf("%s %s login. data: %s", s.conn.RemoteAddr(), s.user, string(data))

    s.rsps <- s.formatMsg(succLoginRsp)

    // add the uesr into groupMgr and userchanMgr
    for _, groupId := range bean.Groups {
        group := groupMgr.GetGroup(groupId)
        group.AddUser(s.user.Uid, s.rsps)
        s.groups[groupId] = group
    }
    userchanMgr.AddUser(s.user.Uid, s.rsps)
    return nil
}

func (s *Session) HandleHeartBeat(data []byte) error {
    s.rsps <- s.formatMsg(succHeartBeatRsp)
    return nil
}

func (s *Session) HandleUpdateUserInfo(data []byte) error {
    var bean UpdateUserInfoBean
    err := json.Unmarshal(data, &bean)
    if err != nil {
        log.Printf("%s %s Unmarshal UpdateUserInfo msg error. %s.  msg:%s", s.conn.RemoteAddr(), s.user, err, string(data))
        return err
    }

    s.user.Nick = bean.Nick
    s.user.Extra = bean.Extra

    s.rsps <- s.formatMsg(succUpdateUserInfoRsp)
    return nil
}

func (s *Session) HandleJoinGroup(data []byte) error {
    var bean JoinGroupBean
    err := json.Unmarshal(data, &bean)
    if err != nil {
        log.Printf("%s %s Unmarshal JoinGroup msg error. %s.  msg:%s", s.conn.RemoteAddr(), s.user, err, string(data))
        return err
    }

    if _, ok := s.groups[bean.Group]; !ok {
        group := groupMgr.GetGroup(bean.Group)
        group.AddUser(s.user.Uid, s.rsps)
        s.groups[bean.Group] = group
    }

    s.rsps <- s.formatMsg(succJoinGroupRsp)
    return nil
}

func (s *Session) HandleLeaveGroup(data []byte) error {
    var bean LeaveGroupBean
    err := json.Unmarshal(data, &bean)
    if err != nil {
        log.Printf("%s %s Unmarshal LeaveGroup msg error. %s.  msg:%s", s.conn.RemoteAddr(), s.user, err, string(data))
        return err
    }

    group, ok := s.groups[bean.Group]
    if ok {
        group.RemoveUser(s.user.Uid)
        delete(s.groups, bean.Group)
    }

    s.rsps <- s.formatMsg(succLeaveGroupRsp)
    return nil
}

func (s *Session) HandleGroupChat(data []byte) error {
    var bean GroupChatBean
    err := json.Unmarshal(data, &bean)
    if err != nil {
        log.Printf("%s %s Unmarshal groupchat msg error. %s.  msg:%s", s.conn.RemoteAddr(), s.user, err, string(data))
        return err
    }

    group, ok := s.groups[bean.Group]
    if !ok {
        log.Printf("%s %s push group msg which group not be registered", s.conn.RemoteAddr(), s.user)
        return dummyError
    }

    log.Printf("%s %s group[%s] chat: %s", s.conn.RemoteAddr(), s.user, bean.Group, bean.Msg)

    s.rsps <- s.formatMsg(succGroupChatRsp)

    var p GroupChatPushBean
    p.Action = ACTION_GROUPCHAT_PUSH
    p.From = *s.user
    p.Group = bean.Group
    p.Msg = bean.Msg

    group.PushMsgToGroup(s.formatMsg(&p), s.user.Uid)
    return nil
}
