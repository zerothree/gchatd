package main

import (
    "errors"
)

const (
    MAX_MSG_LEN     = 5000
    MAX_RSPCHAN_LEN = 100
)

const (
    // request action
    ACTION_LOGIN           = "Login"
    ACTION_HEARTBEAT       = "HeartBeat"
    ACTION_UPDATE_USERINFO = "UpdateUserInfo"
    ACTION_JOINGROUP       = "JoinGroup"
    ACTION_LEAVEGROUP      = "LeaveGroup"
    ACTION_GROUPCHAT       = "GroupChat"

    // response action
    ACTION_LOGIN_RSP           = "login_rsp"
    ACTION_HEARTBEAT_RSP       = "heartbeat_rsp"
    ACTION_UPDATE_USERINFO_RSP = "update_userinfo_rsp"
    ACTION_JOINGROUP_RSP       = "joingroup_rsp"
    ACTION_LEAVEGROUP_RSP      = "leavegroup_rsp"
    ACTION_GROUPCHAT_RSP       = "groupchat_rsp"

    // push action
    ACTION_GROUPCHAT_PUSH   = "groupchat_push"
    ACTION_PRIVATECHAT_PUSH = "privatechat_push"

    ERRCODE_SUCCESS               = 0
    ERRCODE_KICKED_SAMEUSER_LOGIN = 1
    ERRMSG_SUCCESS                = "success"
)

var (
    dummyError            = errors.New("error")
    succLoginRsp          = &LoginRspBean{RspBaseBean: *NewSuccRspBaseBean(ACTION_LOGIN_RSP)}
    succHeartBeatRsp      = &HeartBeatRspBean{RspBaseBean: *NewSuccRspBaseBean(ACTION_HEARTBEAT_RSP)}
    succUpdateUserInfoRsp = &UpdateUserInfoRspBean{RspBaseBean: *NewSuccRspBaseBean(ACTION_UPDATE_USERINFO_RSP)}
    succJoinGroupRsp      = &JoinGroupRspBean{RspBaseBean: *NewSuccRspBaseBean(ACTION_JOINGROUP_RSP)}
    succLeaveGroupRsp     = &LeaveGroupRspBean{RspBaseBean: *NewSuccRspBaseBean(ACTION_LEAVEGROUP_RSP)}
    succGroupChatRsp      = &GroupChatRspBean{RspBaseBean: *NewSuccRspBaseBean(ACTION_GROUPCHAT_RSP)}
)
