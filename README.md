# hardware-agent
hardware agent 

基于小米运维开源的open-falcon，硬件监控专用agent。

##采集的metric列表：
  custom.agent.alive
  custom.data
  
##配置说明
配置文件请参照cfg.example.json，修改该文件名为cfg.json，将该文件里的IP换成实际使用的IP。
```
{
    "debug": true,
    "hostname": "zsf-test-port",                           // endpoint
    "plugin": {
        "enabled": false,
        "dir": "./plugin",
        "git": "https://github.com/open-falcon/plugin.git",
        "logs": "./logs"
    },
    "heartbeat": {
        "enabled": true,
        "addr": "127.0.0.1:6030",
        "interval": 60,
        "timeout": 1000
    },
    "transfer": {
        "enabled": true,
        "addrs": [
            "127.0.0.1:8433",
            "127.0.0.1:8433"
        ],
        "interval": 60,
        "timeout": 1000
    },
    "http": {
        "enabled": false,
        "listen": ":1988",
        "backdoor": false
    },
    "filepath": [
    "D:/myspace/smarteye-windows-agent/hello.exe"
    ]，                                            // 指定需要执行的脚本指令
   "exectimeout": 2000                                // 执行单个脚本文件超时时间ms
}
```

