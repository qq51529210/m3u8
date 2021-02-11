package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/qq51529210/m3u8"
)

func main() {
	if len(os.Args) > 2 {
		dir := filepath.Join(filepath.Dir(os.Args[0]), "m3u8")
		fmt.Println(dir)
		err := m3u8.SimpleDownload(os.Args[1], dir, 5)
		if err != nil {
			fmt.Println(err)
		}
	}
}
