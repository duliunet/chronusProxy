# chronusProxy
TCP UDP UART Port Forward
通过config.json中的配置可做到多协议相互转发
例如:
"Proxy": {
    "ProxyHub": {
        "TCP": {":6666":"192.168.1.100:6666"},  #即为主机192.168.1.100的6666端口转发到本机的6666端口，以下同理只不过协议不同
        "UDP": {},
        "TCP2UDP": {},
        "UDP2TCP": {},
        "UART2UDP": {}
    }
}