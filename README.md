# Go energy workspace

workspace : https://github.com/energye/workspace.git

# Module 3.0

| Module name | Repository addr                           | Desc |
|-------------|-------------------------------------------|------|
| lcl         | https://github.com/energye/lcl.git        |      |
| cef         | https://github.com/energye/cef.git        |      |
| wv          | https://github.com/energye/wv.git         |      |
| wk          | https://github.com/energye/wk.git         |      |
| energy      | https://github.com/energye/energy.git     |      |
| examples    | https://github.com/energye/examples.git   |      |
| assetserve  | https://github.com/energye/assetserve.git |      |
| lib         | https://github.com/energye/lib.git        |      |
| gen         | https://gitee.com/energye/gen.git         |      |


# Go workspace 

Go Version >= 1.18

- 初始化 go workspace

`go work init`

- 模块添加到工作区

`go work use [module name]`

go build -ldflags="-H windowsgui -w -s"
go build -ldflags="-w -s"