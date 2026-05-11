# AGENT DIRECTIVES

- Always consult `ROADMAP.md` and `TODO.md` before picking the next task.
- Be sure to `export PROTOC=/usr/bin/protoc` and `export PKG_CONFIG_PATH=/usr/lib/x86_64-linux-gnu/pkgconfig` before running `./script/presubmit` or `cargo test`.
- All `npm`/`yarn` commands should use `corepack enable` if there are yarn hash mismatch issues.
- Do NOT use `git push` directly as it's blocked by the environment. Use the provided submission tools for handoff.
- Make sure to update `CHANGELOG.md` and bump the version in `VERSION.md` on major milestones.
- Remember that `Screenshot::data` is an `Arc<[u8]>` (as of latest upstream), not a `Vec<u8>`. Use `.into()` to convert from `Vec<u8>` to `Arc<[u8]>`, and `.to_vec()` to extract the vector.
