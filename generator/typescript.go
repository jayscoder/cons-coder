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
		
		// 生成类型定义
		for _, group := range constants.Groups {
			code.WriteString(g.generateTypeDefinitions(group))
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

// generateGroupClass 生成常量组类
func (g *TypeScriptGenerator) generateGroupClass(group *parser.ConstantGroup, projectLabel string) string {
	var code strings.Builder
	
	className := parser.ToJavaName(group.Name)
	
	// 类注释
	code.WriteString("/**\n")
	code.WriteString(fmt.Sprintf(" * %s - %s\n", group.Label, projectLabel))
	code.WriteString(" */\n")
	code.WriteString(fmt.Sprintf("export class %s {\n", className))
	
	// 按字母顺序排序常量
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToTypeScriptName(constants[i].Name) < parser.ToTypeScriptName(constants[j].Name)
	})
	
	// 常量定义
	for _, constant := range constants {
		constName := parser.ToTypeScriptName(constant.Name)
		value := parser.FormatValue(constant.Value, constant.Type, "typescript")
		comment := constant.Label
		if comment == "" {
			comment = constant.Label
		}
		
		code.WriteString(fmt.Sprintf("  /** %s */\n", comment))
		code.WriteString(fmt.Sprintf("  public static readonly %s = %s;\n", constName, value))
	}
	
	// 私有构造函数
	code.WriteString("\n  // 私有构造函数，防止实例化\n")
	code.WriteString("  private constructor() {\n")
	code.WriteString("    throw new Error('常量类不应被实例化');\n")
	code.WriteString("  }\n")
	
	// 生成方法
	code.WriteString("\n")
	code.WriteString(g.generateGetAllValues(group))
	code.WriteString("\n")
	code.WriteString(g.generateGetAllKeys(group))
	code.WriteString("\n")
	code.WriteString(g.generateGetKeyValuePairs(group))
	code.WriteString("\n")
	code.WriteString(g.generateFormatValue(group))
	code.WriteString("\n")
	code.WriteString(g.generateIsValid(group))
	code.WriteString("\n")
	code.WriteString(g.generateFromString(group))
	
	
	code.WriteString("}\n")
	
	return code.String()
}

// generateGetAllValues 生成获取所有值的方法
func (g *TypeScriptGenerator) generateGetAllValues(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	className := parser.ToJavaName(group.Name)
	tsType := parser.GetTypeScriptType(group.Constants[0].Type)
	
	code.WriteString("  /**\n")
	code.WriteString(fmt.Sprintf("   * 获取所有%s常量值\n", group.Label))
	code.WriteString("   * @returns 所有常量值的数组\n")
	code.WriteString("   */\n")
	code.WriteString(fmt.Sprintf("  public static getAllValues(): readonly %s[] {\n", tsType))
	code.WriteString("    return [\n")
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToTypeScriptName(constants[i].Name) < parser.ToTypeScriptName(constants[j].Name)
	})
	
	for _, constant := range constants {
		constName := parser.ToTypeScriptName(constant.Name)
		code.WriteString(fmt.Sprintf("      %s.%s,\n", className, constName))
	}
	
	code.WriteString("    ] as const;\n")
	code.WriteString("  }\n")
	
	return code.String()
}

// generateGetAllKeys 生成获取所有键的方法
func (g *TypeScriptGenerator) generateGetAllKeys(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	code.WriteString("  /**\n")
	code.WriteString(fmt.Sprintf("   * 获取所有%s常量键名\n", group.Label))
	code.WriteString("   * @returns 所有常量键名的数组\n")
	code.WriteString("   */\n")
	code.WriteString("  public static getAllKeys(): readonly string[] {\n")
	code.WriteString("    return [\n")
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToTypeScriptName(constants[i].Name) < parser.ToTypeScriptName(constants[j].Name)
	})
	
	for _, constant := range constants {
		constName := parser.ToTypeScriptName(constant.Name)
		code.WriteString(fmt.Sprintf("      '%s',\n", constName))
	}
	
	code.WriteString("    ] as const;\n")
	code.WriteString("  }\n")
	
	return code.String()
}

// generateGetKeyValuePairs 生成获取键值对的方法
func (g *TypeScriptGenerator) generateGetKeyValuePairs(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	className := parser.ToJavaName(group.Name)
	tsType := parser.GetTypeScriptType(group.Constants[0].Type)
	
	code.WriteString("  /**\n")
	code.WriteString("   * 获取键值对映射\n")
	code.WriteString("   * @returns 键值对映射对象\n")
	code.WriteString("   */\n")
	code.WriteString(fmt.Sprintf("  public static getKeyValuePairs(): Readonly<Record<string, %s>> {\n", tsType))
	code.WriteString("    return {\n")
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToTypeScriptName(constants[i].Name) < parser.ToTypeScriptName(constants[j].Name)
	})
	
	for _, constant := range constants {
		constName := parser.ToTypeScriptName(constant.Name)
		code.WriteString(fmt.Sprintf("      %s: %s.%s,\n", constName, className, constName))
	}
	
	code.WriteString("    } as const;\n")
	code.WriteString("  }\n")
	
	return code.String()
}

