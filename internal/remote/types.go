package remote

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/andreistan26/top-string/internal/core"
	pqueue "github.com/andreistan26/top-string/internal/priority_queue"
)

type ConnWrapper struct {
    net.Conn
}

type MessageType uint16
type ExitCode uint16

type Message interface {
    Serialize() []byte
    GetType() MessageType
}

type BaseMessage struct {
    Status  MessageType
    Code    ExitCode
}

type TopQueryMessage struct {
    Base BaseMessage 
    TopCount uint16
}

type HashMessage struct {
    Base BaseMessage
    Payload core.FileHash
}

type TopResultsMessage struct {
    Base BaseMessage
    Result pqueue.FileHash
}

const (
    SUCCESS ExitCode = 1
    FAIL = 0
)

const (
    TOP_QUERY MessageType = 0

    HASH_SEND      = 1
    ALL_HASH_SEND  = 2
    // server -> client
    RESULTS = 4
    RESULTS_SENT = 5

    ERROR = 6
)

func (baseMessage BaseMessage) Serialize() ([]byte) {
    buffer := make([]byte, 0)

    log.Printf("Trying to serialize %v", baseMessage)
    
    buffer = binary.LittleEndian.AppendUint16(buffer, uint16(baseMessage.Status)) 
    buffer = binary.LittleEndian.AppendUint16(buffer, uint16(baseMessage.Code))
    
    log.Printf("Serialized base msg into: %v", buffer)
    return buffer
}

func (topQueryMessage TopQueryMessage) Serialize() ([]byte) {
    buffer := topQueryMessage.Base.Serialize()

    log.Printf("Trying to serialize %v", topQueryMessage)
    
    buffer = binary.LittleEndian.AppendUint16(buffer, uint16(topQueryMessage.TopCount))
    log.Printf("Serialized top q msg into: %v", buffer)
    return buffer
}

func (hashMessage HashMessage) Serialize() ([]byte) {
    buffer := hashMessage.Base.Serialize()

    log.Printf("Trying to serialize %v", hashMessage)
    
    buffer = append(buffer, hashMessage.Payload.Hash[:]...)
    buffer = append(buffer, []byte(hashMessage.Payload.Filename)...)
    
    return buffer
}

func (topResultsMessage TopResultsMessage) Serialize() ([]byte) {
    buffer := topResultsMessage.Base.Serialize()
    
    log.Printf("Trying to serialize %v", topResultsMessage)

    buffer = binary.LittleEndian.AppendUint16(buffer, uint16(-topResultsMessage.Result.Priority))
    buffer = append(buffer, topResultsMessage.Result.Value[:]...)

    return buffer
}

func Deserialize(buffer []byte) (Message) {
    base := BaseMessage{
        Status: MessageType(binary.LittleEndian.Uint16(buffer[0:2])),
        Code: ExitCode(binary.LittleEndian.Uint16(buffer[2:4])),
    }

    switch base.Status {
    case TOP_QUERY:
        return TopQueryMessage {
            Base: base,
            TopCount: uint16(binary.LittleEndian.Uint16(buffer[4:6])),
        }
    
    case HASH_SEND:
        return HashMessage {
            Base: base,
            Payload: core.FileHash {
                Hash: [16]byte(buffer[4:20]),
                Filename: string(buffer[20:]),
            },
        }
    case RESULTS:
        return TopResultsMessage {
            Base: base,
            Result: pqueue.FileHash {
                Priority: int(binary.LittleEndian.Uint16(buffer[4:6])),
                Value: string(buffer[6:]),
            },
        }
    }
    

    return base
}

func (conn ConnWrapper) GetRequest() (Message, error){
    buffer := make([]byte, 256)
    n, err := conn.Read(buffer)
    
    if err != nil {
        log.Printf("Got an error when trying to read from client\n")
        return BaseMessage{
            Code: FAIL,
            Status: ERROR,
        }, err
    }

    log.Printf("Read %d bytes from client of %v \n", n, buffer[0:n]);
    
    msg := Deserialize(buffer[0:n])

    return msg, err
}

func (conn ConnWrapper) SendResponse (msg Message) {
    buffer := msg.Serialize()

    _, err := conn.Write(buffer)

    if err != nil {
        log.Printf("Got an error when trying to send response : %v\n", err)
    } else {
        log.Printf("Sent an response %v\n", buffer)
    }
}

func (baseMessage BaseMessage) GetType() (MessageType) {
    return baseMessage.Status
}

func (hashMessage HashMessage) GetType() (MessageType) {
    return hashMessage.Base.Status
}

func (topQueryMessage TopQueryMessage) GetType() (MessageType) {
    return topQueryMessage.Base.Status
}

func (topResultsMessage TopResultsMessage) GetType() (MessageType) {
    return topResultsMessage.Base.Status
}
