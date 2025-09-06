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

// PythonGenerator Python代码生成器
type PythonGenerator struct {
	BaseGenerator
}

// NewPythonGenerator 创建Python生成器
func NewPythonGenerator(config Config) *PythonGenerator {
	return &PythonGenerator{
		BaseGenerator: BaseGenerator{Config: config},
	}
}

// Generate 生成Python代码
func (g *PythonGenerator) Generate(constants *parser.ConstantsFile) error {
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
		// class模式 - 生成类
		// 导入
		code.WriteString("from typing import List, Dict, Optional, Any\n")
		code.WriteString("\n\n")

		// 生成每个常量组的类
		for _, group := range constants.Groups {
			code.WriteString(g.generateGroupClass(group, constants.Label))
			code.WriteString("\n\n")
		}
	}

	// 写入文件
	outputPath := g.GetOutputFilePath(constants.FileName)
	return os.WriteFile(outputPath, []byte(code.String()), 0644)
}

// generateConstGroup 生成const模式的常量组
func (g *PythonGenerator) generateConstGroup(group *parser.ConstantGroup, projectLabel string) string {
	var code strings.Builder

	// 生成注释
	code.WriteString(fmt.Sprintf("# %s %s - %s\n", group.Name, group.Label, projectLabel))
	
	// 按字母顺序排序常量
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToPythonName(constants[i].Name) < parser.ToPythonName(constants[j].Name)
	})
	
	// 生成常量定义
	for _, constant := range constants {
		constName := fmt.Sprintf("%s_%s", strings.ToUpper(group.Name), strings.ToUpper(constant.Name))
		value := parser.FormatValue(constant.Value, constant.Type, "python")
		comment := constant.Label
		code.WriteString(fmt.Sprintf("%s = %s  # %s\n", constName, value, comment))
	}
	
	return code.String()
}

// generateGroupClass 生成常量组类
func (g *PythonGenerator) generateGroupClass(group *parser.ConstantGroup, projectLabel string) string {
	var code strings.Builder

	className := parser.ToGoName(group.Name)

	// 类定义和文档字符串
	code.WriteString(fmt.Sprintf("class %s:\n", className))
	code.WriteString(fmt.Sprintf(`    """%s - %s`, group.Label, projectLabel))
	code.WriteString("\n    \n")
	code.WriteString(fmt.Sprintf("    项目: %s\n", projectLabel))
	code.WriteString(fmt.Sprintf("    常量组: %s\n", group.Label))
	code.WriteString(`    """`)
	code.WriteString("\n\n")

	// 按字母顺序排序常量
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToPythonName(constants[i].Name) < parser.ToPythonName(constants[j].Name)
	})

	// 常量定义
	code.WriteString("    # 常量定义 (按字母顺序排列)\n")
	for _, constant := range constants {
		constName := parser.ToPythonName(constant.Name)
		value := parser.FormatValue(constant.Value, constant.Type, "python")
		comment := constant.Label
		if comment == "" {
			comment = constant.Label
		}

		// 计算对齐空格
		spaces := 20 - len(constName)
		if spaces < 1 {
			spaces = 1
		}

		code.WriteString(fmt.Sprintf("    %s = %s%s# %s\n",
			constName, value, strings.Repeat(" ", spaces), comment))
	}

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

	return code.String()
}

// generateGetAllValues 生成获取所有值的方法
func (g *PythonGenerator) generateGetAllValues(group *parser.ConstantGroup) string {
	var code strings.Builder

	// className := parser.ToGoName(group.Name)

	code.WriteString("    @classmethod\n")
	code.WriteString(fmt.Sprintf("    def get_all_values(cls) -> List[%s]:\n",
		parser.GetPythonType(group.Constants[0].Type)))
	code.WriteString(fmt.Sprintf(`        """获取所有%s常量值"""`, group.Label))
	code.WriteString("\n        return [")

	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToPythonName(constants[i].Name) < parser.ToPythonName(constants[j].Name)
	})

	for i, constant := range constants {
		if i > 0 {
			code.WriteString(", ")
		}
		code.WriteString(fmt.Sprintf("cls.%s", parser.ToPythonName(constant.Name)))
	}

	code.WriteString("]\n")

	return code.String()
}

// generateGetAllKeys 生成获取所有键的方法
func (g *PythonGenerator) generateGetAllKeys(group *parser.ConstantGroup) string {
	var code strings.Builder

	code.WriteString("    @classmethod\n")
	code.WriteString("    def get_all_keys(cls) -> List[str]:\n")
	code.WriteString(fmt.Sprintf(`        """获取所有%s常量键名"""`, group.Label))
	code.WriteString("\n        return [")

	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToPythonName(constants[i].Name) < parser.ToPythonName(constants[j].Name)
	})

	for i, constant := range constants {
		if i > 0 {
			code.WriteString(", ")
		}
		code.WriteString(fmt.Sprintf(`"%s"`, parser.ToPythonName(constant.Name)))
	}

	code.WriteString("]\n")

	return code.String()
}

