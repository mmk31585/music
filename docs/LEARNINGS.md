# Muse — Learnings Log

> Append-only log of non-trivial bugs, surprises, and fixes discovered during development.

| Date | Symptom | Cause | Fix | Rule |
|------|---------|-------|-----|------|
| 2026-10-15 | Frontend build fails with "Error parsing JavaScript expression: Unexpected token" in multiple Vue files | Multi-line JavaScript expressions in `@click` handlers (e.g., `@click="var = value; other = value"`) are not valid in Vue templates | Extract multi-line expressions to methods in `<script setup>` and call the method from template (e.g., `@click="handleClick()"`) | Vue template expressions must be single-line; use methods for complex logic |
| 2026-10-15 | Frontend build fails with "Illegal '/' in tags" in SearchOverlay.vue | Text interpolation `{{ suggestion }}` was placed as an attribute on the button element instead of as child content | Move `{{ suggestion }}` between opening and closing button tags | Vue template syntax: content goes between tags, attributes go on opening tag |
| 2026-10-15 | Frontend build fails due to template syntax errors in 5 files | Multiple Vue files had multi-line expressions in event handlers | Fixed all 5 files by extracting to methods | Run `npm run build` early to catch template syntax errors || 2026-10-15 | Import cycle between auth and permissions packages | Auth middleware imported permissions for role hierarchy lookup, permissions imported auth for auth middleware | Refactored to use hierarchy fetcher function pattern and context keys | Break import cycles by using function interfaces and context keys instead of direct package imports |
