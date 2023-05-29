# garbanzo

![image](https://github.com/kijimaD/garbanzo/assets/11595790/f4843b2c-1ec2-486c-9ce2-f2c5e76cacaf)

## install

```
$ go install github.com/kijimaD/garbanzo@main
```

## how to use

Need GitHub Token(notification scope). Token is used to fetch user notification.

```
$ GH_TOKEN=xxx garbanzo
```

and, access http://localhost:8080

## docker run

```
$ docker run -v "$PWD/":/work -w /work --rm -it ghcr.io/kijimad/garbanzo:latest
```

## image

![image](docs/20230528-structure.drawio.svg)

![image](docs/20230529-store.drawio.svg)

## Reference

WebSocketまわりは[O'Reilly Japan \- Go言語によるWebアプリケーション開発](https://www.oreilly.co.jp/books/9784873117522/)のチャットルームのコードを参考にした。
