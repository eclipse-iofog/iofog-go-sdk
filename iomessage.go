package container_sdk_go

import (
	"encoding/binary"
	"errors"
)

const IOMESSAGE_VERSION = 4

type IoMessage struct {
	ID               string `json:"id"`
	Tag              string `json:"tag"`
	GroupId          string `json:"groupid"`
	SequenceNumber   int    `json:"sequencenumber"`
	SequenceTotal    int    `json:"sequencetotal"`
	Priority         int    `json:"priority"`
	Timestamp        int64  `json:"timestamp"`
	Publisher        string `json:"publisher"`
	AuthID           string `json:"authid"`
	AuthGroup        string `json:"authgroup"`
	Version          int    `json:"version"`
	ChainPosition    int64  `json:"chainposition"`
	Hash             string `json:"hash"`
	PreviousHash     string `json:"previoushash"`
	Nonce            string `json:"nonce"`
	DifficultyTarget int    `json:"difficultytarget"`
	InfoType         string `json:"infotype"`
	InfoFormat       string `json:"infoformat"`
	ContextData      []byte `json:"contextdata"`
	ContentData      []byte `json:"contentdata"`
}

func (msg *IoMessage) DecodeBinary(data []byte) error {

	msg.Version = int(binary.BigEndian.Uint16(data[:2]))
	if msg.Version != IOMESSAGE_VERSION {
		return errors.New("Incompatible IoMessage version")
	}

	dataPos := uint32(33)
	var nextLength uint32
	nextLength = uint32(data[2])
	if nextLength != 0 {
		msg.ID = string(data[dataPos : dataPos+nextLength])
		dataPos += nextLength
	}
	nextLength = uint32(binary.BigEndian.Uint16(data[3:5]))
	if nextLength != 0 {
		msg.Tag = string(data[dataPos : dataPos+nextLength])
		dataPos += nextLength
	}
	nextLength = uint32(data[5])
	if nextLength != 0 {
		msg.GroupId = string(data[dataPos : dataPos+nextLength])
		dataPos += nextLength
	}
	nextLength = uint32(data[6])
	if nextLength != 0 {
		b := make([]byte, 4)
		for i := uint32(0); i < 4 && i < nextLength; i++ {
			b[i] = data[dataPos+nextLength-uint32(i)-1]
		}
		msg.SequenceNumber = int(binary.LittleEndian.Uint32(b))
		dataPos += nextLength
	}
	nextLength = uint32(data[7])
	if nextLength != 0 {
		b := make([]byte, 4)
		for i := uint32(0); i < 4 && i < nextLength; i++ {
			b[i] = data[dataPos+nextLength-uint32(i)-1]
		}
		msg.SequenceTotal = int(binary.LittleEndian.Uint32(b))
		dataPos += nextLength
	}
	nextLength = uint32(data[8])
	if nextLength != 0 {
		b := make([]byte, 4)
		for i := uint32(0); i < 4 && i < nextLength; i++ {
			b[i] = data[dataPos+nextLength-uint32(i)-1]
		}
		msg.Priority = int(binary.LittleEndian.Uint32(b))
		dataPos += nextLength
	}
	nextLength = uint32(data[9])
	if nextLength != 0 {
		b := make([]byte, 8)
		for i := uint32(0); i < 8 && i < nextLength; i++ {
			b[i] = data[dataPos+nextLength-uint32(i)-1]
		}
		msg.Timestamp = int64(binary.LittleEndian.Uint64(b))
		dataPos += nextLength
	}
	nextLength = uint32(data[10])
	if nextLength != 0 {
		msg.Publisher = string(data[dataPos : dataPos+nextLength])
		dataPos += nextLength
	}
	nextLength = uint32(binary.BigEndian.Uint16(data[11:13]))
	if nextLength != 0 {
		msg.AuthID = string(data[dataPos : dataPos+nextLength])
		dataPos += nextLength
	}
	nextLength = uint32(binary.BigEndian.Uint16(data[13:15]))
	if nextLength != 0 {
		msg.AuthGroup = string(data[dataPos : dataPos+nextLength])
		dataPos += nextLength
	}
	nextLength = uint32(data[15])
	if nextLength != 0 {
		b := make([]byte, 8)
		for i := uint32(0); i < 8 && i < nextLength; i++ {
			b[i] = data[dataPos+nextLength-uint32(i)-1]
		}
		msg.ChainPosition = int64(binary.LittleEndian.Uint64(b))
		dataPos += nextLength
	}
	nextLength = uint32(binary.BigEndian.Uint16(data[16:18]))
	if nextLength != 0 {
		msg.Hash = string(data[dataPos : dataPos+nextLength])
		dataPos += nextLength
	}
	nextLength = uint32(binary.BigEndian.Uint16(data[18:20]))
	if nextLength != 0 {
		msg.PreviousHash = string(data[dataPos : dataPos+nextLength])
		dataPos += nextLength
	}
	nextLength = uint32(binary.BigEndian.Uint16(data[20:22]))
	if nextLength != 0 {
		msg.Nonce = string(data[dataPos : dataPos+nextLength])
		dataPos += nextLength
	}
	nextLength = uint32(data[22])
	if nextLength != 0 {
		b := make([]byte, 4)
		for i := uint32(0); i < 4 && i < nextLength; i++ {
			b[i] = data[dataPos+nextLength-uint32(i)-1]
		}
		msg.DifficultyTarget = int(binary.LittleEndian.Uint32(b))
		dataPos += nextLength
	}
	nextLength = uint32(data[23])
	if nextLength != 0 {
		msg.InfoType = string(data[dataPos : dataPos+nextLength])
		dataPos += nextLength
	}
	nextLength = uint32(data[24])
	if nextLength != 0 {
		msg.InfoFormat = string(data[dataPos : dataPos+nextLength])
		dataPos += nextLength
	}
	nextLength = binary.BigEndian.Uint32(data[25:29])
	if nextLength != 0 {
		msg.ContextData = data[dataPos : dataPos+nextLength]
		dataPos += nextLength
	}
	nextLength = binary.BigEndian.Uint32(data[29:33])
	if nextLength != 0 {
		msg.ContentData = data[dataPos : dataPos+nextLength]
		dataPos += nextLength
	}
	return nil
}

