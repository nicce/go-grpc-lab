#!/bin/bash
#
# Examines the commits that have been added, since the latest tag, and
# calculates what the next semantic version should be based on the
# conventional commits.
#

# Define the valid message format regex.
VALID_MESSAGE_FORMAT="^([a-z]+)(\([a-z]+\))?(!)?: .*$"

# Define the valid tag format regex.
VALID_TAG_FORMAT="^v([0-9]+)\.([0-9]+)\.([0-9]+)$"

# Extract information about the latest existing tag.
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null)

# Validate the latest tag
if [[ -n "$LATEST_TAG" && ! $LATEST_TAG =~ $VALID_TAG_FORMAT ]]; then
    exit 1
fi

# Fetch commits based on the latest tag or from the beginning if no tag is available
COMMITS=$(git log ${LATEST_TAG:+$LATEST_TAG..HEAD} --no-merges --pretty=format:"%H")

TAG_MAJOR=${BASH_REMATCH[1]:-0}
TAG_MINOR=${BASH_REMATCH[2]:-0}
TAG_PATCH=${BASH_REMATCH[3]:-0}

# Track the changes that have been made.
SHOULD_BUMP_MAJOR=false
SHOULD_BUMP_MINOR=false
SHOULD_BUMP_PATCH=false

# Check each commit to calculate the next semantic version.
for COMMIT in $COMMITS; do
    # Get the commit message of the current commit.
    MESSAGE=$(git log -1 "${COMMIT}" --pretty=format:"%s")

    # Validate message and collect information about the current commit.
    if [[ "$MESSAGE" =~ $VALID_MESSAGE_FORMAT ]]; then
        COMMIT_TYPE="${BASH_REMATCH[1]}${BASH_REMATCH[2]}"
        COMMIT_BREAKING_CHANGE=${BASH_REMATCH[3]}

        # Breaking changes trumps all other commits.
        if [[ -n "$COMMIT_BREAKING_CHANGE" ]]; then
            SHOULD_BUMP_MAJOR=true

            break
        fi

        # Signal that one or more new features have been detected.
        if [[ "$COMMIT_TYPE" =~ ^feat ]]; then
            SHOULD_BUMP_MINOR=true
        fi

        # Signal that one or more new fixes, updated dependencies or performance improvements have been detected.
        if [[ "$COMMIT_TYPE" =~ ^fix ]] || [[ "$COMMIT_TYPE" =~ ^build\(deps\) ]] || [[ "$COMMIT_TYPE" =~ ^perf ]]; then
            SHOULD_BUMP_PATCH=true
        fi
    fi
done

# Decide what the next semantic version should be.
if [[ $SHOULD_BUMP_MAJOR = true ]]; then
    echo "v$((TAG_MAJOR + 1)).0.0"
elif [[ $SHOULD_BUMP_MINOR = true ]]; then
    echo "v$TAG_MAJOR.$((TAG_MINOR + 1)).0"
elif [[ $SHOULD_BUMP_PATCH = true ]]; then
    echo "v$TAG_MAJOR.$TAG_MINOR.$((TAG_PATCH + 1))"
fi
