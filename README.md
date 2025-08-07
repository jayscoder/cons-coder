# 常量代码生成器 (Constants Code Generator)

一个强大的多语言常量代码生成工具，支持从XML配置文件生成Python、Go、Java、Swift、Kotlin等语言的常量定义代码。

## 特性

- 🌍 **多语言支持**: 支持Python、Go、Java、Swift、Kotlin等主流编程语言
- 📝 **XML配置**: 使用简洁的XML格式定义常量
- 🔧 **丰富的工具方法**: 自动生成验证、格式化、转换等实用方法
- 🌐 **多语言标签**: 支持中文、英文、日文等多语言标签
- ⚡ **高效生成**: 一次配置，多语言输出

## 安装与构建

### 前置要求
- Go 1.16+

### 构建
```bash
go build -o cons-coder main.go
```

## 使用方法

### 基本命令
```bash
# 生成Python代码
cons-coder --dir ./data --output ./python-codes --lang python

# 生成Go代码  
cons-coder --dir ./data --output ./go-codes --lang go

# 生成Go代码并指定包名
cons-coder --dir ./data --output ./go-codes --lang go --package constants

# 生成Java代码
cons-coder --dir ./data --output ./java-codes --lang java

# 生成Swift代码
cons-coder --dir ./data --output ./swift-codes --lang swift

# 生成Kotlin代码
cons-coder --dir ./data --output ./kotlin-codes --lang kotlin

# 生成Typescript代码
cons-coder --dir ./data --output ./typescript-codes --lang typescript

# 生成Javascript代码
cons-coder --dir ./data --output ./javascript-codes --lang javascript
```

### 参数说明
| 参数 | 说明 | 必填 | 默认值 |
|------|------|------|--------|
| `--dir` | XML配置文件目录 | 是 | - |
| `--output` | 输出代码目录 | 是 | - |
| `--lang` | 目标语言 (python/go/java/swift/kotlin) | 是 | - |
| `--package` | Go语言包名、Java语言包名、Kotlin语言包名 | 否 | constants/com.example.constants |

## XML数据格式

### 基本结构
```xml
<?xml version='1.0' encoding='utf-8'?>
<constants label="项目或模块描述">
    <constant_group_name label="常量组描述">
        <constant_name type="数据类型" label="显示标签" desc="详细描述" value="常量值" />
        <!-- 更多常量定义... -->
    </constant_group_name>
    <!-- 更多常量组... -->
</constants>
```

### 文件组织规则

- **XML文件名**：XML文件名（不含扩展名）将作为生成代码的文件名
- **多文件支持**：支持多个XML文件，每个文件生成对应的代码文件
- **命名规范**：建议使用小写字母和下划线命名XML文件

**示例文件结构：**
```
data/
├── user.xml          # 生成 user.py / user.go / User.java 等
├── group.xml         # 生成 group.py / group.go / Group.java 等
└── admin.xml         # 生成 admin.py / admin.go / Admin.java 等
```

### 完整示例

**文件：data/app.xml**
```xml
<?xml version='1.0' encoding='utf-8'?>
<constants label="其他">
    <app_sid_type label="散弹号类型">
        <user_preset type="int" label="用户预设" desc="用户预设" value="1" />
        <group_camp type="int" label="大本营预设" desc="大本营预设" value="2" />
        <group_battle type="int" label="对线群预设" desc="对线群预设" value="3" />
        <dynamic type="int" label="动态预设" desc="动态预设" value="4" />
        <world type="int" label="世界预设" desc="世界预设" value="5" />
        <user_custom type="int" label="用户自定义" desc="用户自定义" value="6" />
    </app_sid_type>
    <app_user_token_type label="用户token类型">
        <register type="int" label="注册" desc="注册" value="0" />
        <login type="int" label="登录" desc="登录" value="1" />
        <unbind_old_phone type="int" label="解绑旧手机号" desc="解绑旧手机号" value="2" />
        <bind_new_phone type="int" label="绑定新手机号" desc="绑定新手机号" value="3" />
        <access type="int" label="访问App" desc="访问App" value="4" />
    </app_user_token_type>
</constants>
```

**文件：data/user.xml**
```xml
<?xml version='1.0' encoding='utf-8'?>
<constants label="用户相关">
    <suspend_handle_source label="封禁处理来源">
        <manager type="int" label="管理员" desc="由管理员手动执行的封禁操作" value="0" />
        <osscr type="int" label="OSS内容安全审核" desc="阿里云OSS内容安全自动审核触发" value="1" />
        <system type="int" label="系统自动" desc="系统根据规则自动执行的封禁" value="2" />
        <app_official type="int" label="官方应用" desc="官方应用程序触发的封禁" value="3" />
    </suspend_handle_source>
    
    <user_status label="用户状态">
        <normal type="int" label="正常" desc="用户账号状态正常" value="1" />
        <suspended type="int" label="已封禁" desc="用户账号已被封禁" value="2" />
        <deleted type="int" label="已删除" desc="用户账号已被删除" value="3" />
    </user_status>
</constants>
```

### 支持的数据类型
- `int`: 整型数值
- `string`: 字符串类型
- `float`: 浮点型数值
- `bool`: 布尔类型

## 输出代码结构

### 文件生成规则

1. **Go语言** (`--lang go`)
   - 生成文件：`{xml文件名}.go`
   - 默认包名：`cons`
   - 自定义包名：通过 `--package` 参数指定
   - 示例：`app.xml` → `app.go`

2. **Python语言** (`--lang python`)
   - 生成文件：`{xml文件名}.py` 
   - 自动生成：`__init__.py` (导入所有常量类)
   - 示例：`app.xml` → `app.py`

3. **其他语言**
   - Java：`{XML文件名CamelCase}.java`
   - Swift：`{XML文件名CamelCase}.swift`
   - Kotlin：`{XML文件名CamelCase}.kt`

### 代码头部注释

