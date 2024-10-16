#!/bin/bash
#
# Checks if the commits, on the current branch, contains
# any breaking changes.
#

# Exit immediately if a command exits with a non-zero status.
set -e

# Get the name of the current branch. We're only interested in commits that's
# been added on that specific branch.
BRANCH_NAME=$(git branch --show-current)

# Get the commit SHAs for the commits on the branch, excluding any merge commits
# that might exist.
if [[ -n "$CI" ]]; then
    COMMITS=$(git log --no-merges --format="%H")
else
    COMMITS=$(git log main.."$BRANCH_NAME" --no-merges --format="%H")
fi

# Checks each commit for breaking changes.
for COMMIT in $COMMITS; do
    # Get the commit message of the current commit.
    MESSAGE=$(git log -1 "${COMMIT}" --pretty=format:"%s")

    if [[ "$MESSAGE" =~ !: ]]; then
        echo "⚠️ These changes contain breaking changes"

        break
    fi
done
