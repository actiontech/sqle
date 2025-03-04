//go:build enterprise
// +build enterprise

package ai

import (
	"embed"
	"encoding/xml"
	"fmt"
	"math"
	"path/filepath"
	"regexp"
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
	/*
		在提取子节点时，需要对原始XML文件的内容进行一些处理，以确保生成的Markdown文档的排版和显示正确。
		具体来说，我们需要对原始XML文件中的内容进行以下处理：
		1. 跳过检查流程描述，这部分内容不需要在Markdown文档中显示，因此可以跳过。
		2. 将行首的缩进转换为Markdown的代码块缩进，以确保生成的Markdown文档的代码块显示正确。
		3. 将规则场景的名称添加到场景标题中，将其他属性添加到标签中。
	*/
	processNode = func(node XMLNode, level int) {
		// 跳过检查流程描述
		if node.XMLName.Local == "检查流程描述" {
			return
		}
		if node.Content != "" || len(node.Children) > 0 {
			nodes = append(nodes, ProcessedNode{
				Title:   node.XMLName.Local,
				Content: removeMinIndent(node.Content),
				Level:   level,
			})
		}
		labels := []string{}
		for _, attr := range node.Attrs {
			// 名称的属性，添加到场景标题中
			if attr.Name.Local == "名称" {
				for idx := range nodes {
					if nodes[idx].Title == "场景" {
						nodes[idx].Title = fmt.Sprintf("%s:%s", nodes[idx].Title, attr.Value)
					}
				}
			} else {
				// 其他属性，添加到标签中
				labels = append(labels, fmt.Sprintf("%s:%s", attr.Name.Local, attr.Value))
			}
		}
		if len(labels) > 0 {
			// 添加标签，标签格式：```label [label1,label2]```
			nodes = append(nodes, ProcessedNode{
				Content: fmt.Sprintf("```label [%s]```", strings.Join(labels, ",")),
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

/*
在Markdown中，当行首缩进较多时，Markdown会将其视为代码块。
这可能会导致在某些情况下，行首的缩进被误认为是代码块的一部分，从而影响文档的排版和显示。
为了解决这个问题，我们可以使用正则表达式来去除行首的最小缩进。
这样，即使行首缩进不一致，Markdown也能正确地解析和渲染文档。
*/
func removeMinIndent(input string) string {
	// 分割成行
	lines := strings.Split(input, "\n")

	// 查找最小缩进（忽略零缩进的行）
	minIndent := math.MaxInt32
	re := regexp.MustCompile(`^( +)`)

	for _, line := range lines {
		// 忽略空行
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// 使用正则查找行首空格
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			leadingSpaces := len(matches[1])
			// 更新最小缩进: 忽略零缩进的行
			if leadingSpaces > 0 && leadingSpaces < minIndent {
				minIndent = leadingSpaces
			}
		}
	}

	// 如果没有找到有效的缩进或全部是零缩进
	if minIndent == math.MaxInt32 {
		return input
	}

	// 去除每行的最小缩进
	indentRegex := regexp.MustCompile(fmt.Sprintf(`^( {%d})`, minIndent))
	processedLines := make([]string, len(lines))

	for i, line := range lines {
		if len(strings.TrimSpace(line)) == 0 || !strings.HasPrefix(line, " ") {
			// 保持空行和无缩进行不变
			processedLines[i] = line
		} else {
			// 移除最小缩进
			processedLines[i] = indentRegex.ReplaceAllString(line, "")
		}
	}

	return strings.Join(processedLines, "\n")
}

// GenerateMarkdown 根据节点列表生成Markdown文档
func GenerateMarkdown(nodes []ProcessedNode) string {
	var builder strings.Builder
	for _, node := range nodes {
		// 生成标题
		if node.Title != "" {
			prefix := strings.Repeat("#", node.Level+1)
			builder.WriteString(fmt.Sprintf("%s %s\n", prefix, node.Title))
		}

		// 添加内容（如果有）
		if node.Content != "" {
			builder.WriteString(node.Content + "\n")
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
