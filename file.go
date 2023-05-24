package zutils

import (
	"io/ioutil"
	"os"
	"time"
)

//遍历文件夹，生成文件名集合
func ListDir(folder string) (fslice []string) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if !file.IsDir() {
			fslice = append(fslice, file.Name())
		}
	}
	return fslice
}

// 判断文件是否存在
func IsPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//获取指定目录下包含指定类型文件的所有文件夹,包含子目录
func ListAllDirs(dirPth string) (s []string, err error) {
	var dirs []string
	var tdirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			ListAllDirs(dirPth + PthSep + fi.Name())
		} else {
			tdirs = append(tdirs, dirPth+PthSep)
		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := ListAllDirs(table)
		for _, temp1 := range temp {
			//files = append(files, temp1)
			tdirs = append(tdirs, temp1)
		}
	}
	return tdirs, nil
}

//获取指定目录下的所有文件,包含子目录下的文件
func ListAllfiles(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			ListAllfiles(dirPth + PthSep + fi.Name())
		} else {
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := ListAllfiles(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}
	return files, nil
}

//目录切片去重
func RemoveSliceRepeat(s []string) (newslice []string) {
	newslice = make([]string, 0)
	for i := 0; i < len(s); i++ {
		repeat := false
		for j := i + 1; j < len(s); j++ {
			if s[i] == s[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newslice = append(newslice, s[i])
		}
	}
	return
}

// GetFileModTime 获取文件修改时间 返回unix时间戳
func GetFileModTime(path string) (time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		//fmt.Println("open file error")
		return time.Now(), err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		//fmt.Println("stat fileinfo error")
		return time.Now(), err
	}

	return fi.ModTime(), nil
}
