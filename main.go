package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"path/filepath"
)

func main() {

	template := `
import 'package:flutter/widgets.dart';
/// 代码由程序自动生成。请不要对此文件做任何修改。
class IconFont {
  IconFont._();
  static const font_name = 'IconFont';
{icon_codes}
}
`
	fmt.Printf("Welcome to the iconfont2dart transformation tool. \n")

	filePath := readPath()
	bytes, e := ioutil.ReadFile(filePath)
	if e != nil {
		panic(e)
	}
	if len(bytes) == 0 {
		fmt.Println("iconfont.css is empty")
		return
	}
	iconfont := strings.Replace(string(bytes), "\n", "", -1)
	iconfont = strings.Replace(iconfont, " ", "", -1)
	compile := regexp.MustCompile(`\.icon-(.+?):before{content:"\\(\w+?)";}`)
	allStringSubmatch := compile.FindAllStringSubmatch(iconfont, -1)

	if len(allStringSubmatch) == 0 {
		fmt.Println(filePath + " is a valid iconfont.css")
		return
	}
	var constValSlice []string

	for _, item := range allStringSubmatch {
		iconName := strings.ToLower(strings.Replace(item[1],"-", "_", -1))
		constVal := fmt.Sprintf("  static const IconData %s = const IconData(0x%s, fontFamily: font_name);", iconName, item[2])
		constValSlice = append(constValSlice, constVal)

	}

	iconCodes := strings.Join(constValSlice, "\n")
	realPath, e := filepath.Abs(filepath.Dir(filePath))
	if e != nil {
		panic(e)
	}

	code := strings.Replace(template,"{icon_codes}", iconCodes,-1)
	if e := ioutil.WriteFile(filepath.Join(realPath,"icon_font.dart"), []byte(code), os.ModePerm);e!=nil{
		fmt.Println(code)
	}else{
		fmt.Println("ok,icon_font.dart file in ["+realPath+"] directory！")
	}

}

func readPath() string {
begin:
	filePath := ""
	fmt.Printf("Please enter iconfont.css path: \n")
	//Scanln 扫描来自标准输入的文本，将空格分隔的值依次存放到后续的参数内，直到碰到换行。
	if n, err := fmt.Scanln(&filePath); n == 0 || err != nil {
		fmt.Println("'iconfont.css' path is vaild")
		goto begin
	}

	file, e := os.Open(filePath)
	if e!= nil && os.IsNotExist(e) {
		fmt.Println(filePath + " is not exist")
		goto begin
	}
	e = file.Close()

	return filePath
}
