package main

import (
	"encoding/json"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"bufio"
)

const SERIAL_FILE = "./serial"
const SECRET_FILE = "./secret"

const PASSWORD = "abcd.1234"

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	if len(os.Args) == 2{
		cmd := os.Args[1]
		if cmd == "h" || cmd == "help" || cmd == "?"{
			fmt.Printf("secret   <pwd> \t<serial>\n")
			fmt.Printf("register <pwd> \n")
		}
	} else if len(os.Args) > 2 {
		pwd := os.Args[2]
		cmd := os.Args[1]

		if pwd == PASSWORD {
			switch cmd {
			case "help":
				fmt.Printf("secret   <pwd> \t<serial>\n")
				fmt.Printf("register <pwd> \n")
			case "register":
				register()
				startServe()
			case "secret":
				serial := os.Args[3]
				fmt.Printf("%s\n", getSecret(serial))

			default:
				fmt.Printf("unknown command : %s \n", cmd)
			}
		} else {
			fmt.Printf("what did you type in to amuse me ?\n")
		}
	} else {

		v := verify()

		if v {
			startServe()
		} else {
			fmt.Printf("服务还未注册！\n")
			fmt.Printf("服务序列号：%s \n", getSerial())

			for {
				fmt.Printf("请输入注册码 > ")
				input := bufio.NewScanner(os.Stdin)
				input.Scan()
				secret := input.Text()

				if secret == getSecret(getSerial()) {
					writeFile(SECRET_FILE, secret)
					fmt.Printf("服务注册成功！\n")
					startServe()
				} else {
					fmt.Printf("注册码不正确！\n")
				}
			}
		}

	}
}

func verifyHandler(res http.ResponseWriter, req *http.Request) {
	v := verify()
	if v {
		t, _ := json.Marshal(Result{Code: 200, Message: "ok"})
		fmt.Fprintf(res, "%s", t)
	} else {
		f, _ := json.Marshal(Result{Code: 401, Message: "服务尚未注册！"})
		fmt.Fprintf(res, "%s", f)
	}
}

/*
启动Web服务
*/
func startServe() {
	fmt.Printf("the service has been started.")
	http.HandleFunc("/verify", verifyHandler)
	http.ListenAndServe("localhost:4587", nil)
}

/*
注册秘钥
*/
func register() {
	serial := getSerial()
	secret := getSecret(serial)

	err := writeFile(SECRET_FILE, secret)
	if err == nil {
		// 注册成功
		fmt.Println("register complete.")
	} else {
		// 注册失败
		fmt.Println(serial)
	}
}

/*
验证秘钥
*/
func verify() bool {
	exist, _ := pathExists(SECRET_FILE)
	if exist {
		secret := readFile(SECRET_FILE)
		serial := getSerial()

		if secret == getSecret(serial) {
			// 秘钥结果很完美
			return true
		} else {
			// 秘钥内容不正确
			return false
		}
	} else {
		// 秘钥文件不存在
		return false
	}
}

/*
byte 转字符串
*/
func convert(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, "-")
}

/*
字符串转MD5
*/
func md5V1(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

/*
用序列号计算秘钥
*/
func getSecret(serial string) string {
	return md5V1(serial + "-seeobject")
}

/*
获取本机序列号
*/
func getSerial() string {
	exist, _ := pathExists(SERIAL_FILE)

	result := ""
	if exist {
		result = readFile(SERIAL_FILE)
	} else {
		b := make([]byte, 5)
		rand.Read(b)
		serial := convert(b)
		err := writeFile(SERIAL_FILE, serial)
		if err == nil {
			result = serial
		} else {
			fmt.Println(err)
		}
	}
	return result
}

/*
读文件
*/
func readFile(name string) string {
	if contents, err := ioutil.ReadFile(name); err == nil {
		//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
		result := strings.Replace(string(contents), "\n", "", 1)
		return result
	} else {
		return ""
	}
}

/*
写文件
*/
func writeFile(name, content string) error {
	data := []byte(content)
	return ioutil.WriteFile(name, data, 0644)
}

/*
判断文件是否存在
*/
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
