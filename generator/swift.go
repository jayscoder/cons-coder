package generator

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"cons-coder/parser"
)

// SwiftGenerator Swift代码生成器
type SwiftGenerator struct {
	BaseGenerator
}

// NewSwiftGenerator 创建Swift生成器
func NewSwiftGenerator(config Config) *SwiftGenerator {
	return &SwiftGenerator{
		BaseGenerator: BaseGenerator{Config: config},
	}
}

// Generate 生成Swift代码
func (g *SwiftGenerator) Generate(constants *parser.ConstantsFile) error {
	var code strings.Builder
	
	// 导入
	code.WriteString("import Foundation\n\n")
	
	// 文件头注释
	code.WriteString(g.GetFileHeader(constants))
	code.WriteString("\n")
	
	// 生成每个常量组
	for i, group := range constants.Groups {
		if i > 0 {
			code.WriteString("\n")
		}
		code.WriteString(g.generateGroup(group, constants.Label))
	}
	
	// 写入文件
	outputPath := g.GetOutputFilePath(constants.FileName)
	return os.WriteFile(outputPath, []byte(code.String()), 0644)
}

// generateGroup 生成常量组
func (g *SwiftGenerator) generateGroup(group *parser.ConstantGroup, projectLabel string) string {
	// 所有类型都使用enum形式
	if len(group.Constants) > 0 {
		return g.generateEnumGroup(group, projectLabel)
	}
	
	// 字符串类型仍使用struct形式（作为备份，但实际不会使用）
	return g.generateStructGroup(group, projectLabel)
}

