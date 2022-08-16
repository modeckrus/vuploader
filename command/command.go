package command

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/modeckrus/vuploader/changelog"
	"github.com/modeckrus/vuploader/vuploader"
	log "github.com/sirupsen/logrus"
)

type Commander struct {
	NeedChangeLog bool
	ChangeLog     changelog.ChangeLog
	Uploader      *vuploader.TelegramUploader
	Message       string
	Path          string
}

func NewCommander(needChangeLog bool, changeLog changelog.ChangeLog, message string, uploader *vuploader.TelegramUploader, path string) Commander {
	return Commander{
		NeedChangeLog: needChangeLog,
		ChangeLog:     changeLog,
		Message:       message,
		Uploader:      uploader,
		Path:          path,
	}
}

func (t Commander) ParseCommand(command string) error {
	var err error
	switch command {
	case "clean":
		err = t.Clean()
	case "build":
		err = t.Build()
	case "prod_android":
		_, err = t.ProdAndroid()
	case "dev_android":
		_, err = t.DevAndroid()
	case "stg_android":
		_, err = t.StgAndroid()
	case "prod_upload":
		err = t.ProdUpload()
	case "dev_upload":
		err = t.DevUpload()
	case "stg_upload":
		err = t.StgUpload()
	case "upload_all":
		err = t.UploadAll()
	case "test_rename":
		err = t.TestRename()
	case "test_send":
		err = t.TestSend()
	case "get_chat_id":
		err = t.GetChatID()
	case "just_upload":
		err = t.JustUpload()
	}
	return err
}
func (t Commander) JustUpload() error {
	tasker := NewTasker(t.Path)
	err := tasker.cd("./build/app/outputs/flutter-apk/")
	if err != nil {
		return err
	}
	pathes, err := tasker.getFiles()
	if err != nil {
		return err
	}
	resultPathes := []string{}
	for i := 0; i < len(pathes); i++ {
		path := pathes[i]
		if strings.Contains(path, "prod_") {
			resultPathes = append(resultPathes, path)
		}
		if strings.Contains(path, "dev_") {
			resultPathes = append(resultPathes, path)
		}
		if strings.Contains(path, "stg_") {
			resultPathes = append(resultPathes, path)
		}
	}
	message := t.Message

	if t.NeedChangeLog {
		if message != "" {
			message += "\n"
		}
		message += t.ChangeLog.LastVersion()
	}
	err = t.Uploader.UploadFiles(resultPathes, message)
	return err
}

