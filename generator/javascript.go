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

// JavaScriptGenerator JavaScript代码生成器
type JavaScriptGenerator struct {
	BaseGenerator
}

// NewJavaScriptGenerator 创建JavaScript生成器
func NewJavaScriptGenerator(config Config) *JavaScriptGenerator {
	return &JavaScriptGenerator{
		BaseGenerator: BaseGenerator{Config: config},
	}
}

// Generate 生成JavaScript代码
func (g *JavaScriptGenerator) Generate(constants *parser.ConstantsFile) error {
	var code strings.Builder
	
	// 文件头注释
	code.WriteString(g.GetFileHeader(constants))
	code.WriteString("\n")
	
	// 生成每个常量组的类
	for _, group := range constants.Groups {
		code.WriteString(g.generateGroupClass(group, constants.Label))
		code.WriteString("\n")
	}
	
	// 导出
	code.WriteString("// 导出所有常量类\n")
	code.WriteString("module.exports = {\n")
	for i, group := range constants.Groups {
		if i > 0 {
			code.WriteString(",\n")
		}
		className := parser.ToJavaName(group.Name)
		code.WriteString(fmt.Sprintf("  %s", className))
	}
	code.WriteString("\n};\n")
	
	// 写入文件
	outputPath := g.GetOutputFilePath(constants.FileName)
	return os.WriteFile(outputPath, []byte(code.String()), 0644)
}

// generateGroupClass 生成常量组类
func (g *JavaScriptGenerator) generateGroupClass(group *parser.ConstantGroup, projectLabel string) string {
	var code strings.Builder
	
	className := parser.ToJavaName(group.Name)
	
	// 类注释
	code.WriteString("/**\n")
	code.WriteString(fmt.Sprintf(" * %s - %s\n", group.Label, projectLabel))
	code.WriteString(" */\n")
	code.WriteString(fmt.Sprintf("class %s {\n", className))
	
	// 静态常量定义
	code.WriteString("  static {\n")
	
	// 按字母顺序排序常量
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToJavaScriptName(constants[i].Name) < parser.ToJavaScriptName(constants[j].Name)
	})
	
	// 常量定义
	for _, constant := range constants {
		constName := parser.ToJavaScriptName(constant.Name)
		value := parser.FormatValue(constant.Value, constant.Type, "javascript")
		comment := constant.Label
		if comment == "" {
			comment = constant.Desc
		}
		
		code.WriteString(fmt.Sprintf("    /** %s */\n", comment))
		code.WriteString(fmt.Sprintf("    this.%s = %s;\n", constName, value))
	}
	
	code.WriteString("  }\n\n")
	
	// 私有构造函数
	code.WriteString("  // 私有构造函数，防止实例化\n")
	code.WriteString("  constructor() {\n")
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
	code.WriteString("\n")
	code.WriteString(g.generateGetDescription(group))
	
	code.WriteString("}\n")
	
	return code.String()
}

// generateGetAllValues 生成获取所有值的方法
func (g *JavaScriptGenerator) generateGetAllValues(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	code.WriteString("  /**\n")
	code.WriteString(fmt.Sprintf("   * 获取所有%s常量值\n", group.Label))
	code.WriteString("   * @returns {Array} 所有常量值的数组\n")
	code.WriteString("   */\n")
	code.WriteString("  static getAllValues() {\n")
	code.WriteString("    return [\n")
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToJavaScriptName(constants[i].Name) < parser.ToJavaScriptName(constants[j].Name)
	})
	
	for _, constant := range constants {
		constName := parser.ToJavaScriptName(constant.Name)
		code.WriteString(fmt.Sprintf("      this.%s,\n", constName))
	}
	
	code.WriteString("    ];\n")
	code.WriteString("  }\n")
	
	return code.String()
}

