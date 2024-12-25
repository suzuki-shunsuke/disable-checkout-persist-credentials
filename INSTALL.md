# Install

disable-checkout-persist-credentials is written in Go. So you only have to install a binary in your `PATH`.

There are some ways to install disable-checkout-persist-credentials.

1. [Homebrew](#homebrew)
1. [Scoop](#scoop)
1. [aqua](#aqua)
1. [GitHub Releases](#github-releases)
1. [Build an executable binary from source code yourself using Go](#build-an-executable-binary-from-source-code-yourself-using-go)

## Homebrew

You can install disable-checkout-persist-credentials using [Homebrew](https://brew.sh/).

```sh
brew install suzuki-shunsuke/disable-checkout-persist-credentials/disable-checkout-persist-credentials
```

## Scoop

You can install disable-checkout-persist-credentials using [Scoop](https://scoop.sh/).

```sh
scoop bucket add suzuki-shunsuke https://github.com/suzuki-shunsuke/scoop-bucket
scoop install disable-checkout-persist-credentials
```

## aqua

[aqua-registry >= v4.284.0](https://github.com/aquaproj/aqua-registry/releases/tag/v4.284.0)

You can install disable-checkout-persist-credentials using [aqua](https://aquaproj.github.io/).

```sh
aqua g -i suzuki-shunsuke/disable-checkout-persist-credentials
```

## Build an executable binary from source code yourself using Go

```sh
git clone https://github.com/suzuki-shunsuke/disable-checkout-persist-credentials
cd disable-checkout-persist-credentials
go install ./cmd/disable-checkout-persist-credentials
```

> [!WARNING]
> Unfortunately, `go install github.com/suzuki-shunsuke/disable-checkout-persist-credentials/cmd/disable-checkout-persist-credentials@latest` doesn't work because disable-checkout-persist-credentials uses a replace directive in go.mod.
> 
> ```console
> $ go install github.com/suzuki-shunsuke/disable-checkout-persist-credentials/cmd/disable-checkout-persist-credentials@latest
> go: github.com/suzuki-shunsuke/disable-checkout-persist-credentials/cmd/disable-checkout-persist-credentials@latest (in github.com/suzuki-shunsuke/disable-checkout-persist-credentials@v0.1.0):
> 	The go.mod file for the module providing named packages contains one or
> 	more replace directives. It must not contain directives that would cause
> 	it to be interpreted differently than if it were the main module.
> ```

## GitHub Releases

You can download an asset from [GitHub Releases](https://github.com/suzuki-shunsuke/disable-checkout-persist-credentials/releases).
Please unarchive it and install a pre built binary into `$PATH`. 

### Verify downloaded assets from GitHub Releases

You can verify downloaded assets using some tools.

1. [GitHub CLI](https://cli.github.com/)
1. [slsa-verifier](https://github.com/slsa-framework/slsa-verifier)
1. [Cosign](https://github.com/sigstore/cosign)

### 1. GitHub CLI

You can install GitHub CLI by aqua.

```sh
aqua g -i cli/cli
```

```sh
version=v0.1.0
asset=disable-checkout-persist-credentials_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/disable-checkout-persist-credentials "$version" -p "$asset"
gh attestation verify "$asset" \
  -R suzuki-shunsuke/disable-checkout-persist-credentials \
  --signer-workflow suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml
```

### 2. slsa-verifier

You can install slsa-verifier by aqua.

```sh
aqua g -i slsa-framework/slsa-verifier
```

```sh
version=v0.1.0
asset=disable-checkout-persist-credentials_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/disable-checkout-persist-credentials "$version" -p "$asset" -p multiple.intoto.jsonl
slsa-verifier verify-artifact "$asset" \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/suzuki-shunsuke/disable-checkout-persist-credentials \
  --source-tag "$version"
```

### 3. Cosign

You can install Cosign by aqua.

```sh
aqua g -i sigstore/cosign
```

```sh
version=v0.1.0
checksum_file="disable-checkout-persist-credentials_${version#v}_checksums.txt"
asset=disable-checkout-persist-credentials_darwin_arm64.tar.gz
gh release download "$version" \
  -R suzuki-shunsuke/disable-checkout-persist-credentials \
  -p "$asset" \
  -p "$checksum_file" \
  -p "${checksum_file}.pem" \
  -p "${checksum_file}.sig"
cosign verify-blob \
  --signature "${checksum_file}.sig" \
  --certificate "${checksum_file}.pem" \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  "$checksum_file"
cat "$checksum_file" | sha256sum -c --ignore-missing
```
