package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/modeckrus/vuploader/changelog"
	"github.com/modeckrus/vuploader/command"
	"github.com/modeckrus/vuploader/vuploader"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token         string `json:"token" yaml:"token"`
	ChatId        int64  `json:"chat_id" yaml:"chat_id"`
	NeedChangeLog bool   `json:"need_change_log" yaml:"need_change_log"`
}

var configPath = "vuploader.yml"
var commentary = ""
var arg = ""
var path = ""

func init() {
	flag.StringVar(&configPath, "c", "vuploader.yml", "config path")
	flag.StringVar(&commentary, "m", "", "commentary")
	flag.StringVar(&arg, "e", "clean", "Command")
	initalPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&path, "p", initalPath, "Path")
	flag.Parse()
}

func main() {
	// cmd := exec.Command("flutter", "clean")
	// cmd.Stdin = os.Stdin
	// cmd.Stderr = os.Stderr
	// stdOut, err := cmd.StdoutPipe()
	// if err != nil {
	// 	fmt.Println("Error 1")
	// }
	// if err := cmd.Start(); err != nil {
	// 	fmt.Println("Error 2")
	// }
	// bytes, err := ioutil.ReadAll(stdOut)
	// if err != nil {
	// 	fmt.Println("Error 3")
	// }
	// if err := cmd.Wait(); err != nil {
	// 	fmt.Println("Error 4")
	// 	if exitError, ok := err.(*exec.ExitError); ok {
	// 		fmt.Printf("Exit code is %d\n", exitError.ExitCode())
	// 	}
	// }
	// fmt.Println(string(bytes))
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	yaml.Unmarshal(configBytes, &config)
	uploader, err := vuploader.NewTelegramUploader(config.Token, config.ChatId)
	if err != nil {
		log.Fatal(fmt.Errorf("Error while create NewTelegramUploader: \n%e", err))
	}
	changeLogPath := filepath.Join(path, "CHANGELOG.json")
	changeLog, err := changelog.NewChangelog(changeLogPath)
	if err != nil {
		log.Fatal(fmt.Errorf("Error while parse ChangeLog: \n%e", err))
	}
	log.Println(changeLog.LastVersion())
	commader := command.NewCommander(config.NeedChangeLog, changeLog, commentary, uploader, path)

	err = commader.ParseCommand(arg)
	if err != nil {
		log.Fatal(fmt.Errorf("Error while execute command(%s): \n%e", arg, err))
	}
}
