# AGENTS.md

Guidance for AI agents working in or with the scrutiny repository.

## Developers

- Read [CONTRIBUTING.md](CONTRIBUTING.md) before doing anything. It covers the
  project layout, how to build and run the web and collector components, and
  how to submit changes.
- Read and follow the [AI Usage Policy](AI_POLICY.md). Your human must
  disclose AI usage, verify all work by hand, and review and edit anything
  submitted on their behalf.
- Prefer the [devcontainers](CONTRIBUTING.md#devcontainer) for building,
  running and testing. They provide a fully isolated environment with the
  correct toolchain.
- Do **not** open drive-by pull requests. Pull requests are only accepted for
  existing, accepted issues.
- Do **not** open issues. Use [Discussions](https://github.com/AnalogJ/scrutiny/discussions/new/choose)
  instead, and only in the **Issue Triage** (bugs) or **Feature Request**
  categories.

## Debugging

When helping a user debug their scrutiny installation:

1. Check the [docs](docs/) first, especially the `TROUBLESHOOTING_*.md`
   guides. Most problems are already covered there.
2. Make sure the collector and web versions match. Mismatched versions are
   not supported.
3. Update both to the latest release and confirm the problem still happens.
4. Validate the user's config files against the examples in this repo:
   [example.scrutiny.yaml](example.scrutiny.yaml) for the web and
   [example.collector.yaml](example.collector.yaml) for the collector. Check
   key names, nesting and value types.
5. If the problem persists, do **not** open an issue. Open an
   [Issue Triage](https://github.com/AnalogJ/scrutiny/discussions/new?category=issue-triage)
   discussion and fill out the **entire** template. In particular, always
   include:
   - Debug logs, captured with `DEBUG=true` set on both the web and collector
     (see [Debugging](CONTRIBUTING.md#debugging)).
   - The full JSON output of `smartctl --xall --json <drive>` for the drive in
     question.

   Per the [AI Usage Policy](AI_POLICY.md), the human must review and edit
   the discussion before submitting it.