生成的代码文件头部将包含以下信息：
- XML文件描述（来自根标签 `label` 属性）
- 原始XML文件名和路径
- XML文件最后修改时间
- 代码生成时间
- 生成工具版本信息

### Go语言示例

**生成的 `app.go` 文件：**

```go
package cons

import (
	"fmt"
	"time"
)

/*
 * 其他
 * 
 * 源文件: app.xml
 * 最后修改: 2024-01-15 14:30:25
 * 生成时间: 2024-01-15 16:45:10
 * 生成工具: cons-coder v1.0.0
 */

// AppSidType 散弹号类型 - 其他
type appSidTypeCons struct {
	UserPreset   int // 用户预设
	GroupCamp    int // 大本营预设
	GroupBattle  int // 对线群预设
	Dynamic      int // 动态预设
	World        int // 世界预设
	UserCustom   int // 用户自定义
}

// AppSidType 常量实例
var AppSidType = appSidTypeCons{
	UserPreset:  1,
	GroupCamp:   2,
	GroupBattle: 3,
	Dynamic:     4,
	World:       5,
	UserCustom:  6,
}

// AppUserTokenType 用户token类型 - 其他
type appUserTokenTypeCons struct {
	Register        int // 注册
	Login          int // 登录
	UnbindOldPhone int // 解绑旧手机号
	BindNewPhone   int // 绑定新手机号
	Access         int // 访问App
}

// AppUserTokenType 常量实例
var AppUserTokenType = appUserTokenTypeCons{
	Register:        0,
	Login:          1,
	UnbindOldPhone: 2,
	BindNewPhone:   3,
	Access:         4,
}

// 为每个常量组添加方法...
// (省略详细的方法实现，参考上面的完整示例)
```

### Python语言示例

**生成的 `app.py` 文件：**

```python
"""
其他

源文件: app.xml
最后修改: 2024-01-15 14:30:25
生成时间: 2024-01-15 16:45:10  
生成工具: cons-coder v1.0.0
"""

from typing import List, Dict, Optional


class AppSidType:
    """散弹号类型 - 其他
    
    项目: 其他
    常量组: 散弹号类型
    """

    # 常量定义 (按字母顺序排列)
    DYNAMIC = 4        # 动态预设
    GROUP_BATTLE = 3   # 对线群预设
    GROUP_CAMP = 2     # 大本营预设
    USER_CUSTOM = 6    # 用户自定义
    USER_PRESET = 1    # 用户预设
    WORLD = 5          # 世界预设

    @classmethod
    def get_all_values(cls) -> List[int]:
        """获取所有散弹号类型常量值"""
        return [cls.DYNAMIC, cls.GROUP_BATTLE, cls.GROUP_CAMP, 
                cls.USER_CUSTOM, cls.USER_PRESET, cls.WORLD]

    @classmethod  
    def get_all_keys(cls) -> List[str]:
        """获取所有散弹号类型常量键名"""
        return ["DYNAMIC", "GROUP_BATTLE", "GROUP_CAMP", 
                "USER_CUSTOM", "USER_PRESET", "WORLD"]

    # 省略其他方法...


class AppUserTokenType:
    """用户token类型 - 其他
    
    项目: 其他
    常量组: 用户token类型
    """

    # 常量定义 (按字母顺序排列)
    ACCESS = 4            # 访问App
    BIND_NEW_PHONE = 3    # 绑定新手机号
    LOGIN = 1             # 登录
    REGISTER = 0          # 注册
    UNBIND_OLD_PHONE = 2  # 解绑旧手机号

    @classmethod
    def get_all_values(cls) -> List[int]:
        """获取所有用户token类型常量值"""
        return [cls.ACCESS, cls.BIND_NEW_PHONE, cls.LOGIN, 
                cls.REGISTER, cls.UNBIND_OLD_PHONE]

    # 省略其他方法...
```

**生成的 `__init__.py` 文件：**

```python
"""
常量包初始化文件

生成时间: 2024-01-15 16:45:10
生成工具: cons-coder v1.0.0
"""

# 导入所有常量类
from .app import AppSidType, AppUserTokenType
from .user import SuspendHandleSource, UserStatus  
from .group import GroupType, GroupStatus

# 导出所有常量类
__all__ = [
    # app.py 中的常量
    'AppSidType',
    'AppUserTokenType',
    
    # user.py 中的常量  
    'SuspendHandleSource',
    'UserStatus',
    
    # group.py 中的常量
    'GroupType', 
    'GroupStatus',
]

# 版本信息
__version__ = '1.0.0'
__generator__ = 'cons-coder'
```

### Java语言示例

**生成的 `App.java` 文件：**

```java
package com.example.constants;

import java.util.*;

/**
 * 其他
 * 
 * 源文件: app.xml
 * 最后修改: 2024-01-15 14:30:25
 * 生成时间: 2024-01-15 16:45:10
 * 生成工具: cons-coder v1.0.0
 */

/**
 * 散弹号类型 - 其他
 */
public final class AppSidType {
    /** 动态预设 */
    public static final int DYNAMIC = 4;
    /** 对线群预设 */
    public static final int GROUP_BATTLE = 3;
    /** 大本营预设 */
    public static final int GROUP_CAMP = 2;
    /** 用户自定义 */
    public static final int USER_CUSTOM = 6;
    /** 用户预设 */
    public static final int USER_PRESET = 1;
    /** 世界预设 */
    public static final int WORLD = 5;

    // 私有构造函数，防止实例化
    private AppSidType() {
        throw new AssertionError("常量类不应被实例化");
    }

    // 省略方法实现...
}

/**
 * 用户token类型 - 其他
 */
public final class AppUserTokenType {
    /** 访问App */
    public static final int ACCESS = 4;
    /** 绑定新手机号 */
    public static final int BIND_NEW_PHONE = 3;
    /** 登录 */
    public static final int LOGIN = 1;
    /** 注册 */
    public static final int REGISTER = 0;
    /** 解绑旧手机号 */
    public static final int UNBIND_OLD_PHONE = 2;

    // 私有构造函数，防止实例化
    private AppUserTokenType() {
        throw new AssertionError("常量类不应被实例化");
    }

    // 省略方法实现...
}
```

