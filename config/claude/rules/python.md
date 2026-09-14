---
paths:
  - "**/*.py"
  - "**/pyproject.toml"
---

# Python

- `uv` for everything: `uv init`, `uv add`, `uv run`, `uv sync`. Never `pip install`
  into an ambient interpreter, never `requirements.txt`.
- `ruff` is both the linter and the formatter. `ty` for type checking.
- Annotate every public signature. Modern syntax: `str | None`, `list[str]`, no `typing.List`.
- Dataclasses for data you own, `pydantic` only where you are parsing untrusted input.
- `pathlib`, not `os.path`. f-strings, not `%` or `.format`.
- Catch the narrowest exception that can occur. A bare `except:` is a bug.
- `pytest` with plain `assert`. Parametrize instead of writing the same test five times.
