package main

import (
	"testing"
	"os"
)

func TestGetMd5FromFile(t *testing.T){

	t.Run("md5Verify", func(t *testing.T){
		f,_ := os.Open("resource/testMd5.txt")
		md5Str := getMd5FromFile(f)
		if md5Str != "a906449d5769fa7361d7ecc6aa3f6d28"{
			t.Fail()
		}
	})
}
