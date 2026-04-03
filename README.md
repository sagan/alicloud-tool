# 薅阿里云羊毛辅助工具

参考资料: [使用阿里云CDT 享受每月200G/4元高速流量](https://www.nodeseek.com/post-572284-1)。

这个程序是上面帖子里 Python 脚本的 Go 移植版本。

需要在阿里云里创建 API 用户 ( https://ram.console.alibabacloud.com/users/create ) 并赋予以下权限： AliyunECSFullAccess, AliyunCDTFullAccess 。

## Usage

```
#alicloud-tool -h
A command line tool specifically designed to start/stop ECS instances based on CDT traffic thresholds.

Note: All settings can also be configured using environment variables.
The corresponding environment variables are prefixed with "ALIBABA_CLOUD_"
and use underscores instead of hyphens (e.g., ALIBABA_CLOUD_ACCESS_KEY_ID).

Usage:
  alicloud-tool [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  ecs-start   Start ECS instance manually
  ecs-stop    Stop ECS instance manually
  ecs-watch   Watch CDT traffic and start/stop ECS
  help        Help about any command

Flags:
      --access-key-id string       Alibaba Cloud Access Key ID (Env: ALIBABA_CLOUD_ACCESS_KEY_ID)   
      --access-key-secret string   Alibaba Cloud Access Key Secret (Env: ALIBABA_CLOUD_ACCESS_KEY_SECRET)
  -h, --help                       help for alicloud-tool
      --instance-id string         Alibaba Cloud ECS Instance ID (Env: ALIBABA_CLOUD_INSTANCE_ID)   
      --region string              Alibaba Cloud Region (Env: ALIBABA_CLOUD_REGION) (default "cn-hongkong")

Use "alicloud-tool [command] --help" for more information about a command.
```
