# nassort

A natural directory assorter

## Installation

```
$ go get -u github.com/syucream/nassort
```

## Usage

* Write definitions in yaml file.

```
# for example
$ vim nassort.yaml
-
  dst: moved
  src:
    - contains:
      - target
```

* Before directory tree ... :

```
.
├── dst
├── nassort
├── nassort.yaml
└── src
    └── foo_target_bar
        └── file01
```

* Run nassort

```
$ ./nassort -src src/ -dst dst/ -f nassort.yaml
```

* Then files moved based on nassort.yaml:

```
.
├── dst
│   └── moved
│       └── foo_target_bar
│           └── file01
├── nassort
├── nassort.yaml
└── src
```
