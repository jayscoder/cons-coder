package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Constant 表示单个常量定义
type Constant struct {
	Name  string      // 常量名称
	Type  string      // 数据类型 (int, string)
	Label string      // 中文标签/注释
	Value interface{} // 常量值
}

// ConstantGroup 表示一组常量
type ConstantGroup struct {
	Name      string      // 组名称
	Label     string      // 组描述
	Constants []*Constant // 常量列表
}

// ConstantsFile 表示解析后的完整文件信息
type ConstantsFile struct {
	FileName     string           // 文件名（不含扩展名）
	FilePath     string           // 原始文件路径
	Label        string           // 文件描述
	Groups       []*ConstantGroup // 常量组列表
	LastModified time.Time        // 文件最后修改时间
}

// ParseYAMLFile 解析单个YAML文件
func ParseYAMLFile(filePath string) (*ConstantsFile, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 提取文件名（不含扩展名）
	fileName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	// 解析YAML并提取注释
	label, constants, err := parseYAMLWithComments(data)
	if err != nil {
		return nil, fmt.Errorf("解析YAML失败: %w", err)
	}

	// 创建常量组
	group := &ConstantGroup{
		Name:      fileName,
		Label:     label,
		Constants: constants,
	}

	return &ConstantsFile{
		FileName:     fileName,
		FilePath:     filePath,
		Label:        label,
		Groups:       []*ConstantGroup{group},
		LastModified: fileInfo.ModTime(),
	}, nil
}

// parseYAMLWithComments 解析YAML文件并提取注释
func parseYAMLWithComments(data []byte) (string, []*Constant, error) {
	lines := strings.Split(string(data), "\n")
	
	var label string
	var constants []*Constant
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// 跳过空行
		if line == "" {
			continue
		}
		
		// 提取文件标签（第一行注释）
		if strings.HasPrefix(line, "#") && label == "" {
			label = strings.TrimSpace(strings.TrimPrefix(line, "#"))
			continue
		}
		
		// 解析常量行
		if !strings.HasPrefix(line, "#") && strings.Contains(line, ":") {
			constant, err := parseConstantLine(line)
			if err != nil {
				continue // 跳过无效行
			}
			constants = append(constants, constant)
		}
	}
	
	return label, constants, nil
}

// parseConstantLine 解析单行常量定义
func parseConstantLine(line string) (*Constant, error) {
	// 分割键值对和注释
	parts := strings.Split(line, "#")
	if len(parts) < 2 {
		return nil, fmt.Errorf("缺少注释")
	}
	
	// 解析键值对
	kvPart := strings.TrimSpace(parts[0])
	commentPart := strings.TrimSpace(parts[1])
	
	// 分割键和值
	kvParts := strings.SplitN(kvPart, ":", 2)
	if len(kvParts) != 2 {
		return nil, fmt.Errorf("无效的键值对格式")
	}
	
	name := strings.TrimSpace(kvParts[0])
	valueStr := strings.TrimSpace(kvParts[1])
	
	// 移除引号
	valueStr = strings.Trim(valueStr, `"'`)
	
	// 推断类型和解析值
	var value interface{}
	var dataType string
	
	// 尝试解析为整数
	if intVal, err := strconv.Atoi(valueStr); err == nil {
		value = intVal
		dataType = "int"
	} else {
		// 默认为字符串
		value = valueStr
		dataType = "string"
	}
	
	return &Constant{
		Name:  name,
		Type:  dataType,
		Label: commentPart,
		Value: value,
	}, nil
}

// ToGoName 将下划线命名转换为Go风格的驼峰命名
func ToGoName(name string) string {
	caser := cases.Title(language.English)
	parts := strings.Split(name, "_")
	for i, part := range parts {
		parts[i] = caser.String(strings.ToLower(part))
	}
	return strings.Join(parts, "")
}

// ToPythonName 将名称转换为Python风格的大写下划线命名
func ToPythonName(name string) string {
	return strings.ToUpper(name)
}

