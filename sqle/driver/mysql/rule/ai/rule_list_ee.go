//go:build enterprise
// +build enterprise

package ai

import (
	"embed"
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"
)

//go:embed rule_knowledge_source_files
var ruleFiles embed.FS

const ruleListDir = "rule_knowledge_source_files"

// ReadEmbeddedXMLFiles 读取嵌入的 rule_knowledge_source_files 目录中的所有 XML 文件
func ReadEmbeddedXMLFiles() (map[string]string, error) {
	files := make(map[string]string)
	dirEntries, err := ruleFiles.ReadDir(ruleListDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".xml") {
			filePath := filepath.Join(ruleListDir, entry.Name())
			data, err := ruleFiles.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("read file [%v] failed: %v", filePath, err)
			}
			files[entry.Name()] = string(data)
		}
	}
	return files, nil
}

// XMLNode 定义XML节点结构
type XMLNode struct {
	XMLName  xml.Name
	Content  string     `xml:",chardata"`
	Children []XMLNode  `xml:",any"`
	Attrs    []xml.Attr `xml:",any,attr"`
}

// ProcessedNode 存储处理后的节点信息
type ProcessedNode struct {
	Title   string
	Content string
	Level   int
}

// ParseXMLContent 解析XML内容为XMLNode结构
func ParseXMLContent(xmlContent string) (*XMLNode, error) {
	var root XMLNode
	if err := xml.Unmarshal([]byte(xmlContent), &root); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %v", err)
	}
	return &root, nil
}

// ExtractNodes 从指定节点开始提取所有子节点
func ExtractNodes(node *XMLNode, targetTag string) []ProcessedNode {
	var nodes []ProcessedNode

	// 递归处理节点
	var processNode func(node XMLNode, level int)
	processNode = func(node XMLNode, level int) {
		// 如果节点有内容，添加到结果中
		content := node.Content
		if content != "" || len(node.Children) > 0 {
			nodes = append(nodes, ProcessedNode{
				Title:   node.XMLName.Local,
				Content: content,
				Level:   level,
			})
		}

		// 处理子节点
		for _, child := range node.Children {
			processNode(child, level+1)
		}
	}

	// 找到目标节点并开始处理
	var findAndProcess func(node XMLNode)
	findAndProcess = func(node XMLNode) {
		if node.XMLName.Local == targetTag {
			// 找到目标节点，开始处理其子节点
			for _, child := range node.Children {
				processNode(child, 1)
			}
			return
		}

		// 继续在子节点中查找目标节点
		for _, child := range node.Children {
			findAndProcess(child)
		}
	}

	findAndProcess(*node)
	return nodes
}

// GenerateMarkdown 根据节点列表生成Markdown文档
func GenerateMarkdown(nodes []ProcessedNode) string {
	var builder strings.Builder

	for _, node := range nodes {
		if node.Title == "检查流程描述" {
			continue
		}
		// 生成标题
		prefix := strings.Repeat("#", node.Level)
		builder.WriteString(fmt.Sprintf("%s %s\n", prefix, node.Title))

		// 添加内容（如果有）
		if node.Content != "" {
			builder.WriteString(node.Content + "\n\n")
		}
	}

	return builder.String()
}

// ConvertXMLToMarkdown 主函数，组合上述功能
func ConvertXMLToMarkdown(xmlContent string, targetTag string) (string, error) {
	// 1. 解析XML
	root, err := ParseXMLContent(xmlContent)
	if err != nil {
		return "", err
	}

	// 2. 提取节点
	nodes := ExtractNodes(root, targetTag)

	// 3. 生成Markdown
	markdown := GenerateMarkdown(nodes)

	return markdown, nil
}

func GetAIRulesKnowledge() (map[string]string, error) {
	res := make(map[string]string)
	xmlFiles, err := ReadEmbeddedXMLFiles()
	if err != nil {
		return nil, err
	}

	for fileName, content := range xmlFiles {
		ruleName := strings.TrimSuffix(fileName, "_MySQL.xml")
		md, err := ConvertXMLToMarkdown(content, "规则场景")
		if err != nil {
			return nil, fmt.Errorf("error converting %s to Markdown: %v", fileName, err)
		}
		res[ruleName] = md
	}

	return res, nil
}
