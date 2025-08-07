package generator

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"cons-coder/parser"
)

// KotlinGenerator Kotlin代码生成器
type KotlinGenerator struct {
	BaseGenerator
}

// NewKotlinGenerator 创建Kotlin生成器
func NewKotlinGenerator(config Config) *KotlinGenerator {
	return &KotlinGenerator{
		BaseGenerator: BaseGenerator{Config: config},
	}
}

// Generate 生成Kotlin代码
func (g *KotlinGenerator) Generate(constants *parser.ConstantsFile) error {
	var code strings.Builder
	
	// 包声明
	code.WriteString(fmt.Sprintf("package %s\n\n", g.Config.PackageName))
	
	// 文件头注释
	code.WriteString(g.GetFileHeader(constants))
	code.WriteString("\n")
	
	// 生成每个常量组
	for i, group := range constants.Groups {
		if i > 0 {
			code.WriteString("\n")
		}
		code.WriteString(g.generateObject(group, constants.Label))
	}
	
	// 写入文件
	outputPath := g.GetOutputFilePath(constants.FileName)
	return os.WriteFile(outputPath, []byte(code.String()), 0644)
}

// generateObject 生成常量组对象
func (g *KotlinGenerator) generateObject(group *parser.ConstantGroup, projectLabel string) string {
	var code strings.Builder
	
	objectName := parser.ToKotlinName(group.Name)
	
	// 对象注释
	code.WriteString("/**\n")
	code.WriteString(fmt.Sprintf(" * %s - %s\n", group.Label, projectLabel))
	code.WriteString(" */\n")
	code.WriteString(fmt.Sprintf("object %s {\n", objectName))
	
	// 按字母顺序排序常量
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToKotlinConstantName(constants[i].Name) < parser.ToKotlinConstantName(constants[j].Name)
	})
	
	// 常量定义
	for _, constant := range constants {
		constName := parser.ToKotlinConstantName(constant.Name)
		kotlinType := parser.GetKotlinType(constant.Type)
		value := parser.FormatValue(constant.Value, constant.Type, "kotlin")
		comment := constant.Label
		if comment == "" {
			comment = constant.Desc
		}
		
		code.WriteString(fmt.Sprintf("    /** %s */\n", comment))
		code.WriteString(fmt.Sprintf("    const val %s: %s = %s\n", constName, kotlinType, value))
	}
	
	// 生成方法
	code.WriteString("\n")
	code.WriteString(g.generateGetAllValues(group))
	code.WriteString("\n")
	code.WriteString(g.generateGetAllKeys(group))
	code.WriteString("\n")
	code.WriteString(g.generateGetKeyValuePairs(group))
	code.WriteString("\n")
	code.WriteString(g.generateFormat(group))
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
func (g *KotlinGenerator) generateGetAllValues(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	kotlinType := parser.GetKotlinType(group.Constants[0].Type)
	
	code.WriteString("    /** 获取所有常量值 */\n")
	code.WriteString(fmt.Sprintf("    fun getAllValues(): List<%s> {\n", kotlinType))
	code.WriteString("        return listOf(")
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToKotlinConstantName(constants[i].Name) < parser.ToKotlinConstantName(constants[j].Name)
	})
	
	for i, constant := range constants {
		if i > 0 {
			code.WriteString(", ")
		}
		code.WriteString(parser.ToKotlinConstantName(constant.Name))
	}
	
	code.WriteString(")\n")
	code.WriteString("    }\n")
	
	return code.String()
}

// generateGetAllKeys 生成获取所有键的方法
func (g *KotlinGenerator) generateGetAllKeys(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	code.WriteString("    /** 获取所有常量键名 */\n")
	code.WriteString("    fun getAllKeys(): List<String> {\n")
	code.WriteString("        return listOf(")
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToKotlinConstantName(constants[i].Name) < parser.ToKotlinConstantName(constants[j].Name)
	})
	
	for i, constant := range constants {
		if i > 0 {
			code.WriteString(", ")
		}
		code.WriteString(fmt.Sprintf(`"%s"`, parser.ToKotlinConstantName(constant.Name)))
	}
	
	code.WriteString(")\n")
	code.WriteString("    }\n")
	
	return code.String()
}

// generateGetKeyValuePairs 生成获取键值对的方法
func (g *KotlinGenerator) generateGetKeyValuePairs(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	kotlinType := parser.GetKotlinType(group.Constants[0].Type)
	
	code.WriteString("    /** 获取键值对映射 */\n")
	code.WriteString(fmt.Sprintf("    fun getKeyValuePairs(): Map<String, %s> {\n", kotlinType))
	code.WriteString("        return mapOf(\n")
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToKotlinConstantName(constants[i].Name) < parser.ToKotlinConstantName(constants[j].Name)
	})
	
	for _, constant := range constants {
		constName := parser.ToKotlinConstantName(constant.Name)
		code.WriteString(fmt.Sprintf(`            "%s" to %s,`, constName, constName))
		code.WriteString("\n")
	}
	
	code.WriteString("        )\n")
	code.WriteString("    }\n")
	
	return code.String()
}

