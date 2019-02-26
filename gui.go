package main

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"path/filepath"
	"os"
	"strings"
	"log"
)

func main() {

	err := ui.Main(func() {
		window := ui.NewWindow("指纹图片一键处理V1.17", 500, 500, true)
		box := ui.NewVerticalBox()

		chooseButton := ui.NewButton("选择文件")
		checkbox := ui.NewCheckbox("批量处理整个文件夹")
		hbox := ui.NewHorizontalBox()
		hbox.Append(chooseButton, true)
		hbox.Append(checkbox, false)
		hbox.SetPadded(true)

		black := 60
		white := 50
		labelB := ui.NewLabel("黑")
		spinboxB := ui.NewSpinbox(0, 100)
		spinboxB.SetValue(black)
		sliderB := ui.NewSlider(0, 100)
		sliderB.SetValue(black)

		labelW := ui.NewLabel("白")
		spinboxW := ui.NewSpinbox(0, 100)
		spinboxW.SetValue(white)
		sliderW := ui.NewSlider(0, 100)
		sliderW.SetValue(white)

		vbox1 := ui.NewVerticalBox()
		hboxB := ui.NewHorizontalBox()
		hboxB.Append(labelB, false)
		hboxB.Append(spinboxB, true)
		hboxB.Append(sliderB, true)
		hboxB.SetPadded(true)
		hboxW := ui.NewHorizontalBox()
		hboxW.Append(labelW, false)
		hboxW.Append(spinboxW, true)
		hboxW.Append(sliderW, true)
		hboxW.SetPadded(true)

		vbox1.Append(hboxB, false)
		vbox1.Append(hboxW, false)
		vbox1.SetPadded(true)

		label := ui.NewLabel("待处理的文件")
		fileList := ui.NewEntry()
		fileList.SetReadOnly(false)
		vbox2 := ui.NewVerticalBox()
		vbox2.Append(label, true)
		vbox2.Append(fileList, true)
		vbox2.SetPadded(false)

		progressBar := ui.NewProgressBar()
		progressBar.SetValue(0)
		startButton := ui.NewButton("开始")
		vbox3 := ui.NewVerticalBox()
		vbox3.Append(progressBar, true)
		vbox3.Append(startButton, true)
		vbox3.SetPadded(true)

		//box.Append(ui.NewLabel("1.选择要处理的文件\n"+
		//	"2.选择是否要批量处理整个文件夹\n"+
		//	"3.调整参数或者使用默认值\n"+
		//	"4.点击开始\n"+
		//	"5.查看相关文件\n"), true)
		box.Append(hbox, true)
		box.Append(vbox1, true)
		box.Append(ui.NewHorizontalSeparator(), false)
		box.Append(vbox2, true)
		box.Append(ui.NewHorizontalSeparator(), false)
		box.Append(vbox3, true)
		box.SetPadded(true)

		/****************** logic ******************/
		spinboxB.OnChanged(func(*ui.Spinbox) {
			sliderB.SetValue(spinboxB.Value())
		})
		sliderB.OnChanged(func(*ui.Slider) {
			spinboxB.SetValue(sliderB.Value())
		})
		spinboxW.OnChanged(func(*ui.Spinbox) {
			sliderW.SetValue(spinboxW.Value())
		})
		sliderW.OnChanged(func(*ui.Slider) {
			spinboxW.SetValue(sliderW.Value())
		})

		fileNames := make([]string, 1)
		chooseButton.OnClicked(func(*ui.Button) {
			fileNames = showFiles(ui.OpenFile(window), fileList, checkbox)
		})
		checkbox.OnToggled(func(checkbox *ui.Checkbox) {
			fileNames = showFiles(fileNames[0], fileList, checkbox)
		})
		startButton.OnClicked(func(*ui.Button) {
			log.Println(fileNames)
			if fileNames[0] != "" {
				progress := 0.0
				for i, filePath := range fileNames {
					progress = float64(i+1.0) / float64(len(fileNames))
					p := int(progress * 100)
					progressBar.SetValue(p)
					ext := filepath.Ext(filePath)
					newPath := filePath[0:len(filePath)-len(ext)] + ".bmp"
					log.Println(newPath)
					Process(filePath, newPath, float64(spinboxB.Value())/100, float64(spinboxW.Value())/100)
				}
			}
			ui.MsgBox(window, "提示", "处理完成")
		})
		/****************** logic ******************/

		window.SetChild(box)
		window.SetMargined(true)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		ui.OnShouldQuit(func() bool {
			window.Destroy()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}

func isSupported(path string) bool {
	ext := filepath.Ext(strings.ToLower(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".bmp"
}

func showFiles(fileName string, fileList *ui.Entry, checkbox *ui.Checkbox) []string {
	fileNames := make([]string, 1)
	if fileName == "" {
		fileList.SetText("文件未选择")
		fileNames[0] = fileName
	} else if !isSupported(fileName) {
		fileList.SetText("文件格式不支持")
		fileNames[0] = fileName
	} else {
		if checkbox.Checked() {
			fileNames = make([]string, 0)
			err := filepath.Walk(filepath.Dir(fileName), func(path string, info os.FileInfo, err error) error {
				if isSupported(path) {
					fileNames = append(fileNames, path)
				}
				return nil
			})
			if err != nil {
				panic(err)
			}
		} else {
			fileNames[0] = fileName
		}
		names := ""
		for _, v := range fileNames {
			names += v + "; "
		}
		fileList.SetText(names)
	}
	return fileNames
}
