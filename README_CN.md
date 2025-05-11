[简体中文](README_CN.md) | **[English](README.md)**

---

![GitHub License](https://img.shields.io/github/license/yhlooo/go-tetris)
[![GitHub Release](https://img.shields.io/github/v/release/yhlooo/go-tetris)](https://github.com/yhlooo/go-tetris/releases/latest)
[![release](https://github.com/yhlooo/go-tetris/actions/workflows/release.yaml/badge.svg)](https://github.com/yhlooo/go-tetris/actions/workflows/release.yaml)

# go-tetris

该项目是一个 go 实现的 Tetris 库，同时也是可在浏览器（基于 [Wasm](https://webassembly.org/) ）或终端运行的 Tetris 游戏。

## 玩！

### 在浏览器

![web-ui](docs/img/web-ui.png)

访问 [Tetris](https://yhlooo.github.io/go-tetris/) 立即开始。

### 在终端

![tty-ui](docs/img/tty-ui.png)

可通过下载二进制、 Docker 、从源码编译三种方式之一安装和运行 Tetris ： 

**下载二进制：**

通过 [Releases](https://github.com/yhlooo/go-tetris/releases) 页面下载可执行二进制，解压并将其中 `tetris` 文件放置到任意 `$PATH` 目录下。

然后执行：

```bash
tetris
```

**使用 Docker ：**

```bash
docker run --rm --name tetris -it ghcr.io/yhlooo/tetris:latest
```

**从源码编译：**

```bash
# 下载源码并编译
go install github.com/yhlooo/go-tetris/cmd/tetris@latest
# 运行
$(go env GOPATH)/bin/tetris
```

## 构建该项目

**终端版：**

```bash
go run ./cmd/tetris
```

**浏览器版：**

```bash
GOOS=js GOARCH=wasm go build -o web/app.wasm ./cmd/tetris-wasm && go run ./cmd/tetris-wasm
```

然后访问 <http://localhost:8000> 。

## 构建你自己的 Tetris

该项目不仅是一个可玩的 Tetris 游戏，它同时是一个易于被集成的 Tetris 库。你可以使用它构建你自己的 Tetris 游戏，参考接口 [Tetris](pkg/tetris/tetris.go#L9) 。

## 致谢

- [Hard Drop Tetris Wiki](https://harddrop.com/wiki/Tetris_Wiki) : 提供了关于 Tetris 详细的机制说明
- [rivo/tview](https://github.com/rivo/tview) : 提供了强大的基于终端的 UI
- [maxence-charriere/go-app](https://github.com/maxence-charriere/go-app) : 提供了基于 [Wasm](https://webassembly.org/) 的 Web UI
