package main

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	goos := runtime.GOOS
	dir, err := homedir.Dir()
	if err != nil {
		return
	}

	var files []string
	switch goos {
	case "darwin":
		files, err = findDarwin(dir)
		if err != nil {
			fmt.Println("Config file not found")
			return
		}
		break
	case "windows":
		files, err = findWindows(dir)
		if err != nil {
			fmt.Println("Config file not found")
			return
		}
		break
	default:
		panic("No support OS:" + goos)
	}
	if len(files) == 0 {
		fmt.Println("Nothing change")
		return
	}
	for _, f := range files {
		fmt.Println("Removing " + f)
		err := os.Remove(f)
		if err != nil {
			fmt.Println("Remove failed " + f)
			continue
		}
	}
	fmt.Println("Done.")
}

func findDarwin(dir string) (data []string, err error) {
	var cmd *exec.Cmd
	match := fmt.Sprintf("'%s/Library/Application\\ Support/JetBrains/*/eval/*.evaluation.key'", dir)
	cmd = exec.Command("/bin/bash", "-c", "ls '"+match+"'")
	var output bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return
	}
	files := output.String()
	data = strings.Split(files, "\n")
	return
}

func findWindows(dir string) (data []string, err error) {
	files, _ := ioutil.ReadDir(dir)
	for _, fi := range files {
		if fi.IsDir() {
			if strings.Contains(fi.Name(), ".PhpStorm") || strings.Contains(fi.Name(), ".GoLand") {
				path := filepath.Join(dir, fi.Name(), "config", "eval")
				files, _ = ioutil.ReadDir(path)
				for _, fi2 := range files {
					if !fi2.IsDir() && strings.HasSuffix(fi2.Name(), ".evaluation.key") {
						data = append(data, filepath.Join(dir, fi.Name(), "config", "eval", fi2.Name()))
					}
				}
			}
		}
	}
	return
}
