# 操作系统

## Linux概念原理

- 协程的优势
  - 线程调度 会让操作系统进行调度产生开销；协程是线程内部调度节约
  - 协程不需要权限切换，协程和线程都需要上下文切换
  - 协程能得到自己控制挂起和恢复，能够更好的优化
  - 但是协程需要自己控制挂起yield，因为没有抢占，如果一直站着cpu可能不会被强制踢下来
- 虚拟内存相关
  - MMU的寻址
    - 比如C++ 一个指针 64位 8字节
      - 63-48 16位置是一些权限控制
      - 47-39 是一级页表的索引
      - 39-30 二级页表的索引
      - 29-21 三级页表的索引
      - 21-12 四级页表的索引
      - 12-0 页框内部的偏移量 4KB 刚好是 4的12次方 

## Linux Shell 命令行

## 虚拟化技术



