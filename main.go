package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"vs_export/sln"
)

// getFileType 检测文件类型（不区分大小写）
// 返回 "sln", "vcxproj" 或空字符串（不支持的类型）
func getFileType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".sln":
		return "sln"
	case ".vcxproj":
		return "vcxproj"
	default:
		return ""
	}
}

func main() {
	path := flag.String("s", "", "sln or vcxproj file path")
	configuration := flag.String("c", "Debug|Win32",
		"Configuration, [configuration|platform], default Debug|Win32")
	flag.Parse()

	if *path == "" {
		usage()
		os.Exit(1)
	}

	var solution sln.Sln
	var err error

	fileType := getFileType(*path)
	switch fileType {
	case "vcxproj":
		// 直接解析单个 .vcxproj 文件
		absPath, err := filepath.Abs(*path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		pro, err := sln.NewProject(*path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// 构造只包含单个项目的 Sln 实例
		// SolutionDir 设置为项目文件所在目录
		solution = sln.Sln{
			SolutionDir: filepath.Dir(absPath),
			ProjectList: []sln.Project{pro},
		}
	case "sln":
		// 使用现有的 .sln 解析逻辑
		solution, err = sln.NewSln(*path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "Error: Only .sln and .vcxproj files are supported")
		os.Exit(1)
	}
	cmdList, err := solution.CompileCommandsJson(*configuration)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	js, err := json.Marshal(cmdList)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", js[:])
	ioutil.WriteFile("compile_commands.json", js[:], 0644)
}


func usage() {
	var echo = `Usage: %s -s <path> -c <configuration>

Where:
            -s   path                        sln or vcxproj filename
            -c   configuration               project configuration,eg Debug|Win32.
                                             default Debug|Win32
	`
	echo = fmt.Sprintf(echo, filepath.Base(os.Args[0]))
	fmt.Println(echo)
}
