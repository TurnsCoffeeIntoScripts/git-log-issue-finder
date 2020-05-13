# G.L.I.F. script examples

## Content
1. [Diff between the two latest release versions](#script1)
2. [Diff between the two latest release candidates](#script2)
3. [Diff between the two latest builds](#script3)
4. [Diff between latest build and fixed tag](#script4)

Consider the following versions (from earliest to latest):
1. 1.0.0
2. 1.0.0-rc.1
3. 1.0.1
4. 1.1.0
5. 1.1.0-rc.1
6. 1.1.0-rc.2
7. 1.2.0
8. vrelease_13
9. 1.4.0
10. 1.4.0-build.1
11. 1.4.0-rc.1
11. 1.4.0-build.2  
12. 1.4.0-rc.2
13. 1.4.0-rc.3
14. 1.4.0-build.3
15. 1.4.0-build.4
16. 1.4.0-build.5

## <a name="script1" href="script1">Diff between the two latest release versions</a>
This script will perform a difference between the git logs of the two latest release versions.
The format of the versions in this case follows the semantic versionning.

The script will perform the diff from version 1.2.0 to version 1.4.0.
```
set repopath "."

let repo = initRepo();
let version = "$.$.$";

extractTags(repo, version);

let from = getLatestTag(repo, 1);
let to = getLatestTag(repo, 0);

diff(repo, from, to)
```

## <a name="script2" href="script2">Diff between the two latest release candidates</a>
This script will perform a difference between the git logs of the two latest release candidates. 
The format of the versions in this case follows the semantic versionning.

The script will perform the diff from version 1.4.0-rc.2 to 1.4.0-rc.3.
```
set repopath ".";
print("Using repo path: " + whichRepo());
let repo = initRepo();
let version="$.$.$-rc.$";

extractTags(repo, version);

let from = getLatestTag(repo, 1);
let to = getLatestTag(repo, 0);

diff(repo, from, to);
```

## <a name="script3" href="script3">Diff between the two latest builds</a>
This script will perform a difference between the git logs of the two latest builds.
The format of the versions in this case follows the semantic versionning.

The script will perform the diff from version 1.4.0-build.4 to 1.4.0-build.5.
```
set repopath ".";
print("Using repo path: " + whichRepo());
let repo = initRepo();
let version="$.$.$-build.$";

extractTags(repo, version);

let from = getLatestTag(repo, 1);
let to = getLatestTag(repo, 0);

diff(repo, from, to);
```

## <a name="script4" href="script4">Diff between latest build and fixed tag</a>
This script will perform a difference between the git logs of the latest build and a fixed tag.

The script will perform the diff from version 1.4.0 to 1.4.0-build.5.
```
set repopath ".";
print("Using repo path: " + whichRepo());
let repo = initRepo();
let version="$.$.$-build.$";

extractTags(repo, version);

let from = getTag(repo, "1.4.0");
let to = getLatestTag(repo, 0);

diff(repo, from, to);
```