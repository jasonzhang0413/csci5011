package main

import (
    "log"
    "os"
    "flag"
    "encoding/json"
    "io/ioutil"
    "fmt"
    //"bufio"
)

var Width int
var height int

func foo(x byte) byte { return x + 1 }
func bar(y byte) byte { return y * 2 }

func ReadByte() byte {
    b1 := make([]byte, 1)
    for {
        n, _ := os.Stdin.Read(b1)
        if n == 1 {
            return b1[0]
        }
    }
}
func WriteByte(b byte) {
    b1 := []byte{b}
    for {
        n, _ := os.Stdout.Write(b1)
        if n == 1 {
            return
        }
    }
}

func ReadBytes() ([]byte, error) {

    dat, err := ioutil.ReadFile("./test.json")

    fmt.Println(string(dat))

    var state State
    json.Unmarshal(dat, &state)

    data, err := ioutil.ReadAll(os.Stdin)

    if err != nil {
        return data, err
    }

    /*body := make([]byte, 0, 4*1024)

    n, err := os.Stdin.Read(body)
    if err != nil {
        return body, err
    }

    log.Printf("Read input length ", n)*/

    return data, nil
}

func main() {

    flag.IntVar(&Width, "width", 7, "The width of grid for connect four game, default 7")
    flag.IntVar(&height, "height", 6, "The height of grid for connect four game, default 6")


    data, err := ReadBytes()
    if err != nil {
        log.Printf(err.Error())
    }

    //bodyStr := string(body)

    state := State{}
    json.Unmarshal(data, &state)

    var res byte
    for {
        fn := ReadByte()
        log.Println("fn=", fn)
        arg := ReadByte()
        log.Println("arg=", arg)
        if fn == 1 {
            res = foo(arg)
        } else if fn == 2 {
            res = bar(arg)
        } else if fn == 0 {
            return //exit
        } else {
            res = fn //echo
        }
        WriteByte(1)
        WriteByte(res)
    }
}
