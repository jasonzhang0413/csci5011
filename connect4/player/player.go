package main

import (
    "log"
    "os"
    "flag"
    "encoding/json"
    //"fmt"
    "math/rand"
    "bufio"
)

var width int
var height int

var file1 *os.File

type State struct {
    Grid [][]int `json:"grid"`
}

type Request struct {
    Move int `json:"move"`
}

func main() {
    flag.IntVar(&width, "width", 7, "The width of grid for connect four game, default 7")
    flag.IntVar(&height, "height", 6, "The height of grid for connect four game, default 6")
    flag.Parse()

    file1, _ = os.Create("player_player1.txt")

    defer file1.Close()

    file1.WriteString("temp\n")

    for {
        state, err := GetState()
        if err != nil {
            log.Println("error reading " + err.Error())
        }

        moveIndex := MakeValidMove(state)
        request := &Request{Move: moveIndex}

        enc := json.NewEncoder(os.Stdout)
        enc.Encode(request)
    }

}

func GetState() (*State, error) {

    /*jsonBody := []byte(`{"grid":[[0,0,0,0,0,0],[0,0,0,0,1,1],[0,0,0,0,0,2],[0,0,0,0,0,0],[0,0,0,0,0,2],[0,0,0,0,0,1],[0,0,0,0,0,0]]}`)

    fmt.Println(string(jsonBody))*/

    file1.WriteString("ready to read\n")
    reader := bufio.NewReader(os.Stdin)
    data, err := reader.ReadBytes('\n')

    var state State

    //dec := json.NewDecoder(os.Stdin)
    //dec.Decode(&state)

    //err := json.Unmarshal(jsonBody, &state)
    err = json.Unmarshal(data, &state)

    if err != nil {
        return &state, err
    }

    return &state, nil
}

func MakeValidMove(state *State) int {
    var moveIndex int
    for {
        moveIndex = rand.Intn(width)
        if ValidateMove(state, moveIndex) {
            break
        }
    }

    return moveIndex
}

func ValidateMove(state *State, moveIndex int) bool {
    if state.Grid[moveIndex][0] == 0 {
        return true
    } else {
        return false
    }
}
