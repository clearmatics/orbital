#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
ethdir="$workspace/src/gitlab.clearmatics.net"
if [ ! -L "$ethdir/mobius/ringtool" ]; then
    mkdir -p "$ethdir"
    cd "$ethdir"
    mkdir mobius
    cd mobius
    ln -s ../../../../../. ringtool
    cd "$root"
fi

# Set up the environment to use the workspace.
# Also add Godeps workspace so we build using canned dependencies.
GOPATH="$ethdir/mobius/ringtool/Godeps/_workspace:$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$ethdir/mobius/ringtool"
PWD="$ethdir/mobius/ringtool"

# Launch the arguments with the configured environment.
exec "$@"
