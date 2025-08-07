package generator

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"cons-coder/parser"
)

// JavaGenerator Java代码生成器
type JavaGenerator struct {
	BaseGenerator
}

// NewJavaGenerator 创建Java生成器
func NewJavaGenerator(config Config) *JavaGenerator {
	return &JavaGenerator{
		BaseGenerator: BaseGenerator{Config: config},
	}
}

// Generate 生成Java代码
func (g *JavaGenerator) Generate(constants *parser.ConstantsFile) error {
	var code strings.Builder
	
	// 包声明
	code.WriteString(fmt.Sprintf("package %s;\n\n", g.Config.PackageName))
	
	// 导入
	code.WriteString("import java.util.*;\n\n")
	
	// 文件头注释
	code.WriteString(g.GetFileHeader(constants))
	code.WriteString("\n")
	
	// 文件类
	className := parser.ToJavaName(constants.FileName)
	code.WriteString(fmt.Sprintf("public final class %s {\n\n", className))
	
	// 私有构造函数
	code.WriteString(fmt.Sprintf("\t// 私有构造函数，防止实例化\n"))
	code.WriteString(fmt.Sprintf("\tprivate %s() {\n", className))
	code.WriteString(fmt.Sprintf("\t\tthrow new AssertionError(\"工具类不应被实例化\");\n"))
	code.WriteString(fmt.Sprintf("\t}\n\n"))
	
	// 生成每个常量组作为内部类
	for i, group := range constants.Groups {
		if i > 0 {
			code.WriteString("\n")
		}
		code.WriteString(g.generateGroupClass(group, constants.Label))
	}
	
	code.WriteString("}\n")
	
	// 写入文件
	outputPath := g.GetOutputFilePath(constants.FileName)
	return os.WriteFile(outputPath, []byte(code.String()), 0644)
}

// generateGroupClass 生成常量组内部类
func (g *JavaGenerator) generateGroupClass(group *parser.ConstantGroup, projectLabel string) string {
	var code strings.Builder
	
	className := parser.ToJavaName(group.Name)
	
	// 类注释
	code.WriteString(fmt.Sprintf("\t/**\n"))
	code.WriteString(fmt.Sprintf("\t * %s - %s\n", group.Label, projectLabel))
	code.WriteString(fmt.Sprintf("\t */\n"))
	code.WriteString(fmt.Sprintf("\tpublic static final class %s {\n", className))
	
	// 按字母顺序排序常量
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToJavaConstantName(constants[i].Name) < parser.ToJavaConstantName(constants[j].Name)
	})
	
	// 常量定义
	for _, constant := range constants {
		constName := parser.ToJavaConstantName(constant.Name)
		javaType := parser.GetJavaType(constant.Type)
		value := parser.FormatValue(constant.Value, constant.Type, "java")
		comment := constant.Label
		if comment == "" {
			comment = constant.Desc
		}
		
		code.WriteString(fmt.Sprintf("\t\t/** %s */\n", comment))
		code.WriteString(fmt.Sprintf("\t\tpublic static final %s %s = %s;\n", javaType, constName, value))
	}
	
	// 私有构造函数
	code.WriteString(fmt.Sprintf("\n\t\t// 私有构造函数，防止实例化\n"))
	code.WriteString(fmt.Sprintf("\t\tprivate %s() {\n", className))
	code.WriteString(fmt.Sprintf("\t\t\tthrow new AssertionError(\"常量类不应被实例化\");\n"))
	code.WriteString(fmt.Sprintf("\t\t}\n"))
	
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
	
	code.WriteString("\t}\n")
	
	return code.String()
}

// generateGetAllValues 生成获取所有值的方法
func (g *JavaGenerator) generateGetAllValues(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	javaType := parser.GetJavaType(group.Constants[0].Type)
	boxedType := getBoxedType(javaType)
	
	code.WriteString("\t\t/**\n")
	code.WriteString("\t\t * 获取所有常量值\n")
	code.WriteString("\t\t * @return 所有常量值的列表\n")
	code.WriteString("\t\t */\n")
	code.WriteString(fmt.Sprintf("\t\tpublic static List<%s> getAllValues() {\n", boxedType))
	code.WriteString("\t\t\treturn Arrays.asList(")
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToJavaConstantName(constants[i].Name) < parser.ToJavaConstantName(constants[j].Name)
	})
	
	for i, constant := range constants {
		if i > 0 {
			code.WriteString(", ")
		}
		code.WriteString(parser.ToJavaConstantName(constant.Name))
	}
	
	code.WriteString(");\n")
	code.WriteString("\t\t}\n")
	
	return code.String()
}

