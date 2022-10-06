package localmsg

import (
	"encoding/json"
	"github.com/yunduansing/gtools/gen"
)

type MsgQueue struct {

}
//同步发送
func (m *MsgQueue) SendAsync(msg MsgData)  {

}
//异步发送
func (m *MsgQueue) SendSync(msg MsgData)  {

}

func New() *MsgQueue {
	return &MsgQueue{}
}

// local msg type
type MsgType int

//local msg data struct
type MsgData struct {
	Type MsgType
	Data interface{} //msg data
	Timestamp int64 //unit timestamp
}

func (msg MsgData) String() (string,error) {
	bs,err:=json.Marshal(msg)
	if err!=nil{
		return "", err
	}
	return gen.ByteToString(bs),nil
}

func (msg MsgData) GetData(val interface{}) error {
	bs,err:=json.Marshal(msg)
	if err!=nil{
		return err
	}
	err=json.Unmarshal(bs,val)
	return err
}