### 原有的Python完整示例

```python
    """封禁处理来源 - 用户相关
    
    项目: 用户相关
    常量组: 封禁处理来源
    """

    # 常量定义 (按字母顺序排列)
    APP_OFFICIAL = 3  # 官方应用
    MANAGER = 0       # 管理员  
    OSSCR = 1         # OSS内容安全审核
    SYSTEM = 2        # 系统自动

    @classmethod
    def get_all_values(cls) -> list[int]:
        """获取所有封禁处理来源常量值"""
        return [cls.APP_OFFICIAL, cls.MANAGER, cls.OSSCR, cls.SYSTEM]

    @classmethod  
    def get_all_keys(cls) -> list[str]:
        """获取所有封禁处理来源常量键名"""
        return ["APP_OFFICIAL", "MANAGER", "OSSCR", "SYSTEM"]

    @classmethod
    def get_key_value_pairs(cls) -> dict[str, int]:
        """获取键值对字典"""
        return {
            "APP_OFFICIAL": cls.APP_OFFICIAL,
            "MANAGER": cls.MANAGER, 
            "OSSCR": cls.OSSCR,
            "SYSTEM": cls.SYSTEM,
        }

    @classmethod
    def format_value(cls, value: int, lang: str = 'zh') -> str:
        """根据值和语言格式化封禁处理来源的标签
        
        Args:
            value: 常量值
            lang: 语言代码 ('zh', 'en', 'ja')
            
        Returns:
            格式化后的标签，找不到时返回 'Unknown(value)'
        """
        labels = {
            'zh': {
                cls.APP_OFFICIAL: '官方应用',
                cls.MANAGER: '管理员',
                cls.OSSCR: 'OSS内容安全审核', 
                cls.SYSTEM: '系统自动',
            },
            'en': {
                cls.APP_OFFICIAL: 'Official App',
                cls.MANAGER: 'Manager',
                cls.OSSCR: 'OSS Content Security',
                cls.SYSTEM: 'System Auto',
            },
            'ja': {
                cls.APP_OFFICIAL: '公式アプリ',
                cls.MANAGER: '管理者', 
                cls.OSSCR: 'OSSコンテンツセキュリティ',
                cls.SYSTEM: 'システム自動',
            },
        }

        if lang in labels and value in labels[lang]:
            return labels[lang][value]

        # 默认返回英文，如果英文也没有则返回数值
        if 'en' in labels and value in labels['en']:
            return labels['en'][value]

        return f'Unknown({value})'

    @classmethod
    def is_valid(cls, value: int) -> bool:
        """验证值是否为有效的封禁处理来源常量"""
        return value in cls.get_all_values()

    @classmethod
    def from_string(cls, key: str) -> int | None:
        """从字符串键名获取封禁处理来源常量值
        
        Args:
            key: 常量键名
            
        Returns:
            常量值，找不到时返回 None
        """
        mapping = cls.get_key_value_pairs()
        return mapping.get(key)

    @classmethod
    def get_description(cls, value: int) -> str:
        """获取常量值的详细描述"""
        descriptions = {
            cls.APP_OFFICIAL: '官方应用程序触发的封禁',
            cls.MANAGER: '由管理员手动执行的封禁操作',
            cls.OSSCR: '阿里云OSS内容安全自动审核触发',
            cls.SYSTEM: '系统根据规则自动执行的封禁',
        }
        return descriptions.get(value, f'未知常量值: {value}')
```

### Go

```go
package constants

import "fmt"

// SuspendHandleSource 封禁处理来源 - 用户相关
type suspendHandleSourceCons struct {
	AppOfficial int // 官方应用
	Manager     int // 管理员
	Osscr       int // OSS内容安全审核  
	System      int // 系统自动
}

// SuspendHandleSource 常量实例
var SuspendHandleSource = suspendHandleSourceCons{
	AppOfficial: 3,
	Manager:     0,
	Osscr:       1,
	System:      2,
}

// AllValues 返回所有SuspendHandleSource的值
func (s suspendHandleSourceCons) AllValues() []int {
	return []int{s.AppOfficial, s.Manager, s.Osscr, s.System}
}

// AllKeys 返回所有SuspendHandleSource的键名
func (s suspendHandleSourceCons) AllKeys() []string {
	return []string{"AppOfficial", "Manager", "Osscr", "System"}
}

// KeyValuePairs 返回键值对映射
func (s suspendHandleSourceCons) KeyValuePairs() map[string]int {
	return map[string]int{
		"AppOfficial": s.AppOfficial,
		"Manager":     s.Manager,
		"Osscr":       s.Osscr,
		"System":      s.System,
	}
}

// Format 根据值和语言格式化SuspendHandleSource的标签
func (s suspendHandleSourceCons) Format(value int, lang string) string {
	labels := map[string]map[int]string{
		"zh": {
			s.AppOfficial: "官方应用",
			s.Manager:     "管理员",
			s.Osscr:       "OSS内容安全审核",
			s.System:      "系统自动",
		},
		"en": {
			s.AppOfficial: "Official App",
			s.Manager:     "Manager", 
			s.Osscr:       "OSS Content Security",
			s.System:      "System Auto",
		},
		"ja": {
			s.AppOfficial: "公式アプリ",
			s.Manager:     "管理者",
			s.Osscr:       "OSSコンテンツセキュリティ", 
			s.System:      "システム自動",
		},
	}

	if langMap, exists := labels[lang]; exists {
		if label, exists := langMap[value]; exists {
			return label
		}
	}

	// 默认返回英文
	if enMap, exists := labels["en"]; exists {
		if label, exists := enMap[value]; exists {
			return label
		}
	}

	return fmt.Sprintf("Unknown(%v)", value)
}

// IsValid 检查值是否为有效的SuspendHandleSource常量
func (s suspendHandleSourceCons) IsValid(value int) bool {
	allValues := s.AllValues()
	for _, v := range allValues {
		if v == value {
			return true
		}
	}
	return false
}

// FromString 从字符串键名获取SuspendHandleSource常量值
func (s suspendHandleSourceCons) FromString(key string) (int, bool) {
	mapping := s.KeyValuePairs()
	value, exists := mapping[key]
	return value, exists
}

// GetDescription 获取常量值的详细描述
func (s suspendHandleSourceCons) GetDescription(value int) string {
	descriptions := map[int]string{
		s.AppOfficial: "官方应用程序触发的封禁",
		s.Manager:     "由管理员手动执行的封禁操作", 
		s.Osscr:       "阿里云OSS内容安全自动审核触发",
		s.System:      "系统根据规则自动执行的封禁",
	}
	
	if desc, exists := descriptions[value]; exists {
		return desc
	}
	return fmt.Sprintf("未知常量值: %v", value)
}
```

