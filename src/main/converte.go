package main

import (
	"os"
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"io/ioutil"
	"strings"
)

type fileList struct {
	path string
	fileInfo os.FileInfo
}
func readDir(dir string) (data []fileList, err error) {
	//判断文件或目录是否存在
  file, err := os.Stat(dir)
	if err != nil {
		return data, err
	}

	//如果不是目录，直接返回文件信息
  if !file.IsDir() {
    data = append(data, fileList{path.Dir(dir) + "/", file})
    return data, err
  }

	fmt.Println(dir)
  fileInfo, err := ioutil.ReadDir(dir)
  if err != nil {
    fmt.Println(fileInfo)
    return data, err
  }

  //目录为空
  if len(fileInfo) == 0 {
    return
  }

  for _, v := range fileInfo {
    if v.IsDir() {
      if subDir, err := readDir(dir + "/" + v.Name()); err != nil {
        return data, err
      } else {
        data = append(data, subDir...)
      }
    } else {
      data = append(data, fileList{strings.TrimRight(dir, "/") + "/", v})
    }
  }
  return data, err
}

func convert(mp4in string, mp3out string) {
	cmd := exec.Command("ffmpeg", "-i", mp4in, mp3out)
	cmd.Stdout = os.Stdout
	out := cmd.Run()
	fmt.Println(out)
}

func main() {
	pwd, b := os.Getwd()
	fmt.Println(pwd)
	fmt.Println(b)
	filelist,err := readDir(pwd)
	if err != nil {
		fmt.Println("读取文件时发生错误")
		return 
	}
	for k, v := range filelist {
		fi := v.path + v.fileInfo.Name()
		fmt.Println(k, fi)
		dir, file := filepath.Split(fi)
		ext := filepath.Ext(file)
		fmt.Println("ext name:", ext)
		if ext != ".mp4" {
			fmt.Println("扩展名不是mp4，不处理：", ext)
			continue
		}

		mp4in := fi
		mp3out := filepath.Join(dir, filepath.Base(file) + ".mp3")
		fmt.Println("mp4in: ", mp4in)
		fmt.Println("mp3out: ", mp3out)

		convert(mp4in, mp3out)
	}
}