func (t Commander) GetChatID() error {
	return t.Uploader.ChatIdGetter()
}
func (t Commander) TestSend() error {
	message := t.Message
	if message != "" {
		message += "\n"
	}
	message += t.ChangeLog.LastVersion()
	_, err := t.Uploader.SendMessage(message)
	return err
}
func (t Commander) TestRename() error {
	tasker := NewTasker(t.Path)
	err := tasker.cd("./build/app/outputs/flutter-apk/")
	if err != nil {
		return err
	}
	newName, err := tasker.rename("app-production-release.apk", "prod_2.1.84+140.apk")
	if err != nil {
		return err
	}
	log.Infof("newName: %s", newName)
	return nil
}
func (t Commander) ProdUpload() error {
	path, err := t.ProdAndroid()
	if err != nil {
		return err
	}

	message := t.Message

	if t.NeedChangeLog {
		if message != "" {
			message += "\n"
		}
		message += t.ChangeLog.LastVersion()
	}
	pathes := []string{path}
	err = t.Uploader.UploadFiles(pathes, message)
	return err
}
func (t Commander) StgUpload() error {
	path, err := t.StgAndroid()
	if err != nil {
		return err
	}

	message := t.Message

	if t.NeedChangeLog {
		if message != "" {
			message += "\n"
		}
		message += t.ChangeLog.LastVersion()
	}
	pathes := []string{path}
	err = t.Uploader.UploadFiles(pathes, message)
	return err
}
func (t Commander) DevUpload() error {
	path, err := t.DevAndroid()
	if err != nil {
		return err
	}

	message := t.Message

	if t.NeedChangeLog {
		if message != "" {
			message += "\n"
		}
		message += t.ChangeLog.LastVersion()
	}
	pathes := []string{path}
	err = t.Uploader.UploadFiles(pathes, message)
	return err
}
func (t Commander) UploadAll() error {
	pathes := make([]string, 3)
	path, err := t.ProdAndroid()
	if err != nil {
		return err
	}
	pathes = append(pathes, path)

	path, err = t.DevAndroid()
	if err != nil {
		return err
	}
	pathes = append(pathes, path)
	path, err = t.StgAndroid()
	if err != nil {
		return err
	}
	pathes = append(pathes, path)

	message := t.Message

	if t.NeedChangeLog {
		if message != "" {
			message += "\n"
		}
		message += t.ChangeLog.LastVersion()
	}
	err = t.Uploader.UploadFiles(pathes, message)
	return err
}
func (t Commander) ProdAndroid() (string, error) {
	tasker := NewTasker(t.Path)
	err := tasker.runCommandWithOutput("flutter build apk --dart-define=flavor=production --flavor=production --release -v")
	if err != nil {
		return "", err
	}
	tasker.cd("./build/app/outputs/flutter-apk/")
	oldName := "app-production-release.apk"
	newName := fmt.Sprintf("prod_%s.apk", t.ChangeLog.NumberVersion())
	newPath, err := tasker.rename(oldName, newName)
	if err != nil {
		return "", err
	}
	return newPath, err
}
func (t Commander) DevAndroid() (string, error) {
	tasker := NewTasker(t.Path)
	err := tasker.runCommandWithOutput("flutter build apk --dart-define=flavor=development --flavor=development --release -v")
	if err != nil {
		return "", err
	}
	tasker.cd("./build/app/outputs/flutter-apk/")
	oldName := "app-development-release.apk"
	newName := fmt.Sprintf("dev_%s.apk", t.ChangeLog.NumberVersion())
	newPath, err := tasker.rename(oldName, newName)
	if err != nil {
		return "", err
	}
	return newPath, err
}
func (t Commander) StgAndroid() (string, error) {
	tasker := NewTasker(t.Path)

	err := tasker.runCommandWithOutput("flutter build apk --dart-define=flavor=staging --flavor=staging --release -v")
	if err != nil {
		return "", err
	}
	tasker.cd("./build/app/outputs/flutter-apk/")
	oldName := "app-staging-release.apk"
	newName := fmt.Sprintf("stg_%s.apk", t.ChangeLog.NumberVersion())
	newPath, err := tasker.rename(oldName, newName)
	if err != nil {
		return "", err
	}
	return newPath, err
}
func (t Commander) Intl() error {
	tasker := NewTasker(t.Path)

	err := tasker.runCommandWithOutput("flutter pub global run intl_utils:generate")
	if err != nil {
		return err
	}
	return nil
}
func (t Commander) Build() error {
	tasker := NewTasker(t.Path)

	err := tasker.runCommandWithOutput("flutter pub run build_runner build --delete-conflicting-outputs")
	if err != nil {
		return err
	}
	return nil
}
func (t Commander) Clean() error {
	tasker := NewTasker(t.Path)
	err := tasker.runCommandWithOutput("flutter clean")
	if err != nil {
		return err
	}
	err = tasker.runCommandWithOutput("flutter pub get")
	if err != nil {
		return err
	}
	contains := strings.Contains(runtime.GOOS, "darwin/")
	contains = contains || strings.Contains(runtime.GOOS, "ios/")
	if !contains {
		return nil
	}
	err = tasker.cd("./ios")
	if err != nil {
		return err
	}
	err = tasker.runCommandWithOutput("pod install")
	if err != nil {
		return err
	}
	err = tasker.cd("../macos")
	if err != nil {
		return err
	}
	err = tasker.runCommandWithOutput("pod install")
	if err != nil {
		return err
	}
	return nil
}

type Tasker struct {
	Dir string
}

func NewTasker(path string) *Tasker {
	return &Tasker{
		Dir: path,
	}
}
func (t *Tasker) getFiles() ([]string, error) {
	files, err := ioutil.ReadDir(t.Dir)
	if err != nil {
		return nil, err
	}
	pathes := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		pathes = append(pathes, filepath.Join(t.Dir, file.Name()))
	}
	return pathes, nil
}
func (t *Tasker) cd(path string) error {
	path = filepath.Join(t.Dir, path)
	fs, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}
	if !fs.IsDir() {
		return errors.New("not directory")
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}
	t.Dir = path
	return nil
}
func (t *Tasker) rename(oldName, newName string) (string, error) {
	oldPath := filepath.Join(t.Dir, oldName)
	newPath := filepath.Join(t.Dir, newName)
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return "", err
	}
	return newPath, nil
}
func (t *Tasker) runCommandWithOutput(task string) error {
	splitted := strings.Split(task, " ")
	executable := splitted[0]
	args := splitted[1:]
	command := exec.Command(executable, args...)
	// command.Stderr = os.Stderr
	// command.Stdout = os.Stdout
	command.Dir = t.Dir
	command.Env = os.Environ()
	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}
	// stderr, err := command.StderrPipe()
	// if err != nil {
	// 	return err
	// }
	// go scanOutput(task, true, stderr)
	go scanOutput(task, false, stdout)
	log.Println(task)
	err = command.Start()
	if err != nil {
		return err
	}

	err = command.Wait()
	if errors.Is(err, os.ErrProcessDone) {
		return nil
	}
	return err
}

func scanOutput(name string, isErr bool, output io.ReadCloser) {
	scanner := bufio.NewScanner(output)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		if isErr {
			log.Error(fmt.Sprintf("[%s] %s", name, m))
		}
		log.Info(fmt.Sprintf("[%s] %s", name, m))
	}
}
