package main

import (
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

const SERIAL_FILE = "/home/seeobject/公共的/serial"
const SECRET_FILE = "/home/seeobject/公共的/secret"

func main() {
	fmt.Println("hello")

	ser := getSerial()
	fmt.Println(ser)

	// 设置 gin 的模式（调试模式：DebugMode, 发行模式：ReleaseMode）
	//gin.SetMode(gin.DebugMode)
	// 创建一个不包含中间件的路由器
	//r := gin.Default()
	//r.Run()

	http.HandleFunc("/verify", verifyHandler)
	http.ListenAndServe("localhost:4587", nil)
}

func verifyHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "URL.Path = %q \n", req.URL.Path)

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
func writeFile(name, content string) {
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
		fmt.Println("PathExists err == nil")
		return true, nil
	}
	if os.IsNotExist(err) {
		fmt.Println("PathExists os.IsNotExist(err)")
		return false, nil
	}
	return false, err
}
