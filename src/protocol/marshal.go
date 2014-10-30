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

func UnMarshal(data []byte) (bean *ReqBaseBean, err error) {
    if data[len(data)-1] != '\n' {
        return
    }
    bean = &ReqBaseBean{}
    err = json.Unmarshal(data, bean)
    return
}
