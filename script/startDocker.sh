#!/bin/bash

config_file="./config.yaml"

# 检查配置文件是否存在
if [ -f $config_file ]; then
    # 读取YAML文件中的配置项
    userName=$(yq '.user.userName' $config_file)
    cudaVersion=$(yq '.user.cudaVersion' $config_file)

    # 使用读取到的配置项
    echo "当前人员: $userName"
    echo "cuda版本: $cudaVersion"
else
    echo "配置文件不存在或无法读取。"
fi


# 检查数据卷文件夹是否存在
folder_path="/data1/$userName"
if [ ! -d "$folder_path" ]; then
    # 如果文件夹不存在，则创建它
    mkdir -p "$folder_path"
    echo "文件夹已创建：$folder_path"
else
    echo "文件夹已存在：$folder_path"
fi

# 启动docker容器
docker run -it -d -P --name $userName --gpus 'device=0' -v $folder_path:/data 921image:latest