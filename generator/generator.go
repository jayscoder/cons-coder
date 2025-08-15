package generator

import (
	"fmt"
	"path/filepath"
	"time"

	"cons-coder/parser"
)

// Config 生成器配置
type Config struct {
	Language    string // 目标语言
	OutputDir   string // 输出目录
	PackageName string // 包名
	HeaderComment string // 头部注释
	Version     string // 生成器版本
}

// Generator 代码生成器接口
type Generator interface {
	Generate(constants *parser.ConstantsFile) error
	GenerateIndex(allConstants []*parser.ConstantsFile) error
}

// BaseGenerator 基础生成器，包含通用功能
type BaseGenerator struct {
	Config Config
}

// New 创建对应语言的生成器
func New(config Config) Generator {
	switch config.Language {
	case "python":
		return NewPythonGenerator(config)
	case "go":
		return NewGoGenerator(config)
	case "java":
		return NewJavaGenerator(config)
	case "swift":
		return NewSwiftGenerator(config)
	case "kotlin":
		return NewKotlinGenerator(config)
	case "typescript":
		return NewTypeScriptGenerator(config)
	case "javascript":
		return NewJavaScriptGenerator(config)
	default:
		return nil
	}
}

// GetOutputFileName 获取输出文件名
func (g *BaseGenerator) GetOutputFileName(fileName string) string {
	switch g.Config.Language {
	case "python":
		return fileName + ".py"
	case "go":
		return fileName + ".go"
	case "java":
		return parser.ToJavaName(fileName) + ".java"
	case "swift":
		return parser.ToJavaName(fileName) + ".swift"
	case "kotlin":
		return parser.ToJavaName(fileName) + ".kt"
	case "typescript":
		return fileName + ".ts"
	case "javascript":
		return fileName + ".js"
	default:
		return fileName
	}
}

// GetOutputFilePath 获取完整的输出文件路径
func (g *BaseGenerator) GetOutputFilePath(fileName string) string {
	outputFileName := g.GetOutputFileName(fileName)
	return filepath.Join(g.Config.OutputDir, outputFileName)
}

// FormatGenerationTime 格式化生成时间
func FormatGenerationTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// GetFileHeader 获取文件头部注释
func (g *BaseGenerator) GetFileHeader(constants *parser.ConstantsFile) string {
	var commentStart, commentLine, commentEnd string
	
	switch g.Config.Language {
	case "python":
		commentStart = `"""`
		commentLine = ""
		commentEnd = `"""`
	case "go", "java", "swift", "kotlin", "typescript", "javascript":
		commentStart = "/*"
		commentLine = " * "
		commentEnd = " */"
	}

	header := fmt.Sprintf("%s\n", commentStart)
	
	// 添加自定义头部注释
	if g.Config.HeaderComment != "" {
		if commentLine != "" {
			header += fmt.Sprintf("%s%s\n", commentLine, g.Config.HeaderComment)
			header += fmt.Sprintf("%s\n", commentLine)
		} else {
			header += fmt.Sprintf("%s\n", g.Config.HeaderComment)
			header += "\n"
		}
	}
	
	if commentLine != "" {
		header += fmt.Sprintf("%s%s\n", commentLine, constants.Label)
		header += fmt.Sprintf("%s\n", commentLine)
		header += fmt.Sprintf("%s源文件: %s\n", commentLine, filepath.Base(constants.FilePath))
		header += fmt.Sprintf("%s最后修改: %s\n", commentLine, FormatGenerationTime(constants.LastModified))
		header += fmt.Sprintf("%s生成时间: %s\n", commentLine, FormatGenerationTime(time.Now()))
		header += fmt.Sprintf("%s生成工具: cons-coder v%s\n", commentLine, g.Config.Version)
	} else {
		header += fmt.Sprintf("%s\n", constants.Label)
		header += "\n"
		header += fmt.Sprintf("源文件: %s\n", filepath.Base(constants.FilePath))
		header += fmt.Sprintf("最后修改: %s\n", FormatGenerationTime(constants.LastModified))
		header += fmt.Sprintf("生成时间: %s\n", FormatGenerationTime(time.Now()))
		header += fmt.Sprintf("生成工具: cons-coder v%s\n", g.Config.Version)
	}
	header += fmt.Sprintf("%s\n", commentEnd)
	
	return header
}