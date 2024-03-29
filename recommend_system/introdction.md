## 1 用户数据

### 1.1 日志

- 用户行为汇总会话日志 （Session log）

- 服务端生成展示日志（impression log）

- 还有点击日志（click log）

### 1.2 用户行为

- 显性反馈（explicit feedback）
  
  - 商品评价等

- 隐形反馈（implicit feedback）
  
  - 点击
  
  - 停留
  
  - 展现

## 2 推荐系统标准

### 2.1 长尾理论

在互联网理论下 8/2 理论 更加倾斜                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  

### 2.2 推荐系统经典场景

- 电商

- 电影/视频网站

- 个性化音乐网站

- 社交网络(facebook)

- 个性化服务

- 基于位置的服务

- 个性化邮件

- 个性化广告

### 2.3 推荐系统评测

- 推荐系统参与方
  
  - 用户
  
  - 网站
  
  - 内容提供方
  
  - 广告主

- 推荐系统实验方法
  
  - 离线实验
    
    - 记录日志，生成标准数据集
    
    - 可以+1 检测模型效果
  
  - 用户调查
    
    - 用户给推荐的物品打分
      
      - 优点是非常直接，但是成本很高；且选取的用户 分布需要符合所有用户的分布
  
  - 在线实验
    
    - 在线ab实验 可以在线及时H评估算法 但是要统计层 不同层互不干扰。

- 评价指标
  
  - 用户满意度
  
  - 预测准确度
    
    - 评分预测
    
    - TopN推荐
  
  - 覆盖率
    
    - 描述一个推荐系统对长尾的发掘能力，纯模型或者算法很容易陷入只推荐 特别热的物品。就是陷入马太效应，只推荐用户喜欢的/不能推荐用户以前没看过，但是可能会喜欢的。导致类别基尼系数上升。
  
  - 多样性
    
    - 感觉就是与覆盖率类似，就是不要陷入一个类目
  
  - 实时性
    
    - 新闻/微博 有很强的实时性
  
  - 健壮性
    
    - 防止攻击
  
  - 商业目标
    
    - 赚钱

- 测评维度
  
  主要还是三种类目：
  
  - 用户维度 用户喜欢看，喜欢点
  
  - 物品维度 被买的多
  
  - 事件维度 主要还是实时性

## 3 利用用户行为数据

## 2.1 用户活跃度和物品流行度的分布

- 大部分都是Power law 长尾分布

### 2.2 用户活跃度和武物品流行度的关系

- 基于用户行为数据设计的推荐算法   一般成为协同过滤算法
  
  - 基于领域的方法
  
  - 隐语义模型
  
  - 基于图的随机游算法

### 2.3 基于领域的算法

- 基于用户的协同过滤算法
  
  - 找到和目标用户类似的用户集合
  
  - 找到这个集合中用户喜欢的，且目标用户么有看过的商品。
    
    - 相似度预余弦公式
    
    - 然后计算用户的相似度矩阵
  
  - 算法评估
    
    - 准确性和召回率
    
    - 流行度
    
    - 覆盖率： 多种推荐

- 基于物品的协同过滤算法
  
  - 缺点就是随着用户数目越来越多，计算用户兴趣的举证越来越苦难

### 2.4 基于图的模型

- 传统算法
  
  - 聚类
  
  - 随机   游走

- 深度学习

## 3 推荐系统冷启动

### 3.1 冷启分类

- 用户冷启动   
  
  - 刚来没有用户特征

- 物品冷启动
  
  - 新物品没有特征

- 系统冷启
  
  - 新系统没有特征

- 提供非个性化的推荐

## 4 利用用户标签

标签一般分成两种

- UGC user general content 用户产生

- PGC 专家产生

## 5 基于上下文的推荐

基于上下文，还有地理位置推荐

## 6 基于社交网络的推荐

基于领域的推荐算法
