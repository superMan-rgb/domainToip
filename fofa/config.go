package fofa

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	configFileName = "config.ini"
	emailKey       = "email"
	keyKey         = "key"
)

// 创建文件
func IsConfig() {
	_, err := os.Stat(configFileName)
	if os.IsNotExist(err) {
		fmt.Printf("配置文件 %s 不存在，正在创建...\n", configFileName)

		err := createConfigFile()
		if err != nil {
			fmt.Println("无法创建配置文件:", err)
			return
		}

		fmt.Printf("配置文件 %s 已创建，请修改其中的参数.\n", configFileName)
		return
	}
}

// 检查配置文件是否存在
func Config() (string, string) {

	// 读取配置文件
	config, err := readConfigFile()
	if err != nil {
		fmt.Println("无法读取配置文件:", err)
		return "", ""
	}

	// 获取email和key参数
	email, ok := config[emailKey]
	if !ok {
		fmt.Printf("配置文件中缺少 %s 参数\n", emailKey)
		return "", ""
	}

	key, ok := config[keyKey]
	if !ok {
		fmt.Printf("配置文件中缺少 %s 参数\n", keyKey)
		return "", ""
	}
	return email, key
}

func createConfigFile() error {
	file, err := os.Create(configFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	writer.WriteString(emailKey + " = \n")
	writer.WriteString(keyKey + " = \n")

	return nil
}

func readConfigFile() (map[string]string, error) {
	config := make(map[string]string)

	file, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}
