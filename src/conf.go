package main

type Conf struct {
    Port        int `json:"port"`
    RecvTimeout int `json:"recv_time_out"` //seconds
    SendTimeout int `json:"send_time_out"` //seconds
}
