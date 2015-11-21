gat
===========

Graphical `cat` command

```sh
gat ritsu.png
```

![](https://cloud.githubusercontent.com/assets/931554/11317166/b0b4a2ce-9066-11e5-8341-d536b22b656a.png)

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
gat -cell="**" [imagefile] # output will be constructed with this text
gat -picker="center" [imagefile] # pick center color of output cell
gat -debug [imagefile] # with indexing cells
```

# issues

- https://github.com/otiai10/gat/issues

# advanced projects

- https://github.com/otiai10/amesh

# thanks

- https://github.com/fatih/color/blob/master/color.go
- https://gist.github.com/MicahElliott/719710
- http://hiroki.jp/2012/06/17/4398/
- http://www.m-bsys.com/linux/echo-color-1
- http://qiita.com/mollifier/items/40d57e1da1b325903659
- http://d.hatena.ne.jp/zariganitosh/20150224/escape_sequence
- http://vorfee.hatenablog.jp/entry/2015/03/17/173635
