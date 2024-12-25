# disable-checkout-persist-credentials

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

## Why?

- https://github.com/suzuki-shunsuke/ghalint/blob/main/docs/policies/013.md

actions/checkout's input persist-credentials should be false.

## How To Use

```sh
disable-checkout-persist-credentials [<workflow file> ...]
```

By default, `\.github/workflows/.*\.ya?ml` is changed.

## Known Issues

1. The number of spaces before `#` comment is changed

Before:

```
timeout-minutes: 30   # test
```

After

```
timeout-minutes: 30 # test
```

2. `persistent-credentials: true` isn't fixed

## LICENSE

[MIT](LICENSE)
