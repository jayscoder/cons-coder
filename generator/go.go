package generator

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"cons-coder/parser"
)

// GoGenerator Go代码生成器
type GoGenerator struct {
	BaseGenerator
}

// NewGoGenerator 创建Go生成器
func NewGoGenerator(config Config) *GoGenerator {
	return &GoGenerator{
		BaseGenerator: BaseGenerator{Config: config},
	}
}

// Generate 生成Go代码
func (g *GoGenerator) Generate(constants *parser.ConstantsFile) error {
	var code strings.Builder
	
	// 文件头注释 - 必须在package声明之前
	code.WriteString(g.GetFileHeader(constants))
	code.WriteString("\n")
	
	// 包声明
	code.WriteString(fmt.Sprintf("package %s\n\n", g.Config.PackageName))
	
	if g.Config.Mode == "const" {
		// const模式不需要导入fmt
		// 生成每个常量组
		for _, group := range constants.Groups {
			code.WriteString(g.generateConstGroup(group, constants.Label))
			code.WriteString("\n")
		}
	} else {
		// class模式
		// 导入
		code.WriteString("import (\n")
		code.WriteString("\t\"fmt\"\n")
		code.WriteString(")\n\n")
		
		// 生成每个常量组
		for _, group := range constants.Groups {
			code.WriteString(g.generateGroup(group, constants.Label))
			code.WriteString("\n")
		}
	}
	
	// 写入文件
	outputPath := g.GetOutputFilePath(constants.FileName)
	return os.WriteFile(outputPath, []byte(code.String()), 0644)
}

// generateConstGroup 生成const模式的常量组
func (g *GoGenerator) generateConstGroup(group *parser.ConstantGroup, _ string) string {
	var code strings.Builder
	
	// 生成常量组
	code.WriteString("const (\n")
	
	// 按字母顺序排序常量
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToGoName(constants[i].Name) < parser.ToGoName(constants[j].Name)
	})
	
	// 生成常量定义
	for _, constant := range constants {
		constName := fmt.Sprintf("%s_%s", strings.ToUpper(group.Name), strings.ToUpper(constant.Name))
		value := parser.FormatValue(constant.Value, constant.Type, "go")
		comment := constant.Label
		code.WriteString(fmt.Sprintf("\t%s = %s // %s\n", constName, value, comment))
	}
	code.WriteString(")\n")
	
	return code.String()
}

// generateGroup 生成常量组
func (g *GoGenerator) generateGroup(group *parser.ConstantGroup, _ string) string {
	var code strings.Builder
	
	structName := toCamelCase(group.Name) + "Cons"
	
	// 生成结构体类型定义
	code.WriteString(fmt.Sprintf("type %s struct {\n", structName))
	
	// 按字母顺序排序常量
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToGoName(constants[i].Name) < parser.ToGoName(constants[j].Name)
	})
	
	// 生成结构体字段
	for _, constant := range constants {
		fieldName := parser.ToGoName(constant.Name)
		fieldType := parser.GetGoType(constant.Type)
		comment := constant.Label
		code.WriteString(fmt.Sprintf("\t%s %s // %s\n", fieldName, fieldType, comment))
	}
	code.WriteString("}\n\n")
	
	// 生成常量实例
	groupName := parser.ToGoName(group.Name)
	code.WriteString(fmt.Sprintf("// %s 常量实例\n", groupName))
	code.WriteString(fmt.Sprintf("var %s = %s{\n", groupName, structName))
	for _, constant := range constants {
		fieldName := parser.ToGoName(constant.Name)
		value := parser.FormatValue(constant.Value, constant.Type, "go")
		code.WriteString(fmt.Sprintf("\t%s: %s,\n", fieldName, value))
	}
	code.WriteString("}\n\n")
	
	// 生成方法
	code.WriteString(g.generateAllValues(group, structName))
	code.WriteString("\n")
	code.WriteString(g.generateAllKeys(group, structName))
	code.WriteString("\n")
	code.WriteString(g.generateKeyValuePairs(group, structName))
	code.WriteString("\n")
	code.WriteString(g.generateFormat(group, structName))
	code.WriteString("\n")
	code.WriteString(g.generateIsValid(group, structName))
	code.WriteString("\n")
	code.WriteString(g.generateFromString(group, structName))
	
	return code.String()
}

// toCamelCase 转换为小驼峰命名
func toCamelCase(name string) string {
	parts := strings.Split(name, "_")
	for i, part := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(part)
		} else {
			parts[i] = strings.Title(strings.ToLower(part))
		}
	}
	return strings.Join(parts, "")
}

// generateAllValues 生成获取所有值的方法
func (g *GoGenerator) generateAllValues(group *parser.ConstantGroup, structName string) string {
	var code strings.Builder
	
	groupName := parser.ToGoName(group.Name)
	valueType := parser.GetGoType(group.Constants[0].Type)
	
	code.WriteString(fmt.Sprintf("// AllValues 返回所有%s的值\n", groupName))
	code.WriteString(fmt.Sprintf("func (s %s) AllValues() []%s {\n", structName, valueType))
	code.WriteString(fmt.Sprintf("\treturn []%s{", valueType))
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToGoName(constants[i].Name) < parser.ToGoName(constants[j].Name)
	})
	
	for i, constant := range constants {
		if i > 0 {
			code.WriteString(", ")
		}
		code.WriteString(fmt.Sprintf("s.%s", parser.ToGoName(constant.Name)))
	}
	
	code.WriteString("}\n}\n")
	
	return code.String()
}

