package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func init() {
	flag.Parse()
}

// コンテキストの情報を取得
func getContext() string {
	context, err := exec.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		fmt.Println(err)
	}
	return string(context)
}

func displayInfo(context, command string) {
	fmt.Printf("#%s\n", strings.Repeat("-", 20))
	fmt.Printf("# Context: %s", context)
	fmt.Printf("# Command: %s\n", command)
	fmt.Printf("#%s\n", strings.Repeat("-", 20))
}

// コマンドを引数から取得する
func getCommand() []string {
	var cmd []string
	cmd = append(cmd, flag.Args()...)

	fmt.Println(cmd)

	if len(cmd) != 0 {
		cmd = append(cmd[:1], cmd[0:]...)
		cmd[0] = "kubectl"
	} else {
		cmd = append(cmd, "kubectl")
	}

	return cmd
}

// StdInで読み込んだ結果をファイルに保存
func readStdin() *os.File {
	tmpFile, err := ioutil.TempFile("", "tmp-")

	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	for {
		stdin, err := os.Stdin.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		tmpFile.Write(([]byte)(string(buf[:stdin])))
		if err == io.EOF {
			break
		}
	}
	return tmpFile
}

// 実行前の確認
func askForConfirmation() bool {
	// promptの出力
	fmt.Print("Are you sure? [Y/n]: ")

	// /dev/ttyを開く
	tty, err := os.Open("/dev/tty")
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	// bufferにキーボードからの入力を出力
	buf := make([]byte, 5)
	n, err := tty.Read(buf)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	answer := strings.ToLower(strings.TrimSpace(string(buf[:n])))

	switch answer {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Println("Input (Y)es or (N)o")
		return askForConfirmation()
	}
}

func execCommand(command string) string {
	out, err := exec.Command("sh", "-c", command).CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}

	return string(out)
}

func main() {
	stdinFileInfo, _ := os.Stdin.Stat()
	//前にパイプなし
	if (stdinFileInfo.Mode() & os.ModeCharDevice) != 0 {
		command := strings.Join(getCommand(), " ")
		displayInfo(getContext(), command)

		// 後ろにパイブが続く場合には、確認は行わない
		stdout, _ := os.Stdout.Stat()
		if stdout.Mode()&os.ModeNamedPipe == 0 {
			// 後ろにパイブがない場合
			if askForConfirmation() {
				fmt.Println(execCommand(command))
			}
		} else {
			// 後ろにパイブがある場合
			fmt.Println(execCommand(command))
		}
	} else {
		// 前にパイプあり
		// 標準入力からの読み取り
		command := getCommand()
		displayInfo(getContext(), strings.Join(getCommand(), " "))

		// パイプで渡された処理は一時ファイルに保存
		tmpFile := readStdin()
		// 最後の「-」をパイプで渡された内容のファイル名(tmpRile)に置換
		command[len(command)-1] = tmpFile.Name()

		if askForConfirmation() {
			commandStr := strings.Join(command, " ")
			fmt.Println(execCommand(commandStr))
		}
		// tmpFileのcloseと削除
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()
	}
}
