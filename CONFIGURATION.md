# 配置指南

本文档介绍如何配置 Interlace Money Go SDK 的认证信息和环境设置。

## 📋 目录

1. [获取 API 凭证](#获取-api-凭证)
2. [环境变量配置](#环境变量配置)
3. [代码中配置](#代码中配置)
4. [最佳实践](#最佳实践)
5. [故障排查](#故障排查)

---

## 🔑 获取 API 凭证

### 步骤 1: 注册开发者账号

访问 [Interlace 开发者平台](https://developer.interlace.money) 并注册账号。

### 步骤 2: 创建应用

1. 登录开发者控制台
2. 创建新应用
3. 获取您的 **Client ID** 和 **Client Secret**

### 步骤 3: 选择环境

- **Sandbox（沙盒）**: 用于开发和测试，不涉及真实资金
  - 基础 URL: `https://api-sandbox.interlace.money`
  - 推荐用于开发阶段

- **Production（生产）**: 用于正式环境，涉及真实交易
  - 基础 URL: `https://api.interlace.money`
  - 仅在通过测试后使用

---

## 🌍 环境变量配置

### 方法 1: 使用 .env 文件（推荐）

1. **复制示例文件**:
   ```bash
   cp .env.example .env
   ```

2. **编辑 .env 文件**:
   ```bash
   # 必填
   INTERLACE_CLIENT_ID=your-actual-client-id
   INTERLACE_CLIENT_SECRET=your-actual-client-secret
   
   # 可选
   INTERLACE_ENVIRONMENT=sandbox
   TEST_ACCOUNT_ID=your-test-account-id
   ```

3. **在代码中加载环境变量**:

   安装环境变量加载库:
   ```bash
   go get github.com/joho/godotenv
   ```

   在代码中使用:
   ```go
   package main
   
   import (
       "log"
       "os"
       
       "github.com/joho/godotenv"
       interlace "github.com/difyz9/interlace-go-sdk/pkg"
   )
   
   func main() {
       // 加载 .env 文件
       err := godotenv.Load()
       if err != nil {
           log.Fatal("Error loading .env file")
       }
       
       // 从环境变量读取配置
       clientID := os.Getenv("INTERLACE_CLIENT_ID")
       clientSecret := os.Getenv("INTERLACE_CLIENT_SECRET")
       
       // 创建配置
       config := interlace.DefaultConfig()
       config.ClientID = clientID
       
       // 创建客户端
       client := interlace.NewClient(config)
       
       // 使用客户端...
   }
   ```

### 方法 2: 系统环境变量

在终端中设置环境变量:

**Linux/macOS**:
```bash
export INTERLACE_CLIENT_ID="your-client-id"
export INTERLACE_CLIENT_SECRET="your-client-secret"
export INTERLACE_ENVIRONMENT="sandbox"
```

**Windows (PowerShell)**:
```powershell
$env:INTERLACE_CLIENT_ID="your-client-id"
$env:INTERLACE_CLIENT_SECRET="your-client-secret"
$env:INTERLACE_ENVIRONMENT="sandbox"
```

**Windows (CMD)**:
```cmd
set INTERLACE_CLIENT_ID=your-client-id
set INTERLACE_CLIENT_SECRET=your-client-secret
set INTERLACE_ENVIRONMENT=sandbox
```

### 方法 3: Docker 环境变量

在 `docker-compose.yml` 中:
```yaml
version: '3.8'
services:
  app:
    build: .
    environment:
      - INTERLACE_CLIENT_ID=${INTERLACE_CLIENT_ID}
      - INTERLACE_CLIENT_SECRET=${INTERLACE_CLIENT_SECRET}
      - INTERLACE_ENVIRONMENT=sandbox
    env_file:
      - .env
```

---

## 💻 代码中配置

### 配置方式 1: 使用默认配置

```go
package main

import (
    "context"
    "log"
    "os"
    
    interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
    clientID := os.Getenv("INTERLACE_CLIENT_ID")
    
    // 快速设置（使用默认沙盒配置）
    client, tokenData, err := interlace.QuickSetup(clientID, nil)
    if err != nil {
        log.Fatalf("Setup failed: %v", err)
    }
    
    log.Printf("Authenticated! Token expires in %d seconds", tokenData.ExpiresIn)
}
```

### 配置方式 2: 自定义配置

```go
package main

import (
    "context"
    "time"
    
    interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
    // 创建自定义配置
    config := &interlace.Config{
        BaseURL:   "https://api-sandbox.interlace.money",
        ClientID:  "your-client-id",
        UserAgent: "MyApp/1.0.0",
        Timeout:   60 * time.Second,
    }
    
    // 创建客户端
    client := interlace.NewClient(config)
    
    // 认证
    ctx := context.Background()
    tokenData, err := client.Authenticate(ctx, config.ClientID)
    if err != nil {
        // 处理错误
    }
    
    // 使用客户端...
}
```

### 配置方式 3: 环境切换

```go
package main

import (
    "os"
    
    interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func getConfig() *interlace.Config {
    env := os.Getenv("INTERLACE_ENVIRONMENT")
    
    var config *interlace.Config
    
    if env == "production" {
        // 生产环境配置
        config = interlace.ProductionConfig()
    } else {
        // 默认使用沙盒环境
        config = interlace.SandboxConfig()
    }
    
    // 设置 Client ID
    config.ClientID = os.Getenv("INTERLACE_CLIENT_ID")
    
    return config
}

func main() {
    config := getConfig()
    client := interlace.NewClient(config)
    
    // 使用客户端...
}
```

---

## ✅ 最佳实践

### 1. 安全性

- ✅ **永远不要在代码中硬编码凭证**
- ✅ **使用环境变量或密钥管理服务**
- ✅ **将 .env 文件添加到 .gitignore**
- ✅ **为不同环境使用不同的 Client ID**
- ✅ **定期轮换 API 密钥**

示例 `.gitignore`:
```
.env
.env.local
.env.*.local
```

### 2. 环境管理

- 开发环境使用 **Sandbox**
- 测试环境使用 **Sandbox** 或专用测试账户
- 生产环境使用 **Production** 配置

### 3. 错误处理

```go
func authenticateClient(clientID string) (*interlace.Client, error) {
    if clientID == "" {
        return nil, fmt.Errorf("INTERLACE_CLIENT_ID environment variable is not set")
    }
    
    client, _, err := interlace.QuickSetup(clientID, nil)
    if err != nil {
        return nil, fmt.Errorf("authentication failed: %w", err)
    }
    
    return client, nil
}
```

### 4. 配置验证

```go
func validateConfig(config *interlace.Config) error {
    if config.ClientID == "" {
        return fmt.Errorf("client ID is required")
    }
    if config.BaseURL == "" {
        return fmt.Errorf("base URL is required")
    }
    if config.Timeout == 0 {
        config.Timeout = 30 * time.Second // 设置默认超时
    }
    return nil
}
```

### 5. 测试账户管理

为测试创建专用账户：

```go
func setupTestEnvironment() (*interlace.Client, string, error) {
    clientID := os.Getenv("INTERLACE_CLIENT_ID")
    testAccountID := os.Getenv("TEST_ACCOUNT_ID")
    
    client, _, err := interlace.QuickSetup(clientID, nil)
    if err != nil {
        return nil, "", err
    }
    
    // 如果没有测试账户，创建一个
    if testAccountID == "" {
        ctx := context.Background()
        account, err := client.Account.RegisterGolangTest(ctx)
        if err != nil {
            return nil, "", err
        }
        testAccountID = account.ID
    }
    
    return client, testAccountID, nil
}
```

---

## 🔧 故障排查

### 问题 1: "client ID is required" 错误

**原因**: Client ID 未设置

**解决方案**:
1. 检查 `.env` 文件是否存在
2. 确认 `INTERLACE_CLIENT_ID` 已设置
3. 检查环境变量是否正确加载

```go
// 调试代码
clientID := os.Getenv("INTERLACE_CLIENT_ID")
if clientID == "" {
    log.Fatal("INTERLACE_CLIENT_ID is not set")
}
log.Printf("Using client ID: %s...", clientID[:8])
```

### 问题 2: 认证失败

**原因**: Client ID 不正确或已过期

**解决方案**:
1. 验证 Client ID 是否正确
2. 检查是否使用了正确的环境（Sandbox vs Production）
3. 确认账户状态正常

### 问题 3: 网络连接问题

**原因**: 无法连接到 API 服务器

**解决方案**:
1. 检查网络连接
2. 验证 BaseURL 是否正确
3. 检查防火墙设置
4. 增加超时时间:

```go
config := interlace.DefaultConfig()
config.Timeout = 60 * time.Second // 增加到 60 秒
```

### 问题 4: Token 过期

**原因**: Access Token 已过期

**解决方案**:
实现 token 刷新逻辑:

```go
func refreshTokenIfNeeded(client *interlace.Client, tokenData *interlace.OAuthTokenData) error {
    // 检查 token 是否即将过期（例如，剩余时间少于 5 分钟）
    if time.Now().Unix() + 300 > tokenData.ExpiresIn {
        ctx := context.Background()
        newToken, err := client.OAuth.RefreshAccessToken(ctx, tokenData.RefreshToken)
        if err != nil {
            return err
        }
        
        client.SetAccessToken(newToken.AccessToken)
        tokenData = newToken
    }
    
    return nil
}
```

---

## 📚 相关文档

- [Interlace API 文档](https://developer.interlace.money)
- [SDK README](./README.md)
- [示例代码](./examples/)

---

## 🆘 获取帮助

如果您遇到配置问题:

1. 查看 [GitHub Issues](https://github.com/difyz9/interlace-go-sdk/issues)
2. 阅读 [Interlace API 文档](https://developer.interlace.money)
3. 联系技术支持

