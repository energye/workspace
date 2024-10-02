# Go energy workspace

workspace : https://github.com/energye/workspace.git

# Module 3.0

| Module name | Repository addr                              | Use depend | Desc                                                     |
|-------------|----------------------------------------------|------------|----------------------------------------------------------|
| lcl         | https://github.com/energye/lcl.git           | ✓          | LCL basic library                                        |
| cef         | https://github.com/energye/cef.git           | ✓          | CEF basic library                                        |
| wv          | https://github.com/energye/wv.git            | ✓          | Webview2 basic library                                   |
| wk          | https://github.com/energye/wk.git            | ✓          | Webkit basic library                                     |
| energy      | https://github.com/energye/energy.git        | ✓          | Energy framework                                         |
| examples    | https://github.com/energye/examples.git      | ✓          | All examples                                             |
| assetserve  | https://github.com/energye/assetserve.git    | ✓          | Built-in http static resource service                    |
| lib         | https://github.com/energye/lib.git           | x          | Binary dynamic link library                              |
| gen         | https://gitee.com/energye/gen.git            | x          | Code generation                                          |
| gitpp       | https://gitee.com/snxamdf/gitpp.git          | x          | Git auto pull push cmd                                   |
| doc-api     | https://gitee.com/energye/energy-doc-api.git | x          | Server API, Website publishing and binary file reception |
| energy-doc  | https://github.com/energye/energye.github.io | x          | Energy DOC                                               |
| workflows   | https://github.com/energye/workflows.git     | x          | Workflows  Automatic publishing                          |


# Go workspace 

Go Version >= 1.20

- 初始化 go workspace

`go work init`

- 将使用的模块添加到工作区

`go work use [module name]`

`go build -ldflags="-H windowsgui -w -s"`

`go build -ldflags="-w -s"`