// generateFormatValue 生成格式化值的方法
func (g *TypeScriptGenerator) generateFormatValue(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	className := parser.ToJavaName(group.Name)
	tsType := parser.GetTypeScriptType(group.Constants[0].Type)
	
	code.WriteString("  /**\n")
	code.WriteString(fmt.Sprintf("   * 根据值格式化%s的标签\n", group.Label))
	code.WriteString("   * @param value 常量值\n")
	code.WriteString("   * @returns 格式化后的标签，找不到时返回 'Unknown(value)'\n")
	code.WriteString("   */\n")
	code.WriteString(fmt.Sprintf("  public static formatValue(value: %s): string {\n", tsType))
	code.WriteString(fmt.Sprintf("    const labels: Record<%s, string> = {\n", tsType))
	for _, constant := range group.Constants {
		constName := parser.ToTypeScriptName(constant.Name)
		label := constant.Label
		if label == "" {
			label = constant.Name
		}
		code.WriteString(fmt.Sprintf("      [%s.%s]: '%s',\n", className, constName, label))
	}
	code.WriteString("    };\n\n")
	code.WriteString("    if (labels[value] !== undefined) {\n")
	code.WriteString("      return labels[value];\n")
	code.WriteString("    }\n\n")
	code.WriteString("    return `Unknown(${value})`;\n")
	code.WriteString("  }\n")
	
	return code.String()
}

// generateIsValid 生成验证方法
func (g *TypeScriptGenerator) generateIsValid(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	className := parser.ToJavaName(group.Name)
	tsType := parser.GetTypeScriptType(group.Constants[0].Type)
	
	// 生成类型守卫
	typeGuard := fmt.Sprintf("value is typeof %s.%s", className, parser.ToTypeScriptName(group.Constants[0].Name))
	for i := 1; i < len(group.Constants); i++ {
		typeGuard += fmt.Sprintf(" | typeof %s.%s", className, parser.ToTypeScriptName(group.Constants[i].Name))
	}
	
	code.WriteString("  /**\n")
	code.WriteString(fmt.Sprintf("   * 验证值是否为有效的%s常量\n", group.Label))
	code.WriteString("   * @param value 要验证的值\n")
	code.WriteString("   * @returns 是否为有效常量\n")
	code.WriteString("   */\n")
	code.WriteString(fmt.Sprintf("  public static isValid(value: %s): %s {\n", tsType, typeGuard))
	code.WriteString(fmt.Sprintf("    return %s.getAllValues().includes(value);\n", className))
	code.WriteString("  }\n")
	
	return code.String()
}

// generateFromString 生成从字符串获取值的方法
func (g *TypeScriptGenerator) generateFromString(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	className := parser.ToJavaName(group.Name)
	tsType := parser.GetTypeScriptType(group.Constants[0].Type)
	
	code.WriteString("  /**\n")
	code.WriteString(fmt.Sprintf("   * 从字符串键名获取%s常量值\n", group.Label))
	code.WriteString("   * @param key 常量键名\n")
	code.WriteString("   * @returns 常量值，找不到时返回 undefined\n")
	code.WriteString("   */\n")
	code.WriteString(fmt.Sprintf("  public static fromString(key: string): %s | undefined {\n", tsType))
	code.WriteString(fmt.Sprintf("    const mapping = %s.getKeyValuePairs();\n", className))
	code.WriteString("    return mapping[key];\n")
	code.WriteString("  }\n")
	
	return code.String()
}


// generateTypeDefinitions 生成类型定义
func (g *TypeScriptGenerator) generateTypeDefinitions(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	className := parser.ToJavaName(group.Name)
	
	// 生成值类型
	code.WriteString(fmt.Sprintf("// 类型定义 - 使用后缀避免与其他常量类型重名\n"))
	code.WriteString(fmt.Sprintf("export type %sValue = ", className))
	
	for i, constant := range group.Constants {
		if i > 0 {
			code.WriteString(" | ")
		}
		code.WriteString(fmt.Sprintf("typeof %s.%s", className, parser.ToTypeScriptName(constant.Name)))
	}
	code.WriteString(";\n\n")
	
	// 生成键类型
	code.WriteString(fmt.Sprintf("export type %sKey = ", className))
	
	for i, constant := range group.Constants {
		if i > 0 {
			code.WriteString(" | ")
		}
		code.WriteString(fmt.Sprintf("'%s'", parser.ToTypeScriptName(constant.Name)))
	}
	code.WriteString(";\n")
	
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