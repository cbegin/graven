# Freezing Dependencies

Graven's freeze and unfreeze commands are a mechanism to avoid having to
commit a massive vendor directory to your git repository, while at the 
same time protecting you from losing access to key library dependencies
in the event of a leftpad scenario.

## Option 1: Committing the .freezer directory

This may seem odd to many idiomatic Go programmers.
And others will trigger on concerns surrounding "binaries in Github". 
Graven takes steps to ensure that the frozen files are as efficient for Git
as possible. It does not compress them, ensuring that Git's binary diff 
algorithm and ultimate compression (packing) will handle the files efficiently.
So you can check in the .freezer directory with minimal impact to your repository, 
and keep it cleaner than if you had committed potentially tens of thousands of 
lines of vendored code. 

Still worried about binaries in Git? See this video for more: 

* https://www.youtube.com/watch?v=rALm7BCCY-0#t=6m28s

## Option 2: Use Maven Repo like Nexus or Artifactory

You don't need to commit the .freezer directory. If you use a
Maven compatible repository such as Nexus or Artifactory, Graven will 
translate the Go import path and revision into maven coordinates and 
can push and pull those dependencies from your own private repo, to ensure
you're never leftpadded. In this case the .freezer directory acts like a 
cache and can be .gitignored. 

## Do I have to freeze my dependencies to use Graven?

No.

The freeze and unfreeze commands are entirely optional and separate 
from primary graven workflow (see below). So if you enjoy cluttering your repo 
with thousands of lines of third party code (likely more than your own app has) 
in your vendor directory, you can continue to do so! Graven has no core
dependency on the freeze and unfreeze features.

### But if you choose to try them out...

Graven's freeze and unfreeze commands understand three vendor file formats.
It automatically selects the vendor file to use in priority order as follows:

* Govendor: `vendor/vendor.json`
* Glide: `glide.lock`
* Dep: `Gopkg.lock`

## Usage

```
$ graven freeze
$ graven unfreeze
```

## Diff Note

Graven does as much as possible to ensure that there aren't any noisy diffs in
the frozen archives. But once in a while a file date or other metadata will
cause the dependency to change. Watch your editor and operating system for
such consequences (like on Mac, simply navigating the directory and dropping
a .DS_Store file).

