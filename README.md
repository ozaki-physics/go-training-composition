# go-training-composition
Go 言語のプロジェクトを作るとき どのようなディレクトリ構成にするか

## 目的 Overview
Go でアプリを作るときに シンプルなディレクトリ構成を勉強する

## インストール方法 Install

## 環境 Requirement
- Docker

## 使い方 Usage
```bash
$ docker-compose build
$ docker-compose up
$ docker-compose exec go_training bash
$ docker-compose down
```

`go.mod`を更新したときは build し直してもいいかも

```bash
/go/src/github.com/ozaki-physics/go-training-composition# go mod init $REPOSITORY
```
## 参考文献 References
[Standard Go Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_ja.md)  
[Go の公式 github](https://github.com/golang/go)
