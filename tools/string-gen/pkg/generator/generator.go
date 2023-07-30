package generator

import (
	"fmt"
	"math/rand"
	"os"
    "path/filepath"
    "io"
    "time"
)

type StringGeneratorOpts struct {
    Path    string

    // bytes
    SizeRange []uint

    // chance of a string to repeate
    RepeateChance float32

    Count   uint
    Prefix  string
}

const (
    MAX_CHAR = 126
    MIN_CHAR = 31
)

func GenerateRandomChars(size uint) []byte {
    byteArr := make([]byte, size)
    var i uint = 0
    for ; i < size; i++ {
        byteArr[i] = byte(rand.Uint64() % (MAX_CHAR - MIN_CHAR) + MIN_CHAR)
    }

    return byteArr
}

func StringPrefix(prefix string, idx uint) string {
    return fmt.Sprintf("%s_%d", prefix, idx)
}

// This will generate random strings 
func GenerateStrings(options StringGeneratorOpts) error {
    if _, err := os.Stat(options.Path); err != nil {
        return err
    }

    rand.Seed(time.Now().UnixNano())

    var stringIdx uint
    for stringIdx = 0; stringIdx != options.Count; stringIdx++ {
        currentFile, err := os.Create(filepath.Join(options.Path, StringPrefix(options.Prefix, stringIdx)))
        if err != nil {
            panic(err)
        }

        repeateRand := rand.Float32()
        if repeateRand < options.RepeateChance && stringIdx != 0 {
            // we copy a previous string
            copyStrIdx := rand.Int63n(int64(stringIdx)) 
            copyStrPrefix := StringPrefix(options.Prefix, uint(copyStrIdx))
            copyStrPath := filepath.Join(options.Path, copyStrPrefix)
            fmt.Print(copyStrPath)

            copyFile, err := os.Open(copyStrPath)
            if err != nil {
                panic(err)
            }
            
            _, err = io.Copy(currentFile, copyFile)
            if err != nil {
                panic(err)
            }

            copyFile.Close()
            currentFile.Close()

            continue
        }
    
        stringSize := uint(rand.Uint64()) % (options.SizeRange[1] - options.SizeRange[0]) + options.SizeRange[0]
        _, err = currentFile.Write(GenerateRandomChars(stringSize))
        if err != nil {
            panic(err)
        }

        currentFile.Close()
    }

    return nil
}