// generateEnumGroup 生成enum形式的常量组（适用于整数和字符串类型）
func (g *SwiftGenerator) generateEnumGroup(group *parser.ConstantGroup, projectLabel string) string {
	var code strings.Builder
	
	enumName := parser.ToJavaName(group.Name)
	
	// 判断类型
	isIntType := len(group.Constants) > 0 && group.Constants[0].Type == "int"
	rawType := "String"
	if isIntType {
		rawType = "Int"
	}
	
	// 枚举注释
	code.WriteString(fmt.Sprintf("/// %s - %s\n", group.Label, projectLabel))
	code.WriteString(fmt.Sprintf("public enum %s: %s, CaseIterable, Codable, Identifiable, CustomStringConvertible {\n", enumName, rawType))
	
	// 按字母顺序排序常量
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return strings.ToUpper(constants[i].Name) < strings.ToUpper(constants[j].Name)
	})
	
	// 枚举case定义
	for _, constant := range constants {
		caseName := parser.ToSwiftName(constant.Name)
		value := constant.Value
		comment := constant.Label
		if comment == "" {
			comment = constant.Desc
		}
		
		code.WriteString(fmt.Sprintf("    /// %s\n", comment))
		if isIntType {
			code.WriteString(fmt.Sprintf("    case %s = %s\n", caseName, value))
		} else {
			// 字符串类型需要加引号
			code.WriteString(fmt.Sprintf("    case %s = \"%s\"\n", caseName, value))
		}
	}
	
	// 添加Identifiable协议的实现
	if isIntType {
		code.WriteString("\n    public var id: Int { rawValue }\n")
	} else {
		code.WriteString("\n    public var id: String { rawValue }\n")
	}
	
	// 添加CustomStringConvertible协议的实现
	code.WriteString("    \n")
	code.WriteString("    public var description: String {\n")
	code.WriteString("        format(lang: \"en\")\n")
	code.WriteString("    }\n")
	
	// 生成label属性
	code.WriteString("\n")
	code.WriteString("    /// 获取中文标签\n")
	code.WriteString("    public var label: String {\n")
	code.WriteString("        switch self {\n")
	for _, constant := range constants {
		caseName := parser.ToSwiftName(constant.Name)
		label := constant.Label
		if label == "" {
			label = constant.Name
		}
		code.WriteString(fmt.Sprintf("        case .%s: return \"%s\"\n", caseName, label))
	}
	code.WriteString("        }\n")
	code.WriteString("    }\n")
	
	// 生成详细描述属性
	code.WriteString("\n")
	code.WriteString("    /// 获取详细描述\n")
	code.WriteString("    public var detailDescription: String {\n")
	code.WriteString("        switch self {\n")
	for _, constant := range constants {
		caseName := parser.ToSwiftName(constant.Name)
		desc := constant.Desc
		if desc == "" {
			desc = constant.Label
		}
		code.WriteString(fmt.Sprintf("        case .%s: return \"%s\"\n", caseName, desc))
	}
	code.WriteString("        }\n")
	code.WriteString("    }\n")
	
	// 生成格式化方法
	code.WriteString("\n")
	code.WriteString("    /// 根据语言获取标签\n")
	code.WriteString("    /// - Parameter lang: 语言代码 (\"zh\", \"en\")\n")
	code.WriteString("    /// - Returns: 格式化后的标签\n")
	code.WriteString("    public func format(lang: String = \"en\") -> String {\n")
	code.WriteString("        switch lang {\n")
	code.WriteString("        case \"zh\":\n")
	code.WriteString("            return label\n")
	code.WriteString("        case \"en\":\n")
	code.WriteString("            switch self {\n")
	for _, constant := range constants {
		caseName := parser.ToSwiftName(constant.Name)
		// 简单的英文转换
		enLabel := strings.ReplaceAll(constant.Name, "_", " ")
		enLabel = strings.Title(strings.ToLower(enLabel))
		code.WriteString(fmt.Sprintf("            case .%s: return \"%s\"\n", caseName, enLabel))
	}
	code.WriteString("            }\n")
	code.WriteString("        default:\n")
	code.WriteString("            return label\n")
	code.WriteString("        }\n")
	code.WriteString("    }\n")
	
	// 生成从字符串创建的静态方法
	code.WriteString("\n")
	code.WriteString("    /// 从字符串键名创建枚举\n")
	code.WriteString("    /// - Parameter key: 常量键名（大写）\n")
	code.WriteString("    /// - Returns: 枚举值，找不到时返回nil\n")
	code.WriteString("    public static func fromString(_ key: String) -> Self? {\n")
	code.WriteString("        switch key {\n")
	for _, constant := range constants {
		caseName := parser.ToSwiftName(constant.Name)
		// 对于fromString，仍使用原始的键名（可能是snake_case等）
		originalKey := constant.Name
		code.WriteString(fmt.Sprintf("        case \"%s\": return .%s\n", originalKey, caseName))
	}
	code.WriteString("        default: return nil\n")
	code.WriteString("        }\n")
	code.WriteString("    }\n")
	
	code.WriteString("}\n")
	
	return code.String()
}

// generateStructGroup 生成struct形式的常量组（适用于字符串类型）
func (g *SwiftGenerator) generateStructGroup(group *parser.ConstantGroup, projectLabel string) string {
	var code strings.Builder
	
	structName := parser.ToJavaName(group.Name)
	
	// 结构体注释
	code.WriteString(fmt.Sprintf("/// %s - %s\n", group.Label, projectLabel))
	code.WriteString(fmt.Sprintf("public struct %s {\n", structName))
	
	// 按字母顺序排序常量
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToSwiftName(constants[i].Name) < parser.ToSwiftName(constants[j].Name)
	})
	
	// 常量定义
	for _, constant := range constants {
		constName := parser.ToSwiftName(constant.Name)
		swiftType := parser.GetSwiftType(constant.Type)
		value := parser.FormatValue(constant.Value, constant.Type, "swift")
		comment := constant.Label
		if comment == "" {
			comment = constant.Desc
		}
		
		code.WriteString(fmt.Sprintf("    /// %s\n", comment))
		code.WriteString(fmt.Sprintf("    public static let %s: %s = %s\n", constName, swiftType, value))
	}
	
	// 私有初始化器
	code.WriteString("    \n")
	code.WriteString("    // 私有初始化器，防止实例化\n")
	code.WriteString("    private init() {}\n")
	
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
func (g *SwiftGenerator) generateGetAllValues(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	swiftType := parser.GetSwiftType(group.Constants[0].Type)
	
	code.WriteString("    /// 获取所有常量值\n")
	code.WriteString(fmt.Sprintf("    public static func getAllValues() -> [%s] {\n", swiftType))
	code.WriteString("        return [")
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToSwiftName(constants[i].Name) < parser.ToSwiftName(constants[j].Name)
	})
	
	for i, constant := range constants {
		if i > 0 {
			code.WriteString(", ")
		}
		code.WriteString(parser.ToSwiftName(constant.Name))
	}
	
	code.WriteString("]\n")
	code.WriteString("    }\n")
	
	return code.String()
}

