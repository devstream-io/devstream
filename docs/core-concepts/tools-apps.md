# Tools and Apps

## 1 tools

DevStream treats everything as a concept named _Tool_:

- Each _Tool_ corresponds to a DevStream plugin, which can either install, configure, or integrate some DevOps tools.
- Each _Tool_ has its Name, InstanceID, and Options.
- Each _Tool_ can have its dependencies, specified by the `dependsOn` keyword.

The dependency `dependsOn` is an array of strings, each element being a dependency.

Each dependency is named in the format of "TOOL_NAME.INSTANCE_ID".

---

## 2 apps

Sometimes, you have to define multiple _Tools_ for a single app/microservice. For example, for a web application, you might need to specify the following tools:

- repository scaffolding
- continuous integration
- continuous deployment

If you have multiple apps to manage, you'd have to create many _Tools_ in the config, which can be tedious and hard to read.

To manage multiple apps/microservices more easily, DevStream has another level of abstraction called _Apps_. You can define everything within one app (like the aforementioned repository scaffolding, CI, CD, etc.) with only a few config lines, making the config much easier to read and manage.

Under the hood, DevStream would still convert your _Apps_ configuration into _Tools_ definition, but you do not have to worry about it.

---

## 3 pipelineTemplates

A DevStream _App_ can refer to one or multiple elements of pipelineTemplates, which are mainly CI/CD definitions. In this way, the _Apps_ definition can be shorter, sharing common CI/CD pipelines between multiple microservices.
