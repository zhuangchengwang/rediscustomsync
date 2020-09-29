# rediscustomsync

作者将其定义为一个脚本工具，用于机器之间特定 Redis 的 key 的同步。
解决的问题有：

- 全量同步问题
- 机器 A 可以访问机器 B，机器 A 可以访问机器 C，但是机器 B 无法访问机器 C


支持功能
 - `keysfile` 一行一个，全部同步
 - `pattern` 模式匹配同步