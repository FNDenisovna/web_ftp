package main

import (
	"fmt"
	"strings"
)

func getListDirView(fs []string, ds []string, fullDir string, customDir string) (content string) {
	contentList := make([]string, 0, 3)
	contentList = append(contentList, fmt.Sprintf("<h1>Содержимое дериктории %v:</h1>", fullDir))
	contentList = append(contentList, getRootView())
	if len(fs) == 0 && len(ds) == 0 {
		contentList = append(contentList, fmt.Sprintf("<h2>Пусто!</h2>"))
	}

	if len(fs) > 0 {
		contentList = append(contentList, fmt.Sprintf("<h2>Файлы:</h2>"))
		for _, row := range fs {
			//content += fmt.Sprintf("<a>%v</a><br>", row)
			contentList = append(contentList, getFileView(customDir, row))
		}
	}
	if len(ds) > 0 {
		contentList = append(contentList, fmt.Sprintf("<h2>Папки:</h2>"))
		for _, row := range ds {
			contentList = append(contentList, getDirView(customDir, row))
		}
	}

	contentList = append(contentList, getInputFields([]string{"name", "new_name"}, customDir, "command"))
	content = strings.Join(contentList, "")
	return
}

func getDirView(path string, name string) (content string) {
	content = fmt.Sprintf("<a href=\"http://localhost:%v/ls?dir=%v/%v\"\">%v</a><br>", port, path, name, name)
	return
}

func getFileView(path string, name string) (content string) {
	content = fmt.Sprintf("<a href=\"http://localhost:%v/file?dir=%v&name=%v\"\">%v</a><br>", port, path, name, name)
	return
}

func getRootView() (content string) {
	content = fmt.Sprintf("<a href=\"http://localhost:%v/ls\"\">Вернуться в корень</a><br>", port)
	return
}

func getInputFields(nameList []string, path string, command string) (content string) {
	content += fmt.Sprintf("<h2>Можно что-то изменить:</h2>")
	content += fmt.Sprintf("<form action=\"%v\" method=\"GET\">", command)
	for _, name := range nameList {
		content += fmt.Sprintf("<input autofocus value=\"\" name=\"%s\" type=\"text\">", name)
	}
	content += fmt.Sprintf("<input type=\"hidden\" value=\"%v\" name=\"dir\">", path)
	content += fmt.Sprintf("<input type=\"submit\" value=\"Create|Delete|Rename\" /></form><br>")
	content += fmt.Sprintf("<form action=\"%v\" method=\"POST\" enctype=\"multipart/form-data\">", command)
	content += fmt.Sprintf("<input type=\"hidden\" value=\"%v\" name=\"dir\">", path)
	content += fmt.Sprintf("<input type=\"file\" name=\"file\"/><input type=\"submit\" value=\"Load\" /></form>")
	return
}

func getButtonView(command string) (path string, content string) {
	content = fmt.Sprintf("<a href=\"http://localhost:%v/%s?dir=%v\"\">Create|Delete|Rename</a><br>", port, command, path)
	return
}
