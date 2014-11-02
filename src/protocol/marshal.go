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

func UnMarshalReqBase(data []byte, bean *ReqBaseBean) (beanLen int, err error) {
    for i := 0; i < len(data); i++ {
        if data[i] == '\n' {
            beanLen = i + 1
            break
        }
    }
    if beanLen == 0 {
        err = ErrDataNotEnough
        return
    }

    err = json.Unmarshal(data[:beanLen], bean)
    return
}

func UnmarshalReq(data []byte, bean interface{}) (err error) {
    err = json.Unmarshal(data, bean)
    return
}
