package main

import (
    "fmt"
    "log"
    "os/exec"
    "runtime"
    "time"
    "os"
)

var cout1 chan []byte = make(chan []byte)
var cin1 chan []byte = make(chan []byte)

var cout2 chan []byte = make(chan []byte)
var cin2 chan []byte = make(chan []byte)

var cmd1, cmd2 *exec.Cmd

func Foo(x byte) byte { return call_port1([]byte{1, x}) }
func Foo2(x byte) byte { return call_port2([]byte{1, x}) }
func Bar(y byte) byte { return call_port1([]byte{2, y}) }
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

    cmd1 = exec.Command("../player/player")
    stdin1, err := cmd1.StdinPipe()
    if err != nil {
        log.Fatal(err)
    }
    stdout1, err := cmd1.StdoutPipe()
    if err != nil {
        log.Fatal(err)
    }
    file1, err := os.Create("player1.txt")
    if err != nil {
        panic(err)
    }
    // close fo on exit and check for its returned error
    defer func() {
        if err := file1.Close(); err != nil {
            panic(err)
        }
    }()

    if err := cmd1.Start(); err != nil {
        log.Fatal(err)
    }

    cmd2 = exec.Command("../player/player")
    stdin2, err := cmd2.StdinPipe()
    if err != nil {
        log.Fatal(err)
    }
    stdout2, err := cmd2.StdoutPipe()
    if err != nil {
        log.Fatal(err)
    }
    file2, err := os.Create("player2.txt")
    if err != nil {
        panic(err)
    }
    // close fo on exit and check for its returned error
    defer func() {
        if err := file1.Close(); err != nil {
            panic(err)
        }
    }()

    if err := cmd2.Start(); err != nil {
        log.Fatal(err)
    }

    defer stdin1.Close()
    defer stdout1.Close()
    defer stdin2.Close()
    defer stdout2.Close()

    for {
        select {
        case s := <-cout1:
            // Write to file for audit before write to stdin
            stdin1.Write(s)
            buf := make([]byte, 2)
            runtime.Gosched()
            time.Sleep(100 * time.Millisecond)
            stdout1.Read(buf)
            // Read from stdout and write to file for audit before put into channel
            file1.Write(buf)
            cin1 <- buf
        case s := <-cout2:
            // Write to file for audit before write to stdin
            file2.Write(s)
            stdin2.Write(s)
            buf := make([]byte, 2)
            runtime.Gosched()
            time.Sleep(100 * time.Millisecond)
            stdout2.Read(buf)
            // Read from stdout and write to file for audit before put into channel
            file2.Write(buf)
            cin2 <- buf
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
    
    cmd1.Process.Kill()
    cmd2.Process.Kill()

}
