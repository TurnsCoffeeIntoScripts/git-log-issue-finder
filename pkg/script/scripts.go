package script

// DiffLatestSemverWithLatestBuilds is a predefined script
// It performs a diff between the two latest builds on the latest version (MAJOR.MINOR.PATCH-build.RC)
var DiffLatestSemverWithLatestBuilds = `
set repopath ".";
print("Using repo path: " + whichRepo());
let repo = initRepo();
let version="$.$.$-build.$";

extractTags(repo, version);

let to = getLatestTag(repo, 0);
let from = getLatestTag(repo, 1);

diff(repo, from, to);
`

// DiffLatestSemverWithLatestRCs is a predefined script
// It performs a diff between the two latest RCs on the latest version (MAJOR.MINOR.PATCH-rc.RC)
var DiffLatestSemverWithLatestRCs = `
set repopath ".";
print("Using repo path: " + whichRepo());
let repo = initRepo();
let version="$.$.$-rc.$";

extractTags(repo, version);

let to = getLatestTag(repo, 0);
let from = getLatestTag(repo, 1);

diff(repo, from, to);
`

// DiffLatestSemver is a predefined script
// It performs a diff between the two latest release (production) version (MAJOR.MINOR.PATCH)
var DiffLatestSemver = `
set repopath "."
print("Using repo path: " + whichRepo());
let repo = initRepo();
let version="$.$.$";

extractTags(repo, version);

let to = getLatestTag(repo, 0);
let from = getLatestTag(repo, 1);

diff(repo, from, to);
`
