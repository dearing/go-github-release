# go-github-release
git tool to cut a releaser on github with assets

Tool create a new release at Github for a given owner and repos using a Github Personal Token and the upload assets from a target directory to this release. Finally we give the option to update the message body or publish the release.

## global install

```
go install github.com/dearing/go-github-release
```

## go tool install
```
go get -tool github.com/dearing/go-github-release@latest
```

## usage

```

```

# dev notes

- https://api.github.com/repos/OWNER/REPO/releases
- https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#create-a-release

```

curl -L \
  -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer <YOUR-TOKEN>" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/OWNER/REPO/releases \
  -d '{"tag_name":"v1.0.0","target_commitish":"master","name":"v1.0.0","body":"Description of the release","draft":false,"prerelease":false,"generate_release_notes":false}'

```