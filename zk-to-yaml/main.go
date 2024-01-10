package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"gopkg.in/yaml.v3"
)

func main() {
	// 定义命令行参数
	zkAddress := flag.String("zookeeper", "", "Zookeeper address")
	zkPath := flag.String("path", "/", "Zookeeper path")
	outputFilePath := flag.String("output", "output.yaml", "Output file path")
	flag.Parse()

	// 校验命令行参数
	if *zkAddress == "" {
		log.Fatal("Zookeeper address is required")
	}

	// 连接到Zookeeper
	conn, _, err := zk.Connect([]string{*zkAddress}, 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to Zookeeper: %v", err)
	}
	defer conn.Close()

	// 递归从Zookeeper中读取数据
	data, err := readFromZookeeper(conn, *zkPath)
	if err != nil {
		log.Fatalf("Failed to read from Zookeeper: %v", err)
	}

	// 将数据导出为YAML文件
	err = exportToYAML(data, *outputFilePath)
	if err != nil {
		log.Fatalf("Failed to export to YAML: %v", err)
	}

	fmt.Printf("Successfully exported data to %s\n", *outputFilePath)
}

func readFromZookeeper(conn *zk.Conn, path string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	children, _, err := conn.Children(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get children for path %s: %v", path, err)
	}

	for _, child := range children {
		var childPath = ""
		if path == "/" {
			childPath = path + child
		} else {
			childPath = path + "/" + child
		}
		childData, s, err := conn.Get(childPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get data for path %s: %v", childPath, err)
		}

		if len(childData) > 0 {
			str := string(childData[:s.DataLength-1])
			if strings.Contains(str, ",") {
				data[child] = strings.Split(str, ",")
			} else {
				data[child] = str
			}

		} else {
			subData, err := readFromZookeeper(conn, childPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read data for path %s: %v", childPath, err)
			}
			data[child] = subData
		}
	}

	return data, nil
}

func exportToYAML(data map[string]interface{}, filePath string) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to YAML: %v", err)
	}

	err = os.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write YAML file: %v", err)
	}

	return nil
}
