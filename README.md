# domainToip
复刻棱眼，对收集到的域名进行反查ip的功能，再将ip整理成c段

2023-08-23更新<br/>
新增fofa查询功能，可对输入的域名进行子域名查询，并查询其ip的具体信息，最后将其聚合为c段<br/>
在第一次运行该程序时会在当前路径下自动创建config.ini文件，如果需要使用fofa功能，需要填写配置文件<br/>
将-f参数改为fofa查询，-l为本地文件查询<br/>
## 2023-08-23 update

`./domainToip_windows_amd64.exe`

![image](https://github.com/lvshu811/domainToip/assets/82709360/6a4c71be-0d7c-428f-a0fd-1c5d4a8419da)

`./domainToip_windows_amd64.exe -f xxx.com`

![image](https://github.com/lvshu811/domainToip/assets/82709360/5a46c4a7-3824-4294-a63a-dbe152d7a934)

## 使用方法

`./domaintoip.exe -h` //可查看帮助

![image](https://github.com/lvshu811/domainToip/assets/82709360/9e0ec581-263d-4775-9133-a6cc76d58c2f)

`./domaintoip.exe -d www.xxx.com` //可对单个域名进行反查

![image](https://github.com/lvshu811/domainToip/assets/82709360/c8134de6-d813-46b7-95dd-8aee9246ee72)

`./domaintoip.exe -f domain.txt` //可对文件内的域名进行批量查询，并简单判断是否存在cdn，若存在则不显示

![image](https://github.com/lvshu811/domainToip/assets/82709360/fc2e0a1c-ef35-48b8-8fc3-45027283f187)
