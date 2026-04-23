# Repository Guidelines

## Project Structure & Module Organization
`backend/` contains the Go API, schedulers, channel adapters, MCP integration, and GORM models. Keep server entrypoints in `backend/main.go`, HTTP code in `backend/api/`, business logic in `backend/engine/`, and integrations in `backend/ai/`, `backend/channels/`, and `backend/notifications/`. `frontend/` is a Vue 3 + Vite SPA; views live in `src/views/`, shared UI in `src/components/`, state in `src/stores/`, and routing/i18n in `src/router/` and `src/i18n/`. `docs/` hosts the VitePress documentation, `docker/` contains Nginx and SSL assets, and `scripts/` holds release helpers.

## Build, Test, and Development Commands
Use Docker for full-stack local runs:

- `cp .env.example .env && docker compose up -d --build` builds Nginx, app, and MySQL.
- `cd backend && go test ./...` runs backend unit and integration tests.
- `cd frontend && npm install && npm run dev` starts the SPA on port `3000` with `/api` proxied to `localhost:8080`.
- `cd frontend && npm run build` type-checks with `vue-tsc` and creates the production bundle.
- `cd docs && npm install && npm run docs:dev` serves the docs site locally.

## Coding Style & Naming Conventions
Follow Go defaults: run `gofmt`, keep package names lowercase, and use `_test.go` for tests. In the frontend, use Vue SFCs with `<script setup lang="ts">`, PascalCase component filenames such as `OnboardingWizard.vue`, and camelCase for stores and helpers. Match existing file naming by feature, for example `frontend/src/views/Jobs/JobDetail.vue`. Preserve strict TypeScript settings in `frontend/tsconfig.app.json`; do not add unused locals or parameters.

## Testing Guidelines
Backend tests sit beside implementation files, for example `backend/engine/analyzer_test.go`; prefer table-driven Go tests where practical. Frontend tests use Vitest under `frontend/src/__tests__/` with `*.spec.ts` names. Run the relevant suite before opening a PR, and rebuild with `docker compose up -d --build` when changes affect runtime wiring, config, or containers.

## Commit & Pull Request Guidelines
Recent history uses short Conventional Commit-style subjects such as `fix: ...` and `docs: ...`. Keep commits focused and imperative. PRs should include a brief description, mark the change type, and confirm local testing; this matches `.github/PULL_REQUEST_TEMPLATE.md`. Add screenshots for UI or docs visual changes, and link the related issue when one exists.

## Security & Configuration Tips
Never commit populated `.env` files or secrets. Generate strong values for `JWT_SECRET` and `ENCRYPTION_KEY`, and use `.env.example` as the source of truth for required configuration.
