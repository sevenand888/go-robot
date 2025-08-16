#!/bin/bash
##########################################
# File Name:rebotv1.sh
# Version:V1.0
# Author:zhangpeng
# 0rganization:linuxjk.cn
# Desc:定时发送AI提示词工程有关文章的脚本
###########################################

# 在脚本开头添加环境变量打印
#echo "===== DEBUG START =====" >> /robotv1/logs/cron.log
#echo "PATH: $PATH" >> /robotv1/logs/cron.log
#echo "PWD: $(pwd)" >> /robotv1/logs/cron.log
#env >> /robotv1/logs/cron.log
#vars
file=/robotv1/main.go

# 设置工作目录
cd /robotv1 || exit 1

# 执行程序并记录日志
{
  echo "===== START: $(date) ====="
#go run $file
./robotv1
  echo "===== END: $(date) Exit Code: $? ====="
} >> /robotv1/logs/cron.log 2>&1
