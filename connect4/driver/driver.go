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
    "strings"
)

var cout1 chan []byte = make(chan []byte)
var cin1 chan []byte = make(chan []byte)

var cout2 chan []byte = make(chan []byte)
var cin2 chan []byte = make(chan []byte)

var width, height int
var tournamentTimes int
var player1Cmd, player1Args, player2Cmd, player2Args string
var winCounter1, winCounter2, drawCounter int
var cmd1, cmd2 *exec.Cmd

type State struct {
    Grid [][]int `json:"grid"` //[width][height]
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

    args1 := strings.Split(player1Args, ",")
    args1 = append(args1, fmt.Sprintf("%s%d", "--width=", width))
    args1 = append(args1, fmt.Sprintf("%s%d", "--height=", height))
    args1 = append(args1, fmt.Sprintf("%s%d", "--player=", 1))
    cmd1 = exec.Command(player1Cmd, args1...)
    //cmd1 = exec.Command("/usr/local/go/bin/go", "run", "/Users/jzhang201/go/src/csci5011/connect4/player/player.go", "--width=7", "--height=6", "--player=1")
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

    args2 := strings.Split(player2Args, ",")
    args2 = append(args2, fmt.Sprintf("%s%d", "--width=", width))
    args2 = append(args2, fmt.Sprintf("%s%d", "--height=", height))
    args2 = append(args2, fmt.Sprintf("%s%d", "--player=", 2))
    cmd2 = exec.Command(player2Cmd, args2...)
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
        if err := file2.Close(); err != nil {
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
    flag.StringVar(&player1Cmd, "player1Cmd", "../player/player.mac", "The command to invoke player1 program")
    flag.StringVar(&player1Args, "player1Args", "", "The arguments to invoke player1 program")
    flag.StringVar(&player2Cmd, "player2Cmd", "../player/player.mac", "The command to invoke player2 program")
    flag.StringVar(&player2Args, "player2Args", "", "The arguments to invoke player2 program")
    flag.IntVar(&tournamentTimes, "tournament", 1, "Tournament mode, number of games")
    flag.Parse()

    go start()
    runtime.Gosched()

    for i := 1; i <= tournamentTimes; i++ {
        state := StartNewGame()
        var moveRequest Request

        for {
            if i%2 == 1 {
                //Player 1
                bytes, err := json.Marshal(state)
                if err != nil {
                    fmt.Println("Fail to marshal state " + string(bytes))
                }
                request := call_port1(append(bytes, '\n'))
                json.Unmarshal(request, &moveRequest)

                fmt.Printf("move request made by player 1, column index %d\n", moveRequest)
                if ValidateMove(state, moveRequest.Move) {
                    rowIndex := MakeMove(state, moveRequest.Move, 1)
                    fmt.Println("Current state after player 1 moved")
                    fmt.Println(state.Grid)
                    if checkWinning(state.Grid, moveRequest.Move, rowIndex, 1) {
                        // player 1 win, start new game
                        fmt.Println("player 1 wins")
                        winCounter1++
                        break
                    } else if checkDraw(state.Grid) {
                        // draw, start new game
                        fmt.Println("Draw after player 1 moved")
                        drawCounter++
                        break
                    }
                    fmt.Println("player 1 does not win after move")
                } else {
                    //player 2 win
                    winCounter2++
                    break
                }

                //Player 2
                bytes, err = json.Marshal(state)
                if err != nil {
                    fmt.Println("Fail to marshal state " + string(bytes))
                }
                request = call_port2(append(bytes, '\n'))
                json.Unmarshal(request, &moveRequest)

                fmt.Printf("move request made by player 2, column index %d\n", moveRequest)
                if ValidateMove(state, moveRequest.Move) {
                    rowIndex := MakeMove(state, moveRequest.Move, 2)
                    fmt.Println("Current state after player 2 moved")
                    fmt.Println(state.Grid)
                    if checkWinning(state.Grid, moveRequest.Move, rowIndex, 2) {
                        // player 2 win, start new game
                        fmt.Println("player 2 wins")
                        winCounter2++
                        break
                    } else if checkDraw(state.Grid) {
                        // draw, start new game
                        fmt.Println("Draw after player 2 moved")
                        drawCounter++
                        break
                    }
                    fmt.Println("player 2 does not win after move")
                } else {
                    //player 1 win
                    winCounter1++
                    break
                }
            } else {
                //Player 2
                bytes, err := json.Marshal(state)
                if err != nil {
                    fmt.Println("Fail to marshal state " + string(bytes))
                }
                request := call_port2(append(bytes, '\n'))
                json.Unmarshal(request, &moveRequest)

                fmt.Printf("move request made by player 2, column index %d\n", moveRequest)
                if ValidateMove(state, moveRequest.Move) {
                    rowIndex := MakeMove(state, moveRequest.Move, 2)
                    fmt.Println("Current state after player 2 moved")
                    fmt.Println(state.Grid)
                    if checkWinning(state.Grid, moveRequest.Move, rowIndex, 2) {
                        // player 2 win, start new game
                        fmt.Println("player 2 wins")
                        winCounter2++
                        break
                    } else if checkDraw(state.Grid) {
                        // draw, start new game
                        fmt.Println("Draw after player 2 moved")
                        drawCounter++
                        break
                    }
                    fmt.Println("player 2 does not win after move")
                } else {
                    //player 1 win
                    winCounter1++
                    break
                }

                //Player 1
                bytes, err = json.Marshal(state)
                if err != nil {
                    fmt.Println("Fail to marshal state " + string(bytes))
                }
                request = call_port1(append(bytes, '\n'))
                json.Unmarshal(request, &moveRequest)

                fmt.Printf("move request made by player 1, column index %d\n", moveRequest)
                if ValidateMove(state, moveRequest.Move) {
                    rowIndex := MakeMove(state, moveRequest.Move, 1)
                    fmt.Println("Current state after player 1 moved")
                    fmt.Println(state.Grid)
                    if checkWinning(state.Grid, moveRequest.Move, rowIndex, 1) {
                        // player 1 win, start new game
                        fmt.Println("player 1 wins")
                        winCounter1++
                        break
                    } else if checkDraw(state.Grid) {
                        // draw, start new game
                        fmt.Println("Draw after player 1 moved")
                        drawCounter++
                        break
                    }
                    fmt.Println("player 1 does not win after move")
                } else {
                    //player 2 win
                    winCounter2++
                    break
                }
            }


        }
    }

    fmt.Printf("Result for %d times of game play\n", tournamentTimes)
    fmt.Printf("Player 1 wins %d times\n", winCounter1)
    fmt.Printf("Player 2 wins %d times\n", winCounter2)
    fmt.Printf("Draw happens %d times\n", drawCounter)

    //optional, make sure the player program is not left as orphan process
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

func ValidateMove(state *State, columnIndex int) bool {
    if state.Grid[columnIndex][0] == 0 {
        return true
    } else {
        return false
    }
}

func MakeMove(state *State, columnIndex int, playerValue int) int {
    // rowIndex is top to bottom
    rowIndex := height - 1
    for ; rowIndex >= 0; rowIndex-- {
        if state.Grid[columnIndex][rowIndex] == 0 {
            state.Grid[columnIndex][rowIndex] = playerValue
            break
        }
    }

    return rowIndex
}

func checkWinning(grid [][]int, columnIndex int, rowIndex int, playerValue int) bool {
    fmt.Printf("Row index %d\nStart to check winner against %d\n", rowIndex, playerValue)
    // columnIndex is left to right, rowIndex is from top to bottom
    return checkColumn(grid, columnIndex, rowIndex, playerValue) ||
        checkRow(grid, columnIndex, rowIndex, playerValue) ||
        checkSlashDiagonal(grid, columnIndex, rowIndex, playerValue) ||
        checkBackslashDiagonal(grid, columnIndex, rowIndex, playerValue)
}

func checkColumn(grid [][]int, columnIndex int, rowIndex int, playerValue int) bool {
    win := false
    if rowIndex + 4 < height &&
        grid[columnIndex][rowIndex+1] == playerValue &&
        grid[columnIndex][rowIndex+2] == playerValue &&
        grid[columnIndex][rowIndex+3] == playerValue {
            win = true
    }
    return win
}

func checkRow(grid [][]int, columnIndex int, rowIndex int, playerValue int) bool {
    win := false
    minCol := 0
    if columnIndex - 3 > minCol {
        minCol = columnIndex - 3
    }
    maxCol := width - 1
    if columnIndex + 3 < maxCol {
        maxCol = columnIndex + 3
    }

    for i := minCol; i <= maxCol - 3; i++ {
        if grid[i][rowIndex] == playerValue &&
            grid[i+1][rowIndex] == playerValue &&
            grid[i+2][rowIndex] == playerValue &&
            grid[i+3][rowIndex] == playerValue {
                win = true
                break
        }
    }

    return win
}

func checkSlashDiagonal(grid [][]int, columnIndex int, rowIndex int, playerValue int) bool {
    // rowIndex is top to bottom, figure out the lower left starting cell and upper right end cell
    // the distance check is make sure we are not running out of grid/index
    win := false
    leftDistance := 3
    if columnIndex - 0 < leftDistance {
        leftDistance = columnIndex - 0
    }
    lowerDistance := 3
    if height - 1 - rowIndex < lowerDistance {
        lowerDistance = height - 1 - rowIndex
    }
    lowerLeftDistance := leftDistance
    if lowerDistance < lowerLeftDistance {
        lowerLeftDistance = lowerDistance
    }

    rightDistance := 3
    if width - 1 - columnIndex < rightDistance {
        rightDistance = width - 1 - columnIndex
    }
    upperDistance := 3
    if rowIndex - 0 < upperDistance {
        upperDistance = rowIndex - 0
    }
    upperRightDistance := rightDistance
    if upperDistance < upperRightDistance {
        upperRightDistance = upperDistance
    }

    for i, j := columnIndex - lowerLeftDistance, rowIndex + lowerLeftDistance; i <= columnIndex + upperRightDistance - 3; i, j = i+1, j-1 {
        if grid[i][j] == playerValue &&
            grid[i+1][j-1] == playerValue &&
            grid[i+2][j-2] == playerValue &&
            grid[i+3][j-3] == playerValue {
            win = true
            break
        }
    }

    return win
}

func checkBackslashDiagonal(grid [][]int, columnIndex int, rowIndex int, playerValue int) bool {
    // rowIndex is top to bottom, figure out the upper left starting cell and lower right end cell
    // the distance check is make sure we are not running out of grid/index
    win := false
    leftDistance := 3
    if columnIndex - 0 < leftDistance {
        leftDistance = columnIndex - 0
    }
    upperDistance := 3
    if rowIndex - 0 < upperDistance {
        upperDistance = rowIndex - 0
    }
    upperLeftDistance := leftDistance
    if upperDistance < upperLeftDistance {
        upperLeftDistance = upperDistance
    }

    rightDistance := 3
    if width - 1 - columnIndex < rightDistance {
        rightDistance = width - 1 - columnIndex
    }
    lowerDistance := 3
    if height - 1 - rowIndex < lowerDistance {
        lowerDistance = height - 1 - rowIndex
    }
    lowerRightDistance := rightDistance
    if lowerDistance < lowerRightDistance {
        lowerRightDistance = lowerDistance
    }

    for i, j := columnIndex - upperLeftDistance, rowIndex - upperLeftDistance; i <= columnIndex + lowerRightDistance - 3; i, j = i+1, j+1 {
        if grid[i][j] == playerValue &&
            grid[i+1][j+1] == playerValue &&
            grid[i+2][j+2] == playerValue &&
            grid[i+3][j+3] == playerValue {
            win = true
            break
        }
    }

    return win
}

func checkDraw(grid [][]int) bool {
    // it is a draw if grid is fully filled(the trick is checking top row) without winner
    draw := true
    for i := 0; i < width; i++ {
        if grid[i][0] == 0 {
            draw = false
            break
        }
    }
    return draw
}
