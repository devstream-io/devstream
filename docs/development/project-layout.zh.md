# 项目结构

参见 [`标准 Go 项目布局`](https://github.com/golang-standards/project-layout) 以了解更多背景信息。

更多关于命名和组织包以及其他代码结构的建议如下:
* [GopherCon EU 2018: Peter Bourgon - Best Practices for Industrial Programming](https://www.youtube.com/watch?v=PTE4VJIdHPg)
* [GopherCon Russia 2018: Ashley McNamara + Brian Ketelsen - Go best practices.](https://www.youtube.com/watch?v=MzTcsI6tn-0)
* [GopherCon 2017: Edward Muller - Go Anti-Patterns](https://www.youtube.com/watch?v=ltqV6pDKZD8)
* [GopherCon 2018: Kat Zien - How Do You Structure Your Go Apps](https://www.youtube.com/watch?v=oL6JBUk6tj0)

一篇关于面向包的设计和架构分层的中文文章
* [面向包的设计和架构分层](https://github.com/danceyoung/paper-code/blob/master/package-oriented-design/packageorienteddesign.md)

## 目录

### `/cmd`

存放这个项目的主要程序。

每个程序的目录名称应与你想拥有的可执行文件的名称相匹配（例如，`/cmd/myapp`）。

不要在应用程序目录中放大量的代码。如果你认为这些代码可以被复用的，那么它应该住在`/pkg`目录中。如果代码不能被复用，或者你不希望它别人被复用，那么请把这些代码放在`/internal`目录下。你会惊讶于别人会做什么，所以请明确自己的意图。

通常会有一个`main`函数，导入并调用`/internal`和`/pkg`目录中的代码，而不是其他函数。

### `/internal`

存放私人程序和库的代码。这里是你不希望别人在其应用程序或库中导入的代码。请注意，这种布局模式是由 Go 编译器强制执行的。更多细节请参见 Go 1.4 [`release notes`](https://golang.org/doc/go1.4#internalpackages) 。请注意，你并不局限于顶层的`internal`目录。你可以在项目树的任何一级有一个以上的`internal`目录。

你可以选择性的给你的内部包添加一些额外的结构，以分离你的共享和非共享的内部代码。这不是必须的（特别是对于较小的项目），但可以很好的在视觉上展示包的用途。程序代码可以放在`/internal/app`目录下（例如，`/internal/app/myapp`），被其他程序复用的代码可以放在`/internal/pkg`目录下（例如，`/internal/pkg/myprivlib`）。

### `/pkg`

存放可以被外部应用程序使用的库代码（例如，`/pkg/mypubliclib`）。 其他项目会导入这些库，并且希望它们能够正常工作，所以在你把代码放在这里之前请三思而后行 :-) 请注意，`internal`目录是确保你的私有包不被他人导入的更好的方法，因为它是由 Go 强制执行的。 `/pkg`目录是一个很好的方式，可以明确地传达该目录中的代码可以被他人安全使用。Travis Jeffery 发表的[`I'll take pkg over internal`](https://travisjeffery.com/b/2019/11/i-ll-take-pkg-over-internal/)博文对`pkg`和`internal`目录进行了很好的概述，其内容包含了如何更好的使用它们。

当你的根目录包含大量的非 Go 组件和目录时，它也是一种将 Go 代码分组在一个地方的方法，使其更容易运行各种 Go 工具（正如在这些会谈中提到的：来自 GopherCon EU 2018 的[`Best Practices for Industrial Programming`](https://www.youtube.com/watch?v=PTE4VJIdHPg)、[GopherCon 2018: Kat Zien - How Do You Structure Your Go Apps](https://www.youtube.com/watch?v=oL6JBUk6tj0)以及[GoLab 2018 - Massimiliano Pippi - Project layout patterns in Go](https://www.youtube.com/watch?v=3gQa1LWwuzk)）。

如果你的项目非常小，而且额外的嵌套层次并没有什么意义（除非你真的想这样做:-)），不使用它是可以的。可以在你的项目越来越大，根目录变得杂乱的时候（特别是当你有很多非 Go 应用组件的时候）再进行考虑。

`pkg`目录的由来：旧的 Go 源代码曾经使用`pkg`来存放它的包，然后社区中的各种 Go 项目开始复制这种模式（更多背景请参见[`这里`](https://twitter.com/bradfitz/status/1039512487538970624) Brad Fitzpatrick 的 tweet）。

### `/vendor`

存放程序依赖（手动管理或通过依赖管理工具，如新的内置[`Go Modules`](https://github.com/golang/go/wiki/Modules)功能）。`go mod vendor`命令将为你创建`/vendor`目录。请注意，如果你没有使用 Go 1.14，你可能需要在你的`go build`命令行添加`mod=vendor`，因为该命令是默认开启的。

如果你正在编写一个库，请不要提交程序的依赖关系。

请注意，自从[`1.13`](https://golang.org/doc/go1.13#modules) 版本 Go 也启用了模块代理功能（默认使用[`https://proxy.golang.org`](https://proxy.golang.org)作为其模块代理服务器）。 在[`这里`](https://blog.golang.org/module-mirror-launch)阅读更多关于它的信息，看看它是否符合你所有的要求。如果符合，那么就不需要`vendor`目录。

## 通用程序目录

### `/hack`

[`/hack`](https://github.com/devstream-io/devstream/blob/main/hack/README.md)目录包含许多脚本，确保 DevStream 的持续发展。

### `/build`

打包和持续集成。

把你的云（AMI）、容器（Docker）和操作系统（deb、rpm、pkg）软件包配置和脚本放在`/build/package`目录下。

把你的 CI（travis, circle, drone）配置和脚本放在`/build/ci`目录下。请注意，一些 CI 工具（例如Travis CI）对其配置文件的位置非常挑剔。试着把配置文件放在`/build/ci`目录中，把它们链接到 CI 工具所期望的位置（如果可能的话）。

### `/test`

额外的外部测试程序和测试数据。你可以随心所欲地构建`/test`目录。对于较大的项目，有一个数据子目录是很有必要的。例如，你可以创建`/test/data`或`/test/testdata`目录，如果你需要 Go 忽略该目录中的内容。请注意，Go 也会忽略以". "或"_"开头的目录或文件，所以你在如何命名你的测试数据目录方面有更多灵活性。

## 其他目录

### `/docs`

存放设计和用户文件（除了 godoc 生成的文件之外）。

### `/examples`

存放应用程序与库的示例。

## 不应该存在的目录

### `/src`

有些 Go 项目确实有一个`src`文件夹，但这通常发生在开发人员来自 Java 世界的时候，在那里这是一种常见的模式。请尽可能不要采用这种 Java 模式。请不要将你的 Go 代码或 Go 项目看起来像 Java :-)

不要把项目级的`/src`目录与[`如何编写Go代码`](https://golang.org/doc/code.html)中描述的 Go 用于其工作区的`/src`目录混淆。`$GOPATH`环境变量指向你的（当前）工作空间（在非 Windows 系统上默认指向`$HOME/go`）。这个工作空间包括顶层的`/pkg`、`/bin`和`/src`目录。你的实际项目实际上是`/src`下的一个子目录，所以如果你的项目中有`/src`目录，项目路径看起来会是这样的：`/some/path/to/workspace/src/your_project/src/your_code.go`。请注意，在 Go 1.11中，你的项目有可能在你的`GOPATH`之外，但这不意味着使用这种布局模式是一个好主意。

## Badges

* [Go Report Card](https://goreportcard.com/) - 它会采用`gofmt`、`go vet`、`gocyclo`、`golint`、`ineffassign`、`license`和`misspell`来扫描你的代码。将`github.com/golang-standards/project-layout`替换为你的项目参照。

    [![Go Report Card](https://goreportcard.com/badge/github.com/golang-standards/project-layout?style=flat-square)](https://goreportcard.com/report/github.com/golang-standards/project-layout)

* [Pkg.go.dev](https://pkg.go.dev) - Pkg.go.dev 是一个新的 Go 语言的 package 和 module 的文档中心。你可以用以下方法创建一个[badge](https://pkg.go.dev/badge)。

    [![PkgGoDev](https://pkg.go.dev/badge/github.com/golang-standards/project-layout)](https://pkg.go.dev/github.com/golang-standards/project-layout)

* Release - 它将显示你的项目的最新版本号。改变 github 链接以指向你的项目。

    [![Release](https://img.shields.io/github/release/golang-standards/project-layout.svg?style=flat-square)](https://github.com/golang-standards/project-layout/releases/latest)