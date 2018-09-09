package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
	"bytes"
	"math/rand"
	"net/http"
	"os/exec"
	"strconv"
)

type proverResponse struct {
	Error bool `json:"error"`
}

type verificationRequest struct {
	Proof string `json:"proof"`
	Hash  string `json:"hash"`
}

type proofResponse struct {
	Proof string `json:"proof"`
	Hash  string `json:"hash"`
}

// type Battlefield struct {
// 	Field  [][]int `json:"field"`
// 	Key1   string  `json:"key_1"`
// 	Key2   string  `json:"key_2"`
// 	Key3   string  `json:"key_3"`
// 	Points []Point `json:"points"`
// }

func ProveHander(w http.ResponseWriter, r *http.Request) {
	// if cmd, e := exec.Run("/bin/ls", nil, nil, exec.DevNull, exec.Pipe, exec.MergeWithStdout); e == nil {
	//     b, _ := ioutil.ReadAll(cmd.Stdout)
	//     println("output: " + string(b))
	// }
	// dataJson := `[[1],[2],[3]]`
	var arr [][]int
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&arr)
	if err != nil {
		writeError(w)
		return
	}
	log.Printf("Unmarshaled: %v", arr)
	if len(arr) != 10 {
		writeError(w)
		return
	}
	for i := 0; i < len(arr); i++ {
		if len(arr[i]) != 10 {
			writeError(w)
			return
		}
	}
	fullString := ""
	for i := 0; i < len(arr); i++ {
		substr := ""
		for j := 0; j < len(arr[i]); j++ {
			substr = substr + strconv.Itoa(arr[i][j])
		}
		fullString = fullString + substr
	}

	log.Println(fullString)
	stringForProver := "b" + fullString
	salt := rand.Uint64()
	saltString := strconv.FormatUint(salt, 16)
//	out, err := exec.Command("./battleship -p " + stringForProver + " " + saltString).Output()
//	log.Println(string(out))
	//if cmd, e := exec.Run("./battleship -p", nil, nil, exec.DevNull, exec.Pipe, exec.MergeWithStdout); e == nil {
       	//	b, _ := ioutil.ReadAll(cmd.Stdout)
        //	fmt.Println("output: " + string(b))
        //}
	cmd := exec.Command("./battleship", "-p", stringForProver, saltString)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Println(err)
		writeError(w)
		return
	}
	fmt.Printf("in all caps: %q\n", out.String())
	testVerifier()
	fullFile, err := ioutil.ReadFile("proof.txt")
	if err != nil {
		writeError(w)
		return
	}
	fullContent := string(fullFile)
	writeResponse(w, fullContent, "")
}

func testVerifier() {
        cmd := exec.Command("./battleship", "-v")
        var out bytes.Buffer
        cmd.Stdout = &out
        err := cmd.Run()
        if err != nil {
                log.Println(err)
                return
        }
        fmt.Printf("in all caps: %q\n", out.String())

        outString := out.String()
        log.Println(outString)
        length := len(outString)
        result := outString[length-5 : length-1]
        fmt.Println(result)

}

func writeError(w http.ResponseWriter) {
	resp := proverResponse{true}

	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func writeResponse(w http.ResponseWriter, proof, hash string) {
	resp := proofResponse{proof, hash}

	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
