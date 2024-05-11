# shelf

shelf is PDF viewer. It manage bibliographic information on a file basis.

## install

```
go install github.com/kijimaD/shelf@main
```

## how to use

Shelf generate book config files, and serve PDF files.

```
$ ls
aa.pdf
bb.pdf

# generate shelf file schema(warn: renaming pdf occurs!)
$ shelf gen .

$ ls
20240311T213221070837339_aa.pdf
20240311T213221080079018_bb.pdf
index.toml

# start shelf server
$ shelf web
```

and open http://localhost:8020/

my real use repository: https://github.com/kijimaD/mypdfs

## image

![](./images/top.png)

## help

```
$ shelf print

web             start shelf app server          web
gen             generate files in directory     generate [DIRECTORY]
gensingle       generate file                   gensingle [FILE]
extract         extract image                   extract [PDF]
print           print all command description   print
```
