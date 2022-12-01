# Version 1

前缀: /v1/api

# 用户相关 user

注册 /user/register
post:
{
"username": "",
"password": "",
"email": "",
"chinese_name": "",
"phone": "",
}
get:

登录 /user/login
post:
{
"username": "",
"password": ""
}
get:

登出 /user/logout
post:

鉴权 /user/authenticate
{
"username":"",
"password":""
}

更新 /user/update
{
//update message
"...":"..."
}

删除指定用户 /user/delete/{uid: int}

获取指定用户信息 /user/{uid: int}

获取所有用户的信息 /user

# 项目相关 project

创建 /project/create

更新 /project/update

删除 /project/delete

获取指定项目的信息 /project/{pid: int}

获取所有项目的信息 /project/

# 计划相关  test_plan

创建 /test_plan/create

更新 /test_plan/update

删除 /test_plan/delete

获取指定计划的信息 /test_plan/{tpid: int}

获取所有计划的信息 /test_plan/

# 任务相关 task

创建 /task/create

更新 /task/update

删除 /task/delete

执行 /task/execute

暂停 /task/pause

恢复 /task/resume

重置 /task/reset

结束 /task/terminate

取消 /task/cancel

获取指定任务的信息 /task/{tid: int}

获取所有任务的信息 /task/
