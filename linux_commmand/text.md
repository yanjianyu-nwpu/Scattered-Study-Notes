# Text Tool

## 0 Cat

### 0.0 概述

命令描述：cat命令用于查看内容较少的纯文本文件

cat [opt] [文件]

### 0.1 选项

-n --number 显示行号

-b --number-nonblank 显示行号，但是不对空白行进行编号

-s --squeeze-blank 当遇见连续两行以上的空白行，只显示一条



## 1 stat

命令描述，用来显示文件的详细信息 包括inode atime mtime ctime 等



## 2 wc

-l 只显示行数 -w 只显示单词书 -c 只显示字节数

-w 只显示单词数

-c 只显示字节数

## 3 file

file 表示用于辨识文件类型

file [参数] [文件]

-b 列出辨识结果时 不显示文件名称

-c 详细显示执行过程



## 4 Grep

grep 用于查找符合条件的字符串

grep全程是 Global Regular Expression Print 表示全局正则表达式版本，他能使用正则表达式搜索文本，并把匹配的行打印出来



在shell 脚本中，grep 通过返回一个状态值来表示搜索的状态



grep [参数] [正则表达式] [文件]

### 4.1 参数

-c --count 计算符合样式的列数

-d recurse 或者-r 指定是目录而非文件

-e 【范本样式】指定字符串作为查找文件内容的样式

-E或--extended-regexp 将样式视为固定字符串的列表

-F或--fixed-regexp 将样式视为固定字符串的列表

-G --basic-regexp 将样式视为普通的表示法来使用

-i 或 --ignore-case 忽略大小写

-n 或 --line-number 在显示符合样式的哪一行之前，标示出改行的列数编号

-v 或 --revert-match 显示不包含匹配文本的所有行

### 4.2 例子

ps -ef | grep sshd 这里用的是pipe可以联合使用