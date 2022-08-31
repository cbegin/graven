Deprecated

# graven

Graven is a build management tool for Go projects. It takes light
cues from projects like Maven and Leiningen, but given Go's much
simpler environment and far different take on dependency management,
little is shared beyond the goals.

Want to know more? Read about the [Motivation for Graven](docs/motivation.md).

## Prerequisites

Graven currently requires the following tools to be on your path:

* `go` - the Go build tool, used to compile and test your application.
* `git` - used during the release process to validate the state of your repo,
and tag your repo.
* `docker` - used during build and release of docker images

Of course if you don't plan to use the `release` command commands, you
can still use `graven` just for building, testing and packaging, and thus
would only require the `go` tool.

## Installation

For greater consistency and stability, it is recommended to download a specific release version from:

* https://github.com/cbegin/graven/releases

## Workflow Example

Whether you're starting an entirely new project, or working with existing source,
the workflow should be the same.

Not all of these steps are always necessary. See below for a description of
implied workflow dependencies.

```bash
# Once per project, run the init command and modify
# default project.yaml with relevant names and repos etc.
$ cd ./some/working/directory
$ graven init
$ vi project.yaml

# Typical development cycle
$ graven clean
$ graven build
$ graven test
$ graven package

# When you're ready to release
$ graven repo login --name github
Please type or paste a github token (will not echo):
$ graven release
$ graven bump [major|minor|patch|QUALIFIER]
```

A typical development cycle looks like the following diagram. The `init` command
is run once per project, then `clean`, `build`, `test` and `package` are typically used
throughout the development cycle. Releases occur less frequently, and versions
are bumped after the release.

```
                               +----------+
                             > |  build   | \
                            /  +----------+  \
                           /                  \
                          /                    v
 +----------+    +----------+                 +----------+
 |   init   |--->|   clean  |                 |   test   |
 +----------+    +----------+                 +----------+
                           ^                   /
                            \                 /
                             \ +----------+  /
                              \| package  |<-
                               +----------+
                                     |
                                     v
                               +----------+
                               |  release |
                               +----------+
                                     |
                                     |
                               +-----v----+
                               |   bump   |
                               +----------+
```

## project.yaml example

**Before you have flashbacks of Maven POMs...** note that this is probably the
biggest project.yaml file you'll ever see. Most projects will have shorter,
simpler project.yaml files, and you'll rarely have to modify them once initialized.

The following documented structure (derived from this
very project) will help better understand what you can do with it.

```yaml
# Name, initially derived from parent directory.
name: graven
# Version is typically managed with the graven bump command.
version: 0.6.6
# You can specify a required go version, supporting ranges.
go_version: ">=1.9.1"
# List of artifacts. Each artifact builds one or more packages into executables,
# and provides compiler flags, environment variables, and resources.
# Upon initialization, an artifact config will be generated for darwin, linux and
# windows. You can safely delete any artifacts you don't need. The classifier
# will be used in artifact names and target build directory names.
artifacts:
- classifier: darwin
  # Each target is a package/executable combo with compiler flags and environment variables.
  # Upon initialization, a target will be created for each main package found in the path.
  # You can safely delete any you don't want.
  targets:
  - executable: bin/graven
    package: .
    flags: []
    env: {}
  # archive supports zip and tar.gz (or tgz, as a single dot alias)
  archive: tgz
  resources: []
  # Environment variables, can be set at project, artifact and target level.
  env:
    GOARCH: amd64
    GOOS: darwin
- classifier: linux
  targets:
  - executable: bin/graven
    package: .
    flags: []
    env: {}
  archive: tar.gz
  resources: []
  env:
    GOARCH: amd64
    GOOS: linux
- classifier: win
  targets:
  - executable: graven.exe
    package: .
    flags: []
    env: {}
  archive: zip
  resources: []
  env:
    GOARCH: amd64
    GOOS: windows
# Resources will be included in the packaged archive. Can be overridden at
# artifact level.
resources:
- LICENSE
# Configures a repository for deployment.
# Use graven repo login --name [name] to authenticate
repositories:
  github:
    url: https://api.github.com/
    group: cbegin
    artifact: graven
    type: github
  artifactory:
    url: http://localhost:8081/artifactory/releases/
    group: cbegin
    artifact: graven
    type: maven
  nexus:
    url: http://localhost:8082/nexus/content/repositories/releases/
    group: cbegin
    artifact: graven
    type: maven
  docker:
    url: docker.io
    group: cbegin
    artifact: graven
    type: docker
    file: Dockerfile
```

## Name and Version

```
name: graven
version: 0.6.6
```

The name of your project will initially be set by the parent directory in which it was created when the
`init` command was run. You can change it to whatever you like, but it's recommended to keep it simple,
short and alphabetic. It may be used for generating other values.

The version of your project follows a minimalist set of semantic version practices. That being:

```
M.m.p-Q
```
* M: Major version. Incremented when the software changes significantly and typically in incompatible ways.
* m: Minor version. Incremented when new features are added, and backward compatibility is maintained.
* p: Patch version. Incremented when bugs are fixed. Backward compatibility is typically maintained, but is
* sometimes unavoidably broken.
* Q: Qualifier. This is used to qualify a pre-release build and is typically something like RC1, DEV or TEST.

## Artifacts

```
artifacts:
- classifier: darwin
  targets:
  - executable: bin/graven
    package: .
    flags: []
    env: {}
  archive: tgz
  resources: []
  env:
    GOARCH: amd64
    GOOS: darwin
- ...
```
Each entity listed in the artifacts section represents a distributable artifact that consists of
one or more executables built from a number of packages compiled with specified flags and environment
variables. The resulting binaries are packaged up in an archive format (e.g. zip or tar.gz) and
additional resources can be included in the archive, specified in the resources array. Both resources
and environment variables can be specified at higher levels to avoid duplication. More specific
environment variables will override broader scoped ones. The classifier specifies a suffix for the
artifact that is usually used to indicate the target platform, but can be used to indicate anything
that differs among distributable artifacts.

## Repositories

```
repositories:
  artifactory:
    url: http://localhost:8081/artifactory/releases/
    group: cbegin
    artifact: graven
    type: maven
  ...
```

A repository is where this project will be released. 

Three repository types are currently supported: Github, Docker and Maven (including Nexus and Artifactory).
Their capabilities and settings are summarized in the table below.

| Field | Github | Docker | Maven |
|-------|--------|--------|-------|
| type | github | docker | maven |
| url   | Github API URL | Docker registry URL | Maven release URL |
| group | Owner | Repository | Group ID |
| artifact | Repo | Image Name | Artifact ID |
| file | unused | Dockerfile | unused |

### Authenticating

In order to use a repository for releases, you'll need to authenticate.
To do so, simply call `graven repo login --name [repo-name]`. The credentials you enter
will be stored in your home directory in the .graven.yaml file. Your credentials will be
obfuscated to discourage over-the-shoulder or casual exposure. Even though a strong
encryption algorithm is used, the key is not secure, and thus you should treat it accordingly.


