#!/bin/bash
#
# Validates a commit message against a sub-set of the Conventional Commits
# specification https://www.conventionalcommits.org
#
# Supported features
#  - list of commit types
#  - optional scope
#  - breaking changes are declared using the ! character in the commit subject
#
# Example of valid commit messages
#  - feat: add something
#  - fix(docs): typo in README
#  - feat(api)!: change the parameter order in the api
#  - fix(api): fixed broken import (#123)
#
# Example of invalid commit messages
#  - did some work
#  - feat:add support for something
#  - hack: did my own thing
#  - FEAT: A COOL NEW FEATURE
#

# Exit immediately if a command exits with a non-zero status.
set -e

# Define the list of supported commit types.
SUPPORTED_TYPES="build|chore|ci|docs|feat|fix|perf|refactor|style|test"

# Define the valid message format regex.
VALID_MESSAGE_FORMAT="^($SUPPORTED_TYPES)(\([a-z]+\))?!?: [a-zA-Z0-9\. \/_@()#-]+$"

# Validate message.
if ! [[ "$1" =~ $VALID_MESSAGE_FORMAT ]]; then
    printf "❌ Invalid commit message\n%s" "$1"

    exit 1
fi

