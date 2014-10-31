package main

import (
    "errors"
)

const (
    MAX_RSPCHAN_LEN = 100
)

var (
    dummyError = errors.New("error")

    ErrConnClosedByPeer = errors.New("connection is closed by peer")
)
