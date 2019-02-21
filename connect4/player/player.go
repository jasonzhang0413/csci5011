package main

import (
    "log"
    "os"
    "flag"
    "encoding/json"
    //"io/ioutil"
    "fmt"
    "math/rand"
)

var width int
var height int

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

    jsonBody := []byte(`{"grid":[[0,0,0,0,0,0],[0,0,0,0,1,1],[0,0,0,0,0,2],[0,0,0,0,0,0],[0,0,0,0,0,2],[0,0,0,0,0,1],[0,0,0,0,0,0]]}`)

    fmt.Println(string(jsonBody))

    /*data, err := ioutil.ReadAll(os.Stdin)

    if err != nil {
        return nil, err
    }

    fmt.Println("Read ==>" + string(data))*/
    var state State
    err := json.Unmarshal(jsonBody, &state)

    if err != nil {
        return &state, err
    }
    /*body := make([]byte, 0, 4*1024)

    n, err := os.Stdin.Read(body)
    if err != nil {
        return body, err
    }

    log.Printf("Read input length ", n)*/

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
