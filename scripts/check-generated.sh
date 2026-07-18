#!/usr/bin/env sh
set -eu
task sync
git diff --exit-code -- internal/generated/fx
