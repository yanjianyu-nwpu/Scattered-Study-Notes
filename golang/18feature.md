# 18 新特性

## 1  c++编译

​	go build 命令 --asan flag 支持 c/c++；联合编译。

## 2 Generics 支持

```
func PrintAnything[T any](thing T) {
    fmt.Println(thing)
}
```

简单的泛型，可以有多个模板



也可以指定一些类型 type constrain 约束

```
func A[T any,k int | float64](value T,s k){
	fmt.Println(value,s)
}
```

