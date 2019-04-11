# chronusProxy

TCP UDP UART Port Forward

通过config.json中的配置可做到多协议相互转发
例如:

    "ProxyHub": {
        "TCP": {
            ":6666": "192.168.1.100:6666",
            ":7777": "192.168.1.100:7777"
        },
        "UDP": {},
        "TCP2UDP": {},
        "UDP2TCP": {},
        "UART2UDP": {
            ":55555": {
                "PortName": "/dev/ttyAMA0",
                "BaudRate": 115200,
                "DataBits": 8,
                "StopBits": 1,
                "MinimumReadSize": 4
            }
        }
    }

即为

    -> 主机192.168.1.100的6666端口转发到本机的6666端口
    -> 主机192.168.1.100的7777端口转发到本机的7777端口
    -> 主机RS232或者RS485串口转发到本机UDP协议55555端口