### Java

```java
package com.example.constants;

import java.util.*;

/**
 * 封禁处理来源 - 用户相关
 */
public final class SuspendHandleSource {
    /** 官方应用 */
    public static final int APP_OFFICIAL = 3;
    /** 管理员 */
    public static final int MANAGER = 0;
    /** OSS内容安全审核 */
    public static final int OSSCR = 1;  
    /** 系统自动 */
    public static final int SYSTEM = 2;

    // 私有构造函数，防止实例化
    private SuspendHandleSource() {
        throw new AssertionError("常量类不应被实例化");
    }

    /**
     * 获取所有常量值
     * @return 所有常量值的列表
     */
    public static List<Integer> getAllValues() {
        return Arrays.asList(APP_OFFICIAL, MANAGER, OSSCR, SYSTEM);
    }

    /**
     * 获取所有常量键名
     * @return 所有常量键名的列表
     */
    public static List<String> getAllKeys() {
        return Arrays.asList("APP_OFFICIAL", "MANAGER", "OSSCR", "SYSTEM");
    }

    /**
     * 获取键值对映射
     * @return 键值对映射
     */
    public static Map<String, Integer> getKeyValuePairs() {
        Map<String, Integer> pairs = new HashMap<>();
        pairs.put("APP_OFFICIAL", APP_OFFICIAL);
        pairs.put("MANAGER", MANAGER);
        pairs.put("OSSCR", OSSCR);
        pairs.put("SYSTEM", SYSTEM);
        return Collections.unmodifiableMap(pairs);
    }

    /**
     * 根据值和语言格式化标签
     * @param value 常量值
     * @param lang 语言代码
     * @return 格式化后的标签
     */
    public static String format(int value, String lang) {
        Map<String, Map<Integer, String>> labels = new HashMap<>();
        
        Map<Integer, String> zhLabels = new HashMap<>();
        zhLabels.put(APP_OFFICIAL, "官方应用");
        zhLabels.put(MANAGER, "管理员");
        zhLabels.put(OSSCR, "OSS内容安全审核");
        zhLabels.put(SYSTEM, "系统自动");
        labels.put("zh", zhLabels);
        
        Map<Integer, String> enLabels = new HashMap<>();
        enLabels.put(APP_OFFICIAL, "Official App");
        enLabels.put(MANAGER, "Manager");
        enLabels.put(OSSCR, "OSS Content Security");
        enLabels.put(SYSTEM, "System Auto");
        labels.put("en", enLabels);

        if (labels.containsKey(lang) && labels.get(lang).containsKey(value)) {
            return labels.get(lang).get(value);
        }

        // 默认返回英文
        if (labels.get("en").containsKey(value)) {
            return labels.get("en").get(value);
        }

        return "Unknown(" + value + ")";
    }

    /**
     * 验证值是否有效
     * @param value 要验证的值
     * @return 是否为有效常量
     */
    public static boolean isValid(int value) {
        return getAllValues().contains(value);
    }

    /**
     * 从字符串键名获取常量值
     * @param key 常量键名
     * @return 常量值，找不到时返回null
     */
    public static Integer fromString(String key) {
        return getKeyValuePairs().get(key);
    }

    /**
     * 获取常量值的详细描述
     * @param value 常量值
     * @return 详细描述
     */
    public static String getDescription(int value) {
        Map<Integer, String> descriptions = new HashMap<>();
        descriptions.put(APP_OFFICIAL, "官方应用程序触发的封禁");
        descriptions.put(MANAGER, "由管理员手动执行的封禁操作");
        descriptions.put(OSSCR, "阿里云OSS内容安全自动审核触发");
        descriptions.put(SYSTEM, "系统根据规则自动执行的封禁");
        
        return descriptions.getOrDefault(value, "未知常量值: " + value);
    }
}
```

### Swift

Swift语言生成器统一使用 `enum` 枚举来生成所有常量组：
- **整数类型**：生成 `enum` 枚举，原始值类型为 `Int`
- **字符串类型**：生成 `enum` 枚举，原始值类型为 `String`
- 所有枚举都支持 `CaseIterable`, `Codable`, `Identifiable`, `CustomStringConvertible` 协议
- Swift 保留关键字（如 `default`）会自动使用反引号转义

