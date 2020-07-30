package acn

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
)

type guruErr string

func (ge guruErr) Error() string {
	return string(ge)
}

const (
	ErrGuruAmbiguous guruErr = "ambiguous selection within source file"
)

type GuruDescribeValueTypePos struct {
	ObjPos string `json:"objpos"`
	Desc   string `json:"desc"`
}

type GuruDescribeValue struct {
	Type     string                     `json:"type"`
	ObjPos   string                     `json:"objpos"`
	Typespos []GuruDescribeValueTypePos `json:"typespos"`
}

type GuruDescribeTypeMethod struct {
	Name string `json:"name"`
	Pos  string `json:"pos"`
}

type GuruDescribeType struct {
	Type    string                   `json:"type"`
	NamePos string                   `json:"namepos"`
	NameDef string                   `json:"namedef"`
	Methods []GuruDescribeTypeMethod `json:"methods"`
}

type GuruDescribeResponse struct {
	Desc   string             `json:"identifier"`
	Pos    string             `json:"pos"`
	Detail string             `json:"detail"`
	Value  *GuruDescribeValue `json:"value"`
	Type   *GuruDescribeType  `json:"type"`
}

func RunGuruDescribeQuery(filepath string, position int, resp *GuruDescribeResponse) error {
	cmd := exec.Command("guru", "-json", "describe", filepath+":#"+strconv.Itoa(position))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	processOut := out.Bytes()
	if bytes.Contains(processOut, []byte("ambiguous selection within source file")) {
		return ErrGuruAmbiguous
	} else if bytes.Contains(processOut, []byte("no such file or directory")) {
		return os.ErrNotExist
	} else {
		err := json.Unmarshal(processOut, resp)
		return err
	}
}
