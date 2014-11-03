package main

import (
    "time"

    "./log"
    "./protocol"
)

func (s *Session) HandleLogin(data []byte) error {
    if s.logined {
        log.Notice("%s %s already logined, now login again", s.conn.RemoteAddr(), s.user)
        return dummyError
    }

    var bean protocol.LoginBean
    err := protocol.UnmarshalReq(data, &bean)
    if err != nil {
        log.Notice("%s Unmarshal login msg error. %s.  msg:%s", s.conn.RemoteAddr(), err, string(data))
        return err
    }

    // check uid and token
    // to do ...

    s.user = &bean.UserInfo
    s.logined = true

    log.Info("%s %s login. data: %s", s.conn.RemoteAddr(), s.user, string(data))

    s.WriteMsg(protocol.Marshal(protocol.SuccLoginRsp))

    // add the uesr into groupMgr and userchanMgr
    for _, groupId := range bean.Groups {
        group := groupMgr.GetGroup(groupId)
        group.AddUser(s.user.Uid, s)
        s.groups[groupId] = group
    }
    userMgr.AddUser(s.user.Uid, s)
    return nil
}

func (s *Session) HandleHeartBeat(data []byte) error {
    s.WriteMsg(protocol.Marshal(protocol.SuccHeartBeatRsp))
    return nil
}

func (s *Session) HandleUpdateUserInfo(data []byte) error {
    var bean protocol.UpdateUserInfoBean
    err := protocol.UnmarshalReq(data, &bean)
    if err != nil {
        log.Notice("%s %s Unmarshal UpdateUserInfo msg error. %s.  msg:%s", s.conn.RemoteAddr(), s.user, err, string(data))
        return err
    }

    s.user.Nick = bean.Nick
    s.user.Extra = bean.Extra

    s.WriteMsg(protocol.Marshal(protocol.SuccUpdateUserInfoRsp))
    return nil
}

func (s *Session) HandleJoinGroup(data []byte) error {
    var bean protocol.JoinGroupBean
    err := protocol.UnmarshalReq(data, &bean)
    if err != nil {
        log.Notice("%s %s Unmarshal JoinGroup msg error. %s.  msg:%s", s.conn.RemoteAddr(), s.user, err, string(data))
        return err
    }

    if _, ok := s.groups[bean.Group]; !ok {
        group := groupMgr.GetGroup(bean.Group)
        group.AddUser(s.user.Uid, s)
        s.groups[bean.Group] = group
    }

    s.WriteMsg(protocol.Marshal(protocol.SuccJoinGroupRsp))
    return nil
}

func (s *Session) HandleLeaveGroup(data []byte) error {
    var bean protocol.LeaveGroupBean
    err := protocol.UnmarshalReq(data, &bean)
    if err != nil {
        log.Notice("%s %s Unmarshal LeaveGroup msg error. %s.  msg:%s", s.conn.RemoteAddr(), s.user, err, string(data))
        return err
    }

    group, ok := s.groups[bean.Group]
    if ok {
        group.RemoveUser(s.user.Uid)
        delete(s.groups, bean.Group)
    }

    s.WriteMsg(protocol.Marshal(protocol.SuccLeaveGroupRsp))
    return nil
}

func (s *Session) HandleGroupChat(data []byte) error {
    var bean protocol.GroupChatBean
    err := protocol.UnmarshalReq(data, &bean)
    if err != nil {
        log.Notice("%s %s Unmarshal groupchat msg error. %s.  msg:%s", s.conn.RemoteAddr(), s.user, err, string(data))
        return err
    }

    group, ok := s.groups[bean.Group]
    if !ok {
        log.Notice("%s %s push group msg which group not be registered", s.conn.RemoteAddr(), s.user)
        return dummyError
    }

    log.Info("%s %s group[%s] chat: %s", s.conn.RemoteAddr(), s.user, bean.Group, bean.Msg)

    s.WriteMsg(protocol.Marshal(protocol.SuccGroupChatRsp))

    var push protocol.GroupChatPushBean
    push.Action = protocol.ACTION_GROUPCHAT_PUSH
    push.From = *s.user
    push.Group = bean.Group
    push.Msg = bean.Msg
    push.Ts = int(time.Now().Unix())

    group.PushMsgToGroup(protocol.Marshal(&push), s.user.Uid)
    return nil
}
