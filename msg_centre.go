/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2019/5/19
   Description :
-------------------------------------------------
*/

package zmsgcentre

import (
    "sync"
)

//主题对象储存器
type topicObjectStorage map[string]*topicStruct

//消息中心
type MsgCentre struct {
    rwx    sync.RWMutex
    topics topicObjectStorage
}

//创建一个消息中心(初始化消息主题大小)
func NewMsgCentre() *MsgCentre {
    return &MsgCentre{
        topics: make(topicObjectStorage, InitTopicCapacity),
    }
}

//添加接收器
func (m *MsgCentre) AddReceiver(receiver *Receiver) {
    var name = receiver.topicName
    m.rwx.Lock()
    topic, ok := m.topics[name]
    if !ok {
        topic = newTopicObject(name)
        m.topics[name] = topic
    }
    m.rwx.Unlock()

    topic.AddReceiver(receiver)
}

//移除接收器
func (m *MsgCentre) RemoveReceiver(receiver *Receiver) {
    m.rwx.Lock()
    topic, ok := m.topics[receiver.topicName]
    m.rwx.Unlock()

    if ok {
        topic.RemoveReceiver(receiver)
    }
}

//移除主题, 此操作会同时移除所有绑定该主题的接收器
func (m *MsgCentre) RemoveTopic(topicName string) {
    m.rwx.Lock()
    delete(m.topics, topicName)
    m.rwx.Unlock()
}

//移除所有主题, 此操作会同时移除所有的接收器
func (m *MsgCentre) RemoveAllTopic() {
    m.rwx.Lock()
    m.topics = make(topicObjectStorage, InitTopicCapacity)
    m.rwx.Unlock()
}

//发送消息
func (m *MsgCentre) Send(msg *Message) {
    m.rwx.RLock()
    topic, ok := m.topics[msg.topicName]
    m.rwx.RUnlock()

    if ok {
        topic.Send(msg)
    }
}

//创建一个发送器
func (m *MsgCentre) CreateSender(topicName string) *Sender {
    return &Sender{
        msgCentre: m,
        topicName: topicName,
    }
}
