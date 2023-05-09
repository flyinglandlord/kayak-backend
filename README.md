# kayak-backend

[![Docker Image CI](https://github.com/FreelyCoding/kayak-backend/actions/workflows/docker-image.yml/badge.svg?branch=master)](https://github.com/FreelyCoding/kayak-backend/actions/workflows/docker-image.yml)

学舟刷题软件配套的后端，使用Golang编写

根据`config_demo.yaml`填写相应数据库、Redis等信息，然后重命名为`config.yaml`即可运行

也可以使用Dockerfile构建一个镜像，然后使用`docker-compose`运行