// generateGetAllKeys 生成获取所有键的方法
func (g *JavaScriptGenerator) generateGetAllKeys(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	code.WriteString("  /**\n")
	code.WriteString(fmt.Sprintf("   * 获取所有%s常量键名\n", group.Label))
	code.WriteString("   * @returns {Array<string>} 所有常量键名的数组\n")
	code.WriteString("   */\n")
	code.WriteString("  static getAllKeys() {\n")
	code.WriteString("    return [\n")
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToJavaScriptName(constants[i].Name) < parser.ToJavaScriptName(constants[j].Name)
	})
	
	for _, constant := range constants {
		constName := parser.ToJavaScriptName(constant.Name)
		code.WriteString(fmt.Sprintf("      '%s',\n", constName))
	}
	
	code.WriteString("    ];\n")
	code.WriteString("  }\n")
	
	return code.String()
}

// generateGetKeyValuePairs 生成获取键值对的方法
func (g *JavaScriptGenerator) generateGetKeyValuePairs(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	code.WriteString("  /**\n")
	code.WriteString("   * 获取键值对映射\n")
	code.WriteString("   * @returns {Object} 键值对映射对象\n")
	code.WriteString("   */\n")
	code.WriteString("  static getKeyValuePairs() {\n")
	code.WriteString("    return {\n")
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToJavaScriptName(constants[i].Name) < parser.ToJavaScriptName(constants[j].Name)
	})
	
	for _, constant := range constants {
		constName := parser.ToJavaScriptName(constant.Name)
		code.WriteString(fmt.Sprintf("      %s: this.%s,\n", constName, constName))
	}
	
	code.WriteString("    };\n")
	code.WriteString("  }\n")
	
	return code.String()
}

// generateFormatValue 生成格式化值的方法
func (g *JavaScriptGenerator) generateFormatValue(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	jsType := parser.GetJavaScriptType(group.Constants[0].Type)
	
	code.WriteString("  /**\n")
	code.WriteString(fmt.Sprintf("   * 根据值和语言格式化%s的标签\n", group.Label))
	code.WriteString(fmt.Sprintf("   * @param {%s} value - 常量值\n", jsType))
	code.WriteString("   * @param {string} [lang='zh'] - 语言代码 ('zh', 'en', 'ja')\n")
	code.WriteString("   * @returns {string} 格式化后的标签，找不到时返回 'Unknown(value)'\n")
	code.WriteString("   */\n")
	code.WriteString("  static formatValue(value, lang = 'zh') {\n")
	code.WriteString("    const labels = {\n")
	
	// 中文标签
	code.WriteString("      zh: {\n")
	for _, constant := range group.Constants {
		constName := parser.ToJavaScriptName(constant.Name)
		label := constant.Label
		if label == "" {
			label = constant.Name
		}
		code.WriteString(fmt.Sprintf("        [this.%s]: '%s',\n", constName, label))
	}
	code.WriteString("      },\n")
	
	// 英文标签
	code.WriteString("      en: {\n")
	for _, constant := range group.Constants {
		constName := parser.ToJavaScriptName(constant.Name)
		label := strings.ReplaceAll(constant.Name, "_", " ")
		label = strings.Title(strings.ToLower(label))
		code.WriteString(fmt.Sprintf("        [this.%s]: '%s',\n", constName, label))
	}
	code.WriteString("      },\n")
	
	// 日文标签（示例）
	code.WriteString("      ja: {\n")
	for _, constant := range group.Constants {
		constName := parser.ToJavaScriptName(constant.Name)
		label := constant.Label // 暂时使用中文标签
		if label == "" {
			label = constant.Name
		}
		code.WriteString(fmt.Sprintf("        [this.%s]: '%s',\n", constName, label))
	}
	code.WriteString("      },\n")
	
	code.WriteString("    };\n\n")
	code.WriteString("    const langLabels = labels[lang];\n")
	code.WriteString("    if (langLabels && langLabels[value] !== undefined) {\n")
	code.WriteString("      return langLabels[value];\n")
	code.WriteString("    }\n\n")
	code.WriteString("    // 默认返回英文\n")
	code.WriteString("    const enLabels = labels.en;\n")
	code.WriteString("    if (enLabels && enLabels[value] !== undefined) {\n")
	code.WriteString("      return enLabels[value];\n")
	code.WriteString("    }\n\n")
	code.WriteString("    return `Unknown(${value})`;\n")
	code.WriteString("  }\n")
	
	return code.String()
}

