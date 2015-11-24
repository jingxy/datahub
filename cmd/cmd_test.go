package cmd

import (
	"testing"
)

func TestStopP2P(t *testing.T) {
	if err := StopP2P(); err != nil {
		t.Log("StopP2P():%v", err)
	}
}