// generateGetAllKeys 生成获取所有键的方法
func (g *SwiftGenerator) generateGetAllKeys(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	code.WriteString("    /// 获取所有常量键名\n")
	code.WriteString("    public static func getAllKeys() -> [String] {\n")
	code.WriteString("        return [")
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToSwiftName(constants[i].Name) < parser.ToSwiftName(constants[j].Name)
	})
	
	for i, constant := range constants {
		if i > 0 {
			code.WriteString(", ")
		}
		code.WriteString(fmt.Sprintf(`"%s"`, parser.ToSwiftName(constant.Name)))
	}
	
	code.WriteString("]\n")
	code.WriteString("    }\n")
	
	return code.String()
}

// generateGetKeyValuePairs 生成获取键值对的方法
func (g *SwiftGenerator) generateGetKeyValuePairs(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	swiftType := parser.GetSwiftType(group.Constants[0].Type)
	
	code.WriteString("    /// 获取键值对字典\n")
	code.WriteString(fmt.Sprintf("    public static func getKeyValuePairs() -> [String: %s] {\n", swiftType))
	code.WriteString("        return [\n")
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToSwiftName(constants[i].Name) < parser.ToSwiftName(constants[j].Name)
	})
	
	for _, constant := range constants {
		constName := parser.ToSwiftName(constant.Name)
		code.WriteString(fmt.Sprintf(`            "%s": %s,`, constName, constName))
		code.WriteString("\n")
	}
	
	code.WriteString("        ]\n")
	code.WriteString("    }\n")
	
	return code.String()
}

// generateFormat 生成格式化方法
func (g *SwiftGenerator) generateFormat(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	swiftType := parser.GetSwiftType(group.Constants[0].Type)
	
	code.WriteString("    /// 根据值和语言格式化标签\n")
	code.WriteString("    /// - Parameters:\n")
	code.WriteString("    ///   - value: 常量值\n")
	code.WriteString(`    ///   - lang: 语言代码 (默认: "zh")`)
	code.WriteString("\n")
	code.WriteString("    /// - Returns: 格式化后的标签\n")
	code.WriteString(fmt.Sprintf(`    public static func format(_ value: %s, lang: String = "zh") -> String {`, swiftType))
	code.WriteString("\n")
	code.WriteString(fmt.Sprintf("        let labels: [String: [%s: String]] = [\n", swiftType))
	
	// 中文标签
	code.WriteString(`            "zh": [`)
	code.WriteString("\n")
	for _, constant := range group.Constants {
		constName := parser.ToSwiftName(constant.Name)
		label := constant.Label
		if label == "" {
			label = constant.Name
		}
		code.WriteString(fmt.Sprintf(`                %s: "%s",`, constName, label))
		code.WriteString("\n")
	}
	code.WriteString("            ],\n")
	
	// 英文标签
	code.WriteString(`            "en": [`)
	code.WriteString("\n")
	for _, constant := range group.Constants {
		constName := parser.ToSwiftName(constant.Name)
		label := strings.ReplaceAll(constant.Name, "_", " ")
		label = strings.Title(strings.ToLower(label))
		code.WriteString(fmt.Sprintf(`                %s: "%s",`, constName, label))
		code.WriteString("\n")
	}
	code.WriteString("            ],\n")
	
	// 日文标签（示例）
	code.WriteString(`            "ja": [`)
	code.WriteString("\n")
	for _, constant := range group.Constants {
		constName := parser.ToSwiftName(constant.Name)
		label := constant.Label // 暂时使用中文标签
		if label == "" {
			label = constant.Name
		}
		code.WriteString(fmt.Sprintf(`                %s: "%s",`, constName, label))
		code.WriteString("\n")
	}
	code.WriteString("            ]\n")
	
	code.WriteString("        ]\n")
	code.WriteString("        \n")
	code.WriteString("        if let langLabels = labels[lang], let label = langLabels[value] {\n")
	code.WriteString("            return label\n")
	code.WriteString("        }\n")
	code.WriteString("        \n")
	code.WriteString("        // 默认返回英文\n")
	code.WriteString(`        if let enLabels = labels["en"], let label = enLabels[value] {`)
	code.WriteString("\n")
	code.WriteString("            return label\n")
	code.WriteString("        }\n")
	code.WriteString("        \n")
	code.WriteString(`        return "Unknown(\(value))"`)
	code.WriteString("\n")
	code.WriteString("    }\n")
	
	return code.String()
}

