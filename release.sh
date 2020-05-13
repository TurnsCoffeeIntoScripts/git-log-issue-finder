#!/bin/bash
# Utility script to simplify release procedure.
# Need to have credentials to both github and dockerhub to use.

GREEN='\033[0;32m'
NO_COLOR='\033[0m'

if [[ $# -ne 2 ]]; then
    echo "Invalid input parameter count. Should be 2"
    echo "  1 --> New version"
    echo "  2 --> Increment type (MAJOR,MINOR,PATCH)"
    exit 1
fi

# Description of input parameters:
# $1 ==> version
# $2 ==> increment type (MAJOR,MINOR,PATCH)

# Increment version in README.me
echo "=================================================================="
echo -e ${GREEN}"Incrementing $2 version in README.md"${NO_COLOR}
case $2 in
    MAJOR)
        sed -r -i 's/(Version:\s)([0-9])\.([0-9])\.([0-9])/echo "\1$((\2+1)).0.0"/ge' README.md
        ;;
    MINOR)
        sed -r -i 's/(Version:\s)([0-9])\.([0-9])\.([0-9])/echo "\1\2.$((\3+1)).0"/ge' README.md
        ;;
    PATCH)
        sed -r -i 's/(Version:\s)([0-9])\.([0-9])\.([0-9])/echo "\1\2.\3.$((\4+1))"/ge' README.md
        ;;
    *)
        echo "Invalid increment type (MAJOR,MINOR,PATCH)"
        ;;
esac

# Change the unreleased to current date in changelog.md
echo "=================================================================="
echo -e ${GREEN}"Updating the 'unrealesed' section of the changelog with current date/version"${NO_COLOR}
sed -r -i "s/(## \[Unreleased\])/echo '\1 \n\n## [$1] - $(date +%Y-%m-%d)'/ge" changelog.md

# Update "version.go"
echo "=================================================================="
echo -e ${GREEN}"Updating the versions constant in pkg/version/version.go with the new version"${NO_COLOR}
case $2 in
    MAJOR)
        sed -r -i 's/(major\s=\s)([0-9]+)/echo "\t\1$((\2+1))"/ge' pkg/version/version.go
        sed -r -i 's/(minor\s=\s)([0-9]+)/echo "\tminor = 0"/ge' pkg/version/version.go
        sed -r -i 's/(patch\s=\s)([0-9]+)/echo "\tpatch = 0"/ge' pkg/version/version.go
        ;;
    MINOR)
        sed -r -i 's/(minor\s=\s)([0-9]+)/echo "\t\1$((\2+1))"/ge' pkg/version/version.go
        sed -r -i 's/(patch\s=\s)([0-9]+)/echo "\tpatch = 0"/ge' pkg/version/version.go
        ;;
    PATCH)
        sed -r -i 's/(patch\s=\s)([0-9]+)/echo "\t\1$((\2+1))"/ge' pkg/version/version.go
        ;;
    *)
        echo "Invalid increment type (MAJOR,MINOR,PATCH)"
        ;;
esac

# Add/Commit/Push the changes
echo "=================================================================="
echo -e ${GREEN}"git add/commit/push of README.md and changelog.md"${NO_COLOR}
git add README.md
git add changelog.md
git add pkg/version/version.go
git commit -m"Incrementing version in doc and version.go ($1)"
git push

# Building docker image
echo "=================================================================="
echo -e ${GREEN}"Building Docker image ($1) with following build arguments:"${NO_COLOR}
echo -e ${GREEN}"\tVERSION=$1"${NO_COLOR}
docker build -t turnscoffeeintoscripts/git-log-issue-finder:$1 \
    --build-arg VERSION=$1 \
    .

if [[ $? -eq "0" ]]; then
    # Pushing docker image
    echo "=================================================================="
    echo -e ${GREEN}"Pushing Docker image ($1) to docker hub"${NO_COLOR}
    docker push turnscoffeeintoscripts/git-log-issue-finder:$1
else
    exit 1
fi

if [[ $? -eq "0" ]]; then
    # Tag in git when the image was successfully push
    echo "=================================================================="
    echo -e ${GREEN}"Creating/pushing tag $1 for git repository"${NO_COLOR}
    git tag -a $1 -m "Tagging version $1"
    git push origin $1
else
    exit 1
fi
