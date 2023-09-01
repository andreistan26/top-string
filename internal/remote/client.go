package remote

import (
	"fmt"
	"net"
    "log"

	"github.com/andreistan26/top-string/internal/core"
)

func SendFiles(opts core.SenderOpts, connOpts ConnOpts) (err error) {
    netConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", string(connOpts.Ip), connOpts.Port))

    if err != nil {
        log.Printf("Could not dial %s:%d", string(connOpts.Ip), connOpts.Port) 
    }

    conn := ConnWrapper {
        netConn,
    }

    conn.SendResponse(
        TopQueryMessage {
            Base: BaseMessage{
                Status: TOP_QUERY,
                Code: SUCCESS,
            },
            TopCount: uint16(opts.QueryCount),
        },
    )


    hashes := core.GetHashStream(opts)
    for hash := range hashes {
        conn.SendResponse(
            HashMessage {
                Base: BaseMessage {
                    Status: HASH_SEND,
                    Code: SUCCESS,
                },
                Payload: hash,
            },
        )
    }

    conn.SendResponse(
        BaseMessage{
            Status: ALL_HASH_SEND,
            Code: SUCCESS,
        },
    )
    
    for ;; {
        msg, err := conn.GetRequest()
        if err != nil {
            log.Printf("Got an error when reading requests : %v\n", err.Error())
        }

        if msg.GetType() != RESULTS {
            log.Println("RECEIVED ALL RESULTS")
            break
        }

        result := msg.(TopResultsMessage).Result

        fmt.Printf("%d : %s", result.Priority, result.Value)
    }

    return nil
}
