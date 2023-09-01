package remote

import (
	"errors"
	"fmt"
	"log"
	"net"
    
	"github.com/andreistan26/top-string/internal/core"
)

type Server struct {
    Addr    net.Addr
    Listner net.Listener
}

type ConnOpts struct {
    Port    int
    Ip      net.IP
}

func StartServer(opts *ConnOpts) (server *Server, err error) {
    server = &Server{}
    server.Listner, err = net.Listen("tcp", fmt.Sprintf("%s:%d", string(opts.Ip), opts.Port))
    
    if err != nil {
        log.Printf("Error has occured when starting server : %v\n", err.Error())
    }
    
    log.Printf("Listening on port %s:%d", string(opts.Ip), opts.Port)

    return server, err
}

func (server *Server) Run() error {
    for {
        conn, err := server.Listner.Accept()
        if err != nil {
            log.Printf("Got an error when Accept() : %v\n", err.Error())
        }

        log.Printf("Connection established with : %v\n", conn.RemoteAddr())
        go func() {
            defer conn.Close()
            server.HandleConnection(ConnWrapper{
                conn,
            })
        } ()
    }
}

// Here the server will receive firstly the requirements of the
// client {number of top strings to be counter}
// after which the server will loop until there are no more hashes to be received
func (server *Server) HandleConnection(conn ConnWrapper) error {
    firstMsg, err := conn.GetRequest()
    if err != nil {
        return err
    }

    if firstMsg.GetType() != TOP_QUERY {
        return errors.New("First message should be a TOP_COUNT_QUERRY")
    } else {
        log.Printf("Got message %v\n", firstMsg)
    }

    topCount := firstMsg.(TopQueryMessage).TopCount
    
    hashes := make(chan core.FileHash, 1000)
    
    results := core.CountStrings(hashes, int(topCount))
    
    for ;; {
        msg, err := conn.GetRequest()
        if err != nil {
            log.Printf("Got an error when reading requests : %v\n", err.Error())
        }
        if msg.GetType() != HASH_SEND {
            log.Printf("Got a different message than hash_send")
            break
        }

        hashes <- msg.(HashMessage).Payload
    }

    
    close(hashes)
    for result := range results {
        fmt.Println("got result", result)
        resultMsg := TopResultsMessage{
            Base: BaseMessage {
                Status: RESULTS,
                Code: SUCCESS,
            },
            Result: result,
        }

        conn.SendResponse(resultMsg)
    }

    conn.SendResponse(
        BaseMessage{
            Status: RESULTS_SENT,
            Code: SUCCESS,
        },
    )

    return nil
}
