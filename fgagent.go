package main

/*
golang
web:         www.simonzhang.net
Email:       simon-zzm@163.com
*/

import (
	"archive/zip"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"net/http"
	"time"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IP地址白名单设置
var (
	white_list []string = []string{"192.168.5.28", "["}
	// 是否开启白名单功能。true为开启，false为关闭
	filter_ip_startus bool = false
)

func wlog(log_context string) {
	//fmt.Println("Write file")
	logfile, err := os.OpenFile("fgagent.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//fmt.Println(err)
	if err != nil {
		panic(err)
		os.Exit(-1)
	}
	defer logfile.Close()
    
    // 拼接日志信息，格式为 日期 日志文字
    _log_context := time.Now().Format("2006-01-02 15:04:05")+" "+log_context+"\r\n"
    _,err = logfile.WriteString(_log_context)
    if err != nil {
        panic(err)
    }
}

func unzip(src_zip string) string {
	// 解析解压包名
	dest := strings.Split(src_zip, ".")[0]
	// 打开压缩包
	unzip_file, err := zip.OpenReader(src_zip)
	if err != nil {
		wlog("压缩包损坏")
		return "压缩包损坏"
	}
	defer unzip_file.Close()
	// 创建解压目录
	os.MkdirAll(dest, 0755)
	// gb18030 转 utf-8
	enc := mahonia.NewDecoder("gb18030")
	// 循环解压zip文件
	for _, f := range unzip_file.File {
		rc, err := f.Open()
		if err != nil {
			return "压缩包中文件损坏"
		}
		defer rc.Close()
		path := filepath.Join(dest, enc.ConvertString(f.Name))
		// 判断解压出的是文件还是目录
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			// 创建解压文件
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return "创建本地文件失败"
			}
			defer f.Close()
			// 写入本地d
			_, err = io.Copy(f, rc)
			if err != nil {
				if err != io.EOF {
					return "写入本地失败"
				}
			}
			// end
		}
	}
	wlog("解压完成")
	return "OK"
}

func Upfile(w http.ResponseWriter, req *http.Request) {
	// 检查IP是否在白名单里
	if checkip(req) == false {
		w.Write([]byte("IP限制"))
		return
	}
	req.ParseForm()
	if req.Method == "POST" {
		wlog("开始处理POST信息")

		// 获取上传到的文件
		f, formHeader, err := req.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		defer f.Close()
		wlog("开始接收上传文件完毕。开始获取传输参数。")

		// 获取上传参数部分
		get_file_name := formHeader.Filename
		f_md5 := req.FormValue("md5")
		wlog("上传文件" + get_file_name + " md5:" + f_md5)
		// 创建本地存储文件
		to_file, err := os.Create(get_file_name)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		defer to_file.Close()
		wlog("上传文件" + get_file_name + "完成")
		// 存储到本地，后关闭相关文件
		io.Copy(to_file, f)
		if strings.Split(get_file_name, ".")[1] == "zip" {
			// 解压zip包
			get_unzip_status := unzip(get_file_name)
			if get_unzip_status != "OK" {
				w.Write([]byte("压缩文件处理错误"))
				return
			}
		}
		// 返回相关信息
		w.Write([]byte("OK"))
	} else {
		w.Write([]byte("Error"))
	}
	return
}

func commline(w http.ResponseWriter, req *http.Request) {
	// 检查IP是否在白名单里
	if checkip(req) == false { //; checkip_startus == 0 {
		w.Write([]byte("IP限制"))
		return
	}
	// 获取要执行的命令
	// 参数为cline
	// 格式为http://ip:port/comm/?cline=copy ttt\\test1.txt ttt\\test2.txt
	req.ParseForm()
	if req.Method == "GET" {
		cline := req.FormValue("cline")
		wlog("运行 " + cline)
		//_, err := exec.Command("cmd.exe", "/c", cline).Output()
		cmd := exec.Command("cmd.exe", "/c", cline)
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
		wlog("运行 " + cline + " 完成")
		w.Write([]byte("OK"))
	} else {
		w.Write([]byte("Error"))
	}
}

func Hello(w http.ResponseWriter, req *http.Request) {
	req.ParseForm() //解析参数，默认是不会解析的
	w.Write([]byte("my name is deploy agent!"))
}

// 检查白名单
func checkip(req *http.Request) bool {
	remote_ip := strings.Split(req.RemoteAddr, ":")[0]
	var checkstatus bool = false
	// 判断是否在白名单里
	if filter_ip_startus == true {
		for _, v := range white_list {
			if v == remote_ip {
				checkstatus = true
			}
		}
	} else {
		checkstatus = true
	}
	return checkstatus
}

/*
func test(w http.ResponseWriter, req *http.Request) {
	// 检查IP是否在白名单里
	if checkip(req) == false { //; checkip_startus == 0 {
		w.Write([]byte("IP限制"))
		return
	}
	fmt.Println("test")
}
*/
func main() {
	http.HandleFunc("/", Hello)         //根目录路由
	http.HandleFunc("/upfile/", Upfile) // 上传文件
	http.HandleFunc("/comm/", commline) // 需要执行的命令行
	//http.HandleFunc("/test/", test)          // 测试部分
	err := http.ListenAndServe(":8866", nil) //设置监听的端口
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
