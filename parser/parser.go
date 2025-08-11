package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Constant 表示单个常量定义
type Constant struct {
	XMLName xml.Name `xml:""`
	Name    string   // 从XML标签名提取
	Type    string   `xml:"type,attr"`
	Label   string   `xml:"label,attr"`
	Desc    string   `xml:"desc,attr"`
	Value   string   `xml:"value,attr"`
}

// ConstantGroup 表示一组常量
type ConstantGroup struct {
	XMLName   xml.Name    `xml:""`
	Name      string      // 从XML标签名提取
	Label     string      `xml:"label,attr"`
	Constants []*Constant `xml:",any"`
}

// Constants 表示XML根元素
type Constants struct {
	XMLName xml.Name         `xml:"constants"`
	Label   string           `xml:"label,attr"`
	Groups  []*ConstantGroup `xml:",any"`
}

// ConstantsFile 表示解析后的完整文件信息
type ConstantsFile struct {
	FileName     string           // 文件名（不含扩展名）
	FilePath     string           // 原始文件路径
	Label        string           // 文件描述
	Groups       []*ConstantGroup // 常量组列表
	LastModified time.Time        // 文件最后修改时间
}

// ParseXMLFile 解析单个XML文件
func ParseXMLFile(filePath string) (*ConstantsFile, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 解析XML
	var constants Constants
	if err := xml.Unmarshal(data, &constants); err != nil {
		return nil, fmt.Errorf("解析XML失败: %w", err)
	}

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 提取文件名（不含扩展名）
	fileName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	// 处理常量组和常量名称
	for _, group := range constants.Groups {
		// 设置组名（从XML标签名提取）
		group.Name = group.XMLName.Local

		// 处理组内的常量
		for _, constant := range group.Constants {
			// 设置常量名（从XML标签名提取）
			constant.Name = constant.XMLName.Local
		}
	}

	return &ConstantsFile{
		FileName:     fileName,
		FilePath:     filePath,
		Label:        constants.Label,
		Groups:       constants.Groups,
		LastModified: fileInfo.ModTime(),
	}, nil
}

// ToGoName 将下划线命名转换为Go风格的驼峰命名
func ToGoName(name string) string {
	parts := strings.Split(name, "_")
	for i, part := range parts {
		parts[i] = strings.Title(strings.ToLower(part))
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
	for i, part := range parts {
		parts[i] = strings.Title(strings.ToLower(part))
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
func GetGoType(xmlType string) string {
	switch xmlType {
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
func GetPythonType(xmlType string) string {
	switch xmlType {
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
func GetJavaType(xmlType string) string {
	switch xmlType {
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
func GetSwiftType(xmlType string) string {
	switch xmlType {
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
func GetKotlinType(xmlType string) string {
	switch xmlType {
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
func GetTypeScriptType(xmlType string) string {
	switch xmlType {
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
func GetJavaScriptType(xmlType string) string {
	return GetTypeScriptType(xmlType)
}

// FormatValue 根据类型格式化值
func FormatValue(value string, xmlType string, lang string) string {
	switch lang {
	case "python":
		if xmlType == "string" {
			return fmt.Sprintf(`"%s"`, value)
		}
		if xmlType == "bool" {
			return strings.Title(strings.ToLower(value))
		}
		return value
	case "go":
		if xmlType == "string" {
			return fmt.Sprintf(`"%s"`, value)
		}
		return value
	case "java", "kotlin":
		if xmlType == "string" {
			return fmt.Sprintf(`"%s"`, value)
		}
		return value
	case "swift":
		if xmlType == "string" {
			return fmt.Sprintf(`"%s"`, value)
		}
		return value
	case "typescript", "javascript":
		if xmlType == "string" {
			return fmt.Sprintf(`'%s'`, value)
		}
		return value
	default:
		return value
	}
}
