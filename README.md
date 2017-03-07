_Status: This is a proof of concept. The code was written during about 10 episodes of The 100. I'll post some next steps soon, to bring the code to production level._

# graven

Graven is a build management tool for Go projects. It takes light
cues from projects like Maven and Leiningen, but given Go's much
simpler environment and far different take on dependency management, 
little is shared beyond the goals.

These shared goals are as follows.

## Providing an easy, uniform build system

Go projects are often built with makefiles or bash scripts. As much 
as I like Make, given Go's cross platform capabilities, Make is 
a poor experience for Windows developers. In addition, Make always
suffers from a lack of consistency, and also carries a bit of 
legacy baggage that is overhead that is simply not necessary 
anymore. 

Graven offers a single, consistent artifact: `project.yaml` that
can be used to build a project on Mac, Linux or Windows in a 
consistent and easy way. It also supports creating consistent
archives of deployable executables and any required resources.

The Graven command line interface offers a simple lifecycle 
similar to that which is often implemented in makefiles or
bash scripts: `clean`, `build`, `test`, `package`, `release`.

## Providing project information and guidelines for best practices

For all the greatness that is Go, there is a major gap in practices
surrounding versioning and dependency management. Vendoring has slightly
improved the dependency management, but lacks a single consistent tool.
Furthermore, builds aren't easily repeatable and versions are usually
based on commit hashcodes, rather than intelligently selected semantic
versions that describe capabilities, compatibility and bug fixes. 

Graven supports, automate and encourages proper semantic versioning and 
can freeze vendor dependencies to ensure repeatable builds are possible. 
Graven is opinionated about vendoring tools, and has chosen Govendor as 
its standard. However, it may support other vendoring tools in the future, 
and will embrace any standard tools that eventually come from the Go
project.


## Where things differ

While Graven takes queues from Maven and Leiningen, it also casts out 
the annoying, verbose and repetitive aspects that most developers
agree weigh Maven down. 

So Graven embrases:

* A much simpler build artifact based on light YAML
* Batteries included, no plugins - none are even supported yet, 
which is considered a good thing for now

## TODO

- Unit Tests
- A functional test project suite
- Interfaces around likely integration points (repository support)
- bug in freeze .frozen metadata error
- docker support