```swift
import Foundation

/// 用户状态 - 用户相关（整数类型示例）
public enum UserStatus: Int, CaseIterable, Codable, Identifiable, CustomStringConvertible {
    /// 正常
    case normal = 1
    /// 已封禁
    case suspended = 2
    /// 已删除
    case deleted = 3
    /// 未激活
    case inactive = 4

    public var id: Int { rawValue }
    
    public var description: String {
        format(lang: "en")
    }

    /// 获取中文标签
    public var label: String {
        switch self {
        case .normal: return "正常"
        case .suspended: return "已封禁"
        case .deleted: return "已删除"
        case .inactive: return "未激活"
        }
    }

    /// 获取详细描述
    public var detailDescription: String {
        switch self {
        case .normal: return "用户账号状态正常"
        case .suspended: return "用户账号已被封禁"
        case .deleted: return "用户账号已被删除"
        case .inactive: return "用户账号未激活"
        }
    }

    /// 根据语言获取标签
    /// - Parameter lang: 语言代码 ("zh", "en")
    /// - Returns: 格式化后的标签
    public func format(lang: String = "en") -> String {
        switch lang {
        case "zh":
            return label
        case "en":
            switch self {
            case .normal: return "Normal"
            case .suspended: return "Suspended"
            case .deleted: return "Deleted"
            case .inactive: return "Inactive"
            }
        default:
            return label
        }
    }

    /// 从字符串键名创建枚举
    /// - Parameter key: 常量键名
    /// - Returns: 枚举值，找不到时返回nil
    public static func fromString(_ key: String) -> Self? {
        switch key {
        case "normal": return .normal
        case "suspended": return .suspended
        case "deleted": return .deleted
        case "inactive": return .inactive
        default: return nil
        }
    }
}

/// 管理员统计视图类型 - 系统管理（字符串类型示例，包含关键字转义）
public enum AdminStatsViewType: String, CaseIterable, Codable, Identifiable, CustomStringConvertible {
    /// 日期
    case date = "date"
    /// 时间
    case datetime = "datetime"
    /// 默认（注意：default 是 Swift 关键字，会自动转义）
    case `default` = "default"
    /// 动态
    case dynamic = "dynamic"
    
    public var id: String { rawValue }
    
    public var description: String {
        format(lang: "en")
    }

    /// 获取中文标签
    public var label: String {
        switch self {
        case .date: return "日期"
        case .datetime: return "时间"
        case .`default`: return "默认"
        case .dynamic: return "动态"
        }
    }
    
    // ... 其他方法省略
}
```

### Kotlin

```kotlin
package com.example.constants

/**
 * 封禁处理来源 - 用户相关
 */
object SuspendHandleSource {
    /** 官方应用 */
    const val APP_OFFICIAL = 3
    /** 管理员 */
    const val MANAGER = 0
    /** OSS内容安全审核 */
    const val OSSCR = 1
    /** 系统自动 */
    const val SYSTEM = 2

    /** 获取所有常量值 */
    fun getAllValues(): List<Int> {
        return listOf(APP_OFFICIAL, MANAGER, OSSCR, SYSTEM)
    }

    /** 获取所有常量键名 */
    fun getAllKeys(): List<String> {
        return listOf("APP_OFFICIAL", "MANAGER", "OSSCR", "SYSTEM")
    }

    /** 获取键值对映射 */
    fun getKeyValuePairs(): Map<String, Int> {
        return mapOf(
            "APP_OFFICIAL" to APP_OFFICIAL,
            "MANAGER" to MANAGER,
            "OSSCR" to OSSCR,
            "SYSTEM" to SYSTEM
        )
    }

    /**
     * 根据值和语言格式化标签
     * @param value 常量值
     * @param lang 语言代码 (默认: "zh")
     * @return 格式化后的标签
     */
    fun format(value: Int, lang: String = "zh"): String {
        val labels = mapOf(
            "zh" to mapOf(
                APP_OFFICIAL to "官方应用",
                MANAGER to "管理员",
                OSSCR to "OSS内容安全审核",
                SYSTEM to "系统自动",
            ),
            "en" to mapOf(
                APP_OFFICIAL to "Official App",
                MANAGER to "Manager",
                OSSCR to "OSS Content Security", 
                SYSTEM to "System Auto",
            ),
            "ja" to mapOf(
                APP_OFFICIAL to "公式アプリ",
                MANAGER to "管理者",
                OSSCR to "OSSコンテンツセキュリティ",
                SYSTEM to "システム自動",
            ),
        )

        labels[lang]?.get(value)?.let { return it }

        // 默认返回英文
        labels["en"]?.get(value)?.let { return it }

        return "Unknown($value)"
    }

    /** 
     * 检查值是否为有效常量
     * @param value 要检查的值
     * @return 是否为有效常量
     */
    fun isValid(value: Int): Boolean {
        return getAllValues().contains(value)
    }

    /**
     * 从字符串键名获取常量值
     * @param key 常量键名
     * @return 常量值，找不到时返回null
     */
    fun fromString(key: String): Int? {
        return getKeyValuePairs()[key]
    }

    /**
     * 获取常量值的详细描述
     * @param value 常量值
     * @return 详细描述
     */
    fun getDescription(value: Int): String {
        val descriptions = mapOf(
            APP_OFFICIAL to "官方应用程序触发的封禁",
            MANAGER to "由管理员手动执行的封禁操作",
            OSSCR to "阿里云OSS内容安全自动审核触发",
            SYSTEM to "系统根据规则自动执行的封禁"
        )
        
        return descriptions[value] ?: "未知常量值: $value"
    }
}
```

### Typescript

