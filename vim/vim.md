# Vim

## 索引配置

​	在Rest风格设置 /_settings 或者 {index}/_ _settintgs 设置 一个或者多个索引

```http
PUT http://127.0.0.1:9200/secisland/_settings

{
	“index”：{“number_of_replicas”:4}
}
```

​	在更新分词器，创建索引后添加新的分析器，添加分析器之前必须关闭索引，添加之后关闭

```
POST http://127.0.0.1:9200/secisland/_close
PUT http://127.0.0.1:9200/secisland/_settings
{
	“analysis”{
		“analyzer”{
			“content”：{“type”:"custom","tokenizer":"whitespace"}
		}
	}
}
POST http://127.0.0.1:9200/secisland/_open
```

获取配置 索引中包含很多参数。可以通过命令获取索引的参数配置

```
GET http://127.0.0.1:9200/secisland/_settings
host:port/{index}/_settings
GET http://127.0.0.1:9200/secisland/_settings/name=index.number_*
```

这是拿索引配置

## 索引分析

​	把文本块分析成一个个单独的词(term)，为了后面的倒排做准备，然后标准化为标准形式。提高可搜索性。 都是分析器完成的，分析器有三个功能

- 字符过滤器 字符串经过（character filter）处理，标记化 之前处理字符，能够去出HTML标记，或者转化为“&”或 and
- 分词器：分词器被标记化成独立的词，一个简单的分词器可以根据空格或者逗号把单词分开
- 标记过滤器：每个词都通过所有标记过滤（token filters）处理，它可以修改词 ，例如将“Quick”转为小写怕0，