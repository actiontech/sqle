## 使用方法
```golang

// 传递文件进行解析
p, err := parser.CreateJavaParser("/root/javaexample/test/Test7.java")
if err != nil {
    os.Exit(-1)
}

// 生成java解析器的访问者
v := parser.NewJavaVisitor()

// 从跟节点开始访问，并生成自定义的变量树
a:=p.CompilationUnit()
a.Accept(v)

// 从变量树中根据jdbc运行sql的函数获取sql
fmt.Println(parser.GetSqlsFromVisitor(v))

/*
delete from t1;
select BYTES from user_segments where segment_name =?
*/
```