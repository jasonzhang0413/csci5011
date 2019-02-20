package main

import (
    "fmt"
    "log"
    "os/exec"
    "runtime"
    "time"
)

var cout1 chan []byte = make(chan []byte)
var cin1 chan []byte = make(chan []byte)
var exit1 chan bool = make(chan bool)

var cout2 chan []byte = make(chan []byte)
var cin2 chan []byte = make(chan []byte)
var exit2 chan bool = make(chan bool)

func Foo(x byte) byte { return call_port1([]byte{1, x}) }
func Foo2(x byte) byte { return call_port2([]byte{1, x}) }
func Bar(y byte) byte { return call_port1([]byte{2, y}) }
func Exit1() byte      { return call_port1([]byte{0, 0}) }
func Exit2() byte      { return call_port2([]byte{0, 0}) }
func call_port1(s []byte) byte {
    cout1 <- s
    s = <-cin1
    return s[1]
}

func call_port2(s []byte) byte {
    cout2 <- s
    s = <-cin2
    return s[1]
}

func start() {
    fmt.Println("start")

    cmd1 := exec.Command("../player/player")
    stdin1, err := cmd1.StdinPipe()
    if err != nil {
        log.Fatal(err)
    }
    stdout1, err2 := cmd1.StdoutPipe()
    if err2 != nil {
        log.Fatal(err2)
    }
    if err := cmd1.Start(); err != nil {
        log.Fatal(err)
    }

    cmd2 := exec.Command("../player/player")
    stdin2, err := cmd2.StdinPipe()
    if err != nil {
        log.Fatal(err)
    }
    stdout2, err2 := cmd2.StdoutPipe()
    if err2 != nil {
        log.Fatal(err2)
    }
    if err := cmd2.Start(); err != nil {
        log.Fatal(err)
    }

    defer stdin1.Close()
    defer stdout1.Close()
    defer stdin2.Close()
    defer stdout2.Close()

    b1 := true
    b2 := true
    for b1 || b2 {
        select {
        case s := <-cout1:
            stdin1.Write(s)
            buf := make([]byte, 2)
            runtime.Gosched()
            time.Sleep(100 * time.Millisecond)
            stdout1.Read(buf)
            cin1 <- buf
        case s := <-cout2:
            stdin2.Write(s)
            buf := make([]byte, 2)
            runtime.Gosched()
            time.Sleep(100 * time.Millisecond)
            stdout2.Read(buf)
            cin2 <- buf
        case b1 = <-exit1:
            if b1 {
                fmt.Printf("Exit1")
                //return //os.Exit(0)
            }
        case b2 = <-exit2:
            if b2 {
                fmt.Printf("Exit2")
                //return //os.Exit(0)
            }
        }

    }
}
func main() {
    go start()
    runtime.Gosched()
    fmt.Println("30+1=", Foo(30)) //30+1= 31
    fmt.Println("2*40=", Bar(40)) //2*40= 80

    for i := 1; i <= 10; i++ {
        fmt.Printf("end %d \n", i)
    }
    fmt.Println("30+1=", Foo(100)) //30+1= 31
    fmt.Println("30+1=", Foo2(200)) //30+1= 31
    Exit1()
    //exit1 <- true
    Exit2()
    //exit2 <- true

}
