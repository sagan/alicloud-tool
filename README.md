# 薅阿里云羊毛辅助工具

参考资料: [使用阿里云CDT 享受每月200G/4元高速流量](https://www.nodeseek.com/post-572284-1)。

这个程序是上面帖子里 Python 脚本的 Go 移植版本。仅支持阿里云国际版( https://www.alibabacloud.com/ )。

需要在阿里云里创建 API 用户 ( https://ram.console.alibabacloud.com/users/create ) 并赋予以下权限： AliyunECSFullAccess, AliyunCDTReadOnlyAccess, AliyunBSSReadOnlyAccess 。

然后设置环境变量：

- ALIBABA_CLOUD_ACCESS_KEY_ID & ALIBABA_CLOUD_ACCESS_KEY_SECRET : 创建的 API 用户的 id 和 secret。
- ALIBABA_CLOUD_INSTANCE_ID : 阿里云 ECS 的实例 id。

使用方法：

- `alicloud-tool ecs-watch` : 检测当前 ECS 实例以及 CDT 流量状态。如果本月 200 GB CDT 流量没用完则自动将 ECS 实例开机；否则将其关机。
- `alicloud-tool ecs-start`, `alicloud-tool ecs-stop` : 手动开机 / 关机 ECS 实例。
- `alicloud-tool status` : 显示当前阿里云账户状态，包括 CDT 流量用量以及账单信息。

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
  status      Display current Alibaba Cloud account status info (CDT traffic, Billing)

Flags:
      --access-key-id string       Alibaba Cloud Access Key ID (Env: ALIBABA_CLOUD_ACCESS_KEY_ID)
      --access-key-secret string   Alibaba Cloud Access Key Secret (Env: ALIBABA_CLOUD_ACCESS_KEY_SECRET)
  -h, --help                       help for alicloud-tool
      --instance-id string         Alibaba Cloud ECS Instance ID (Env: ALIBABA_CLOUD_INSTANCE_ID)
      --region string              Alibaba Cloud Region (Env: ALIBABA_CLOUD_REGION) (default "cn-hongkong")
  -v, --version                    version for alicloud-tool

Use "alicloud-tool [command] --help" for more information about a command.
```
