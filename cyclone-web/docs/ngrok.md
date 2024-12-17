## 远程联调工具 ngrok

**场景**：现在开发团队经常处于异地的情况，这个时候相互联调挺麻烦的，通过 ngrok 可以快速解决这个问题！

**原来的解决方案：** 
1. 在杭州启动一个测试服务，然后直接连接杭州的
2. 更新代码，自己本地下载代码部署开发

`缺点:如果遇到联调开发，前后端都在不断的调试修改代码的时候，上面的方案会极大浪费时间`

**Ngrok**：

1.  https://ngrok.com/： 可以将 NAT 后面的本机 localhost 服务暴露给外部访问
2.  例子：后端开发本机启动 java 在  8080 端口， 在命令行运行  ngrok 8080

```
Session Status                online
Version                       2.1.18
Region                        United States (us)
Web Interface                 http://127.0.0.1:4040
Forwarding                    http://6f324930.ngrok.io -> localhost:8080
Forwarding                    https://6f324930.ngrok.io -> localhost:8080

Connections                   ttl     opn     rt1     rt5     p50     p90
                              0       0       0.00    0.00    0.00    0.00
```

3. 外部通过访问 http://6f324930.ngrok.io  就能直接链接开发机
