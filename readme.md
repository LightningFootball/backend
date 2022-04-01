<div align="center">

# EduOJ

</div>

# 项目简介

该项目是一个基于[EduOJ](https://github.com/EduOJ/backend) 的分支项目，为本人的毕设项目。 该项目目的是在原项目基础上添加附加功能和改进。
更多关于项目的开发信息，如代码风格、成员、权限、存储信息等，请参考原项目，恕不赘述。

# 协议

原EduOJ项目使用GPL v3开源协议，因而本项目也使用GPL v3开源协议。

[GNU AFFERO GENERAL PUBLIC LICENSE Version 3](./license.md).

# 部署方式

该项目需要与前端项目配合使用，因此需要部署 [Frontend](https://github.com/LightningFootball/frontend) 部署。

当前后端项目的推荐部署方式为：

1. docker部署存储相关容器
    1. minio/minio
    2. postgres
    3. redis
2. 拉取后端
    ```
    git clone https://github.com/EduOJ/backend.git
    cd backend
    go build .
   ```
3. 修改配置文件
    ```
    cp config.yml.example config.yml
    nano config.yml
    ```
4. 启动后端
    ```
    ./backend serve
    ```