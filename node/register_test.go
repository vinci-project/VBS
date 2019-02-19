package main

import (
	"testing"
)

func TestGetHardwareInformation(t *testing.T){
	_,err := node.getHardwareInformation();
	if err != nil{
		t.Error("Get Failed")
	} else {
		t.Log("Get Passed")
	}
}