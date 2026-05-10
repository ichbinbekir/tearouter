# Tearouter Task List

## 🏗 Architectural Mandates (CRITICAL)
- **Reference:** This project is inspired by Flutter's **GoRouter** package.
- **State Management:** Core Bubble Tea functions (`Init`, `Update`, `View`) MUST NOT use pointer receivers for `tearouter.Model`. This is a deliberate design choice to align with specific state management expectations.
- **Navigation:** Hierarchical path resolution (e.g., `/a/b/c`) automatically builds the model stack to allow natural `Pop` behavior.

## 📁 Project Structure
- `router.go`: Core routing logic and model.
- `messages.go`: Message definitions for redirection and errors.
- `internal/`: Internal implementation details (models, router wrapper).
- `pkg/`: Publicly usable packages.
- `examples/`: Usage examples.
- `cmd/tearouter`: **(ON HOLD)** Planned CLI tool for generating and managing Tearouter projects.

## ✅ Completed Tasks
- [x] Initial project exploration.
- [x] Identified missing sub-route support.
- [x] Implement nested route support in `Route` struct.
- [x] Implement hierarchical path resolution.
- [x] Support automatic stack population for hierarchical routes (Go/Push to `/a/b/c` builds full stack).
- [x] Add an example demonstrating hierarchical sub-routes.
- [x] Refactor and translate `.gemini/STYLE.md` with Go best practices.
- [x] Ensure state management works WITHOUT pointer receivers for `tearouter.Model`.
- [x] Add and verify Middleware support with hierarchical routes.
- [x] Update documentation (README.md, README.tr.md) with new features.

## 🚀 Upcoming Tasks
- [ ] Support for route parameters (e.g. `/user/:id`) - *Optional/Future*
- [ ] **(Low Priority)** Development of the `cmd/tearouter` CLI tool.
