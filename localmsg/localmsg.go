package localmsg

import (
	"encoding/json"
	"github.com/yunduansing/gtools/utils"
)

type MsgQueue struct {
	state int32
	c     chan MsgData
}

// 异步发送
func (m *MsgQueue) SendAsync(msg MsgData) {
	m.c <- msg
}

// 同步发送
func (m *MsgQueue) SendSync(msg MsgData) {

}

func New() *MsgQueue {
	return &MsgQueue{}
}

func (m *MsgQueue) Start() {
	for msg := range m.c {
		go msg.Do()
	}
}

func (m *MsgQueue) Stop() {
	close(m.c)
}

// local msg type
type MsgType int

// local msg data struct
type MsgData struct {
	Type      MsgType
	Data      any   //msg data
	Timestamp int64 //unit timestamp
	Do        func()
}

func (msg MsgData) String() (string, error) {
	bs, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return utils.ByteToString(bs), nil
}

func (msg MsgData) GetData(val interface{}) error {
	bs, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bs, val)
	return err
}