// generateFormat 生成格式化方法
func (g *KotlinGenerator) generateFormat(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	kotlinType := parser.GetKotlinType(group.Constants[0].Type)
	
	code.WriteString("    /**\n")
	code.WriteString("     * 根据值和语言格式化标签\n")
	code.WriteString("     * @param value 常量值\n")
	code.WriteString(`     * @param lang 语言代码 (默认: "zh")`)
	code.WriteString("\n")
	code.WriteString("     * @return 格式化后的标签\n")
	code.WriteString("     */\n")
	code.WriteString(fmt.Sprintf(`    fun format(value: %s, lang: String = "zh"): String {`, kotlinType))
	code.WriteString("\n")
	code.WriteString("        val labels = mapOf(\n")
	
	// 中文标签
	code.WriteString(`            "zh" to mapOf(`)
	code.WriteString("\n")
	for i, constant := range group.Constants {
		if i > 0 {
			code.WriteString(",\n")
		}
		constName := parser.ToKotlinConstantName(constant.Name)
		label := constant.Label
		if label == "" {
			label = constant.Name
		}
		code.WriteString(fmt.Sprintf(`                %s to "%s"`, constName, label))
	}
	code.WriteString(",\n")
	code.WriteString("            ),\n")
	
	// 英文标签
	code.WriteString(`            "en" to mapOf(`)
	code.WriteString("\n")
	for i, constant := range group.Constants {
		if i > 0 {
			code.WriteString(",\n")
		}
		constName := parser.ToKotlinConstantName(constant.Name)
		label := strings.ReplaceAll(constant.Name, "_", " ")
		label = strings.Title(strings.ToLower(label))
		code.WriteString(fmt.Sprintf(`                %s to "%s"`, constName, label))
	}
	code.WriteString(",\n")
	code.WriteString("            ),\n")
	
	// 日文标签（示例）
	code.WriteString(`            "ja" to mapOf(`)
	code.WriteString("\n")
	for i, constant := range group.Constants {
		if i > 0 {
			code.WriteString(",\n")
		}
		constName := parser.ToKotlinConstantName(constant.Name)
		label := constant.Label // 暂时使用中文标签
		if label == "" {
			label = constant.Name
		}
		code.WriteString(fmt.Sprintf(`                %s to "%s"`, constName, label))
	}
	code.WriteString(",\n")
	code.WriteString("            ),\n")
	
	code.WriteString("        )\n\n")
	code.WriteString("        labels[lang]?.get(value)?.let { return it }\n\n")
	code.WriteString("        // 默认返回英文\n")
	code.WriteString(`        labels["en"]?.get(value)?.let { return it }`)
	code.WriteString("\n\n")
	code.WriteString(`        return "Unknown($value)"`)
	code.WriteString("\n")
	code.WriteString("    }\n")
	
	return code.String()
}

// generateIsValid 生成验证方法
func (g *KotlinGenerator) generateIsValid(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	kotlinType := parser.GetKotlinType(group.Constants[0].Type)
	
	code.WriteString("    /** \n")
	code.WriteString("     * 检查值是否为有效常量\n")
	code.WriteString("     * @param value 要检查的值\n")
	code.WriteString("     * @return 是否为有效常量\n")
	code.WriteString("     */\n")
	code.WriteString(fmt.Sprintf("    fun isValid(value: %s): Boolean {\n", kotlinType))
	code.WriteString("        return getAllValues().contains(value)\n")
	code.WriteString("    }\n")
	
	return code.String()
}

// generateFromString 生成从字符串获取值的方法
func (g *KotlinGenerator) generateFromString(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	kotlinType := parser.GetKotlinType(group.Constants[0].Type)
	
	code.WriteString("    /**\n")
	code.WriteString("     * 从字符串键名获取常量值\n")
	code.WriteString("     * @param key 常量键名\n")
	code.WriteString("     * @return 常量值，找不到时返回null\n")
	code.WriteString("     */\n")
	code.WriteString(fmt.Sprintf("    fun fromString(key: String): %s? {\n", kotlinType))
	code.WriteString("        return getKeyValuePairs()[key]\n")
	code.WriteString("    }\n")
	
	return code.String()
}

// generateGetDescription 生成获取描述的方法
func (g *KotlinGenerator) generateGetDescription(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	kotlinType := parser.GetKotlinType(group.Constants[0].Type)
	
	code.WriteString("    /**\n")
	code.WriteString("     * 获取常量值的详细描述\n")
	code.WriteString("     * @param value 常量值\n")
	code.WriteString("     * @return 详细描述\n")
	code.WriteString("     */\n")
	code.WriteString(fmt.Sprintf("    fun getDescription(value: %s): String {\n", kotlinType))
	code.WriteString("        val descriptions = mapOf(\n")
	
	for i, constant := range group.Constants {
		if i > 0 {
			code.WriteString(",\n")
		}
		constName := parser.ToKotlinConstantName(constant.Name)
		desc := constant.Desc
		if desc == "" {
			desc = constant.Label
		}
		code.WriteString(fmt.Sprintf(`            %s to "%s"`, constName, desc))
	}
	
	code.WriteString("\n        )\n")
	code.WriteString("        \n")
	code.WriteString(`        return descriptions[value] ?: "未知常量值: $value"`)
	code.WriteString("\n")
	code.WriteString("    }\n")
	
	return code.String()
}

// GenerateIndex Kotlin不需要生成索引文件
func (g *KotlinGenerator) GenerateIndex(allConstants []*parser.ConstantsFile) error {
	return nil
}