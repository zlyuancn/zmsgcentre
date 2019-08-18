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
    name              string
    rwx               sync.RWMutex
    receivers         receiverObjectStorage
    receiverChain     *list.List
    receiverChanged   bool
    receiverSnapshoot []*Receiver
}

func newTopicObject(name string) *topicStruct {
    return &topicStruct{
        name:            name,
        receivers:       make(receiverObjectStorage),
        receiverChain:   list.New(),
        receiverChanged: true,
    }
}

func (m *topicStruct) snapshoot() {
    m.rwx.Lock()
    defer m.rwx.Unlock()

    if m.receiverChanged {
        m.receiverChanged = false

        var snapshoot = make([]*Receiver, m.receiverChain.Len())
        e := m.receiverChain.Front()
        var index int
        for {
            if e == nil {
                break
            }
            snapshoot[index] = e.Value.(*Receiver)
            index++
            e = e.Next()
        }
        m.receiverSnapshoot = snapshoot
    }
}

//添加接收器
func (m *topicStruct) AddReceiver(receiver *Receiver) {
    var id = receiver.id

    m.rwx.Lock()
    defer m.rwx.Unlock()

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
        m.receiverChanged = true
    }
}

//移除接收器
func (m *topicStruct) RemoveReceiver(receiver *Receiver) {
    m.rwx.Lock()
    defer m.rwx.Unlock()

    element, ok := m.receivers[receiver.id]
    if ok {
        m.receiverChain.Remove(element)
        delete(m.receivers, receiver.id)
        m.receiverChanged = true
    }
}

//发送消息
func (m *topicStruct) Send(msg *Message) {
    m.rwx.RLock()
    defer m.rwx.RUnlock()

    for m.receiverChanged {
        m.rwx.RUnlock()
        m.snapshoot()
        m.rwx.RLock()
    }

    if len(m.receiverSnapshoot) == 0 {
        return
    }

    //遍历接收器链镜像
    for _, r := range m.receiverSnapshoot {
        r.fn(msg)
        if msg.IsStop() {
            return
        }
    }
}
