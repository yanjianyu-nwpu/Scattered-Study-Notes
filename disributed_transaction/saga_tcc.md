# 分布式事务

## 1 Tcc

tcc 就是 try，confirm 和 cancel

### 1.1 try

比如 调用方 会调用多个 被调用方  例如：

- 先try 冻结库存 -》返回冻结成功

- try冻结账户 -》 返回失败

那么调用cancel 取消所有冻结

## 1.2 confirm

try所有的被调用方都成功，然后真的扣库存

如果失败会重试

### 1.3 cancel

 只有在try 阶段失败或异常的情况才会进行cancel处理

### 1.4 幂等

因为可能超时所以需要重发啊，为了因为重复调用而占用资源，幂等是非常重要。

## 2 Saga 方案

### 2.1 简介

和tcc 的不同就是 try的时候直接扣了，如果失败 调用 补偿方法

saga是一阶段方案

tcc是两阶段方案

saga 是一种长事务的解决方案

对于saga 来说很可能是要做到事务的并发控制，还是在业务逻辑层实现并发控制，

```
https://juejin.cn/post/6857520180894351374
```

以上