package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"cons-coder/parser"
)

// TypeScriptGenerator TypeScript代码生成器
type TypeScriptGenerator struct {
	BaseGenerator
}

// NewTypeScriptGenerator 创建TypeScript生成器
func NewTypeScriptGenerator(config Config) *TypeScriptGenerator {
	return &TypeScriptGenerator{
		BaseGenerator: BaseGenerator{Config: config},
	}
}

// Generate 生成TypeScript代码
func (g *TypeScriptGenerator) Generate(constants *parser.ConstantsFile) error {
	var code strings.Builder
	
	// 文件头注释
	code.WriteString(g.GetFileHeader(constants))
	code.WriteString("\n")
	
	if g.Config.Mode == "const" {
		// const模式 - 生成简单常量
		for _, group := range constants.Groups {
			code.WriteString(g.generateConstGroup(group, constants.Label))
			code.WriteString("\n")
		}
	} else {
		// class模式
		// 生成每个常量组的类
		for _, group := range constants.Groups {
			code.WriteString(g.generateGroupClass(group, constants.Label))
			code.WriteString("\n")
		}
	}
	
	// 写入文件
	outputPath := g.GetOutputFilePath(constants.FileName)
	return os.WriteFile(outputPath, []byte(code.String()), 0644)
}

// generateConstGroup 生成const模式的常量组
func (g *TypeScriptGenerator) generateConstGroup(group *parser.ConstantGroup, projectLabel string) string {
	var code strings.Builder
	
	// 生成注释
	code.WriteString(fmt.Sprintf("// %s %s - %s\n", group.Name, group.Label, projectLabel))
	
	// 按字母顺序排序常量
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToTypeScriptName(constants[i].Name) < parser.ToTypeScriptName(constants[j].Name)
	})
	
	// 生成常量定义
	for _, constant := range constants {
		constName := fmt.Sprintf("%s_%s", strings.ToUpper(group.Name), strings.ToUpper(constant.Name))
		value := parser.FormatValue(constant.Value, constant.Type, "typescript")
		comment := constant.Label
		code.WriteString(fmt.Sprintf("export const %s = %s; // %s\n", constName, value, comment))
	}
	
	return code.String()
}

// generateGroupClass 生成常量组（简化格式）
func (g *TypeScriptGenerator) generateGroupClass(group *parser.ConstantGroup, _ string) string {
	var code strings.Builder
	
	className := parser.ToJavaName(group.Name)
	
	// 按字母顺序排序常量
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToTypeScriptName(constants[i].Name) < parser.ToTypeScriptName(constants[j].Name)
	})
	
	// 生成常量定义
	code.WriteString(fmt.Sprintf("export const %s = {\n", className))
	
	// 生成常量值
	for _, constant := range constants {
		fieldName := strings.ToUpper(constant.Name)
		value := parser.FormatValue(constant.Value, constant.Type, "typescript")
		comment := constant.Label
		if comment == "" {
			comment = constant.Name
		}
		code.WriteString(fmt.Sprintf("  /** %s */\n", comment))
		code.WriteString(fmt.Sprintf("  %s: %s,\n", fieldName, value))
	}
	
	code.WriteString("} as const;\n\n")
	
	// 生成类型定义
	code.WriteString(fmt.Sprintf("export type %sValue = typeof %s[keyof typeof %s];\n", className, className, className))
	code.WriteString(fmt.Sprintf("export type %sKey = keyof typeof %s;", className, className))
	
	return code.String()
}



// GenerateIndex 生成TypeScript的index.ts文件
func (g *TypeScriptGenerator) GenerateIndex(allConstants []*parser.ConstantsFile) error {
	var code strings.Builder
	
	// 文件头注释
	code.WriteString("/**\n")
	code.WriteString(" * 常量包索引文件\n")
	code.WriteString(" * \n")
	code.WriteString(fmt.Sprintf(" * 生成时间: %s\n", FormatGenerationTime(time.Now())))
	code.WriteString(fmt.Sprintf(" * 生成工具: cons-coder v%s\n", g.Config.Version))
	code.WriteString(" */\n\n")
	
	// 导出所有文件
	for _, constants := range allConstants {
		code.WriteString(fmt.Sprintf("export * from './%s';\n", constants.FileName))
	}
	
	// 写入文件
	outputPath := filepath.Join(g.Config.OutputDir, "index.ts")
	return os.WriteFile(outputPath, []byte(code.String()), 0644)
}