/*
-------------------------------------------------
   Author :       zlyuan
   date：         2019/5/19
   Description :  消息结构体
-------------------------------------------------
*/

package zmsgcentre

//接受消息的函数
type msgFunction func(msg *Message)

//消息体
type Message struct {
    id        uint64      //全局唯一id
    topicName string      //消息主题
    isStop    bool        //是否被中断
    Body      interface{} //消息内容
}

//创建一个消息
func NewMessage(topicName string, body interface{}) *Message {
    return &Message{
        id:        autoMsgNum.NextNum(),
        topicName: topicName,
        isStop:    false,
        Body:      body,
    }
}

//获取该消息主题
func (m *Message) TopicName() string {
    return m.topicName
}

//获取该消息id
func (m *Message) ID() uint64 {
    return m.id
}

//停止消息传播
func (m *Message) Stop() {
    m.isStop = true
}

//获取该是否被中断
func (m *Message) IsStop() bool {
    return m.isStop
}
