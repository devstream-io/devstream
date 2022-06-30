# 项目布局

参见[`Standard Go Project Layout`](https://github.com/golang-standards/project-layout)以了解更多背景信息。

更多关于命名和组织包以及其他代码结构的推荐阅读。
* [GopherCon EU 2018: Peter Bourgon - 工业编程的最佳实践](https://www.youtube.com/watch?v=PTE4VJIdHPg)
* [GopherCon Russia 2018: Ashley McNamara + Brian Ketelsen - Go最佳实践。](https://www.youtube.com/watch?v=MzTcsI6tn-0)
* [GopherCon 2017: Edward Muller - Go Anti-Patterns](https://www.youtube.com/watch?v=ltqV6pDKZD8)
* [GopherCon 2018: Kat Zien - How Do You Structure Your Go Apps](https://www.youtube.com/watch?v=oL6JBUk6tj0)

一篇关于面向包的设计指南和架构层的中文帖子
* [面向包的设计和架构分层](https://github.com/danceyoung/paper-code/blob/master/package-oriented-design/packageorienteddesign.md)

## 目录

### `/cmd `

本项目的主要应用程序。

每个应用程序的目录名称应该与你想要的可执行文件的名称一致（例如，`/cmd/myapp`）。

不要在应用程序目录中放大量的代码。如果你认为这些代码可以被导入并用于其他项目，那么它应该放在`/pkg`目录中。如果代码不能重复使用，或者你不希望别人复用改代码，那么就把这些代码放在`/internal`目录下。你无法预料别人会如何你的代码，所以要明确说明你的意图!

Go的惯用法是编写一个少而简洁的`main'函数，导入并调用`/internal'和`/pkg'目录中的代码。

### `/internal `

存放不对外暴露的应用程序和库代码。这是你不希望别人在他们的应用程序或库中导入的代码。注意，这种布局模式是由Go编译器本身强制执行的。更多细节请参见 Go 1.4 [`release notes`](https://golang.org/doc/go1.4#internalpackages) 。请注意，你并不局限于只能在顶层的设置`internal`目录。你可以在你的项目树的任何一级设置多个的`internal`目录。

你可以选择为你的内部包添加一些额外的结构，以分离你的共享和非共享的内部代码。这不是必须的（特别是对于较小的项目），但有视觉上的显示包的用途是很好的做法。你的实际应用代码可以放在`/internal/app`目录下（例如，`/internal/app/myapp`），这些应用共享的代码放在`/internal/pkg`目录下（例如，`/internal/pkg/myprivlib`）。

### `/pkg`

可以被外部应用程序使用的库代码（例如，`/pkg/mypubliclib`）。其他项目会导入这些库，期望它们能够正确发挥作用，所以在你把东西放在这里之前要三思 :-) 请注意，`internal`目录是确保你的私有包不能被导入的更好方法，因为它是由Go强制执行的。`/pkg`目录仍然是一个很好的方式，可以明确地传达该目录下的代码可以被他人安全使用。Travis Jeffery的[`I'll take pkg over internal`](https://travisjeffery.com/b/2019/11/i-ll-take-pkg-over-internal/)博文提供了关于`pkg`和`internal`目录的很好的概述，以及何时使用它们可能是有意义的。

当你的根目录包含很多非Go的组件和目录时，这也是一种将Go代码集中在一个地方的方法也是一种不错的实践，这样使你更容易运行各种Go工具（正如这些讲座中提到的。来自GopherCon EU 2018的[`Best Practices for Industrial Programming`](https://www.youtube.com/watch?v=PTE4VJIdHPg)，[GopherCon 2018: Kat Zien - How Do You Structure Your Go Apps](https://www.youtube.com/watch?v=oL6JBUk6tj0)和[GoLab 2018 - Massimiliano Pippi - Project layout patterns in Go](https://www.youtube.com/watch?v=3gQa1LWwuzk)）。

如果你的应用程序项目真的很小，而且额外的嵌套层次并没有增加多少负担（除非你真的想这么做:-)），不使用它也没关系。当你的项目变得足够大，你的根目录变得相当繁杂混乱时，再考虑一下优化结构也不迟（特别是当你有很多非Go应用组件时）。

`pkg`目录的由来。旧的Go源代码曾经使用`pkg`来存放它的包，然后社区中的各种Go项目开始引用这种模式（更多背景请参见[`this`](https://twitter.com/bradfitz/status/1039512487538970624) Brad Fitzpatrick的tweet）。

### `/vendor `

应用程序的依赖关系（手动管理或通过你喜欢的依赖管理工具，如新的内置[`Go Modules`](https://github.com/golang/go/wiki/Modules)功能）。`go mod vendor`命令将为你创建`/vendor`目录。注意，如果你使用的不是Go 1.14，你可能需要在你的`go build`命令中添加`mod=vendor`标志，因为Go 1.14默认是打开的。

如果你正在构建一个库，不要提交你的应用程序依赖。

请注意，从[`1.13`](https://golang.org/doc/go1.13#modules)开始，Go也启用了模块代理功能（默认使用[`https://proxy.golang.org`](https://proxy.golang.org)作为其模块代理服务器）。阅读更多关于它的信息[`here`](https://blog.golang.org/module-mirror-launch)，看看它是否符合你所有的要求和限制。如果符合，那么你就根本不需要`vendor`目录了。

## 通用的应用目录

### `/hack `

[`/hack`](https://github.com/devstream-io/devstream/blob/main/hack/README.md)目录包含许多用于确保DevStream持续发展的脚本，。

### `/build`

打包和持续集成。

把你的云（AMI）、容器（Docker）、操作系统（deb、rpm、pkg）包、配置和脚本放在`/build/package`目录下。

把你的CI（travis, circle, drone）配置和脚本放在`/build/ci`目录下。注意，一些CI工具（例如Travis CI）对其配置文件的位置有特殊要求。试着把配置文件放在`/build/ci`目录中，把它们链接到CI工具期望的位置（如果有必要的话）。

### `/test `

额外的外部测试应用程序和测试数据。你可以随心所欲地构建`/test`目录。对于较大的项目，是有必要设置一个数据子目录的。例如，你可以设置`/test/data`或`/test/testdata`，让Go忽略该目录中的内容。注意，Go也会忽略以". "或"_"开头的目录或文件，所以你在如何命名你的测试数据目录方面有更大的灵活性。

## 其他目录

### `/docs `

存放设计和用户文档（除了你的godoc生成的文档外）。

### `/examples `

你应用程序的事例和/或公共库。

## 不应该有的目录

### `/src

有些Go项目确实有一个`src`文件夹，但这通常发生在开发人员来自Java世界的时候，在java中这是一种常见的模式。如果你能试图让自己尽量不要采用这种Java模式，你将写出真正正确的Go代码或Go项目:-)

不要把项目级的`/src`目录与[`如何编写Go代码`](https://golang.org/doc/code.html)中描述的Go用于其工作空间的`/src`目录混淆。`$GOPATH`环境变量指向你的（当前）工作空间（在非windows系统上默认指向`$HOME/go`）。这个工作空间包括顶层的`/pkg`、`/bin`和`/src`目录。你的实际项目最终会成为`/src`下的一个子目录，所以如果你的项目中有`/src`目录，项目路径会是这样的：`/some/path/to/workspace/src/your_project/src/your_code.go`。注意，在Go 1.11中，你的项目有可能在你的`GOPATH`之外，但这种布局模式并不是一个好的方式。


## 徽章

* [Go Report Card](https://goreportcard.com/) - 它将用`gofmt`、`go vet`、`gocyclo`、`golint`、`ineffassign`、`license`和`misspell`扫描你的代码。将`github.com/golang-standards/project-layout`替换为你的项目。

    [！[Go报告卡](https://goreportcard.com/badge/github.com/golang-standards/project-layout?style=flat-square)](https://goreportcard.com/report/github.com/golang-standards/project-layout)

* [Pkg.go.dev](https://pkg.go.dev) - Pkg.go.dev是一个全新的Go包搜索和文档的站点。你可以使用[徽章生成工具](https://pkg.go.dev/badge)创建一个徽章。

    [！[PkgGoDev](https://pkg.go.dev/badge/github.com/golang-standards/project-layout)](https://pkg.go.dev/github.com/golang-standards/project-layout)

* 发布 - 它将显示你的项目的最新版本号。通过修改github链接来指向你的项目。

    [![Release](https://img.shields.io/github/release/golang-standards/project-layout.svg?style=flat-square)](https://github.com/golang-standards/project-layout/releases/latest)
