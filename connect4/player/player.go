package main

import (
    "log"
    "os"
    "flag"
    "encoding/json"
    "math/rand"
    "bufio"
	//"fmt"
    "time"
)

var width int
var height int
var player int

type State struct {
    Grid [][]int `json:"grid"` //[width][height]
}

type Request struct {
    Move int `json:"move"`
}

func main() {
    flag.IntVar(&width, "width", 7, "The width of grid for connect four game, default 7")
    flag.IntVar(&height, "height", 6, "The height of grid for connect four game, default 6")
    flag.IntVar(&player, "player", 1, "The player number, default 1")
    flag.Parse()

    rand.Seed(time.Now().UnixNano())

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
