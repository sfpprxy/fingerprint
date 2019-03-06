package main

import (
	"log"
	"os"
	"strconv"
	"path/filepath"
)

/*
	命令行版本：
	参数0：文件路径或者目录路径
	参数1：黑色值
	参数2：白色值
	参数3：一个写死的随机字符串密码(防止随意盗用程序)

	例子1：/Users/joe/Documents 60 40 qw5e4r6t
	例子2：D:\工作\图片\指纹1.jpg 50 30 qw5e4r6t

	例子1批量处理这个目录下的所有位图文件，例子2处理这个目录下指定的单个文件

	健壮性要求：
	1. 对不合法的输入进行检测并给出相应提示
	2. 密码错误无提示直接退出程序
	3. 继续保持对OS X和Windows的兼容性
*/
func HandleInput(args []string)  {
	var path, token, b, w string
	if len(args) != 4 {
		log.Fatalln("err code: 1 - invalid arg number") // log.Fatalln will exit program

	}
	// exam token
	token = args[3]
	if token != "foobar" {
		log.Fatalln("wrong token") // no return and exit
	}
	// exam path
	path = args[0]
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatalln("err code: 2 - path not valid")
	}
	// exam B & W
	b = args[1]
	w = args[2]
	err3 := "err code: 3 - black arg not valid"
	ib, err := strconv.ParseInt(b, 0, 64)
	if err != nil {
		log.Fatalln(err3)
	}
	if ib < 0 || ib > 100 {
		log.Fatalln(err3)
	}
	iw, err := strconv.ParseInt(w, 0, 64)
	if err != nil {
		log.Fatalln(err3)
	}
	if iw < 0 || iw > 100 {
		log.Fatalln(err3)
	}

	// pass exam and start process
	switch mode := fi.Mode(); {
	case mode.IsDir():
		log.Println("directory")
		fileNames := make([]string, 0)
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if IsSupported(filePath) {
				fileNames = append(fileNames, filePath)
			}
			return nil
		})
		if err != nil {
			log.Fatalln(err)
		}
		for _, f := range fileNames {
			Process(f, ToBmp(f), float64(ib), float64(iw))
		}
	case mode.IsRegular():
		log.Println("file")
		Process(path, ToBmp(path), float64(ib), float64(iw))
	}

}
