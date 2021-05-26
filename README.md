# monitorStock
this project propose to monitor a stock price, and send an email to appointed mail address  
本工程旨在监控指定的一只股票，当股票价格超过设定的值的时候，会发送一封邮件到指定的邮箱进行提醒  

# 使用方法
* 查找股票代码替换代码中的股票代码sz000586
``` golang
    res, err := client.Get("https://hq.sinajs.cn/?list=sz000586") 
    if err != nil {
        fmt.Println("error occurred!:", err)
        continue
    }
```
* 设置邮件的发件箱、密码、邮箱服务器和收件地址
``` golang
	user := "********@126.com"
	password := "*******"
	host := "smtp.126.com:25"
	to := "******@qq.com"
```
* 编译后执行可执行文件，需要传入一个价格值（此处为7.00），当股价超过这个价格时会触发邮件提醒（半小时之内只会触发一次）
```
    ./monitorStock 7.00
```
* 执行过程如下
```
    邮件通知价格阈值： 7
    汇源通信： 6.990   2021-05-26   15:00:03
    汇源通信： 6.990   2021-05-26   15:00:03
```