// generateAllKeys 生成获取所有键的方法
func (g *GoGenerator) generateAllKeys(group *parser.ConstantGroup, structName string) string {
	var code strings.Builder
	
	groupName := parser.ToGoName(group.Name)
	
	code.WriteString(fmt.Sprintf("// AllKeys 返回所有%s的键名\n", groupName))
	code.WriteString(fmt.Sprintf("func (s %s) AllKeys() []string {\n", structName))
	code.WriteString("\treturn []string{")
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToGoName(constants[i].Name) < parser.ToGoName(constants[j].Name)
	})
	
	for i, constant := range constants {
		if i > 0 {
			code.WriteString(", ")
		}
		code.WriteString(fmt.Sprintf(`"%s"`, parser.ToGoName(constant.Name)))
	}
	
	code.WriteString("}\n}\n")
	
	return code.String()
}

// generateKeyValuePairs 生成获取键值对的方法
func (g *GoGenerator) generateKeyValuePairs(group *parser.ConstantGroup, structName string) string {
	var code strings.Builder
	
	valueType := parser.GetGoType(group.Constants[0].Type)
	
	code.WriteString("// KeyValuePairs 返回键值对映射\n")
	code.WriteString(fmt.Sprintf("func (s %s) KeyValuePairs() map[string]%s {\n", structName, valueType))
	code.WriteString(fmt.Sprintf("\treturn map[string]%s{\n", valueType))
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToGoName(constants[i].Name) < parser.ToGoName(constants[j].Name)
	})
	
	for _, constant := range constants {
		fieldName := parser.ToGoName(constant.Name)
		code.WriteString(fmt.Sprintf("\t\t\"%s\": s.%s,\n", fieldName, fieldName))
	}
	
	code.WriteString("\t}\n}\n")
	
	return code.String()
}

// generateFormat 生成格式化方法
func (g *GoGenerator) generateFormat(group *parser.ConstantGroup, structName string) string {
	var code strings.Builder
	
	groupName := parser.ToGoName(group.Name)
	valueType := parser.GetGoType(group.Constants[0].Type)
	
	code.WriteString(fmt.Sprintf("// Format 根据值格式化%s的标签\n", groupName))
	code.WriteString(fmt.Sprintf("func (s %s) Format(value %s) string {\n", structName, valueType))
	code.WriteString(fmt.Sprintf("\tlabels := map[%s]string{\n", valueType))
	for _, constant := range group.Constants {
		fieldName := parser.ToGoName(constant.Name)
		label := constant.Label
		if label == "" {
			label = constant.Name
		}
		code.WriteString(fmt.Sprintf("\t\ts.%s: \"%s\",\n", fieldName, label))
	}
	code.WriteString("\t}\n\n")
	code.WriteString("\tif label, exists := labels[value]; exists {\n")
	code.WriteString("\t\treturn label\n")
	code.WriteString("\t}\n\n")
	code.WriteString("\treturn fmt.Sprintf(\"Unknown(%v)\", value)\n")
	code.WriteString("}\n")
	
	return code.String()
}

// generateIsValid 生成验证方法
func (g *GoGenerator) generateIsValid(group *parser.ConstantGroup, structName string) string {
	var code strings.Builder
	
	groupName := parser.ToGoName(group.Name)
	valueType := parser.GetGoType(group.Constants[0].Type)
	
	code.WriteString(fmt.Sprintf("// IsValid 检查值是否为有效的%s常量\n", groupName))
	code.WriteString(fmt.Sprintf("func (s %s) IsValid(value %s) bool {\n", structName, valueType))
	code.WriteString("\tallValues := s.AllValues()\n")
	code.WriteString("\tfor _, v := range allValues {\n")
	code.WriteString("\t\tif v == value {\n")
	code.WriteString("\t\t\treturn true\n")
	code.WriteString("\t\t}\n")
	code.WriteString("\t}\n")
	code.WriteString("\treturn false\n")
	code.WriteString("}\n")
	
	return code.String()
}

// generateFromString 生成从字符串获取值的方法
func (g *GoGenerator) generateFromString(group *parser.ConstantGroup, structName string) string {
	var code strings.Builder
	
	groupName := parser.ToGoName(group.Name)
	valueType := parser.GetGoType(group.Constants[0].Type)
	
	code.WriteString(fmt.Sprintf("// FromString 从字符串键名获取%s常量值\n", groupName))
	code.WriteString(fmt.Sprintf("func (s %s) FromString(key string) (%s, bool) {\n", structName, valueType))
	code.WriteString("\tmapping := s.KeyValuePairs()\n")
	code.WriteString("\tvalue, exists := mapping[key]\n")
	code.WriteString("\treturn value, exists\n")
	code.WriteString("}\n")
	
	return code.String()
}


// GenerateIndex Go不需要生成索引文件
func (g *GoGenerator) GenerateIndex(allConstants []*parser.ConstantsFile) error {
	return nil
}