# Motivation for Graven

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

Graven supports, automate and encourages proper semantic versioning. 

## Where things differ

While Graven takes queues from Maven and Leiningen, it also casts out 
the annoying, verbose and repetitive aspects that most developers
agree weigh Maven down. 

So Graven embraces:

* A much simpler build artifact based on light YAML
* Batteries included, no plugins - none are even supported yet, 
which is considered a good thing for now