// generateGetAllKeys 生成获取所有键的方法
func (g *JavaGenerator) generateGetAllKeys(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	code.WriteString("\t\t/**\n")
	code.WriteString("\t\t * 获取所有常量键名\n")
	code.WriteString("\t\t * @return 所有常量键名的列表\n")
	code.WriteString("\t\t */\n")
	code.WriteString("\t\tpublic static List<String> getAllKeys() {\n")
	code.WriteString("\t\t\treturn Arrays.asList(")
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToJavaConstantName(constants[i].Name) < parser.ToJavaConstantName(constants[j].Name)
	})
	
	for i, constant := range constants {
		if i > 0 {
			code.WriteString(", ")
		}
		code.WriteString(fmt.Sprintf(`"%s"`, parser.ToJavaConstantName(constant.Name)))
	}
	
	code.WriteString(");\n")
	code.WriteString("\t\t}\n")
	
	return code.String()
}

// generateGetKeyValuePairs 生成获取键值对的方法
func (g *JavaGenerator) generateGetKeyValuePairs(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	javaType := parser.GetJavaType(group.Constants[0].Type)
	boxedType := getBoxedType(javaType)
	
	code.WriteString("\t\t/**\n")
	code.WriteString("\t\t * 获取键值对映射\n")
	code.WriteString("\t\t * @return 键值对映射\n")
	code.WriteString("\t\t */\n")
	code.WriteString(fmt.Sprintf("\t\tpublic static Map<String, %s> getKeyValuePairs() {\n", boxedType))
	code.WriteString(fmt.Sprintf("\t\t\tMap<String, %s> pairs = new HashMap<>();\n", boxedType))
	
	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToJavaConstantName(constants[i].Name) < parser.ToJavaConstantName(constants[j].Name)
	})
	
	for _, constant := range constants {
		constName := parser.ToJavaConstantName(constant.Name)
		code.WriteString(fmt.Sprintf("\t\t\tpairs.put(\"%s\", %s);\n", constName, constName))
	}
	
	code.WriteString("\t\t\treturn Collections.unmodifiableMap(pairs);\n")
	code.WriteString("\t\t}\n")
	
	return code.String()
}

// generateFormat 生成格式化方法
func (g *JavaGenerator) generateFormat(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	javaType := parser.GetJavaType(group.Constants[0].Type)
	boxedType := getBoxedType(javaType)
	
	code.WriteString("\t\t/**\n")
	code.WriteString("\t\t * 根据值和语言格式化标签\n")
	code.WriteString("\t\t * @param value 常量值\n")
	code.WriteString("\t\t * @param lang 语言代码\n")
	code.WriteString("\t\t * @return 格式化后的标签\n")
	code.WriteString("\t\t */\n")
	code.WriteString(fmt.Sprintf("\t\tpublic static String format(%s value, String lang) {\n", javaType))
	code.WriteString(fmt.Sprintf("\t\t\tMap<String, Map<%s, String>> labels = new HashMap<>();\n", boxedType))
	code.WriteString("\t\t\t\n")
	
	// 中文标签
	code.WriteString(fmt.Sprintf("\t\t\tMap<%s, String> zhLabels = new HashMap<>();\n", boxedType))
	for _, constant := range group.Constants {
		constName := parser.ToJavaConstantName(constant.Name)
		label := constant.Label
		if label == "" {
			label = constant.Name
		}
		code.WriteString(fmt.Sprintf("\t\t\tzhLabels.put(%s, \"%s\");\n", constName, label))
	}
	code.WriteString("\t\t\tlabels.put(\"zh\", zhLabels);\n")
	code.WriteString("\t\t\t\n")
	
	// 英文标签
	code.WriteString(fmt.Sprintf("\t\t\tMap<%s, String> enLabels = new HashMap<>();\n", boxedType))
	for _, constant := range group.Constants {
		constName := parser.ToJavaConstantName(constant.Name)
		label := strings.ReplaceAll(constant.Name, "_", " ")
		label = strings.Title(strings.ToLower(label))
		code.WriteString(fmt.Sprintf("\t\t\tenLabels.put(%s, \"%s\");\n", constName, label))
	}
	code.WriteString("\t\t\tlabels.put(\"en\", enLabels);\n\n")
	
	code.WriteString("\t\t\tif (labels.containsKey(lang) && labels.get(lang).containsKey(value)) {\n")
	code.WriteString("\t\t\t\treturn labels.get(lang).get(value);\n")
	code.WriteString("\t\t\t}\n\n")
	code.WriteString("\t\t\t// 默认返回英文\n")
	code.WriteString("\t\t\tif (labels.get(\"en\").containsKey(value)) {\n")
	code.WriteString("\t\t\t\treturn labels.get(\"en\").get(value);\n")
	code.WriteString("\t\t\t}\n\n")
	code.WriteString("\t\t\treturn \"Unknown(\" + value + \")\";\n")
	code.WriteString("\t\t}\n")
	
	return code.String()
}

