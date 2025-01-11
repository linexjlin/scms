package main

import (
	"os"
	"path/filepath"
	"strings"
)

func listFiles(directory string, depth int, exts string) ([]string, error) {
	var files []string

	// 将扩展名字符串分割成切片
	extList := strings.Split(exts, ",")

	// 遍历目录
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算当前文件的深度
		relPath, err := filepath.Rel(directory, path)
		if err != nil {
			return err
		}
		currentDepth := len(strings.Split(relPath, string(filepath.Separator))) - 1

		// 如果当前深度超过指定深度，则跳过
		if depth != 0 && currentDepth > depth {
			return nil
		}

		// 检查文件扩展名是否匹配
		for _, ext := range extList {
			if matched, _ := filepath.Match(ext, filepath.Ext(path)); matched {
				files = append(files, path)
				break
			}
		}

		return nil
	})

	return files, err
}
