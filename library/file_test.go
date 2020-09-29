package library

import "testing"

func TestReadlines(t *testing.T) {
	filename := "../keysfile"
	lines, err := Readlines(filename)
	if err != nil {
		t.Error(err)
	}else{
		t.Log(lines)
	}
}