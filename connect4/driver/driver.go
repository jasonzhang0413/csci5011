package main

import (
    "fmt"
    "log"
    "os/exec"
    "runtime"
    "time"
    "os"
    "flag"
    "encoding/json"
    "bufio"
)

var cout1 chan []byte = make(chan []byte)
var cin1 chan []byte = make(chan []byte)

var cout2 chan []byte = make(chan []byte)
var cin2 chan []byte = make(chan []byte)

var width, height int
var tournament bool
var cmd1, cmd2 *exec.Cmd

type State struct {
    Grid   [][]int `json:"grid"`
}

type Request struct {
    Move int `json:"move"`
}

func call_port1(bytes []byte) []byte {
    cout1 <- bytes
    bytes = <-cin1
    return bytes
}

func call_port2(bytes []byte) []byte {
    cout2 <- bytes
    bytes = <-cin2
    return bytes
}

func start() {
    fmt.Println("start")

    cmd1 = exec.Command("../player/player", fmt.Sprintf("%s%d", "--width=", width), fmt.Sprintf("%s%d", "--height=", height))
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

    cmd2 = exec.Command("../player/player", fmt.Sprintf("%s%d", "--width=", width), fmt.Sprintf("%s%d", "--height=", height))
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
            file1.Write(s)
            stdin1.Write(s)
            runtime.Gosched()
            time.Sleep(100 * time.Millisecond)
            reader := bufio.NewReader(stdout1)
            data, _ := reader.ReadBytes('\n')
            // Read from stdout and write to file for audit before put into channel
            file1.Write(data)
            cin1 <- data
        case s := <-cout2:
            // Write to file for audit before write to stdin
            file2.Write(s)
            stdin2.Write(s)
            runtime.Gosched()
            time.Sleep(100 * time.Millisecond)
            reader := bufio.NewReader(stdout2)
            data, _ := reader.ReadBytes('\n')
            // Read from stdout and write to file for audit before put into channel
            file2.Write(data)
            cin2 <- data
        }

    }
}
func main() {

    flag.IntVar(&width, "width", 7, "The width of grid for connect four game, default 7")
    flag.IntVar(&height, "height", 6, "The height of grid for connect four game, default 6")
    flag.BoolVar(&tournament, "tournament", false, "Tournament mode")
    flag.Parse()

    go start()
    runtime.Gosched()

    times := 1
    if tournament {
        times = 20
    }

    for i := 1; i <= times; i++ {
        state := StartNewGame()
        var moveRequest Request

        bytes, err := json.Marshal(state)
        if err != nil {
            fmt.Println("Fail to marshal state " + string(bytes))
        }
        request := call_port1(append(bytes, '\n'))
        json.Unmarshal(request, &moveRequest)

        fmt.Println(moveRequest)
        if ValidateMove(state, moveRequest.Move) {
            MakeMove(state, moveRequest.Move, 1)
        }
        fmt.Println("New state")
        fmt.Println(state.Grid)

        bytes, err = json.Marshal(state)
        if err != nil {
            fmt.Println("Fail to marshal state " + string(bytes))
        }
        request = call_port2(append(bytes, '\n'))
        json.Unmarshal(request, &moveRequest)

        fmt.Println(moveRequest)
        if ValidateMove(state, moveRequest.Move) {
            MakeMove(state, moveRequest.Move, 2)
        }
        fmt.Println("New state")
        fmt.Println(state.Grid)

    }


    cmd1.Process.Kill()
    cmd2.Process.Kill()

}

func StartNewGame() *State {
    grid := make([][]int, width)
    for i := range grid {
        grid[i] = make([]int, height)
    }

    initialState := &State{Grid: grid}

    return initialState
}

func ValidateMove(state *State, moveIndex int) bool {
    if state.Grid[moveIndex][0] == 0 {
        return true
    } else {
        return false
    }
}

func MakeMove(state *State, moveIndex int, player int) {
    for i := height - 1; i >= 0; i-- {
        if state.Grid[moveIndex][i] == 0 {
            state.Grid[moveIndex][i] = player
            break
        }
    }
}
