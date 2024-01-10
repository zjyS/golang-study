package main

import (
	"flag"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// 定义命令行参数
	yamlFilePath := flag.String("yaml", "", "Path to YAML file")
	zkAddress := flag.String("zookeeper", "", "Zookeeper address")
	zkPath := flag.String("path", "/", "Zookeeper path")
	flag.Parse()

	// 校验命令行参数
	if *yamlFilePath == "" {
		log.Fatal("YAML file path is required")
	}
	if *zkAddress == "" {
		log.Fatal("Zookeeper address is required")
	}

	// 读取yaml文件
	yamlFile, err := os.ReadFile(*yamlFilePath)
	if err != nil {
		log.Fatalf("Failed to read yaml file: %v", err)
	}

	// 解析yaml文件内容
	data := make(map[interface{}]interface{})
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatalf("Failed to unmarshal yaml file: %v", err)
	}

	// 连接到Zookeeper
	conn, _, err := zk.Connect([]string{*zkAddress}, 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to Zookeeper: %v", err)
	}
	defer conn.Close()

	// 创建路径
	if err := createZookeeperPath(conn, *zkPath); err != nil {
		log.Fatalf("Failed to create Zookeeper path: %v", err)
	}

	// 递归将yaml内容写入到Zookeeper中
	if err := writeToZookeeper(conn, *zkPath, data); err != nil {
		log.Fatalf("Failed to write to Zookeeper: %v", err)
	}

	fmt.Println("Successfully wrote yaml content to Zookeeper")
}

func createZookeeperPath(conn *zk.Conn, path string) error {
	parts := strings.Split(path, "/")
	parts = parts[1:] // Skip the first empty part
	for i := range parts {
		subPath := "/" + strings.Join(parts[:i+1], "/")
		exists, _, err := conn.Exists(subPath)
		if err != nil {
			return fmt.Errorf("failed to check if path exists: %v", err)
		}
		if !exists {
			_, err := conn.Create(subPath, nil, 0, zk.WorldACL(zk.PermAll))
			if err != nil {
				return fmt.Errorf("failed to create path %s: %v", subPath, err)
			}
		}
	}
	return nil
}

func writeToZookeeper(conn *zk.Conn, path string, data map[interface{}]interface{}) error {
	for k, v := range data {
		var nodePath = ""
		if path == "/" {
			nodePath = fmt.Sprint(path, k)
		} else {
			nodePath = fmt.Sprint(path, "/", k)
		}
		switch value := v.(type) {
		case map[string]interface{}:
			if err := createZookeeperPath(conn, nodePath); err != nil {
				return fmt.Errorf("failed to create node %s: %v", nodePath, err)
			}
			if err := writeToZookeeper(conn, nodePath, transferMap(value)); err != nil {
				return err
			}
		case map[interface{}]interface{}:
			if err := createZookeeperPath(conn, nodePath); err != nil {
				return fmt.Errorf("failed to create node %s: %v", nodePath, err)
			}
			if err := writeToZookeeper(conn, nodePath, value); err != nil {
				return err
			}

		case []interface{}:
			var builder strings.Builder
			for i, d := range value {
				builder.WriteString(fmt.Sprint(d))
				if i < len(value)-1 {
					builder.WriteString(",")
				}
			}
			valueBytes, err := yaml.Marshal(builder.String())
			if err != nil {
				return fmt.Errorf("failed to marshal value: %v", err)
			}
			updateNode(conn, nodePath, valueBytes)

		default:
			valueBytes, err := yaml.Marshal(value)
			if err != nil {
				return fmt.Errorf("failed to marshal value: %v", err)
			}
			updateNode(conn, nodePath, valueBytes)
		}
	}
	return nil
}

func transferMap(strMap map[string]interface{}) map[interface{}]interface{} {
	newMap := make(map[interface{}]interface{}, len(strMap))
	for k, v := range strMap {
		newMap[k] = v
	}
	return newMap
}

func updateNode(conn *zk.Conn, nodePath string, data []byte) {
	exists, s, err := conn.Exists(nodePath)
	if err != nil {
		_ = fmt.Errorf("failed to check if path exists: %v", err)
	}
	if !exists {
		if _, err := conn.Create(nodePath, data, 0, zk.WorldACL(zk.PermAll)); err != nil {
			_ = fmt.Errorf("failed to create node %s: %v", nodePath, err)
		} else {
			fmt.Printf("set node %s --> %s", nodePath, data)
		}
	} else {
		if _, err = conn.Set(nodePath, data, s.Version); err != nil {
			_ = fmt.Errorf("failed to set node %s: %v", nodePath, err)
		} else {
			fmt.Printf("update node %s --> %s", nodePath, data)
		}
	}
}
