/*
-------------------------------------------------
   Author :       zlyuan
   date：         2019/5/19
   Description :
-------------------------------------------------
*/

package zmsgcentre

type Sender struct {
    msgCentre *MsgCentre
    topicName string
}

//发送消息
func (m *Sender) Send(body interface{}) {
    m.msgCentre.Send(NewMessage(m.topicName, body))
}
