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

// ファイルの存を在確認する
func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// コンテキストや実行コマンドの情報を表示する
func displayInfo(context, command string, tmpFile string) {
	fmt.Printf("#%s\n", strings.Repeat("-", 20))
	fmt.Printf("# Context: %s", context)
	fmt.Printf("# Command: %s\n", command)
	if fileExists(tmpFile) {
		fmt.Printf("# Manifest: %s\n", tmpFile)
	}
	fmt.Printf("#%s\n", strings.Repeat("-", 20))
}

// コマンドを引数から取得する
func getCommand() []string {
	var cmd []string
	cmd = append(cmd, flag.Args()...)

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
	out, _ := exec.Command("sh", "-c", command).CombinedOutput()
	return string(out)
}

// 文字列に配列に指定した単語が含まれるかを確認する
func confirmSentenceContainWords(sentence string, words []string) bool {
	sentences := strings.Split(sentence, " ")

	for _, word := range words {
		for _, sentence := range sentences[1:] {
			if sentence == word {
				return true
			}
		}
	}
	return false
}

// ファイルに記述されている言葉を配列に変換して返す
func convertFileWordsToArray(assetName string) []string {
	// file, err := Assets.Open(assetName)
	file, err := Assets.Open(assetName)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 64)
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			panic(err)
		}
	}

	words := strings.Split(string(buf), "\n")
	return words
}

func main() {
	// 標準入力からの読み取り
	command := getCommand()
	//除外コマンドのリスト作成
	excludeCommands := convertFileWordsToArray("/config/exclude_commands.conf")

	// 前にパイプなし
	stdinFileInfo, _ := os.Stdin.Stat()
	if (stdinFileInfo.Mode() & os.ModeCharDevice) != 0 {
		commandStr := strings.Join(command, " ")

		// excludeCommandsに含まれないコマンドの場合、確認を実行
		if confirmSentenceContainWords(commandStr, excludeCommands) == false {
			displayInfo(getContext(), commandStr, "")

			// 後ろにパイブが続く場合には、確認は行わない
			stdout, _ := os.Stdout.Stat()
			if stdout.Mode()&os.ModeNamedPipe == 0 {
				// 後ろにパイブがない場合
				if askForConfirmation() {
					fmt.Println(execCommand(commandStr))
				}
			}
		} else {
			// 後ろにパイブがある場合
			fmt.Println(execCommand(commandStr))
		}
	} else {
		// 前にパイプあり

		// パイプで渡された処理は一時ファイルに保存
		os.Setenv("TMPDIR", "/var/tmp")
		tmpFile := readStdin()
		// 最後の「-」をパイプで渡された内容のファイル名(tmpRile)に置換
		commandForInfo := strings.Join(command, " ")

		command[len(command)-1] = tmpFile.Name()
		commandStr := strings.Join(command, " ")

		// excludeCommandsに含まれないコマンドの場合、確認を実行
		if confirmSentenceContainWords(commandStr, excludeCommands) == false {

			displayInfo(getContext(), commandForInfo, tmpFile.Name())
			if askForConfirmation() {
				fmt.Println(execCommand(commandStr))
			}
		}

		// tmpFileのcloseと削除
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()
	}
}
