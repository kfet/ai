# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial extraction from fir (`github.com/kfet/fir`) Phase 5. Portable
  AI primitives: message/content types, `Tool`, `Usage`, `Context`,
  `Model`, streaming events, `StreamFunction`, and retry classification
  in the root `ai` package; `ratelimit/`, `overflow/`, and `jsonparse/`
  subpackages. Portions ported from pi-mono (MIT, Copyright (c) 2025
  Mario Zechner).