// generateIsValid 生成验证方法
func (g *JavaScriptGenerator) generateIsValid(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	jsType := parser.GetJavaScriptType(group.Constants[0].Type)
	
	code.WriteString("  /**\n")
	code.WriteString(fmt.Sprintf("   * 验证值是否为有效的%s常量\n", group.Label))
	code.WriteString(fmt.Sprintf("   * @param {%s} value - 要验证的值\n", jsType))
	code.WriteString("   * @returns {boolean} 是否为有效常量\n")
	code.WriteString("   */\n")
	code.WriteString("  static isValid(value) {\n")
	code.WriteString("    return this.getAllValues().includes(value);\n")
	code.WriteString("  }\n")
	
	return code.String()
}

// generateFromString 生成从字符串获取值的方法
func (g *JavaScriptGenerator) generateFromString(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	jsType := parser.GetJavaScriptType(group.Constants[0].Type)
	
	code.WriteString("  /**\n")
	code.WriteString(fmt.Sprintf("   * 从字符串键名获取%s常量值\n", group.Label))
	code.WriteString("   * @param {string} key - 常量键名\n")
	code.WriteString(fmt.Sprintf("   * @returns {%s|undefined} 常量值，找不到时返回 undefined\n", jsType))
	code.WriteString("   */\n")
	code.WriteString("  static fromString(key) {\n")
	code.WriteString("    const mapping = this.getKeyValuePairs();\n")
	code.WriteString("    return mapping[key];\n")
	code.WriteString("  }\n")
	
	return code.String()
}

// generateGetDescription 生成获取描述的方法
func (g *JavaScriptGenerator) generateGetDescription(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	jsType := parser.GetJavaScriptType(group.Constants[0].Type)
	
	code.WriteString("  /**\n")
	code.WriteString("   * 获取常量值的详细描述\n")
	code.WriteString(fmt.Sprintf("   * @param {%s} value - 常量值\n", jsType))
	code.WriteString("   * @returns {string} 详细描述\n")
	code.WriteString("   */\n")
	code.WriteString("  static getDescription(value) {\n")
	code.WriteString("    const descriptions = {\n")
	
	for _, constant := range group.Constants {
		constName := parser.ToJavaScriptName(constant.Name)
		desc := constant.Desc
		if desc == "" {
			desc = constant.Label
		}
		code.WriteString(fmt.Sprintf("      [this.%s]: '%s',\n", constName, desc))
	}
	
	code.WriteString("    };\n\n")
	code.WriteString("    return descriptions[value] || `未知常量值: ${value}`;\n")
	code.WriteString("  }\n")
	
	return code.String()
}

// GenerateIndex 生成JavaScript的index.js文件
func (g *JavaScriptGenerator) GenerateIndex(allConstants []*parser.ConstantsFile) error {
	var code strings.Builder
	
	// 文件头注释
	code.WriteString("/**\n")
	code.WriteString(" * 常量包索引文件\n")
	code.WriteString(" * \n")
	code.WriteString(fmt.Sprintf(" * 生成时间: %s\n", FormatGenerationTime(time.Now())))
	code.WriteString(fmt.Sprintf(" * 生成工具: cons-coder v%s\n", g.Config.Version))
	code.WriteString(" */\n\n")
	
	// 导入所有文件
	for _, constants := range allConstants {
		code.WriteString(fmt.Sprintf("const %s = require('./%s');\n", 
			constants.FileName, constants.FileName))
	}
	
	code.WriteString("\n// 导出所有常量\n")
	code.WriteString("module.exports = {\n")
	
	// 导出所有常量类
	for _, constants := range allConstants {
		for _, group := range constants.Groups {
			className := parser.ToJavaName(group.Name)
			code.WriteString(fmt.Sprintf("  %s: %s.%s,\n", 
				className, constants.FileName, className))
		}
	}
	
	code.WriteString("};\n")
	
	// 写入文件
	outputPath := filepath.Join(g.Config.OutputDir, "index.js")
	return os.WriteFile(outputPath, []byte(code.String()), 0644)
}