```typescript
/**
 * 其他
 * 
 * 源文件: app.xml
 * 最后修改: 2024-01-15 14:30:25
 * 生成时间: 2024-01-15 16:45:10
 * 生成工具: cons-coder v1.0.0
 */

/**
 * 散弹号类型 - 其他
 */
export class AppSidType {
  /** 动态预设 */
  public static readonly DYNAMIC = 4;
  /** 对线群预设 */
  public static readonly GROUP_BATTLE = 3;
  /** 大本营预设 */
  public static readonly GROUP_CAMP = 2;
  /** 用户自定义 */
  public static readonly USER_CUSTOM = 6;
  /** 用户预设 */
  public static readonly USER_PRESET = 1;
  /** 世界预设 */
  public static readonly WORLD = 5;

  // 私有构造函数，防止实例化
  private constructor() {
    throw new Error('常量类不应被实例化');
  }

  /**
   * 获取所有散弹号类型常量值
   * @returns 所有常量值的数组
   */
  public static getAllValues(): readonly number[] {
    return [
      AppSidType.DYNAMIC,
      AppSidType.GROUP_BATTLE,
      AppSidType.GROUP_CAMP,
      AppSidType.USER_CUSTOM,
      AppSidType.USER_PRESET,
      AppSidType.WORLD,
    ] as const;
  }

  /**
   * 获取所有散弹号类型常量键名
   * @returns 所有常量键名的数组
   */
  public static getAllKeys(): readonly string[] {
    return [
      'DYNAMIC',
      'GROUP_BATTLE', 
      'GROUP_CAMP',
      'USER_CUSTOM',
      'USER_PRESET',
      'WORLD',
    ] as const;
  }

  /**
   * 获取键值对映射
   * @returns 键值对映射对象
   */
  public static getKeyValuePairs(): Readonly<Record<string, number>> {
    return {
      DYNAMIC: AppSidType.DYNAMIC,
      GROUP_BATTLE: AppSidType.GROUP_BATTLE,
      GROUP_CAMP: AppSidType.GROUP_CAMP,
      USER_CUSTOM: AppSidType.USER_CUSTOM,
      USER_PRESET: AppSidType.USER_PRESET,
      WORLD: AppSidType.WORLD,
    } as const;
  }

  /**
   * 根据值和语言格式化散弹号类型的标签
   * @param value 常量值
   * @param lang 语言代码 ('zh', 'en', 'ja')
   * @returns 格式化后的标签，找不到时返回 'Unknown(value)'
   */
  public static formatValue(value: number, lang: 'zh' | 'en' | 'ja' = 'zh'): string {
    const labels: Record<string, Record<number, string>> = {
      zh: {
        [AppSidType.DYNAMIC]: '动态预设',
        [AppSidType.GROUP_BATTLE]: '对线群预设',
        [AppSidType.GROUP_CAMP]: '大本营预设',
        [AppSidType.USER_CUSTOM]: '用户自定义',
        [AppSidType.USER_PRESET]: '用户预设',
        [AppSidType.WORLD]: '世界预设',
      },
      en: {
        [AppSidType.DYNAMIC]: 'Dynamic Preset',
        [AppSidType.GROUP_BATTLE]: 'Group Battle Preset',
        [AppSidType.GROUP_CAMP]: 'Group Camp Preset',
        [AppSidType.USER_CUSTOM]: 'User Custom',
        [AppSidType.USER_PRESET]: 'User Preset',
        [AppSidType.WORLD]: 'World Preset',
      },
      ja: {
        [AppSidType.DYNAMIC]: 'ダイナミックプリセット',
        [AppSidType.GROUP_BATTLE]: 'グループバトルプリセット',
        [AppSidType.GROUP_CAMP]: 'グループキャンププリセット',
        [AppSidType.USER_CUSTOM]: 'ユーザーカスタム',
        [AppSidType.USER_PRESET]: 'ユーザープリセット',
        [AppSidType.WORLD]: 'ワールドプリセット',
      },
    };

    const langLabels = labels[lang];
    if (langLabels && langLabels[value] !== undefined) {
      return langLabels[value];
    }

    // 默认返回英文
    const enLabels = labels.en;
    if (enLabels && enLabels[value] !== undefined) {
      return enLabels[value];
    }

    return `Unknown(${value})`;
  }

  /**
   * 验证值是否为有效的散弹号类型常量
   * @param value 要验证的值
   * @returns 是否为有效常量
   */
  public static isValid(value: number): value is typeof AppSidType.DYNAMIC | typeof AppSidType.GROUP_BATTLE | typeof AppSidType.GROUP_CAMP | typeof AppSidType.USER_CUSTOM | typeof AppSidType.USER_PRESET | typeof AppSidType.WORLD {
    return AppSidType.getAllValues().includes(value);
  }

  /**
   * 从字符串键名获取散弹号类型常量值
   * @param key 常量键名
   * @returns 常量值，找不到时返回 undefined
   */
  public static fromString(key: string): number | undefined {
    const mapping = AppSidType.getKeyValuePairs();
    return mapping[key];
  }

  /**
   * 获取常量值的详细描述
   * @param value 常量值
   * @returns 详细描述
   */
  public static getDescription(value: number): string {
    const descriptions: Record<number, string> = {
      [AppSidType.DYNAMIC]: '动态预设',
      [AppSidType.GROUP_BATTLE]: '对线群预设',
      [AppSidType.GROUP_CAMP]: '大本营预设',
      [AppSidType.USER_CUSTOM]: '用户自定义',
      [AppSidType.USER_PRESET]: '用户预设',
      [AppSidType.WORLD]: '世界预设',
    };

    return descriptions[value] || `未知常量值: ${value}`;
  }
}

/**
 * 用户token类型 - 其他
 */
export class AppUserTokenType {
  /** 访问App */
  public static readonly ACCESS = 4;
  /** 绑定新手机号 */
  public static readonly BIND_NEW_PHONE = 3;
  /** 登录 */
  public static readonly LOGIN = 1;
  /** 注册 */
  public static readonly REGISTER = 0;
  /** 解绑旧手机号 */
  public static readonly UNBIND_OLD_PHONE = 2;

  // 私有构造函数，防止实例化
  private constructor() {
    throw new Error('常量类不应被实例化');
  }

  /**
   * 获取所有用户token类型常量值
   * @returns 所有常量值的数组
   */
  public static getAllValues(): readonly number[] {
    return [
      AppUserTokenType.ACCESS,
      AppUserTokenType.BIND_NEW_PHONE,
      AppUserTokenType.LOGIN,
      AppUserTokenType.REGISTER,
      AppUserTokenType.UNBIND_OLD_PHONE,
    ] as const;
  }

  /**
   * 获取所有用户token类型常量键名
   * @returns 所有常量键名的数组
   */
  public static getAllKeys(): readonly string[] {
    return [
      'ACCESS',
      'BIND_NEW_PHONE',
      'LOGIN',
      'REGISTER',
      'UNBIND_OLD_PHONE',
    ] as const;
  }

  /**
   * 获取键值对映射
   * @returns 键值对映射对象
   */
  public static getKeyValuePairs(): Readonly<Record<string, number>> {
    return {
      ACCESS: AppUserTokenType.ACCESS,
      BIND_NEW_PHONE: AppUserTokenType.BIND_NEW_PHONE,
      LOGIN: AppUserTokenType.LOGIN,
      REGISTER: AppUserTokenType.REGISTER,
      UNBIND_OLD_PHONE: AppUserTokenType.UNBIND_OLD_PHONE,
    } as const;
  }

  /**
   * 根据值和语言格式化用户token类型的标签
   * @param value 常量值
   * @param lang 语言代码 ('zh', 'en', 'ja')
   * @returns 格式化后的标签，找不到时返回 'Unknown(value)'
   */
  public static formatValue(value: number, lang: 'zh' | 'en' | 'ja' = 'zh'): string {
    const labels: Record<string, Record<number, string>> = {
      zh: {
        [AppUserTokenType.ACCESS]: '访问App',
        [AppUserTokenType.BIND_NEW_PHONE]: '绑定新手机号',
        [AppUserTokenType.LOGIN]: '登录',
        [AppUserTokenType.REGISTER]: '注册',
        [AppUserTokenType.UNBIND_OLD_PHONE]: '解绑旧手机号',
      },
      en: {
        [AppUserTokenType.ACCESS]: 'Access App',
        [AppUserTokenType.BIND_NEW_PHONE]: 'Bind New Phone',
        [AppUserTokenType.LOGIN]: 'Login',
        [AppUserTokenType.REGISTER]: 'Register',
        [AppUserTokenType.UNBIND_OLD_PHONE]: 'Unbind Old Phone',
      },
      ja: {
        [AppUserTokenType.ACCESS]: 'アプリアクセス',
        [AppUserTokenType.BIND_NEW_PHONE]: '新しい電話番号バインド',
        [AppUserTokenType.LOGIN]: 'ログイン',
        [AppUserTokenType.REGISTER]: '登録',
        [AppUserTokenType.UNBIND_OLD_PHONE]: '古い電話番号アンバインド',
      },
    };

    const langLabels = labels[lang];
    if (langLabels && langLabels[value] !== undefined) {
      return langLabels[value];
    }

    // 默认返回英文
    const enLabels = labels.en;
    if (enLabels && enLabels[value] !== undefined) {
      return enLabels[value];
    }

    return `Unknown(${value})`;
  }

  /**
   * 验证值是否为有效的用户token类型常量
   * @param value 要验证的值
   * @returns 是否为有效常量
   */
  public static isValid(value: number): value is typeof AppUserTokenType.ACCESS | typeof AppUserTokenType.BIND_NEW_PHONE | typeof AppUserTokenType.LOGIN | typeof AppUserTokenType.REGISTER | typeof AppUserTokenType.UNBIND_OLD_PHONE {
    return AppUserTokenType.getAllValues().includes(value);
  }

  /**
   * 从字符串键名获取用户token类型常量值
   * @param key 常量键名
   * @returns 常量值，找不到时返回 undefined
   */
  public static fromString(key: string): number | undefined {
    const mapping = AppUserTokenType.getKeyValuePairs();
    return mapping[key];
  }

  /**
   * 获取常量值的详细描述
   * @param value 常量值
   * @returns 详细描述
   */
  public static getDescription(value: number): string {
    const descriptions: Record<number, string> = {
      [AppUserTokenType.ACCESS]: '访问App',
      [AppUserTokenType.BIND_NEW_PHONE]: '绑定新手机号',
      [AppUserTokenType.LOGIN]: '登录',
      [AppUserTokenType.REGISTER]: '注册',
      [AppUserTokenType.UNBIND_OLD_PHONE]: '解绑旧手机号',
    };

    return descriptions[value] || `未知常量值: ${value}`;
  }
}

// 类型定义 - 使用后缀避免与其他常量类型重名
export type AppUserTokenTypeValue = typeof AppUserTokenType.ACCESS | 
                                   typeof AppUserTokenType.BIND_NEW_PHONE | 
                                   typeof AppUserTokenType.LOGIN | 
                                   typeof AppUserTokenType.REGISTER | 
                                   typeof AppUserTokenType.UNBIND_OLD_PHONE;

export type AppUserTokenTypeKey = 'ACCESS' | 'BIND_NEW_PHONE' | 'LOGIN' | 'REGISTER' | 'UNBIND_OLD_PHONE';

```
生成的 index.ts 文件（TypeScript 包导出）：
```typescript
export * from './app';
export * from './user';
export * from './group';
export * from './admin';
```



