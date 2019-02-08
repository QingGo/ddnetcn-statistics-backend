## ddnet国服2018年统计数据

to-do list：
* 去重后的新增的分数和国服相应排名，原排名，现排名
* 最常玩的图，排名？
* 最常组队过图的队友，排名？
* 过图时间一天之内哪个小时，一星期那一天，是否需要得到数组用来做图。
* MongoDB缓存已查询结果，通过redis储存未处理，处理中，处理完成的状态，通过MQ限制并发查询数？
* 还是一开始就预处理好？那是放Mysql还是MongoDB？
* 重构减少重复代码？
* vue写前端，或者试试其它框架
* nginx处理静态文件和反向代理，通过docker一键部署代码和所有依赖，试试k8s
* 测试，ci持续集成

部署笔记
下载游戏过图sql文件
wget https://ddnet.tw/stats/ddnet-sql.zip
unzip ddnet-sql.zip
安装mysql或mariadb，root用户登录：
create database ddnet;
CREATE USER 'ddnet'@'%' IDENTIFIED BY 'ddnet';
GRANT ALL PRIVILEGES ON ddnet.* TO 'ddnet'@'%';
FLUSH PRIVILEGES;
source record_teamrace.sql;
source record_maps.sql;
source record_race.sql;
安装go，
安装go的依赖：
go get github.com/gin-gonic/gin
go get github.com/go-sql-driver/mysql
进入代码目录，试运行：
go run main.go
实际运行：
go build main.go
nohup ./main >> web.lob 2>&1 &