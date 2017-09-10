# graven

Graven is a build management tool for Go projects. It takes light
cues from projects like Maven and Leiningen, but given Go's much
simpler environment and far different take on dependency management, 
little is shared beyond the goals.

Want to know more? Read about the [Motivation for Graven](docs/motivation.md)

# Prerequisites

Graven currently requires the following tools to be on your path:

* `go` - the Go build tool, used to compile and test your application.
* `git` - used during the release process to validate the state of your repo, 
and tag your repo.

Of course if you don't plan to use the `release` command commands, you
can still use `graven` just for building, testing and packaging, and thus 
would only require the `go` tool. 

# Installation

If you want to run the latest, you can just `go get` the tool.

```
go get -u github.com/cbegin/graven
```

For greater consistency and stability, you can run a specific relese version from: 

* https://github.com/cbegin/graven/releases

# Workflow Example

Whether your starting an entirely new project, or working with existing source, 
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
$ graven repo --login --name github
Please type or paste a github token (will not echo):
$ graven release
$ graven bump [major|minor|patch|QUALIFIER]
```

A typical development cycle looks like this. The `init` command is run once 
per project, then clean, build, test and package are typically used throughout 
the development cycle. Releases occur less frequently, and versions are bumped 
after the release.

The `freeze` and `unfreeze` commands are optional and on a completely 
independent flow, thus can be executed any time. They are discussed separately
here: [Freezing Dependencies](docs/freezing.md).

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
    +----------+               +----------+                           
    |  freeze  |                     |                                
    +----------+                     v                                
          |                    +----------+                           
          |                    |  release |                           
    +-----v----+               +----------+                           
    | unfreeze |                     |                                
    +----------+                     |                                
                               +-----v----+                           
                               |   bump   |                           
                               +----------+                           
```

# project.yaml

For many projects, only minimal interaction will be needed with `project.yaml`
after initialization. The following documented structure (derived from this
very project) will help better understand what you can do with it. 

```yaml
# Name, initially derived from parent directory. 
name: graven
# Version is typically managed with the graven bump command.
version: 0.6.6
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
    flags: ""
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
    flags: ""
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
    flags: ""
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
# Use graven repo --login --name [name] to authenticate
repositories:
  github:
    url: https://api.github.com/
    group: cbegin
    artifact: graven
    type: github
    # Github repos only support releases
    roles: 
    - release
  artifactory:
    url: http://localhost:8081/artifactory/releases/
    group: cbegin
    artifact: graven
    type: maven
    # Supports both releases and frozen dependencies
    roles: 
    - release
    - dependency
  nexus:
    url: http://localhost:8082/nexus/content/repositories/releases/
    group: cbegin
    artifact: graven
    type: maven
    # Supports both releases and frozen dependencies
    roles: 
    - release
    - dependency
  docker:
    url: docker.io
    group: cbegin
    artifact: graven
    type: docker
    file: Dockerfile
    # Docker repos only support releases
    roles: 
    - release
```
## A Comment about Comments in project.yaml

Currently Graven uses `gopkg.in/yaml.v2` which does not have round trip support for comments
or document structure. Therefore, when Graven rewrites your project file at certain times, 
your comments will be lost. I'll look at resolving this in the near future. Graven probably
doesn't need rich YAML rewriting support, so I can probably get away with a minimal YAML 
parser that preserves structure.