// generateIsValid 生成验证方法
func (g *SwiftGenerator) generateIsValid(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	swiftType := parser.GetSwiftType(group.Constants[0].Type)
	
	code.WriteString("    /// 验证值是否有效\n")
	code.WriteString("    /// - Parameter value: 要验证的值\n")
	code.WriteString("    /// - Returns: 是否为有效常量\n")
	code.WriteString(fmt.Sprintf("    public static func isValid(_ value: %s) -> Bool {\n", swiftType))
	code.WriteString("        return getAllValues().contains(value)\n")
	code.WriteString("    }\n")
	
	return code.String()
}

// generateFromString 生成从字符串获取值的方法
func (g *SwiftGenerator) generateFromString(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	swiftType := parser.GetSwiftType(group.Constants[0].Type)
	
	code.WriteString("    /// 从字符串键名获取常量值\n")
	code.WriteString("    /// - Parameter key: 常量键名\n")
	code.WriteString("    /// - Returns: 常量值，找不到时返回nil\n")
	code.WriteString(fmt.Sprintf("    public static func fromString(_ key: String) -> %s? {\n", swiftType))
	code.WriteString("        return getKeyValuePairs()[key]\n")
	code.WriteString("    }\n")
	
	return code.String()
}

// generateGetDescription 生成获取描述的方法
func (g *SwiftGenerator) generateGetDescription(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	swiftType := parser.GetSwiftType(group.Constants[0].Type)
	
	code.WriteString("    /// 获取常量值的详细描述\n")
	code.WriteString("    /// - Parameter value: 常量值\n")
	code.WriteString("    /// - Returns: 详细描述\n")
	code.WriteString(fmt.Sprintf("    public static func getDescription(_ value: %s) -> String {\n", swiftType))
	code.WriteString(fmt.Sprintf("        let descriptions: [%s: String] = [\n", swiftType))
	
	for _, constant := range group.Constants {
		constName := parser.ToSwiftName(constant.Name)
		desc := constant.Desc
		if desc == "" {
			desc = constant.Label
		}
		code.WriteString(fmt.Sprintf(`            %s: "%s",`, constName, desc))
		code.WriteString("\n")
	}
	
	code.WriteString("        ]\n")
	code.WriteString("        \n")
	code.WriteString(`        return descriptions[value] ?? "未知常量值: \(value)"`)
	code.WriteString("\n")
	code.WriteString("    }\n")
	
	return code.String()
}

// GenerateIndex Swift不需要生成索引文件
func (g *SwiftGenerator) GenerateIndex(allConstants []*parser.ConstantsFile) error {
	return nil
}