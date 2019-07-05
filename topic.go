/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2019/5/19
   Description :  主题
-------------------------------------------------
*/

package zmsgcentre

import (
    "container/list"
    "sync"
)

//接收器对象储存器
type receiverObjectStorage map[uint64]*list.Element

//主题结构
type topicStruct struct {
    name          string
    rwx           *sync.RWMutex
    receivers     receiverObjectStorage
    receiverChain *list.List
}

func newTopicObject(name string) *topicStruct {
    return &topicStruct{
        name:          name,
        rwx:           new(sync.RWMutex),
        receivers:     make(receiverObjectStorage),
        receiverChain: list.New(),
    }
}

//添加接收器
func (m *topicStruct) AddReceiver(receiver *Receiver) {
    var id = receiver.id
    m.rwx.Lock()
    _, ok := m.receivers[id]
    if !ok {
        e := m.receiverChain.Back()
        for {
            if e == nil {
                m.receivers[id] = m.receiverChain.PushFront(receiver)
                break
            }
            if e.Value.(*Receiver).priority <= receiver.priority {
                m.receivers[id] = m.receiverChain.InsertAfter(receiver, e)
                break
            }
            e = e.Prev()
        }
    }
    m.rwx.Unlock()
}

//移除接收器
func (m *topicStruct) RemoveReceiver(receiver *Receiver) {
    m.rwx.Lock()
    element, ok := m.receivers[receiver.id]
    if ok {
        m.receiverChain.Remove(element)
        delete(m.receivers, receiver.id)
    }
    m.rwx.Unlock()
}

//发送消息
func (m *topicStruct) Send(msg *Message) {
    m.rwx.RLock()
    if m.receiverChain.Len() == 0 {
        m.rwx.RUnlock()
        return
    }

    //获取当前接收器链镜像
    var chain = make([]*Receiver, m.receiverChain.Len())
    e := m.receiverChain.Front()
    var index int
    for {
        if e == nil {
            break
        }

        chain[index] = e.Value.(*Receiver)
        index++
        e = e.Next()
    }
    m.rwx.RUnlock() //不用defer

    //遍历接收器链镜像
    for _, r := range chain {
        r.fn(msg)
        if msg.IsStop() {
            return
        }
    }
}
