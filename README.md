# oyishyi-docker

# 0. motivation
This is a learning note of [《自己动手写 Docker》](https://github.com/xianlubird/mydocker).

# 1. Blogs
1. [使用 GoLang 从零开始写一个 Docker（概念篇）-- 《自己动手写 Docker》读书笔记](https://juejin.cn/post/6971335828060504094)
2. [使用 GoLang 从零开始写一个 Docker（容器篇）-- 《自己动手写 Docker》读书笔记](https://juejin.cn/post/6973901434555203598)
3. [使用 GoLang 从零开始写一个 Docker（镜像篇）-- 《自己动手写 Docker》读书笔记](https://juejin.cn/post/6976152015747596301)
4. [使用 GoLang 从零开始写一个 Docker（容器进阶篇/完结篇？）-- 《自己动手写 Docker》读书笔记](https://juejin.cn/post/6978120651676581895)

# 1. How to use
Similar with the official Docker, my docker can only be used on **Linux**. Because it need the features of LXC(and aufs). You can try to run it on wsl(not tested although).

In linux environment, just run the pre-build binary file: `./docker`.  

# 2. Implemented commands
These have exactly the same usages as official Docker.
1. ./docker images
2. ./docker ps
3. ./docker run 
    - --name
    - presudo terminal or detach 
      - -it
      - -d
    - volume
      - -v
    - resource limit
      - -m
      - -cpu
      - cpushare
4. ./docker commit
5. ./docker logs
6. ./docker exec
7. ./docker stop
8. ./docker rm

# 3. How to use other images
The first time you run `./docker images`, you will find only one image named busybox.   
If you want to use other images, you need to follow the below steps:
1. have an installed official docker
2. docker pull the image you want to use
3. run an image as a container(using -d)
4. docker export this container as a tar file
5. move this tar file to the `runtime` folder.
6. have fun
