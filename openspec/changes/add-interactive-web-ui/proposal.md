# Proposal: Add Interactive Web UI

## Why

Currently, `inu web` provides only a RESTful API without a user interface. Users must interact with the anonymization service through API clients like curl or Postman. The `inu interactive` command offers a better workflow with entity memory across multiple restoration cycles, but it's limited to CLI environments.

**Problems:**
- Web mode lacks visual interface, limiting accessibility
- Users unfamiliar with APIs cannot easily use the web service
- Interactive workflow (anonymize once, restore multiple times) is unavailable in web mode
- No visual feedback for entity mappings during restoration

## What

Add a single-page web UI to `inu web` that replicates the interactive command workflow in a browser interface.

**Core Features:**

1. **Anonymize View**
   - Entity type selector (populated from `--entity-types` flag + custom input)
   - Input textarea for text to anonymize
   - Output panel showing anonymized result
   - "Anonymize" button with loading state
   - "Switch to Restore Mode" button (appears after anonymization)

2. **Restore View**
   - Entity mapping display (shows placeholder → original value pairs)
   - Read-only anonymized text panel (left)
   - Editable input textarea (right) for external processing results
   - "Restore" button to de-anonymize current text
   - "Back to Anonymize" button to return to first view

3. **State Management**
   - Client-side storage of entities using sessionStorage
   - Entities persist across view switches within same browser session
   - No server-side session required (stateless API design)

**Implementation:**
- Vanilla HTML/CSS/JavaScript (no framework dependencies)
- Static files served by Gin's static file handler
- Reuses existing `/api/v1/anonymize` and `/api/v1/restore` endpoints
- New route: `GET /` serves UI homepage (no authentication required)
- New route: `GET /api/v1/config` provides entity types configuration (optional)

## Impact

**Specs:**
- **web-api**: ADDED requirements for UI routes, static file serving, and frontend functionality

**Code:**
- `pkg/web/server.go`: Add `GET /` route and static file handler
- `pkg/web/static/`: New directory for HTML/CSS/JS files
- `cmd/inu/commands/web.go`: No changes needed (entity types already configurable)

**Dependencies:**
- None (vanilla JS, no new Go dependencies)

**Breaking Changes:**
- None (API endpoints remain unchanged)

**Documentation:**
- Update README with web UI usage
- Add screenshots of web interface

## Alternatives Considered

1. **React/Vue SPA**: More complex, requires build tooling, overkill for simple UI
2. **Server-rendered templates**: Requires state management on server, breaks stateless API design
3. **Separate UI server**: Additional deployment complexity, unnecessary for single-team project

## Open Questions

1. Should UI routes require authentication? (Recommend: no auth for GET /, keep API auth)
2. Should entity types be configurable at runtime? (Recommend: CLI flag only for v1)
3. Should we support multiple concurrent sessions? (Recommend: not needed, client-side state sufficient)
