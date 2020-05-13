# G.L.I.F. script documentation

1. [Installing glif](#install)
2. [Glif R.E.P.L. interface](#repl)
3. [Git repository management and glif operations](#grm)
4. [Built-in functions](#bif)

## <a href="install" name="install">Installing glif</a>
To use glif on your machine simply clone this repository: https://github.com/TurnsCoffeeIntoScripts/git-log-issue-finder.git.  
Once done, execute the command `make put`, which will build and copy the executable `glif` in `/usr/bin`.

## <a href="repl" name="repl">Glif R.E.P.L. interface</a>
When glif is built and/or installed you can run `glif` or `glif --repl` to launch the repl
(Read-Evaluate-Print-Loop) to test glif scripts.

## <a href="grm" name="grm">Git repository management and glif operations</a>
Since glif scripts' main purpose are to parse git logs and perform a 'diff' between two specific
point in the git history, it is imperative that those scripts are easily able to manage (read interac with)
git repository found on the local file system.  

The following steps are are necessary for the proper management for a git repository.

### Set the repopath
The first step is to initialize the `repopath` which indicate the location (on the local machine)
of the git repository. By default this value is already set within the interpreter to `.`; the
current directory. 

```
set repopath ".";
```

You can set the repopath to any path on your local machine or environment. Here's two examples. The
first one being a Windows type path, and the second one being a Linux type path.
```
set repopath "C:\Development\repo1\";

set repopath "/home/development/repo1/";
```

### Init the repository object
The glif interpreter contains an internal representation of the git repository. To properly
init this object one needs to call the `initRepo` function. This must be done **after** the `repopath`
has been initiated (if different than `.`).
```
let r = initRepo();
```
This repository object (declared with the variable declaration `let`) will be required for most
of the subsequent calls to the builtin functions.

### Extract the relevant tag(s)
Before performing the diff between tags the interpreter needs to know what to look for.
This is done in a few steps. First, a format must be specified. In this string, the dollar signs
`$` represents numbers. In terms of regex it translates to this: `[0-9]+`. So the next example
defines the format of a standard tag following semantic versionning.  

```
let versionRelease = "$.$.$"; 
```

This next one also follows semantic versionning but with the qualifier of a release candidate:
```
let versionRC = "$.$.$-rc.$";
```
A wide variety of format can be specified. Once this version is specified the repository object
inside the interpreter needs to extract and remember all tags matching this format. This is done
by the `extractTags` function call that takes the repository object as a first parameters and 
the format string as the second parameter.
```
let repo = initRepo();
let version = "$.$.$";
extractTags(repo, version);
```
After that point in the execution, the interpreter will have an order list of all tags matching
the specified format. Note that this list is order with the latest element first. So the index `0`
contains the latest, most recent, tag of the git repository. 

### Finding tags
Once the interpreter knows what to look for and extracted the possibles values, the last things it
will need are a starting point (from) and an ending point (to). There's two possible way of defining
these points. The first one uses the order list mentionned in the previous section and is called with
the `getLatestTag` function that take the repository object as a first parameter and an integer index
as the second parameter. For example these two calls are often used in the predefined scripts:
```
let from = getLatestTag(repo, 1);
let to = getLatestTag(repo, 0);
```
The second parameter, the integer index, tells the interpreter how far back it should look
within its internal ordered list of tags. Therefore, by giving `0` it will take the first element of
the list which is the most recent tag.  
By giving `1` it will take the second element in that list which is the previous latest tag.  

This combination is frequent as it allows to perform a diff between the two most latest element
that will match the specified format.

The second method to define the starting/ending point is to use the `getTag` function. With this function
you can explicitly call to any existing tags in the git repository.
```
let tag = getTag(repo, "1.5.1-rc.1");
let anotherTag = getTag(repo, "custom_tag_v3");
```

### Performing the actual diff
The final and most important function call is the one that will execute the "diff" operation of the
git logs. Assuming we have a script that contains these lines:
```
let repo = initRepo();
let format = "$.$.$";

extractTags(repo, format);

let from = getLatestTag(repo, 1);
let to = getLatestTag(repo, 0);
```
The only thing left to do to get the actual diff is to call the following function which takes the 
repository object and two tags as parameters. 
```
diff(repo, from, to)
```
This will print all the issues found in the form of a slice (array). 