// generateGetKeyValuePairs 生成获取键值对的方法
func (g *PythonGenerator) generateGetKeyValuePairs(group *parser.ConstantGroup) string {
	var code strings.Builder

	code.WriteString("    @classmethod\n")
	code.WriteString(fmt.Sprintf("    def get_key_value_pairs(cls) -> Dict[str, %s]:\n",
		parser.GetPythonType(group.Constants[0].Type)))
	code.WriteString(`        """获取键值对字典"""`)
	code.WriteString("\n        return {\n")

	// 按字母顺序排序
	constants := make([]*parser.Constant, len(group.Constants))
	copy(constants, group.Constants)
	sort.Slice(constants, func(i, j int) bool {
		return parser.ToPythonName(constants[i].Name) < parser.ToPythonName(constants[j].Name)
	})

	for _, constant := range constants {
		constName := parser.ToPythonName(constant.Name)
		code.WriteString(fmt.Sprintf(`            "%s": cls.%s,`, constName, constName))
		code.WriteString("\n")
	}

	code.WriteString("        }\n")

	return code.String()
}

// generateFormatValue 生成格式化值的方法
func (g *PythonGenerator) generateFormatValue(group *parser.ConstantGroup) string {
	var code strings.Builder

	valueType := parser.GetPythonType(group.Constants[0].Type)

	code.WriteString("    @classmethod\n")
	code.WriteString(fmt.Sprintf("    def format_value(cls, value: %s) -> str:\n", valueType))
	code.WriteString(fmt.Sprintf(`        """根据值格式化%s的标签`, group.Label))
	code.WriteString("\n        \n")
	code.WriteString("        Args:\n")
	code.WriteString("            value: 常量值\n")
	code.WriteString("            \n")
	code.WriteString("        Returns:\n")
	code.WriteString("            格式化后的标签，找不到时返回 'Unknown(value)'\n")
	code.WriteString(`        """`)
	code.WriteString("\n        labels = {\n")
	for _, constant := range group.Constants {
		constName := parser.ToPythonName(constant.Name)
		label := constant.Label
		if label == "" {
			label = constant.Name
		}
		code.WriteString(fmt.Sprintf("            cls.%s: '%s',\n", constName, label))
	}
	code.WriteString("        }\n\n")
	code.WriteString("        if value in labels:\n")
	code.WriteString("            return labels[value]\n\n")
	code.WriteString("        return f'Unknown({value})'\n")

	return code.String()
}

// generateIsValid 生成验证方法
func (g *PythonGenerator) generateIsValid(group *parser.ConstantGroup) string {
	var code strings.Builder

	valueType := parser.GetPythonType(group.Constants[0].Type)

	code.WriteString("    @classmethod\n")
	code.WriteString(fmt.Sprintf("    def is_valid(cls, value: %s) -> bool:\n", valueType))
	code.WriteString(fmt.Sprintf(`        """验证值是否为有效的%s常量"""`, group.Label))
	code.WriteString("\n        return value in cls.get_all_values()\n")

	return code.String()
}

// generateFromString 生成从字符串获取值的方法
func (g *PythonGenerator) generateFromString(group *parser.ConstantGroup) string {
	var code strings.Builder

	valueType := parser.GetPythonType(group.Constants[0].Type)

	code.WriteString("    @classmethod\n")
	code.WriteString(fmt.Sprintf("    def from_string(cls, key: str) -> Optional[%s]:\n", valueType))
	code.WriteString(fmt.Sprintf(`        """从字符串键名获取%s常量值`, group.Label))
	code.WriteString("\n        \n")
	code.WriteString("        Args:\n")
	code.WriteString("            key: 常量键名\n")
	code.WriteString("            \n")
	code.WriteString("        Returns:\n")
	code.WriteString("            常量值，找不到时返回 None\n")
	code.WriteString(`        """`)
	code.WriteString("\n        mapping = cls.get_key_value_pairs()\n")
	code.WriteString("        return mapping.get(key)\n")

	return code.String()
}


// GenerateIndex 生成Python的__init__.py文件
func (g *PythonGenerator) GenerateIndex(allConstants []*parser.ConstantsFile) error {
	var code strings.Builder

	// 文件头注释
	code.WriteString(`"""
常量包初始化文件

生成时间: `)
	code.WriteString(FormatGenerationTime(time.Now()))
	code.WriteString("\n生成工具: cons-coder v")
	code.WriteString(g.Config.Version)
	code.WriteString("\n")
	code.WriteString(`"""`)
	code.WriteString("\n\n")

	// 导入所有常量类
	code.WriteString("# 导入所有常量类\n")
	for _, constants := range allConstants {
		var classes []string
		for _, group := range constants.Groups {
			classes = append(classes, parser.ToGoName(group.Name))
		}
		if len(classes) > 0 {
			code.WriteString(fmt.Sprintf("from .%s import %s\n",
				constants.FileName, strings.Join(classes, ", ")))
		}
	}

	code.WriteString("\n# 导出所有常量类\n")
	code.WriteString("__all__ = [\n")

	for _, constants := range allConstants {
		if len(constants.Groups) > 0 {
			code.WriteString(fmt.Sprintf("    # %s.py 中的常量\n", constants.FileName))
			for _, group := range constants.Groups {
				code.WriteString(fmt.Sprintf("    '%s',\n", parser.ToGoName(group.Name)))
			}
			code.WriteString("    \n")
		}
	}

	code.WriteString("]\n\n")

	// 版本信息
	code.WriteString("# 版本信息\n")
	code.WriteString(fmt.Sprintf("__version__ = '%s'\n", g.Config.Version))
	code.WriteString("__generator__ = 'cons-coder'\n")

	// 写入文件
	outputPath := filepath.Join(g.Config.OutputDir, "__init__.py")
	return os.WriteFile(outputPath, []byte(code.String()), 0644)
}
