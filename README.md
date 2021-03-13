## 项目名称
> typoraUpImage  
## 配置文件
文件为config.json
```json
{
    "message": "add", //上传的时候的附加信息     
    "branch": "master", // 分支
    "token":"", // api的token信息
    "userName":"", // 你的gitee的用户名
    "repositorie":"picbed", //你的代码仓库的名
    "Folder":"image", // 代码库里的文件夹名
    "BucketDomain": "https://gitee.com/api/v5/repos/" //这个是固定的，可以去看一下api              
}
```


## make的作用

```
内建函数make(T, args)与new(T)的用途不一样。它只用来创建slice，map和channel，并且返回一个初始化的(而不是置零)，类型为T的值（而不是*T）。之所以有所不同，是因为这三个类型的背后引用了使用前必须初始化的数据结构。例如，slice是一个三元描述符，包含一个指向数据（在数组中）的指针，长度，以及容量，在这些项被初始化之前，slice都是nil的。对于slice，map和channel，make初始化这些内部数据结构，并准备好可用的值。
```



## map

```
在json.Unmarshal传参中第二个参数必须为指针类型,因为第二个类型是空接口类型具有普遍适用性,所以map也必续用&取指针.可以理解为map为特殊类型的指针可以达到指针的效果,并且&map和&int用*引用后效果是一样的,&map取到的是存放map本身指针的一个变量的指针,因为map本身的特殊性用一个*就可以引用map本身,这是golang语言内进行处理的.
```

### typora运行命令

> 运行命令必须先转到文件再运行
>
> 例如:windows的命令
>
> cd /d D:\golangCode\typoraUpImage & main.exe