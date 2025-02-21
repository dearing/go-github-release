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
Simple GitHub tool to cut a release and upload assets.

NOTE: An env:GITHUB_TOKEN is required to authenticate with GitHub and its
recommended to use a personal access token scoped to repository.

GitHub releases are fairly simple, its a ccollection of assets with name and
message body that reference a commit, which can be a newly created tag as a
commitish (--target-commitish) or existing tag in which case the target
commitish is ignored. Assets are uploaded to this release based on mime-type
which is inferred from the file extension.

If you are unsure of what workflow to use, the author recommends creating a
tag ahead of release and using that tag as the --tag-name. This way you can
seperate the release process from the tag creation process and ensure your
tag is referencing the proper moment in history. The automatic tag creation
feature is a GitHub conveinence and not related to git.

ex: to draft a new release for an existing tag v1:

[go tool] go-github-release \
        --github-owner myorg \
        --github-repo myrepo \
        --tag-name v1 \
        --asset-dir build

ex: to draft a new release and new tag v1 pointing to abc123:

[go tool] go-github-release \
        --github-owner myorg \
        --github-repo myrepo \
        --tag-name v1 \
        --target-commitish abc123 \
        --asset-dir dist

Tips:
  - drafts that create a tag, won't do so until the release is published
  - if --github-repo is not set, the current directory is used
  - if --release-name is not set, the tag name is used
  - if --target-commitish is not set, main is used
  - a release body supports markdown
  - you can safely edit the release body while the assets are uploading
  - find the latest by name tag: git describe --tags --abbrev=0
  - find the a tag's reference commit: git rev-list -n 1 tags/v0.0.1

Options:

  -asset-dir string
        directory to search for assets to upload (default "build")
  -generate-release-notes
        tell GitHub to generate release notes
  -github-owner string
        github owner
  -github-repo string
        github repo
  -release-body string
        release body contents (markdown supported) (default "no results")
  -release-draft
        set draft flag on release (default true)
  -release-name string
        release name, defaults to tag name
  -release-prerelease
        set prerelease flag on release
  -tag-name string
        release tag to reference or create
  -target-commitish string
        when creating a tag, the commit to reference
  -version
        print version and exit
```