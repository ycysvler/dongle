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
)

const SERIAL_FILE = "./serial"
const SECRET_FILE = "./secret"

const PASSWORD = "abcd.1234"

type Result struct{
    Code int          `json:"code"`
    Message string    `json:"message"`
}

func main() {
	fmt.Println(os.Args)

	if len(os.Args) > 2{
		pwd := os.Args[2]
        cmd := os.Args[1]

		if pwd == PASSWORD{
            switch cmd{
				case "register" : fmt.Println("register")
                default : fmt.Printf("unknown command : %s \n", cmd)
   			}
		}else{
			fmt.Printf("what did you type in to amuse me ? %s", cmd)
		}

	}else{
		fmt.Printf("the service has been started.")
		http.HandleFunc("/verify", verifyHandler)
		http.ListenAndServe("localhost:4587", nil)
	}
    ser := getSerial()
	fmt.Printf("serial : %s",ser)
}

func verifyHandler(res http.ResponseWriter, req *http.Request) {
    v := verify()
    if v {
		t, _ := json.Marshal(Result{Code:200, Message:"ok"})
		fmt.Fprintf(res,"%s", t)
    }else{
		f, _ := json.Marshal(Result{Code:401, Message:"服务尚未注册！"})
		fmt.Fprintf(res,"%s", f)
    }
}

func verify() bool{
	return false
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
		writeFile(SERIAL_FILE, serial)
		result = serial
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
func writeFile(name, content string) string{
	data := []byte(content)

	err := ioutil.WriteFile(name, data, 0644)
	if err == nil {
		fmt.Println("写入文件成功:", content)
	} else {
		fmt.Println("write file", err)
	}
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
