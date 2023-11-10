package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// 读取hosts文件内容
func readHostsFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var content []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return content, nil
}

// 在内存中替换行
func replaceLine(content []string, oldLine, newLine string) []string {
	for i, line := range content {
		if strings.Contains(line, " github.com") {
			content[i] = newLine
			break
		}
	}

	content = append(content, newLine)
	return content
}

// 将修改后的内容写回hosts文件
func writeHostsFile(filePath string, content []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range content {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func main() {

	// 要请求的URL
	url := "https://sites.ipaddress.com/github.com/"

	pattern := `IN\s+A\s+<a\s+href="https://www\.ipaddress\.com/ipv4/(\d+\.\d+\.\d+\.\d+)">`

	// 发送GET请求
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer response.Body.Close()

	// IN  A  <a href="https://www.ipaddress.com/ipv4/140.82.114.4">
	// 读取响应的内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}
	// 编译正则表达式
	regex, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("正则表达式编译失败:", err)
		return
	}

	// 在文本中查找匹配的字符串
	matches := regex.FindStringSubmatch(string(body))
	// 打印匹配结果
	if len(matches) > 1 {
		ipAddress := matches[1]
		fmt.Println("匹配的IPv4地址:", ipAddress)
		hostPath := "C:\\Windows\\System32\\drivers\\etc\\hosts"
		cotent := ipAddress + " github.com"
		// 读取hosts文件内容
		content, err := readHostsFile(hostPath)
		if err != nil {
			fmt.Println("读取hosts文件失败:", err)
			return
		}
		// 在内存中替换行
		newContent := replaceLine(content, " github.com", cotent)
		// 将修改后的内容写回hosts文件
		err = writeHostsFile(hostPath, newContent)
		if err != nil {
			fmt.Println("写入hosts文件失败:", err)
			return
		}

		fmt.Println("成功替换hosts文件中的行.")

	} else {
		fmt.Println("未找到匹配")
	}

}