// generateIsValid 生成验证方法
func (g *JavaGenerator) generateIsValid(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	javaType := parser.GetJavaType(group.Constants[0].Type)
	
	code.WriteString("\t\t/**\n")
	code.WriteString("\t\t * 验证值是否有效\n")
	code.WriteString("\t\t * @param value 要验证的值\n")
	code.WriteString("\t\t * @return 是否为有效常量\n")
	code.WriteString("\t\t */\n")
	code.WriteString(fmt.Sprintf("\t\tpublic static boolean isValid(%s value) {\n", javaType))
	code.WriteString("\t\t\treturn getAllValues().contains(value);\n")
	code.WriteString("\t\t}\n")
	
	return code.String()
}

// generateFromString 生成从字符串获取值的方法
func (g *JavaGenerator) generateFromString(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	javaType := parser.GetJavaType(group.Constants[0].Type)
	boxedType := getBoxedType(javaType)
	
	code.WriteString("\t\t/**\n")
	code.WriteString("\t\t * 从字符串键名获取常量值\n")
	code.WriteString("\t\t * @param key 常量键名\n")
	code.WriteString("\t\t * @return 常量值，找不到时返回null\n")
	code.WriteString("\t\t */\n")
	code.WriteString(fmt.Sprintf("\t\tpublic static %s fromString(String key) {\n", boxedType))
	code.WriteString("\t\t\treturn getKeyValuePairs().get(key);\n")
	code.WriteString("\t\t}\n")
	
	return code.String()
}

// generateGetDescription 生成获取描述的方法
func (g *JavaGenerator) generateGetDescription(group *parser.ConstantGroup) string {
	var code strings.Builder
	
	javaType := parser.GetJavaType(group.Constants[0].Type)
	boxedType := getBoxedType(javaType)
	
	code.WriteString("\t\t/**\n")
	code.WriteString("\t\t * 获取常量值的详细描述\n")
	code.WriteString("\t\t * @param value 常量值\n")
	code.WriteString("\t\t * @return 详细描述\n")
	code.WriteString("\t\t */\n")
	code.WriteString(fmt.Sprintf("\t\tpublic static String getDescription(%s value) {\n", javaType))
	code.WriteString(fmt.Sprintf("\t\t\tMap<%s, String> descriptions = new HashMap<>();\n", boxedType))
	
	for _, constant := range group.Constants {
		constName := parser.ToJavaConstantName(constant.Name)
		desc := constant.Desc
		if desc == "" {
			desc = constant.Label
		}
		code.WriteString(fmt.Sprintf("\t\t\tdescriptions.put(%s, \"%s\");\n", constName, desc))
	}
	
	code.WriteString("\t\t\t\n")
	code.WriteString("\t\t\treturn descriptions.getOrDefault(value, \"未知常量值: \" + value);\n")
	code.WriteString("\t\t}\n")
	
	return code.String()
}

// getBoxedType 获取基本类型的装箱类型
func getBoxedType(primitiveType string) string {
	switch primitiveType {
	case "int":
		return "Integer"
	case "double":
		return "Double"
	case "boolean":
		return "Boolean"
	case "float":
		return "Float"
	case "long":
		return "Long"
	default:
		return primitiveType
	}
}

// GenerateIndex Java不需要生成索引文件
func (g *JavaGenerator) GenerateIndex(allConstants []*parser.ConstantsFile) error {
	return nil
}