## 项目结构

```bash
cons-coder/
├── main.go              # 主程序入口
├── README.md           # 项目文档
├── go.mod              # Go模块配置
├── data/               # XML配置文件目录
│   ├── app.xml         # 应用相关常量
│   ├── user.xml        # 用户相关常量
│   ├── group.xml       # 群组相关常量
│   └── admin.xml       # 管理相关常量
├── output/             # 生成代码输出目录
│   ├── python/         # Python代码
│   │   ├── __init__.py # 包初始化文件 (自动生成)
│   │   ├── app.py      # 对应 app.xml
│   │   ├── user.py     # 对应 user.xml  
│   │   ├── group.py    # 对应 group.xml
│   │   └── admin.py    # 对应 admin.xml
│   ├── go/             # Go代码 (默认包名: cons)
│   │   ├── app.go      # 对应 app.xml
│   │   ├── user.go     # 对应 user.xml
│   │   ├── group.go    # 对应 group.xml
│   │   └── admin.go    # 对应 admin.xml
│   ├── java/           # Java代码
│   │   ├── App.java    # 对应 app.xml
│   │   ├── User.java   # 对应 user.xml
│   │   ├── Group.java  # 对应 group.xml
│   │   └── Admin.java  # 对应 admin.xml
│   ├── swift/          # Swift代码
│   │   ├── App.swift   # 对应 app.xml
│   │   ├── User.swift  # 对应 user.xml
│   │   ├── Group.swift # 对应 group.xml
│   │   └── Admin.swift # 对应 admin.xml
│   └── kotlin/         # Kotlin代码
│       ├── App.kt      # 对应 app.xml
│       ├── User.kt     # 对应 user.xml
│       ├── Group.kt    # 对应 group.xml
│       └── Admin.kt    # 对应 admin.xml
│   ├── typescript/    # Typescript代码
│   │   ├── index.ts    # 对应 index.xml
│   │   ├── app.ts      # 对应 app.xml
│   │   ├── user.ts     # 对应 user.xml
│   │   ├── group.ts    # 对应 group.xml
│   │   └── admin.ts    # 对应 admin.xml
│   └── javascript/     # Javascript代码
│       ├── index.js    # 对应 index.xml
│       ├── app.js      # 对应 app.xml
│       ├── user.js     # 对应 user.xml
│       ├── group.js    # 对应 group.xml
│       └── admin.js    # 对应 admin.xml
└── templates/          # 代码生成模板（文件名是语言名）
    ├── python.tmpl      # 对应python
    ├── go.tmpl          # 对应go
    ├── java.tmpl        # 对应java
    ├── swift.tmpl       # 对应swift
    ├── kotlin.tmpl      # 对应kotlin
    ├── typescript.tmpl  # 对应typescript
    └── javascript.tmpl  # 对应javascript
```

