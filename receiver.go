/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2019/5/19
   Description :  接收器
-------------------------------------------------
*/

package zmsgcentre

//接收器
type Receiver struct {
    id        uint64      //全局唯一id
    topicName string      //主题名
    priority  int32       //优先级(数值小的会先收到消息)
    fn        msgFunction //接收消息的函数
}

//创建一个接收器(主题名, 接收消息的函数, 优先级(数值小的会先收到消息))
func NewReceiver(topicName string, priority int32, fn msgFunction) *Receiver {
    return &Receiver{
        id:        autoReceiverNum.NextNum(),
        topicName: topicName,
        priority:  priority,
        fn:        fn,
    }
}
