# shelf

shelf is PDF viewer. It manage bibliographic information on a file basis.

## install

```
go install github.com/kijimaD/shelf@main
```

## how to use

serve current directory PDF files.

```
$ ls
aa.pdf
bb.pdf

$ shelf gen . # generate shelf file schema(warn: renaming pdf occurs!)

$ ls
20240311T213221070837339_aa.pdf
20240311T213221070837339.toml
20240311T213221080079018_bb.pdf
20240311T213221080079018.toml

$ shelf web   # start server
```

and open http://localhost:8020/

## image

![](./images/top.png)