// ToJavaName 将下划线命名转换为Java风格的驼峰命名（首字母大写）
func ToJavaName(name string) string {
	return ToGoName(name)
}

// ToJavaConstantName 将名称转换为Java常量风格（全大写下划线）
func ToJavaConstantName(name string) string {
	return strings.ToUpper(name)
}

// ToSwiftName 将下划线命名转换为Swift风格的驼峰命名
func ToSwiftName(name string) string {
	parts := strings.Split(name, "_")
	caser := cases.Title(language.English)
	for i, part := range parts {
		parts[i] = caser.String(strings.ToLower(part))
	}
	return strings.Join(parts, "")
}

// ToKotlinName 与Java相同
func ToKotlinName(name string) string {
	return ToJavaName(name)
}

// ToKotlinConstantName 与Java常量相同
func ToKotlinConstantName(name string) string {
	return ToJavaConstantName(name)
}

// ToTypeScriptName 将下划线命名转换为TypeScript风格
func ToTypeScriptName(name string) string {
	return strings.ToUpper(name)
}

// ToJavaScriptName 与TypeScript相同
func ToJavaScriptName(name string) string {
	return ToTypeScriptName(name)
}

// GetGoType 获取Go语言对应的类型
func GetGoType(dataType string) string {
	switch dataType {
	case "int":
		return "int"
	case "string":
		return "string"
	case "float":
		return "float64"
	case "bool":
		return "bool"
	default:
		return "interface{}"
	}
}

// GetPythonType 获取Python语言对应的类型
func GetPythonType(dataType string) string {
	switch dataType {
	case "int":
		return "int"
	case "string":
		return "str"
	case "float":
		return "float"
	case "bool":
		return "bool"
	default:
		return "Any"
	}
}

// GetJavaType 获取Java语言对应的类型
func GetJavaType(dataType string) string {
	switch dataType {
	case "int":
		return "int"
	case "string":
		return "String"
	case "float":
		return "double"
	case "bool":
		return "boolean"
	default:
		return "Object"
	}
}

// GetSwiftType 获取Swift语言对应的类型
func GetSwiftType(dataType string) string {
	switch dataType {
	case "int":
		return "Int"
	case "string":
		return "String"
	case "float":
		return "Double"
	case "bool":
		return "Bool"
	default:
		return "Any"
	}
}

// GetKotlinType 获取Kotlin语言对应的类型
func GetKotlinType(dataType string) string {
	switch dataType {
	case "int":
		return "Int"
	case "string":
		return "String"
	case "float":
		return "Double"
	case "bool":
		return "Boolean"
	default:
		return "Any"
	}
}

// GetTypeScriptType 获取TypeScript语言对应的类型
func GetTypeScriptType(dataType string) string {
	switch dataType {
	case "int", "float":
		return "number"
	case "string":
		return "string"
	case "bool":
		return "boolean"
	default:
		return "any"
	}
}

// GetJavaScriptType 获取JavaScript语言对应的类型（用于JSDoc）
func GetJavaScriptType(dataType string) string {
	return GetTypeScriptType(dataType)
}

// FormatValue 根据类型格式化值
func FormatValue(value interface{}, dataType string, lang string) string {
	caser := cases.Title(language.English)
	valueStr := fmt.Sprintf("%v", value)
	
	switch lang {
	case "python":
		if dataType == "string" {
			return fmt.Sprintf(`"%s"`, valueStr)
		}
		if dataType == "bool" {
			return caser.String(strings.ToLower(valueStr))
		}
		return valueStr
	case "go":
		if dataType == "string" {
			return fmt.Sprintf(`"%s"`, valueStr)
		}
		return valueStr
	case "java", "kotlin":
		if dataType == "string" {
			return fmt.Sprintf(`"%s"`, valueStr)
		}
		return valueStr
	case "swift":
		if dataType == "string" {
			return fmt.Sprintf(`"%s"`, valueStr)
		}
		return valueStr
	case "typescript", "javascript":
		if dataType == "string" {
			return fmt.Sprintf(`'%s'`, valueStr)
		}
		return valueStr
	default:
		return valueStr
	}
}
