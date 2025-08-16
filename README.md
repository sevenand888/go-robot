<img width="1528" height="796" alt="image" src="https://github.com/user-attachments/assets/dd900df8-a02f-41bc-8df7-085dfc678061" /><img width="1528" height="796" alt="image" src="https://github.com/user-attachments/assets/6a6a3de6-53b0-4210-a288-0cf90f136ecd" />本软件用于以企业微信机器人的身份每天自动向企业微信群中发送和AI提示词有关的文章，需要配置的选项有：
1.企业微信机器人地址
2.deepseek的api
前两条更改地址在/robotv1/config/config.yaml
3.文件服务器ip地址/域名（显示每天生成出文章的网页连接）
/robotv1/pkg/wechat/robot.go
文件服务器端口设置（非必须）
/robotv1/main.go
