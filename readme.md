## 配置
在 tgbot/chatdraw.go 的init 方法中 设置奖品、管理员

## 命令

### 帮助
/help 

	regCommand("/sign", "", Sign)

### 签到
签到的用户，才会进入抽奖名单

/sign

### 初始化抽奖
初始化配置的奖品 与签到的用户为抽奖池

/drawinit

### 基本信息
查询中奖信息，剩余奖品及 剩余抽奖名单

/info

### 选择抽奖用户 
选择或随机确定 下一次抽奖名单

/makedrawuser

### 抽奖
只有选定的用户才能抽奖

/luckydraw

### 回滚抽奖
选择撤回一条抽奖信息

/rollbackrecord

### 新增奖品
添加新的奖品

/addawards 奖品 数量

### 新增抽奖名单
签到抽奖的模式下 添加的用户无法抽奖

/adddrawuser name


### 摇摇乐开始
/rollstart

### 摇摇乐
每人只能摇一次

/roll

### 摇摇乐信息
查询本轮摇摇乐信息

/rollinfo






