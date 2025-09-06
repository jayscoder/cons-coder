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

// Swift 保留关键字列表
var swiftKeywords = map[string]bool{
	"associatedtype": true, "class": true, "deinit": true, "enum": true,
	"extension": true, "fileprivate": true, "func": true, "import": true,
	"init": true, "inout": true, "internal": true, "let": true,
	"open": true, "operator": true, "private": true, "protocol": true,
	"public": true, "static": true, "struct": true, "subscript": true,
	"typealias": true, "var": true, "break": true, "case": true,
	"continue": true, "default": true, "defer": true, "do": true,
	"else": true, "fallthrough": true, "for": true, "guard": true,
	"if": true, "in": true, "repeat": true, "return": true,
	"switch": true, "where": true, "while": true, "as": true,
	"catch": true, "false": true, "is": true, "nil": true,
	"rethrows": true, "super": true, "self": true, "Self": true,
	"throw": true, "throws": true, "true": true, "try": true,
	"type": true,
}

// escapeSwiftKeyword 如果是 Swift 关键字，使用反引号转义
func escapeSwiftKeyword(name string) string {
	if swiftKeywords[strings.ToLower(name)] {
		return "`" + name + "`"
	}
	return name
}

// Generate 生成Swift代码
func (g *SwiftGenerator) Generate(constants *parser.ConstantsFile) error {
	var code strings.Builder

	// 导入
	code.WriteString("import Foundation\n\n")

	// 文件头注释
	code.WriteString(g.GetFileHeader(constants))
	code.WriteString("\n")

	if g.Config.Mode == "const" {
		// const模式 - 生成简单常量
		for i, group := range constants.Groups {
			if i > 0 {
				code.WriteString("\n")
			}
			code.WriteString(g.generateConstGroup(group, constants.Label))
		}
	} else {
		// class模式
		// 生成每个常量组
		for i, group := range constants.Groups {
			if i > 0 {
				code.WriteString("\n")
			}
			code.WriteString(g.generateGroup(group, constants.Label))
		}
	}

	// 写入文件
	outputPath := g.GetOutputFilePath(constants.FileName)
	return os.WriteFile(outputPath, []byte(code.String()), 0644)
}

// generateConstGroup 生成const模式的常量组
func (g *SwiftGenerator) generateConstGroup(group *parser.ConstantGroup, projectLabel string) string {
	var code strings.Builder
	
	// 生成注释
	code.WriteString(fmt.Sprintf("// %s %s - %s\n", group.Name, group.Label, projectLabel))
	
	// 按字母顺序排序常量
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToSwiftName(constants[i].Name) < parser.ToSwiftName(constants[j].Name)
	})
	
	// 生成常量定义
	for _, constant := range constants {
		constName := fmt.Sprintf("%s_%s", strings.ToUpper(group.Name), strings.ToUpper(constant.Name))
		valueType := parser.GetSwiftType(constant.Type)
		value := parser.FormatValue(constant.Value, constant.Type, "swift")
		comment := constant.Label
		code.WriteString(fmt.Sprintf("public let %s: %s = %s // %s\n", constName, valueType, value, comment))
	}
	
	return code.String()
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
		// 处理 Swift 关键字
		caseName = escapeSwiftKeyword(caseName)
		value := constant.Value
		comment := constant.Label
		if comment == "" {
			comment = constant.Label
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
	code.WriteString("        label\n")
	code.WriteString("    }\n")

	// 生成label属性
	code.WriteString("\n")
	code.WriteString("    /// 获取中文标签\n")
	code.WriteString("    public var label: String {\n")
	code.WriteString("        switch self {\n")
	for _, constant := range constants {
		caseName := parser.ToSwiftName(constant.Name)
		caseName = escapeSwiftKeyword(caseName)
		label := constant.Label
		if label == "" {
			label = constant.Name
		}
		code.WriteString(fmt.Sprintf("        case .%s: return \"%s\"\n", caseName, label))
	}
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
		caseName = escapeSwiftKeyword(caseName)
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
			comment = constant.Label
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

	code.WriteString("    /// 根据值格式化标签\n")
	code.WriteString("    /// - Parameter value: 常量值\n")
	code.WriteString("    /// - Returns: 格式化后的标签\n")
	code.WriteString(fmt.Sprintf(`    public static func format(_ value: %s) -> String {`, swiftType))
	code.WriteString("\n")
	code.WriteString(fmt.Sprintf("        let labels: [%s: String] = [\n", swiftType))
	for _, constant := range group.Constants {
		constName := parser.ToSwiftName(constant.Name)
		label := constant.Label
		if label == "" {
			label = constant.Name
		}
		code.WriteString(fmt.Sprintf(`            %s: "%s",`, constName, label))
		code.WriteString("\n")
	}
	code.WriteString("        ]\n")
	code.WriteString("        \n")
	code.WriteString("        if let label = labels[value] {\n")
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


// GenerateIndex Swift不需要生成索引文件
func (g *SwiftGenerator) GenerateIndex(allConstants []*parser.ConstantsFile) error {
	return nil
}
