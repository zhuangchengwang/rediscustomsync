package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/guoruibiao/rediscustomsync/models/dao"
	"github.com/guoruibiao/rediscustomsync/models/dao/redisdao"
	"github.com/guoruibiao/rediscustomsync/models/service"
	"os"
	"strings"
)

var appConfig dao.AppConfig
func Init(configPath string) {
	_, err := toml.DecodeFile(configPath, &appConfig)
	if err != nil {
		fmt.Println("加载配置文件失败， 错误：", err.Error())
		os.Exit(0)
	}
	// 初始化读写的 redisdao 配置
	redisdao.Init(appConfig)
}

func main() {
	helpInfo := `
    ____           ___         ______           __                     _____                 
   / __ \___  ____/ (_)____   / ____/_  _______/ /_____  ____ ___     / ___/__  ______  _____
  / /_/ / _ \/ __  / / ___/  / /   / / / / ___/ __/ __ \/ __  __ \    \__ \/ / / / __ \/ ___/
 / _, _/  __/ /_/ / (__  )  / /___/ /_/ (__  ) /_/ /_/ / / / / / /   ___/ / /_/ / / / / /__  
/_/ |_|\___/\__,_/_/____/   \____/\__,_/____/\__/\____/_/ /_/ /_/   /____/\__, /_/ /_/\___/  
                                                                         /____/
    Redis Custom Synchronization, 自定义 Redis 同步工具，支持功能：
      - 无需全量同步的场景，仅同步 conf/app.toml 中的 keysfile 的指定 key（一行一个）
      - 跨机房同步问题，机器 A 可访问机器 B，机器 A 可访问机器 C，但是机器 B 不能访问机器 C，却需要把机器 C上部分 key 同步到机器 B
      - 模式匹配，可支持模式串的 key 同步，优先级低于 keysfile 中指定的 key
    
    注意：
      - keysfile 路径可以在 conf/app.toml 中配置，建议使用绝对路径; 建议末尾预留一个换行
      - conf 文件夹需放置到工具运行的同级目录 或者通过 -c 参数指定绝对路径
      - 无需模式匹配时，patterns 置空即可

    用法：
        rediscustomsync -config ./conf/app.toml

                                                   Author: guoruibiao
`
	fmt.Println(helpInfo)
	var configPath = flag.String("config", "", "配置文件路径")
	flag.Parse()
	if *configPath == "" {
		fmt.Println("conf/app.toml 路径为空")
		os.Exit(0)
	}

	// 初始化配置
	Init(*configPath)

	fmt.Println("sync...")
	// 开始模式匹配 pattern 同步
	success, failed, err := service.TransferPatterns(appConfig)
	if err != nil {
		fmt.Println("同步失败， 错误：", err.Error())
		os.Exit(0)
	}

	fmt.Println("同步结果：")
	fmt.Printf("同步成功: %d 个，详细：\n %s\n", len(success), strings.Join(success, "\n"))
	fmt.Printf("同步失败: %d 个，详细：\n %s\n", len(failed), strings.Join(failed, "\n"))

	// 开始同步 keysfile
	success, failed, err = service.TransferKeysfile(appConfig)
	if err != nil {
		fmt.Println("同步失败， 错误：", err.Error())
		os.Exit(0)
	}

	fmt.Println("同步结果：")
	fmt.Printf("同步成功: %d 个，详细：\n %s\n", len(success), strings.Join(success, "\n"))
	fmt.Printf("同步失败: %d 个，详细：\n %s\n", len(failed), strings.Join(failed, "\n"))

}