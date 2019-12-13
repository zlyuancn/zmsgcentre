/*
-------------------------------------------------
   Author :       zlyuan
   date：         2019/5/19
   Description :  自动id
-------------------------------------------------
*/

package zmsgcentre

import (
    "sync/atomic"
)

type autoNum struct {
    id uint64
}

func (m *autoNum) NextNum() uint64 {
    return atomic.AddUint64(&m.id, 1)
}

var autoMsgNum = &autoNum{}
var autoReceiverNum = &autoNum{}
