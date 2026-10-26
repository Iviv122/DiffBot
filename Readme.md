Diffbot, your process to spot changes on pages

written in go by a go newbie :D

usage:

- create populate file (e.g list.txt) with format {link, repeat time in seconds}

```
https://link1.com 7
https://link2.com 4
https://link3.com 10
```

- then populate database

```
./diffbot populate ./list.txt ./foo.db
```

or create database and then populate

```
./diffbot initdb ./foo.db
./diffbot populate ./list.txt ./foo.db
```

after that you can use it with

```
./diffbot run ./foo.db
```

features:

- simple cli
- automatic database initialization
- descriptive logging
- multithreading

used dependecies:

- github.com/glebarez/go-sqlite
- github.com/nuskey8/go-cliapp
- github.com/pmezard/go-difflib/difflib

```

```
