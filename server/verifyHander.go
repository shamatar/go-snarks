package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
	"bytes"
	"net/http"
	"os/exec"
)

type verificationResponse struct {
	Error bool `json:"error"`
}

func VerifyHander(w http.ResponseWriter, r *http.Request) {
	// if cmd, e := exec.Run("/bin/ls", nil, nil, exec.DevNull, exec.Pipe, exec.MergeWithStdout); e == nil {
	//     b, _ := ioutil.ReadAll(cmd.Stdout)
	//     println("output: " + string(b))
	// }

	// dataJson := `[[1],[2],[3]]`
	var req verificationRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		writeError(w)
		return
	}

	log.Printf("Unmarshaled: %v", req)

	d1 := []byte(req.Proof)
	err = ioutil.WriteFile("proof.txt", d1, 0644)
	if err != nil {
		writeError(w)
		return
	}
	cmd := exec.Command("./battleship", "-v")
        var out bytes.Buffer
        cmd.Stdout = &out
        err = cmd.Run()
        if err != nil {
                log.Println(err)
                writeError(w)
                return
        }
        fmt.Printf("in all caps: %q\n", out.String())

	outString := out.String()
	log.Println(outString)
	length := len(outString)
	result := outString[length-5 : length-1]
	fmt.Println(result)
	if result == "PASS" {
		writeSuccess(w)
		return
	}
	writeError(w)
}

func writeSuccess(w http.ResponseWriter) {
	resp := proverResponse{false}

	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
