## 1. 文件类型检测

- [x] 1.1 在 `main.go` 中添加文件扩展名检测函数，支持不区分大小写的 `.sln` 和 `.vcxproj` 匹配
- [x] 1.2 添加不支持文件类型的错误处理，返回清晰的错误消息说明仅支持 .sln 和 .vcxproj

## 2. .vcxproj 文件处理流程

- [x] 2.1 在 `main.go` 中添加文件类型判断分支逻辑
- [x] 2.2 实现 .vcxproj 分支：调用 `sln.NewProject()` 直接解析项目文件
- [x] 2.3 获取项目文件的绝对路径并提取所在目录作为 `SolutionDir`
- [x] 2.4 手动构造 `sln.Sln` 实例，包含单个项目和正确的 `SolutionDir`
- [x] 2.5 验证 `sln.Sln` 字段可见性，确保可以在 `main.go` 中直接构造

## 3. .sln 文件处理流程保持不变

- [x] 3.1 确保 .sln 分支仍然使用 `sln.NewSln()` 的现有逻辑
- [x] 3.2 验证两种输入类型共用相同的后续处理流程（`CompileCommandsJson()` 和输出逻辑）

## 4. 错误处理完善

- [x] 4.1 添加 .vcxproj 文件不存在时的错误处理
- [x] 4.2 添加 .vcxproj 文件格式错误（无效 XML）时的错误处理
- [x] 4.3 确保配置不存在时的错误消息与 .sln 处理保持一致

## 5. 测试与验证

- [x] 5.1 测试：使用有效的 .vcxproj 文件生成 compile_commands.json
- [x] 5.2 测试：验证生成的编译命令包含正确的 directory、file 和 command 字段
- [x] 5.3 测试：验证 `$(SolutionDir)` 变量在 .vcxproj 输入时被正确替换为项目目录
- [x] 5.4 测试：大写和混合大小写的 .vcxproj 扩展名能被正确识别
- [x] 5.5 测试：不支持的文件扩展名（如 .txt）返回清晰的错误消息
- [x] 5.6 测试：.vcxproj 文件不存在时返回适当的错误
- [x] 5.7 测试：现有 .sln 文件处理逻辑未受影响（回归测试）
- [x] 5.8 测试：同一项目分别通过 .sln 和 .vcxproj 处理，生成的编译命令应一致

## 6. 文档更新

- [x] 6.1 更新 README.md，说明 `-s` 参数现在支持 .vcxproj 文件
- [x] 6.2 在 README.md 中添加使用 .vcxproj 文件的示例
- [x] 6.3 更新 usage 函数的帮助文本，说明支持的文件类型
