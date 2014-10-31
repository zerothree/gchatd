package protocol

import (
    "errors"
)

const (
    MAX_MSG_LEN = 5000
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
    ACTION_LOGIN_RSP           = "Login_rsp"
    ACTION_HEARTBEAT_RSP       = "HeartBeat_rsp"
    ACTION_UPDATE_USERINFO_RSP = "UpdateUserInfo_rsp"
    ACTION_JOINGROUP_RSP       = "JoinGroup_rsp"
    ACTION_LEAVEGROUP_RSP      = "LeaveGroup_rsp"
    ACTION_GROUPCHAT_RSP       = "GroupChat_rsp"

    // push action
    ACTION_GROUPCHAT_PUSH   = "GroupChat_push"
    ACTION_PRIVATECHAT_PUSH = "PrivateChat_push"

    ERRCODE_SUCCESS               = 0
    ERRCODE_KICKED_SAMEUSER_LOGIN = 1
    ERRMSG_SUCCESS                = "success"
)

var (
    ErrDataNotEnough = errors.New("data not enough")

    SuccLoginRsp          = &LoginRspBean{RspBaseBean: *NewSuccRspBaseBean(ACTION_LOGIN_RSP)}
    SuccHeartBeatRsp      = &HeartBeatRspBean{RspBaseBean: *NewSuccRspBaseBean(ACTION_HEARTBEAT_RSP)}
    SuccUpdateUserInfoRsp = &UpdateUserInfoRspBean{RspBaseBean: *NewSuccRspBaseBean(ACTION_UPDATE_USERINFO_RSP)}
    SuccJoinGroupRsp      = &JoinGroupRspBean{RspBaseBean: *NewSuccRspBaseBean(ACTION_JOINGROUP_RSP)}
    SuccLeaveGroupRsp     = &LeaveGroupRspBean{RspBaseBean: *NewSuccRspBaseBean(ACTION_LEAVEGROUP_RSP)}
    SuccGroupChatRsp      = &GroupChatRspBean{RspBaseBean: *NewSuccRspBaseBean(ACTION_GROUPCHAT_RSP)}
)
