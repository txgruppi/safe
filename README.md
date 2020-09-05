# Safe

Safe file storage built on top of [BadgerDB](https://github.com/dgraph-io/badger).

## Be careful with your password

At the moment there is only a non-interative cli interface.
It requires that you type invoke the command with the `--password` option.

Only use this version of the Safe if you know how to pass the `--password` option safely.

## Installing

```
$ go get github.com/txgruppi/safe
```

## Usage

Before anything you have to define 2 things.

1. Where to store the database files.
    
    You must defined in which folder your DB files will be stored.

2. Your database password.
    
    This is the master password for your database, choose a strong passowrd and store it safely.
    
    The password must be 16 (AES-128), 24 (AES-192), or 32 (AES-256) bytes long.

## CLI interface

### Commands

#### `ls` - list files

All entries are listed as `<location> <size> <mime type>`.

Since spaces are not allowed in the `<location>` it is easy to pipe this output to other programs to do more complex tasks.

```
$ safe db --database "./myprivatedb" --password "tJgPEjkXbWm0Gl6x" ls
/math/01-01.jpg 18208 image/jpeg
/math/01.jpg 23987 image/jpeg
/math/02.jpg 19960 image/jpeg
/math/03-01.jpg 29811 image/jpeg
/math/03.jpg 39184 image/jpeg
/math/all.jpg 69956 image/jpeg
```

#### `put` - add files to the database

Several files can be added at once by passing to the command as many `<location> <file in disk>` pair as you want.

Any existing file with `<location>` will be replaced.

```
$ safe db --database "./myprivatedb" --password "tJgPEjkXbWm0Gl6x" put \
    /math/01.jpg ~/Downloads/math/01.jpg \
    /math/02.jpg ~/Downloads/math/02.jpg
```

#### `get` - get files from the database

Several files can be fetched at once by passing to the command as many `<location> <file in disk>` pair as you want.

If a `<location>` is not found in the database the file will be ignored and the next pair will be processed.

```
$ safe db --database "./myprivatedb" --password "tJgPEjkXbWm0Gl6x" get \
    /math/01.jpg ~/Downloads/math/01.jpg \
    /math/02.jpg ~/Downloads/math/02.jpg
```

#### `rm` - delete files from the database

Several files can be deleted at once by passing to the command as many `<location>` values as you want.

If a `<location>` is not found in the database the file will be ignored and the next item will be processed.

```
$ safe db --database "./myprivatedb" --password "tJgPEjkXbWm0Gl6x" rm \
    /math/01.jpg /math/02.jpg 
```

## TODO

- [ ] http(s) interface
- [ ] repl / interactive interface