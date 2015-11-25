package cmd

import (
	"fmt"
	"testing"
)

func TestDp(t *testing.T) {
	fmt.Println("TestDp")
	if err := Dp(false, nil); err != nil {
		t.Error(err.Error())
		fmt.Println("err!=nil")
	} else {
		t.Log("Dp list all datapools ok!")
		fmt.Println("Dp ok!\n")
	}

	args := []string{"dp1"}
	if err := Dp(false, args); err != nil {
		t.Error(err.Error())
	}
}