## 代码特性

### 生成的代码包含以下功能

1. **文件头部注释**: 包含XML文件信息、修改时间、生成时间等元数据
2. **基础常量定义**: 清晰的常量声明和注释
3. **获取方法**: 
   - `getAllValues()` - 获取所有常量值
   - `getAllKeys()` - 获取所有常量键名  
   - `getKeyValuePairs()` - 获取键值对映射
4. **工具方法**:
   - `format()` - 多语言标签格式化
   - `isValid()` - 常量值验证
   - `fromString()` - 从字符串键名获取值
   - `getDescription()` - 获取详细描述
5. **类型安全**: 使用强类型和适当的访问控制
6. **文档完整**: 包含完整的注释和文档
7. **Python特殊支持**: 自动生成 `__init__.py` 导入所有常量类
8. **Typescript特殊支持**: 自动生成 `index.ts` 导入所有常量类
9. **Javascript特殊支持**: 自动生成 `index.js` 导入所有常量类

## 生成示例

假设 `data/` 目录下有以下XML文件：

**文件结构：**
```
data/
├── app.xml      # 应用相关常量
├── user.xml     # 用户相关常量  
└── group.xml    # 群组相关常量
```

**生成Python代码：**
```bash
cons-coder --dir ./data --output ./python-codes --lang python
```

**输出结果：**
```
python-codes/
├── __init__.py          # 自动生成，导入所有常量
├── app.py              # AppSidType, AppUserTokenType 等类
├── user.py             # SuspendHandleSource, UserStatus 等类
└── group.py            # GroupType, GroupStatus 等类
```

**生成Go代码：**
```bash
cons-coder --dir ./data --output ./go-codes --lang go --package constants
```

**输出结果：**
```  
go-codes/
├── app.go              # package constants; AppSidType, AppUserTokenType 等
├── user.go             # package constants; SuspendHandleSource, UserStatus 等
└── group.go            # package constants; GroupType, GroupStatus 等
```

## 最佳实践

1. **XML文件组织**: 
   - 按业务模块组织XML文件 (如：`user.xml`, `order.xml`, `payment.xml`)
   - 每个XML文件包含相关的常量组，避免单一文件过大
   - 使用有意义的文件名，因为它们将成为生成代码的文件名

2. **命名规范**: 
   - **XML文件名**: 使用小写字母和下划线 (如：`user_profile.xml`)
   - **XML元素名**: 使用小写字母和下划线
   - **生成代码**: 自动遵循各语言命名规范

3. **版本管理**: 
   - 将XML配置文件纳入版本控制
   - 在XML文件中添加版本信息或修改说明
   - 定期更新和维护常量定义

4. **持续集成**: 
   - 在构建流程中集成代码生成步骤
   - 设置自动化检查确保XML格式正确
   - 生成代码后运行单元测试验证

5. **多语言支持**:
   - 在XML中为不同语言提供对应的标签
   - 统一管理多语言标签，避免遗漏
   - 考虑使用国际化框架结合生成的常量

6. **Go语言特殊配置**:
   - 使用 `--package` 参数指定合适的包名
   - 建议包名与项目结构保持一致
   - 考虑生成到不同包中以避免循环依赖

## 使用技巧

### 批处理生成
```bash
# 一次性生成多种语言
./scripts/generate-all.sh

# generate-all.sh 内容示例:
#!/bin/bash
cons-coder --dir ./data --output ./output/python --lang python
cons-coder --dir ./data --output ./output/go --lang go --package constants  
cons-coder --dir ./data --output ./output/java --lang java
cons-coder --dir ./data --output ./output/swift --lang swift
cons-coder --dir ./data --output ./output/kotlin --lang kotlin
```

### Python包使用示例
```python
# 导入生成的常量
from constants import AppSidType, UserStatus

# 使用常量
if user_type == AppSidType.USER_PRESET:
    print("这是用户预设类型")

# 验证常量值
if AppSidType.is_valid(user_input):
    label = AppSidType.format_value(user_input, 'zh')
    print(f"用户输入的标签是: {label}")
```

### Go包使用示例
```go
package main

import (
    "fmt"
    "your-project/constants"
)

func main() {
    // 使用常量
    if userType == constants.AppSidType.UserPreset {
        fmt.Println("这是用户预设类型")
    }
    
    // 验证和格式化
    if constants.AppSidType.IsValid(userInput) {
        label := constants.AppSidType.Format(userInput, "zh")
        fmt.Printf("用户输入的标签是: %s\n", label)
    }
}
```

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。