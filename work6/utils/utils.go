package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
)
type Err struct {
	Code int
	Msg string
	Notice string
}

func (e * Err)Error()string{
	return fmt.Sprintf(e.Msg)
}

func ErrInfo(code int,msg,notice string)Err{
	err := Err{
		code,
		msg,
		notice,
	}
	return err
}

func Encode(msg string)([]byte,error){
	var length = int32(len(msg))		// 读取消息的长度
	var pkg = new(bytes.Buffer)
	err := binary.Write(pkg,binary.LittleEndian, length)		// 写入消息头
	if err != nil {
		return nil, err
	}
	err = binary.Write(pkg,binary.LittleEndian,[]byte(msg))		// 写入消息实体
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(),nil
}

func Decode(reader *bufio.Reader)(string,error){
	lengthByte,_:=reader.Peek(4)					// 读取消息的长度
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff,binary.LittleEndian,&length)
	if err != nil {
		return "", err
	}
	// Buffered返回缓冲中现有的可读取的字节数
	if int32(reader.Buffered())<length+4{
		return "", err
	}
	pack := make([]byte,int(4+length))
	_,err = reader.Read(pack)		// 读取真正的消息数据
	if err != nil {
		return "", err
	}
	s := string(pack[4:])
	if err != nil {
		return "", err
	}
	return s,nil
}