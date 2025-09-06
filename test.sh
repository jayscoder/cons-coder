#!/bin/bash

echo "========================================="
echo "Testing CLASS mode (output to output/class)"
echo "========================================="

# Class mode tests (输出到 output/class)
go run main.go --dir ./data --output ./output/class/go --lang go --mode class --package constants
go run main.go --dir ./data --output ./output/class/python --lang python --mode class
go run main.go --dir ./data --output ./output/class/typescript --lang typescript --mode class
go run main.go --dir ./data --output ./output/class/javascript --lang javascript --mode class
go run main.go --dir ./data --output ./output/class/kotlin --lang kotlin --mode class --package com.app.constants
go run main.go --dir ./data --output ./output/class/swift --lang swift --mode class
go run main.go --dir ./data --output ./output/class/java --lang java --mode class --package com.app.constants

echo ""
echo "========================================="
echo "Testing CONST mode (output to output/const)"
echo "========================================="

# Const mode tests (输出到 output/const)
go run main.go --dir ./data --output ./output/const/go --lang go --mode const --package constants
go run main.go --dir ./data --output ./output/const/python --lang python --mode const
go run main.go --dir ./data --output ./output/const/typescript --lang typescript --mode const
go run main.go --dir ./data --output ./output/const/javascript --lang javascript --mode const
go run main.go --dir ./data --output ./output/const/kotlin --lang kotlin --mode const --package com.app.constants
go run main.go --dir ./data --output ./output/const/swift --lang swift --mode const
go run main.go --dir ./data --output ./output/const/java --lang java --mode const --package com.app.constants

echo ""
echo "========================================="
echo "All tests completed!"
echo "Class mode output: ./output/class/"
echo "Const mode output: ./output/const/"
echo "========================================="