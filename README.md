# go-github-release
git tool to cut a releaser on github with assets

Tool create a new release at Github for a given owner and repos using a Github Personal Token and the upload assets from a target directory to this release. Finally we give the option to update the message body or publish the release.

## global install

```
go install github.com/dearing/go-github-release@latest
```
>[!NOTE]
>With Go 1.24+, we can use the tool feature to pin the tool in go.mod instead of installing it in the host's path
## go tool usage
```
go get -tool github.com/dearing/go-github-release@latest
go tool go-github-release --version
go tool go-github-release --init-config
go tool go-github-release
```
## tool maintenance tips
```
      pin => go get -tool github.com/dearing/go-github-release@v1.0.1
   update => go get -tool github.com/dearing/go-github-release@latest
downgrade => go get -tool github.com/dearing/go-github-release@v1.0.0
uninstall => go get -tool github.com/dearing/go-github-release@none
```
---

## usage 

```
Usage: [go tool] go-github-release [options]
  -asset-dir string
        directory to search for assets to upload (default "build")
  -generate-release-notes
        generate release notes
  -github-owner string
        github owner
  -github-repo string
        github repo
  -release-body string
        release body contents (markdown supported) (default "initial release")
  -release-draft
        draft (default true)
  -release-name string
        release name (default "v0.0.1")
  -release-prerelease
        prerelease
  -tag-name string
        release tag name (default "v0.0.1")
  -target-commitish string
        target commitish tie the release to (default "main")
  -version
        print version and exit
```