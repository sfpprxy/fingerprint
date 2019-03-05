package main

import (
	"github.com/disintegration/imaging"
	"hawx.me/code/img/levels"
	"hawx.me/code/img/channel"
	"image"
	"log"
	"os/exec"
	"strings"
	"path/filepath"
	"os"
	"io/ioutil"
)

func Process(srcPath string, newPath string, black float64, white float64) {
	src, err := imaging.Open(srcPath)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	optimize(src, newPath, black, white)
	vectorize(newPath)
}

func vectorize(fileName string) {
	pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Println(pwd)
	files, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Fatalf(err.Error())
	}
	name := ""
	// TODO: 这里通过文件夹名定位potrace程序, 目前是写死的, 应修改代码, 判断当前是哪种操作系统来自行选择win还是mac的potrace
	for _, v := range files {
		if strings.Contains(v.Name(), "potrace") {
			name = v.Name()
			log.Println(name)
		}
	}
	potracePath := filepath.Join(pwd, name, "potrace")
	cmd := exec.Command(potracePath, "-b", "dxf",
		fileName, "-o",
		strings.Replace(fileName, ".bmp", ".dxf", -1))
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Println(string(stdout))
}

func optimize(img image.Image, fileName string, black float64, white float64) {
	img = imaging.Grayscale(img)

	// Resize the cropped image to width = 200px preserving the aspect ratio.
	img = imaging.Resize(img, 1600, 0, imaging.MitchellNetravali)

	img = levels.SetBlack(img, channel.Red, black)
	img = levels.SetBlack(img, channel.Green, black)
	img = levels.SetBlack(img, channel.Blue, black)
	img = levels.SetWhite(img, channel.Red, white)
	img = levels.SetWhite(img, channel.Green, white)
	img = levels.SetWhite(img, channel.Blue, white)

	err := imaging.Save(img, fileName)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}
