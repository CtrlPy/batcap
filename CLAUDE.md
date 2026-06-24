# Git

Use Angular Conventional Commits for all commit messages.
Follow GitHub Flow — work on main for small changes, feature branches for larger ones.

# Releases & Tagging

* **Never forget to release a new version tag** after merging features or key fixes:
  1. Determine the next Semantic Version bump (e.g. `v1.4.1`).
  2. Update the version in `README.md`'s `go install` instructions.
  3. Commit and push the version update.
  4. Create and push the annotated tag (`git tag -a v1.x.x -m "..." && git push origin v1.x.x`).
