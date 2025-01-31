# disable-checkout-persist-credentials

[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/disable-checkout-persist-credentials/main/LICENSE) | [Install](INSTALL.md)

CLI to disable [actions/checkout](https://github.com/actions/checkout)'s persist-credentials.

To modify GitHub Actions workflows, you can run:

```sh
disable-checkout-persist-credentials
```

Then `persist-credentials: false` is patched.

```diff
     steps:
       - uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2
+        with:
+          persist-credentials: false
       - uses: actions/setup-go@3041bf56c941b39c61721a86cd11f3bb1338122a # v5.2.0
         with:
           go-version-file: go.mod
```

You can also specify file paths:

```sh
disable-checkout-persist-credentials .github/workflows/test.yaml .github/workflows/release.yaml
```

This tool supports composite actions too.

```sh
disable-checkout-persist-credentials action.yaml
```

## Why?

- https://github.com/suzuki-shunsuke/ghalint/blob/main/docs/policies/013.md

actions/checkout's input persist-credentials should be false.

## How To Use

```sh
disable-checkout-persist-credentials [<workflow file> ...]
```

By default, `\.github/workflows/.*\.ya?ml` is changed.

## :warning: Known Issues

1. The number of spaces before `#` comment is changed

Before:

```
timeout-minutes: 30   # test
```

After

```
timeout-minutes: 30 # test
```

2. Multiple newlines are modified to a newline

Before:

```
# test


jobs:
```

After:

```
# test

jobs:
```

## LICENSE

[MIT](LICENSE)
