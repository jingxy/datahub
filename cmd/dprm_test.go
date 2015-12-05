package cmd

import (
	"testing"
)

func TestDprm(t *testing.T) {
	var testdpnames []string = []string{"DpUnitTest1", "DpUnitTest2", "DpUnitTest3", "DpUnitTest4"}

	var cargs = []string{testdpnames[0], "--type=file", "--conn=/var/lib/datahub/test"}
	if err := DpCreate(false, cargs); err != nil {
		t.Errorf("DpCreate error: %v", err)
	}
	if err := DpRm(false, []string{testdpnames[0]}); err != nil {
		t.Errorf("DpRm error:%v", err)
	}

	var dargs = []string{testdpnames[1], "--type=file", "--conn=unittest"}
	if err := DpCreate(false, dargs); err != nil {
		t.Errorf("DpCreate error: %v", err)
	}
	if err := DpRm(false, []string{testdpnames[1]}); err != nil {
		t.Errorf("DpRm error:%v", err)
	}

}
