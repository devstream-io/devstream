# Commit 信息
我们尽最大的努力遵守[conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary)规范。

TL;DR：提交信息应当结构如下：

```
<type>[optional scope]: <description>
[optional body]
[optional footer(s)]
```

其中的"type"可以是以下内容：

- `feat`: 实现了一个功能，如例1
- `fix`: 修复一个错误
- `BREAKING CHANGE`: 一项重大更改。或者在`feat`或`fix`的结尾处添加`!`，像`feat!`和`fix!`,如例2和例4
- 其他类型同样被允许，例如：`build`, `chore`, `ci`, `docs`, `style`, `refactor`, `perf`, `test`

"body"和"footer"都是可选项；重大改变既可以写在标题里，如例5；又可以写在脚注里，如例6。举例如下：

1. `feat: send an email to the customer when a product is shipped`
2. `feat!: send an email to the customer when a product is shipped`
3. `feat(api): send an email to the customer when a product is shipped`
4. `feat(api)!: send an email to the customer when a product is shipped`
5. `BREAKING CHANGE: send an email to the customer when a product is shipped`
6. ```
   feat!: send an email to the customer when a product is shipped
   A detailed description in the body.
   BREAKING CHANGE: readdressing the breaking change in the footer.
   ```
<!-- todo -->
