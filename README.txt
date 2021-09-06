金融反欺诈(图数据库项目)
主要实现功能包含：
    对用户行为进行监控数据进行采集处理和分析，最终构建多类型的实体关系网络，可以及时预警潜在风险，为金融行业用户提供风险异常检测和反欺诈行为分析。
    todo:
    1、关系推理
    2、关联度检测
    3、集中度测量
    4、语义分析
    5、团伙发现
    6、可视化展示

#### 工具 #####
1、本地docker搭建nebula
    docker pull vesoft/nebula-dev
    docker run --rm -ti --security-opt seccomp=unconfined -v /opt/software/nebula-graph:/home/nebula -w /home/nebula vesoft/nebula-dev 
    容器内执行
        mkdir build && cd build
        cmake -DENABLE_BUILD_STORAGE=on - DCMAKE_EXPORT_COMPILE_COMMANDS=on ..
        路径： /user/local/nebula
        make install-all
        docker exec -it 26602a3943e49e38e4c21b2758a4a1c8c4ca6ce73093cad2c34088a0e0099fd2 /bin/sh
        
2、go语言实现nebula写入以及查询
    cat requirements | xargs go get
    go mod init nebula.go
#### 语言 ######
go
    版本：go1.16.6
#### 数据库 ####
    1、nebula graph

#### moc数据 ####
go mod init