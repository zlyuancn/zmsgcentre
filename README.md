# zmsgcentre
###### 消息中心

## 获得zmsgcentre
`go get -u github.com/zlyuancn/zmsgcentre`

## 使用zmsgcentre

```go
package main

import (
    "fmt"
    "github.com/zlyuancn/zmsgcentre"
)

func main() {
    a := zmsgcentre.NewMsgCentre()

    testtopic := "test"
    a.AddReceiver(zmsgcentre.NewReceiver(testtopic, 100, func(msg *zmsgcentre.Message) {
        fmt.Println("R1", msg.ID(), msg.Body)
    }))
    a.AddReceiver(zmsgcentre.NewReceiver(testtopic, 10, func(msg *zmsgcentre.Message) {
        fmt.Println("R2", msg.ID(), msg.Body)
    }))

    sender := a.CreateSender(testtopic)
    sender.Send("消息体")
}
```
