---
name: python-pypi-package-builder
description: 'End-to-end skill for building, testing, linting, versioning, and publishing a Python library to PyPI. Covers setuptools, hatchling, flit, poetry backends; PEP 440 versioning; Trusted Publishing (OIDC); and CI/CD.'
---

# Python PyPI Package Builder

A complete guide for building, testing, linting, versioning, and publishing a production-grade Python library to PyPI.

## Quick Navigation
| Section | Covers |
|---------|--------|
| Build Backend Decision | setuptools / hatchling / flit / poetry |
| Folder Structure | src/ vs flat vs monorepo |
| Versioning | PEP 440, semver, dynamic vs static |
| PyPA Flow | Build → Validate → Test → Publish |

## Build Backend Decision
- **setuptools + setuptools_scm**: Git tag versioning, C extensions supported
- **hatchling**: Fast, modern, pure Python (recommended default)
- **flit**: Zero-config for simple packages
- **poetry**: All-in-one dep management + publishing

## Folder Structure
- New projects: `src/` layout (PyPA recommended)
- Small packages: Flat layout
- Monorepo: Namespace packages (PEP 420)

## Versioning
- PEP 440: `N.N.N[a|b|rc]N[.postN][.devN]`
- SemVer: MAJOR.MINOR.PATCH
- Dynamic: `setuptools_scm` from git tags
- Static: `__version__` in `__init__.py`

## Reference Files
- `references/pyproject-toml.md` — Full templates for all 4 backends
- `references/library-patterns.md` — OOP/SOLID, type hints
- `references/testing-quality.md` — pytest, ruff, mypy, pre-commit
- `references/ci-publishing.md` — CI/CD with Trusted Publishing
- `references/versioning-strategy.md` — PEP 440 deep-dive
- `scripts/scaffold.py` — Generate project from template
