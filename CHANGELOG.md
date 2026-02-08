# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Famous Football Player Names 🎖️

This project uses famous football player names (A-Z) as release codenames:

| Letter | Player Name | Country/Era | Tag Name |
| ------ | ----------- | ----------- | -------- |
| A | Ademir | Brazil | `ademir` |
| B | Bobby (Moore) | England | `bobby` |
| C | Cafu | Brazil | `cafu` |
| D | (Alfredo) Di Stéfano | Argentina/Spain | `distefano` |
| E | Eusébio | Portugal | `eusebio` |
| F | Franz (Beckenbauer) | Germany | `franz` |
| G | Garrincha | Brazil | `garrincha` |
| H | (Thierry) Henry | France | `henry` |
| I | (Filippo) Inzaghi | Italy | `inzaghi` |
| J | Jairzinho | Brazil | `jairzinho` |
| K | (Roy) Keane | Ireland | `keane` |
| L | (Frank) Lampard | England | `lampard` |
| M | (Diego) Maradona | Argentina | `maradona` |
| N | Nilton (Santos) | Brazil | `nilton` |
| O | (Jay-Jay) Okocha | Nigeria | `okocha` |
| P | Pelé | Brazil | `pele` |
| Q | (Fabio) Quagliarella | Italy | `quagliarella` |
| R | Romário | Brazil | `romario` |
| S | (Paul) Scholes | England | `scholes` |
| T | (Francesco) Totti | Italy | `totti` |
| U | Uwe (Seeler) | Germany | `uwe` |
| V | (Patrick) Vieira | France | `vieira` |
| W | (Ian) Wright | England | `wright` |
| X | Xabi (Alonso) | Spain | `xabi` |
| Y | (Lev) Yashin | USSR | `yashin` |
| Z | Zico | Brazil | `zico` |

---

## [Unreleased]

### Added

### Changed

### Deprecated

### Removed

### Fixed

### Security

---

## [1.0.0 - Ademir] - 2026-02-06

Initial release. See [README.md](README.md) for complete feature list and documentation.

---

## How to Release

To create a new release, follow these steps in order:

### 1. Update CHANGELOG.md

Move items from the `[Unreleased]` section to a new release section:

```markdown
## [X.Y.Z - PLAYER_NAME] - YYYY-MM-DD

### Added
- New features here

### Changed
- Changes here

### Fixed
- Bug fixes here

### Removed
- Removed features here
```

**Important:** Commit and push this change before creating the tag.

### 2. Create and Push Version Tag

```bash
git tag -a vX.Y.Z-player -m "Release X.Y.Z - Player"
git push origin vX.Y.Z-player
```

Example:

```bash
git tag -a v1.0.0-ademir -m "Release 1.0.0 - Ademir"
git push origin v1.0.0-ademir
```

### 3. Automated CD Workflow

The CD workflow automatically:

- ✅ Validates the player name against the A-Z list
- ✅ Builds and tests the project
- ✅ Publishes Docker images to GHCR with three tags (`:X.Y.Z`, `:player`, `:latest`)
- ✅ Creates a GitHub Release with auto-generated notes from commits

### Pre-Release Checklist

- [ ] CHANGELOG.md updated with release notes
- [ ] CHANGELOG.md changes committed and pushed
- [ ] Tag created with correct format: `vX.Y.Z-player`
- [ ] Player name is valid (A-Z from table above)
- [ ] Tag pushed to trigger CD workflow

---

<!-- Template for new releases:

## [X.Y.Z - PLAYER_NAME] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security vulnerability fixes

-->

---

[unreleased]: https://github.com/nanotaboada/go-samples-gin-restful/compare/v1.0.0-ademir...HEAD
[1.0.0 - Ademir]: https://github.com/nanotaboada/go-samples-gin-restful/releases/tag/v1.0.0-ademir
