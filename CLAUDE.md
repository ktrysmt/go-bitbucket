# go-bitbucket conventions

## Testing policy

Before declaring a task complete, run all three suites:

1. `make test/unit` (or `make test/unit-short` if no network) — package-level
   unit tests in the repo root.
2. `make test/mock` — gomock-based interface tests under `mock_tests/`.
3. `make test/swagger` — contract tests under `tests/` against a Prism mock
   server bound to `:4010`.

Never substitute `make test/ci` (which only runs build + unit-short + mock)
for the full set. Swagger contract coverage is mandatory.

### Running the swagger suite

The Prism mock server is launched via Docker (per `tests/README.md:17`):

```sh
docker run --rm -it -p 4010:4010 stoplight/prism:3 mock -h 0.0.0.0 https://bitbucket.org/api/swagger.json
```

Then in a separate shell:

```sh
make test/swagger
```

If `docker run` does not come up (connection refused, daemon errors), verify
that OrbStack is running before retrying:

```sh
docker info >/dev/null 2>&1 || open -a OrbStack
```

Re-run the `docker run` command once OrbStack reports a healthy daemon.

## Release policy

Before merging any PR and before cutting a release with
`gh release create ... --generate-notes`, the full test set above
(unit + mock + swagger) must pass locally on the head commit. Do not
rely on a green CI alone — `make test/ci` skips swagger contract
coverage. Concretely:

1. Check out the branch (or `master` for a release tag) at the exact
   commit being merged or tagged.
2. Run unit, mock, and swagger as documented above. Start the Prism
   container first; do not skip swagger because Prism is not running.
3. Only after all three are green: `gh pr merge --squash` (per PR) or
   `gh release create vX.Y.Z --generate-notes`.

If a swagger run surfaces a regression that is in scope of the PR or
release, fix it before merging or tagging. Pre-existing swagger
failures unrelated to the change should be called out in the PR
description, not silently shipped.

## Output language

Follow `~/.claude/CLAUDE.md`: chat replies in Japanese, code/identifiers in
English, LLM-facing artifacts (this file included) in English.
