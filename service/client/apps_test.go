package client

import (
	"fmt"
	pprint "github.com/NubeIO/edge/pkg/helpers/print"
	"testing"
)

func TestHost(*testing.T) {

	cli := New("", 0)
	file, err := cli.UploadFile("/home/aidan/Downloads/fileutils-master.zip", "/home/aidan")

	pprint.PrintJOSN(file)
	fmt.Println(err)
}
