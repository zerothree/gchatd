package protocol

import (
    "encoding/json"
    "fmt"
)

func Marshal(msg interface{}) []byte {
    data, err := json.Marshal(msg)
    if err != nil {
        panic(fmt.Sprintf("Marshal logic err: %s", err))
    }

    return data
}

func UnMarshalReqBase(data []byte, bean *ReqBaseBean) (err error) {
    if data[len(data)-1] != '\n' {
        return ErrDataNotEnough
    }
    err = json.Unmarshal(data, bean)
    return
}

func UnMarshalReq(baseBean *ReqBaseBean, data []byte, bean interface{}) (err error) {
    err = json.Unmarshal(data, bean)
    return
}
