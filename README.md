# garbanzo

garbanzo is fast notification viewer!

- RSS feeds
- GitHub notifications

[screencast-localhost_8080-2023.06.08-22_11_14.webm](https://github.com/kijimaD/garbanzo/assets/11595790/4b0c6559-18d0-4f87-9d9d-a04884973a01)

support push notification

![image](https://github.com/kijimaD/garbanzo/assets/11595790/5ce7eab9-efc3-4462-b2cc-059cbbef3dbd)

restriction...

- can't open a link in iframe
- can't open a private link

## Install

#### go install

```
$ go install github.com/kijimaD/garbanzo@main
```

#### brew

```
$ brew install kijimaD/tap/garbanzo
```

## How to use

```
$ garbanzo
```

and, access http://localhost:8080

## [optional] GitHub token

If you want to receive GitHub notifications, require GitHub Personal Access Token(**notification scope**)! Token is used to fetch users notifications.

<img src="https://github.com/kijimaD/garbanzo/assets/11595790/9cabb383-a5a2-484c-8967-0860ad87d5a9" width=800>

## docker run

```
$ docker run --rm -it -p 8080:8080 -p 8081:8081 ghcr.io/kijimad/garbanzo:latest
```

## image

![image](docs/20230528-structure.drawio.svg)

![image](docs/20230529-store.drawio.svg)

## Reference

- [Githubのタイムラインや通知を見るアプリをnode\-webkitで作った \| Web Scratch](https://efcl.info/2014/0430/res3872/)を見て、自分で作ってみようと思った。
- WebSocketまわりは[O'Reilly Japan \- Go言語によるWebアプリケーション開発](https://www.oreilly.co.jp/books/9784873117522/)のチャットルームのコードを参考にした。
