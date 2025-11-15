# Interface Layout Documentation

This document shows the visual layout of the Be My Eyes TUI application.

## Main Interface Overview

The application uses a two-column layout with a 40-60 split:
- **Left Column (40%)**: Status, Videos, and History sections
- **Right Column (60%)**: Details panel
- **Footer**: Key bindings reference

## Layout Diagram

```
┌─ Status ──┐  ┌─── Details ────┐
│ Connected │  │ Video/Query    │
└───────────┘  │ information    │
┌─ Videos ──┐  │ displayed      │
│ • Video A │  │ here based on  │
│   Video B │  │ selection      │
└───────────┘  │                │
┌─ History ─┐  │                │
│ Q: "..."  │  │                │
│ Q: "..."  │  │                │
└───────────┘  └────────────────┘
 Key bindings shown here
```

## Key Features

- **Navigation**: Arrow keys, Tab, Mouse
- **Sections**: Status (connection), Videos (library), History (Q&A)
- **Details**: Context-aware right panel
- **Dialogs**: Question input, Menu, Help, About
- **Footer**: Always visible key bindings

For detailed usage, see [QUICKSTART.md](../QUICKSTART.md)
