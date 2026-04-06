# Frontend

React + Vite + TypeScript dashboard for the Finance Agent project.

## Requirements

- Node.js 18+
- Backend running at `http://localhost:8080`

## Install

```bash
npm install
```

## Run locally

```bash
npm run dev
```

The app will be available at the Vite dev server URL, usually `http://localhost:5173`.

## Build

```bash
npm run build
```

## Notes

- The frontend fetches data from `http://localhost:8080/api`.
- If the backend is unavailable, the UI falls back to bundled mock data.
- Tailwind utility classes are used throughout, and the shared class contracts from `shared/contracts.md` are preserved in the UI markup.