先看实现效果：![6841a98a23b09c8bf1315be10fc6580c](https://github.com/user-attachments/assets/46e136cb-1b72-4b0e-ab82-c722dd701798)
<img width="631" height="1282" alt="image" src="https://github.com/user-attachments/assets/39342ab1-7f66-4742-8b9a-3d73a25837fd" />


本软件用于以企业微信机器人的身份每天自动向企业微信群中发送和AI提示词有关的文章，需要配置的选项有：
1.企业微信机器人地址
2.deepseek的api
<img width="1010" height="535" alt="image" src="https://github.com/user-attachments/assets/8521efb5-2dc1-4e07-ad70-424599219449" />
前两条更改地址在/robotv1/config/config.yaml
3.文件服务器ip地址/域名（显示每天生成出文章的网页连接）
/robotv1/pkg/wechat/robot.go
<img width="1378" height="774" alt="image" src="https://github.com/user-attachments/assets/8cb111a4-8153-4124-8170-c5c3577538f7" />
文件服务器端口设置（非必须）
/robotv1/main.go
<img width="1378" height="774" alt="image" src="https://github.com/user-attachments/assets/4c91e575-bf46-46c7-a735-ea60a94ea432" />



实现项目目录的articles和下载链接里面的http://网址:8080/articles同步的关键：lsyncd，下面是配置文件
<img width="1188" height="678" alt="image" src="https://github.com/user-attachments/assets/477998a9-ea3f-4e25-b85f-d7cbdcc82233" />

保证点击下载链接触发下载动作，而不是触发预览（预览会乱码）
nginx配置：
location ~* \.md$ {
    # 强制所有 .md 文件作为附件下载
    add_header Content-Disposition "attachment";
    
    # 可选：防止浏览器尝试渲染文本
    add_header Content-Type "application/octet-stream";
}

apache（httpd）配置：
<FilesMatch "\.md$">
    # 强制下载
    Header set Content-Disposition "attachment"

    # 可选：禁用文本渲染
    Header set Content-Type "application/octet-stream"
</FilesMatch>
