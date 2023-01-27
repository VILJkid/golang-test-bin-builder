package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/VILJkid/golang-test-bin-builder/helpers"
)

func init() {
	helpers.GetFlags()
}

type platform struct {
	os        string
	shortName string
	extension string
}

func main() {
	found := false
	switch os.Args[1] {
	case "build", "b":
		found = true
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Panicln(err)
		}

		testDir := homeDir + "/Work/api_test_ymls/"
		toolName := "api.test"
		commands := []string{
			"go test -c -o " + toolName,
			"cp " + toolName + " " + testDir,
		}

		for _, cmd := range commands {
			splittedCmd := strings.Split(cmd, " ")
			if err := exec.Command(splittedCmd[0], splittedCmd[1:]...).Run(); err != nil {
				log.Panicln(err)
			}
		}

		log.Println("== New " + toolName + " binary generated and copied to the destination ==")

	case "upload", "u":
		found = true
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Panicln(err)
		}

		dlDir := homeDir + "/Downloads/"
		toolName := "api.test"
		separator := "-"
		arch := "amd64"
		version := helpers.GetFlags().Version.Value.String()

		platforms := []platform{
			{
				os:        "linux",
				shortName: "lnx",
				extension: "",
			},
			{
				os:        "darwin",
				shortName: "dwn",
				extension: "",
			},
			{
				os:        "windows",
				shortName: "win",
				extension: ".exe",
			},
		}

		if err := os.MkdirAll("v"+version, os.ModePerm); err != nil {
			log.Panicln(err)
		}

		wg := new(sync.WaitGroup)
		for _, p := range platforms {
			wg.Add(1)
			go func(p platform) {
				binaryName := fmt.Sprintf("%s%sv%s%s%s%s%s%s", toolName, separator, version, separator, p.shortName, separator, arch, p.extension)
				cmd := "env GOOS=" + p.os + " GOARCH=" + arch + " go test -c -o v" + version + "/" + binaryName
				splittedCmd := strings.Split(cmd, " ")
				if err := exec.Command(splittedCmd[0], splittedCmd[1:]...).Run(); err != nil {
					log.Panicln(err)
				}
				log.Println(binaryName + " binary generated.")
				wg.Done()
			}(p)
		}
		wg.Wait()

		cmd := "cp CHANGELOG.md v" + version + "/"
		splittedCmd := strings.Split(cmd, " ")
		if err := exec.Command(splittedCmd[0], splittedCmd[1:]...).Run(); err != nil {
			log.Panicln(err)
		}

		cmd = "mv v" + version + " " + dlDir
		splittedCmd = strings.Split(cmd, " ")
		if err := exec.Command(splittedCmd[0], splittedCmd[1:]...).Run(); err != nil {
			log.Panicln(err)
		}

		log.Println("== All binaries are ready to upload. ==")
	}

	if !found {
		log.Println("Oops, command not matched")
	}
}
