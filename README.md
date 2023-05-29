# garbanzo

![image](https://github.com/kijimaD/garbanzo/assets/11595790/e4acc4ce-4bc6-45c4-a5c7-0a18a273835f)

## install

```
$ go install github.com/kijimaD/garbanzo@main
```

## docker run

```
$ docker run -v "$PWD/":/work -w /work --rm -it ghcr.io/kijimad/garbanzo:latest
```

## image

![image](docs/20230528-structure.drawio.svg)

![image](docs/20230529-store.drawio.svg)

## Reference

WebSocketまわりは[O'Reilly Japan \- Go言語によるWebアプリケーション開発](https://www.oreilly.co.jp/books/9784873117522/)のチャットルームのコードを参考にした。
