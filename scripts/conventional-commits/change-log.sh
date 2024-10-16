#!/bin/bash
#
# Generates a change log, in markdown format, of the changes that
# has been added between the latest and previous tags.
# Use PREVIOUS_TAG to generate change log between latest and PREVIOUS_TAG.
#
# Example usage:
#  ./scripts/change-log.sh
# or
#  PREVIOUS_TAG=<previous-tag> ./scripts/change-log.sh

# Exit immediately if a command exits with a non-zero status.
set -e

# Extract information about the latest tag.
LATEST_TAG=$(git describe --tags --always --abbrev=0)

# Extract information about the previous tag if it wasn't provided.
if [[ -z "$PREVIOUS_TAG" ]]; then
    PREVIOUS_TAG=$(git describe --tags --always --abbrev=0 "$LATEST_TAG^")
fi

# Regular expression for valid tag (e.g. v1.0.0).
VALID_TAG_FORMAT="^v[0-9]+(\.[0-9]+){2}$"

# Check if the variable is in a valid tag format
if [[ ! $PREVIOUS_TAG =~ $VALID_TAG_FORMAT ]]; then
    echo "PREVIOUS_TAG invalid format: $PREVIOUS_TAG"
    echo "Expected tag in format that fulfills regexp $VALID_TAG_FORMAT"
    exit 1
fi

# Define supported conventional commit types with corresponding emojis.
types=(
    'build|🏗️️'
    'ci|🤖'
    'chore|🧹'
    'docs|📝'
    'feat|🚀'
    'fix|🛠️'
    'perf|⚡'
    'refactor|🔧'
    'style|💄'
    'test|✅'
)

echo "# Changelog #"

# Loop through each type and retrieve commit SHAs.
for fields in "${types[@]}"; do
    # Split the type and emoji by the '|' delimiter.
    IFS=$'|' read -r TYPE EMOJI < <(echo "$fields")

    # Retrieve the commit SHAs for the commits that have been added since the previous tag was created
    # up to the latest tag, and match the commit message against the current type being processed.
    COMMITS=$(git log "$PREVIOUS_TAG".."$LATEST_TAG" --no-merges --pretty=format:"%H" --grep="^$TYPE")

    # If there are no commits for this type, skip to the next one.
    if [[ -z "$COMMITS" ]]; then
        continue
    fi

    # Output a header for this type with its corresponding emoji.
    echo "## $EMOJI $TYPE ##"

    # Loop through each commit and output its message with a bullet point.
    for COMMIT in $COMMITS; do
        MESSAGE=$(git log -1 "$COMMIT" --pretty=format:"%s (%h)")

        echo -n "- "

        # If the commit contains a breaking change, add a warning emoji to the output.
        if [[ "$MESSAGE" =~ !: ]]; then
            echo -n "⚠️ "
        fi

        # Output the commit message.
        echo "$MESSAGE"
    done
done
