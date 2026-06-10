# ADR-0013: Use Player-Themed Semantic Versioning

Date: 2026-06-10

## Status

Accepted

## Context

The project uses Semantic Versioning (MAJOR.MINOR.PATCH). Purely numeric tags
are accurate but forgettable. The project is football-themed and part of a
cross-language comparison set where each repo adopts a different
football-domain naming convention. Several well-known projects use alphabetical
codename conventions (Ubuntu, macOS).

Release tags follow the format `v{MAJOR}.{MINOR}.{PATCH}-{player}`, where the
player name is drawn from the fixed list below, assigned A→Z sequentially:

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

## Decision

Every release tag appends an alphabetically ordered footballer surname to the
Semantic Version. Format: `v{MAJOR}.{MINOR}.{PATCH}-{player}`. Names are drawn
from the fixed list above, assigned A→Z sequentially. The CD pipeline validates
the name before publishing.

## Consequences

**Positive:**

- Release names are memorable and human-friendly.
- Alphabetical ordering provides an implicit sequence visible in `git tag`.
- The naming scheme is deterministic — the next name is always the next letter.
- Reinforces the football theme across the cross-language comparison set.

**Negative:**

- Non-standard tag format; may confuse new contributors unfamiliar with the
  convention.
- The list is finite — 26 slots before the sequence must restart or be
  extended.
- CD validation adds a small amount of pipeline complexity.

**Neutral:**

- The full list is published in `CHANGELOG.md`; the current position is always
  the last released tag.
