gat
===========

Graphical `cat` command

```sh
gat ritsu.png
```

![](https://cloud.githubusercontent.com/assets/931554/11054870/4c0d8206-87b4-11e5-9b23-60d4262686c1.png)

# install

```sh
go get github.com/otiai10/gat/gat # <- gat/gat!
```

# usage

```sh
gat [imagefile]
gat -b [imagefile]    # with border
gat -h 20 [imagefile] # output height will be about 20 rows
gat -w 40 [imagefile] # output width will be about 40 cols
gat -debug [imagefile] # with indexing cells
```

# issues

https://github.com/otiai10/gat/issues

# thanks

- https://github.com/fatih/color/blob/master/color.go
- https://gist.github.com/MicahElliott/719710
- http://hiroki.jp/2012/06/17/4398/
- http://www.m-bsys.com/linux/echo-color-1
- http://qiita.com/mollifier/items/40d57e1da1b325903659
- http://d.hatena.ne.jp/zariganitosh/20150224/escape_sequence
