# Backend of Pixiu

Use Go + Gorm + Sqlite to build the backend of Pixiu.

# 包依赖关系

1. Runtime 层可以依赖以下层：Adapter、Business、Pkg
2. Adapter 层可以依赖以下层：Business、Pkg
3. Business 层可以依赖以下层：Pkg
4. Pkg 层不依赖其他层（基础包）

```mermaid
    Runtime --> Adapter
    Runtime --> Business
    Runtime --> Pkg
    Adapter --> Business
    Adapter --> Pkg
    Business --> Pkg
```

这种分层架构保证了项目的清晰度和可维护性，遵循了从上到下的依赖原则，避免循环依赖。