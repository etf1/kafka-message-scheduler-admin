# Kafka Message Scheduler Admin — Client
React/TypeScript frontend for the Kafka Message Scheduler Admin.
## Prerequisites
- **Node.js** ≥ 22
- **pnpm** (install: `corepack enable && corepack prepare pnpm@latest --activate`)
## Available Scripts
### `pnpm install`
Install dependencies.
### `pnpm start`
Runs the app in development mode with Vite HMR.\
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.
API calls are proxied to `http://localhost:9000` (Go backend).
### `pnpm test`
Runs unit tests with Vitest.
### `pnpm run build`
Builds the app for production to the `build` folder.\
Uses Vite with optimized chunking for best performance.
### `pnpm preview`
Previews the production build locally.
## Tech Stack
- **React 18** with TypeScript
- **Vite 6** (bundler + dev server)
- **Bulma** (CSS framework)
- **i18next** (internationalization — English + French)
- **Vitest** + **Testing Library** (unit tests)
