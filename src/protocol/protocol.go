package protocol

import (
    "fmt"
)

// request base struct
type ReqBaseBean struct {
    Action string `json:"action"`
}

// reponse base struct
type RspBaseBean struct {
    Action  string `json:"action"`
    ErrCode int    `json:"errcode"`
    ErrMsg  string `json:"errmsg"`
}

func NewSuccRspBaseBean(action string) *RspBaseBean {
    return &RspBaseBean{Action: action, ErrCode: ERRCODE_SUCCESS, ErrMsg: ERRMSG_SUCCESS}
}

// push base struct
type PushBaseBean struct {
    Action string `json:"action"`
}

// userinfo base struct
type UserInfoBean struct {
    Uid   string `json:"uid"`
    Nick  string `json:"nick"`
    Extra string `json:"extra"`
}

func (u *UserInfoBean) String() string {
    return fmt.Sprintf("%s|%s", u.Uid, u.Nick)
}

// login request
type LoginBean struct {
    ReqBaseBean
    UserInfo UserInfoBean `json:"userinfo"`
    Token    string       `json:"token"`
    Groups   []string     `json:"groups"`
    Friends  []string     `json:"friends"`
    Ignores  []string     `json:"ignores"`
}

// login response
type LoginRspBean struct {
    RspBaseBean
}

// hearbbeat request
type HeartBeatBean struct {
    ReqBaseBean
    Ping string `json:"ping"`
}

// heartbeat response
type HeartBeatRspBean struct {
    RspBaseBean
}

// update userinfo request
type UpdateUserInfoBean struct {
    ReqBaseBean
    Nick  string `json"nick"`
    Extra string `json:"extra"`
}

// update userinfo request
type UpdateUserInfoRspBean struct {
    RspBaseBean
}

// joingroup request
type JoinGroupBean struct {
    ReqBaseBean
    Group string `json:"group"`
}

// joingroup response
type JoinGroupRspBean struct {
    RspBaseBean
}

// leavegroup request
type LeaveGroupBean struct {
    ReqBaseBean
    Group string `json:"group"`
}

// leavegroup response
type LeaveGroupRspBean struct {
    RspBaseBean
}

// add friend
type AddFriendBean struct {
    ReqBaseBean
    Friend string `json:"friend"`
}

// add friend rsp
type AddFriendRspBean struct {
    RspBaseBean
}

// add ignore
type AddIgnoreBean struct {
    ReqBaseBean
    Ignore string `json:"ignore"`
}

// add ignore response
type AddIgnoreRspBean struct {
    RspBaseBean
}

// groupchat request
type GroupChatBean struct {
    ReqBaseBean
    Group string `json:"group"`
    Msg   string `json:"msg"`
}

// groupchat response
type GroupChatRspBean struct {
    RspBaseBean
}

// private chat
type PrivateChatBean struct {
    ReqBaseBean
    To  string `json:"to"`
    Msg string `json:"msg"`
}

// private chat rsp
type PrivateChatRspBean struct {
    RspBaseBean
}

// groupchat push
type GroupChatPushBean struct {
    PushBaseBean
    From  UserInfoBean `json:"from"`
    Group string       `json:"group"`
    Msg   string       `json:"msg"`
    Ts    int          `json:"ts`
}

// private chat push
type PrivateChatPushBean struct {
    PushBaseBean
    From UserInfoBean `json:"from"`
    Msg  string       `json:"msg"`
    Ts   int          `json:"ts"`
}