func (msg *IoMessage) EncodeBinary() ([]byte, error) {
	if msg.Version != IOMESSAGE_VERSION {
		return nil, errors.New("Incompatible IoMessage version")
	}

	bytesContentData := make([]byte, 1)
	bytesContextData := make([]byte, 1)

	if msg.ContentData != nil {
		bytesContentData = msg.ContentData
	}
	if msg.ContextData != nil {
		bytesContextData = msg.ContextData
	}

	msgHeaderBytes := make([]byte, 0, 33)
	msgBodyBytes := make([]byte, 0, 128)
	versionBytes := make([]byte, 2)
	lenBytes := make([]byte, 4)
	var value []byte
	var n int

	binary.BigEndian.PutUint16(versionBytes, uint16(msg.Version))
	msgHeaderBytes = append(msgHeaderBytes, versionBytes...)

	msgHeaderBytes = append(msgHeaderBytes, byte(len([]byte(msg.ID))))
	msgBodyBytes = append(msgBodyBytes, []byte(msg.ID)...)

	binary.BigEndian.PutUint16(lenBytes[:2], uint16(len([]byte(msg.Tag))))
	msgHeaderBytes = append(msgHeaderBytes, lenBytes[:2]...)
	msgBodyBytes = append(msgBodyBytes, []byte(msg.Tag)...)

	msgHeaderBytes = append(msgHeaderBytes, byte(len([]byte(msg.GroupId))))
	msgBodyBytes = append(msgBodyBytes, []byte(msg.GroupId)...)

	value, n = intToBytesBE(msg.SequenceNumber)
	msgHeaderBytes = append(msgHeaderBytes, byte(n))
	msgBodyBytes = append(msgBodyBytes, value...)

	value, n = intToBytesBE(msg.SequenceTotal)
	msgHeaderBytes = append(msgHeaderBytes, byte(n))
	msgBodyBytes = append(msgBodyBytes, value...)

	value, n = intToBytesBE(msg.Priority)
	msgHeaderBytes = append(msgHeaderBytes, byte(n))
	msgBodyBytes = append(msgBodyBytes, value...)

	value, n = int64ToBytesBE(msg.Timestamp)
	msgHeaderBytes = append(msgHeaderBytes, byte(n))
	msgBodyBytes = append(msgBodyBytes, value...)

	msgHeaderBytes = append(msgHeaderBytes, byte(len([]byte(msg.Publisher))))
	msgBodyBytes = append(msgBodyBytes, []byte(msg.Publisher)...)

	binary.BigEndian.PutUint16(lenBytes[:2], uint16(len([]byte(msg.AuthID))))
	msgHeaderBytes = append(msgHeaderBytes, lenBytes[:2]...)
	msgBodyBytes = append(msgBodyBytes, []byte(msg.AuthID)...)

	binary.BigEndian.PutUint16(lenBytes[:2], uint16(len([]byte(msg.AuthGroup))))
	msgHeaderBytes = append(msgHeaderBytes, lenBytes[:2]...)
	msgBodyBytes = append(msgBodyBytes, []byte(msg.AuthGroup)...)

	value, n = int64ToBytesBE(msg.ChainPosition)
	msgHeaderBytes = append(msgHeaderBytes, byte(n))
	msgBodyBytes = append(msgBodyBytes, value...)

	binary.BigEndian.PutUint16(lenBytes[:2], uint16(len([]byte(msg.Hash))))
	msgHeaderBytes = append(msgHeaderBytes, lenBytes[:2]...)
	msgBodyBytes = append(msgBodyBytes, []byte(msg.Hash)...)

	binary.BigEndian.PutUint16(lenBytes[:2], uint16(len([]byte(msg.PreviousHash))))
	msgHeaderBytes = append(msgHeaderBytes, lenBytes[:2]...)
	msgBodyBytes = append(msgBodyBytes, []byte(msg.PreviousHash)...)

	binary.BigEndian.PutUint16(lenBytes[:2], uint16(len([]byte(msg.Nonce))))
	msgHeaderBytes = append(msgHeaderBytes, lenBytes[:2]...)
	msgBodyBytes = append(msgBodyBytes, []byte(msg.Nonce)...)

	value, n = intToBytesBE(msg.DifficultyTarget)
	msgHeaderBytes = append(msgHeaderBytes, byte(n))
	msgBodyBytes = append(msgBodyBytes, value...)

	msgHeaderBytes = append(msgHeaderBytes, byte(len([]byte(msg.InfoType))))
	msgBodyBytes = append(msgBodyBytes, []byte(msg.InfoType)...)

	msgHeaderBytes = append(msgHeaderBytes, byte(len([]byte(msg.InfoFormat))))
	msgBodyBytes = append(msgBodyBytes, []byte(msg.InfoFormat)...)

	binary.BigEndian.PutUint32(lenBytes[:4], uint32(len(bytesContextData)))
	msgHeaderBytes = append(msgHeaderBytes, lenBytes[:4]...)
	msgBodyBytes = append(msgBodyBytes, bytesContextData...)

	binary.BigEndian.PutUint32(lenBytes[:4], uint32(len(bytesContentData)))
	msgHeaderBytes = append(msgHeaderBytes, lenBytes[:4]...)
	msgBodyBytes = append(msgBodyBytes, bytesContentData...)

	return append(msgHeaderBytes, msgBodyBytes...), nil
}
