# go_skel

Go template repository.

```
git grep -l 'go_skel' | xargs sed -i 's/go_skel/your_repo/g'
git grep -l 'kijimaD' | xargs sed -i 's/kijimaD/your_name/g'
```

## install

```
$ go install github.com/kijimaD/go_skel@main
```

## docker run

```
$ docker run -v "$PWD/":/work -w /work --rm -it ghcr.io/kijimad/go_skel:latest
```
