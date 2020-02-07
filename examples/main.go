package main

import (
	"context"
	"fmt"

	"github.com/winterssy/mxget/pkg/provider/netease"
)

func main() {
	client := netease.Client()
	resp, err := client.GetSong(context.Background(), "36990266")
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
