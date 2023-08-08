package core

import (
	"container/heap"
	"crypto/md5"
	"fmt"
    "io"
    "io/ioutil"
	"os"
	"sync"

	pqueue "github.com/andreistan26/top-string/internal/priority_queue"
)

type SenderOpts struct {
    Paths       []string
}

type ReceiverOpts struct {
    QueryCount  int
}

func HashFile(file *os.File) ([16]byte) {
    h := md5.New()

    if _, err := io.Copy(h, file); err != nil {
        panic(err)
    }
    
    // only way to do this :(
    var buf [16]byte
    copy(buf[:], h.Sum(nil))

    return buf
}

type FileHash struct {
    Hash        [16]byte
    Filename    string
}

const (
    MAX_CONCURRENT_JOBS = 50
    CH_BUFSZ = 1000
)


func HashFiles(filenames <-chan string, out chan<- FileHash){
    waitChan := make(chan struct{}, MAX_CONCURRENT_JOBS)
    wg := &sync.WaitGroup{}
    
    for filename := range filenames {
        wg.Add(1)
        waitChan <- struct{}{}
        go func(filename string) {
            file, err := os.Open(filename)

            if err != nil {
                return
            }

            fileHash := FileHash {
                Hash:       HashFile(file),
                Filename:   filename,
            }
            
            out <- fileHash
            <- waitChan
            wg.Done()
        } (filename)
    }

    wg.Wait()
    close(out)
}


func GetHashStream(opts *SenderOpts) (chan FileHash) {
    c := make(chan string, CH_BUFSZ)
    out := make(chan FileHash, CH_BUFSZ)
    
    go HashFiles(c, out)

    wg := &sync.WaitGroup{}
    // Get all files from directories
    for _, path := range opts.Paths {
        wg.Add(1)
        go func(path string) {
            files, err := ioutil.ReadDir(path)
            if err != nil {
                return 
            }

            for _, file := range files {
                c <- fmt.Sprintf("%s/%s", path, file.Name())
            }

            wg.Done()
        }(path)
    }

    go func() {
        wg.Wait()
        close(c)
    } ()

    return out
}

// Now we need to actually see the top strings 
// we will need to have a cli option for choosing the number of "top" strings that we will
// hold, so we will plug the hashes in a hash-map and using a priority queue we will get the 
// first k top strings

type MapValue struct {
    Filename    string
    Frequency   int
}

func CountStrings(hashes <-chan FileHash, opts *ReceiverOpts) []pqueue.FileHash {
    hashMap := make(map[[16]byte] MapValue)
    fmt.Println("in count strings")
    var pq pqueue.PQueue

    for fileHash := range hashes {
        if value, exists := hashMap[fileHash.Hash]; exists == false {
            hashMap[fileHash.Hash] = MapValue{
                Filename: fileHash.Filename,
                Frequency: 1,
            }
        } else {
            value.Frequency += 1
            hashMap[fileHash.Hash] = value
        }
    }

    if len(hashMap) < opts.QueryCount {
        opts.QueryCount = len(hashMap)
    }

    for _, value := range hashMap {
        if len(pq) < opts.QueryCount {
            pq = append(pq, &pqueue.FileHash {
                Priority: -value.Frequency,
                Value: value.Filename,
                Index: len(pq),
            })

            if len(pq) == opts.QueryCount {
                heap.Init(&pq)
            }

            continue
        } else if -value.Frequency < pq[0].Priority {
            heap.Pop(&pq)
            heap.Push(&pq, &pqueue.FileHash {
                Priority: -value.Frequency,
                Value: value.Filename,
            })
        }
    }
    
    resCount := len(pq)
    results := make([]pqueue.FileHash, resCount)
    for pq.Len() > 0 {
        item := heap.Pop(&pq).(*pqueue.FileHash)
        results[pq.Len()] = *item
    }

    return results
}
