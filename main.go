package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cons-coder/generator"
	"cons-coder/parser"

	flag "github.com/spf13/pflag"
)

// Version information (will be set by build script)
var (
	Version   = "1.0.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
	GitBranch = "unknown"
	GitState  = "unknown"
	GitTag    = ""
)

func main() {
	var (
		dir         string
		output      string
		lang        string
		pkgName     string
		help        bool
		showVersion bool
	)

	flag.StringVarP(&dir, "dir", "d", "", "XML配置文件目录 (必填)")
	flag.StringVarP(&output, "output", "o", "", "输出代码目录 (必填)")
	flag.StringVarP(&lang, "lang", "l", "", "目标语言 (python/go/java/swift/kotlin/typescript/javascript) (必填)")
	flag.StringVarP(&pkgName, "package", "p", "", "包名 (可选，Go/Java/Kotlin语言使用)")
	flag.BoolVarP(&help, "help", "h", false, "显示帮助信息")
	flag.BoolVarP(&showVersion, "version", "v", false, "显示版本信息")

	flag.Parse()

	if showVersion {
		printVersion()
		os.Exit(0)
	}

	if help {
		printHelp()
		os.Exit(0)
	}

	// 验证必填参数
	if dir == "" || output == "" || lang == "" {
		fmt.Println("错误: 缺少必填参数")
		printHelp()
		os.Exit(1)
	}

	// 验证语言参数
	supportedLangs := []string{"python", "go", "java", "swift", "kotlin", "typescript", "javascript"}
	if !contains(supportedLangs, lang) {
		fmt.Printf("错误: 不支持的语言 '%s'\n", lang)
		fmt.Printf("支持的语言: %s\n", strings.Join(supportedLangs, ", "))
		os.Exit(1)
	}

	// 设置默认包名
	if pkgName == "" {
		switch lang {
		case "go":
			pkgName = "cons"
		case "java", "kotlin":
			pkgName = "com.example.constants"
		}
	}

	// 验证目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("错误: 目录 '%s' 不存在\n", dir)
		os.Exit(1)
	}

	// 创建输出目录
	if err := os.MkdirAll(output, 0755); err != nil {
		fmt.Printf("错误: 无法创建输出目录 '%s': %v\n", output, err)
		os.Exit(1)
	}

	// 读取所有XML文件
	xmlFiles, err := filepath.Glob(filepath.Join(dir, "*.xml"))
	if err != nil {
		log.Fatalf("错误: 读取XML文件失败: %v", err)
	}

	if len(xmlFiles) == 0 {
		fmt.Printf("警告: 在目录 '%s' 中没有找到XML文件\n", dir)
		os.Exit(0)
	}

	fmt.Printf("找到 %d 个XML文件\n", len(xmlFiles))

	// 解析所有XML文件
	var allConstants []*parser.ConstantsFile
	for _, xmlFile := range xmlFiles {
		fmt.Printf("正在解析: %s\n", xmlFile)

		constants, err := parser.ParseXMLFile(xmlFile)
		if err != nil {
			log.Printf("警告: 解析文件 '%s' 失败: %v", xmlFile, err)
			continue
		}

		allConstants = append(allConstants, constants)
	}

	if len(allConstants) == 0 {
		fmt.Println("错误: 没有成功解析任何XML文件")
		os.Exit(1)
	}

	// 生成代码
	config := generator.Config{
		Language:    lang,
		OutputDir:   output,
		PackageName: pkgName,
		Version:     Version,
	}

	gen := generator.New(config)

	for _, constants := range allConstants {
		fmt.Printf("正在生成 %s 代码: %s\n", lang, constants.FileName)

		if err := gen.Generate(constants); err != nil {
			log.Printf("警告: 生成文件 '%s' 的代码失败: %v", constants.FileName, err)
			continue
		}
	}

	// 生成索引文件（Python的__init__.py, TypeScript的index.ts等）
	if lang == "python" || lang == "typescript" || lang == "javascript" {
		fmt.Printf("正在生成索引文件...\n")
		if err := gen.GenerateIndex(allConstants); err != nil {
			log.Printf("警告: 生成索引文件失败: %v", err)
		}
	}

	fmt.Println("代码生成完成!")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func printHelp() {
	fmt.Println("常量代码生成器 (Constants Code Generator)")
	fmt.Printf("版本: %s\n\n", Version)
	fmt.Println("用法:")
	fmt.Println("  cons-coder --dir <XML目录> --output <输出目录> --lang <语言> [选项]")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  cons-coder --dir ./data --output ./python-codes --lang python")
	fmt.Println("  cons-coder --dir ./data --output ./go-codes --lang go --package constants")
	fmt.Println()
	fmt.Println("选项:")
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Printf("cons-coder version %s\n", Version)
	
	// Show detailed version info if available
	if BuildTime != "unknown" || GitCommit != "unknown" {
		fmt.Println("\nBuild Information:")
		fmt.Printf("  Build Time:  %s\n", BuildTime)
		fmt.Printf("  Git Commit:  %s\n", GitCommit)
		if GitTag != "" {
			fmt.Printf("  Git Tag:     %s\n", GitTag)
		}
		fmt.Printf("  Git Branch:  %s\n", GitBranch)
		fmt.Printf("  Git State:   %s\n", GitState)
		fmt.Printf("  Go Version:  %s\n", runtime.Version())
		fmt.Printf("  OS/Arch:     %s/%s\n", runtime.GOOS, runtime.GOARCH)
	}
}
