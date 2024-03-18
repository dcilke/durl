# durl

Command line URL utilities

# API

## Decode

Fully decode a URL.

```bash
durl --decode "http%3A%2F%2Fexample.com%2F%3Fq%3Dtest"

http://example.com/?q=test
```

## Encode

Fully encode a URL.

```bash
durl --encode "http://www.example.com/file%20one%26two"

http://www.example.com/file one&two
```

## Password

Extract the password from the URL. Decodes if necessary

```bash
durl --password "http://%3Fam:pa%3Fsword@google.com"

pa?sword
```

# Installation

Install locally via go.

```bash
make install
```
