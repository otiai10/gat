gat
===========

Graphical `cat` command

```sh
# Local file
% gat example.png

# By URL
% gat https://raw.githubusercontent.com/otiai10/gat/master/samples/sample.png
```

| default | iTerm | Sixel |
|:-------:|:------:|:------:|
| <img width="320px" src="https://cloud.githubusercontent.com/assets/931554/11317166/b0b4a2ce-9066-11e5-8341-d536b22b656a.png"> | <img width="320px" src="https://user-images.githubusercontent.com/931554/44513376-26ed8280-a6f8-11e8-83df-f1f877228189.png"> | supported but no image |

# Install

```sh
go get -u github.com/otiai10/gat
```

# Options

```sh
# iTerm and Sixel
gat -S 0.5 [imagefile] # Scale of output image

# Only for cell grid mode
gat -c [imagefile]      # Use cell grid mode
gat -H 20 [imagefile]   # Rows of output
gat -W 40 [imagefile]   # Cols of output
gat -b [imagefile]      # Print border
gat -t="**" [imagefile] # Text to be printed for each cell
gat -debug [imagefile]  # Dump available colors and color for each cell
```

# Thanks

- https://github.com/fatih/color/blob/master/color.go
- https://gist.github.com/MicahElliott/719710
- http://hiroki.jp/2012/06/17/4398/
- http://www.m-bsys.com/linux/echo-color-1
- http://qiita.com/mollifier/items/40d57e1da1b325903659
- http://d.hatena.ne.jp/zariganitosh/20150224/escape_sequence
- http://vorfee.hatenablog.jp/entry/2015/03/17/173635
