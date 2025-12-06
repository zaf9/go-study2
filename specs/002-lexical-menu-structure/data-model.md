# Data Model: Lexical Menu Structure

**Feature**: 002-lexical-menu-structure  
**Date**: 2025-12-04  
**Status**: Complete

## Overview

This feature involves minimal data modeling as it primarily deals with menu navigation structure rather than persistent data. This document defines the menu item structure and the menu hierarchy.

## Entities

### MenuItem (Modified)

**Purpose**: Represents a single menu option with its description and action handler.

**Location**: `main.go`

**Fields**:

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| Description | string | Display text for the menu option | Non-empty |
| Action | func(io.Reader, io.Writer, io.Writer) | Handler function that executes when selected | Non-nil |

**Changes from Current**:
- **Before**: `Action func()`
- **After**: `Action func(io.Reader, io.Writer, io.Writer)`

**Rationale**: Enables menu actions to be interactive by receiving I/O streams, allowing sub-menus to read user input and write output.

---

### Menu Hierarchy

**Structure**: Two-level hierarchy

```
Main Menu (Level 1)
└── 0. Lexical elements
    └── Lexical Elements Sub-Menu (Level 2)
        ├── 0. Comments
        ├── 1. Tokens
        ├── 2. Semicolons
        ├── 3. Identifiers
        ├── 4. Keywords
        ├── 5. Operators
        ├── 6. Integers
        ├── 7. Floats
        ├── 8. Imaginary
        ├── 9. Runes
        ├── 10. Strings
        └── q. 返回上级菜单 (Return to main menu)
```

---

### SubMenuItem (Conceptual)

**Purpose**: Represents a topic in the lexical elements sub-menu.

**Location**: `internal/app/lexical_elements/lexical_elements.go`

**Structure** (implemented as map):

| Key | Description | Handler Function |
|-----|-------------|------------------|
| "0" | Comments | DisplayComments |
| "1" | Tokens | DisplayTokens |
| "2" | Semicolons | DisplaySemicolons |
| "3" | Identifiers | DisplayIdentifiers |
| "4" | Keywords | DisplayKeywords |
| "5" | Operators | DisplayOperators |
| "6" | Integers | DisplayIntegers |
| "7" | Floats | DisplayFloats |
| "8" | Imaginary | DisplayImaginary |
| "9" | Runes | DisplayRunes |
| "10" | Strings | DisplayStrings |
| "q" | 返回上级菜单 | (return from function) |

**Implementation**: `map[string]SubMenuItem` where SubMenuItem has Description and Action fields.

---

## State Transitions

### Main Menu State Machine

```
[Main Menu Display]
    ↓
[User Input]
    ↓
┌───────────────┐
│ Input = "0"?  │
└───────┬───────┘
        │ Yes
        ↓
[Show Lexical Elements Sub-Menu] ←──┐
        ↓                            │
[User Input in Sub-Menu]             │
        ↓                            │
┌───────────────────┐                │
│ Input = "q"?      │                │
└───────┬───────────┘                │
        │ No                         │
        ↓                            │
[Execute Topic Display] ─────────────┘
        │
        │ Yes
        ↓
[Return to Main Menu]
```

### Sub-Menu State Machine

```
[Sub-Menu Display]
    ↓
[User Input]
    ↓
┌────────────────┐
│ Input = "q"?   │
└────┬───────────┘
     │ Yes → Return to Main Menu
     │
     │ No
     ↓
┌────────────────────┐
│ Input in 0-10?     │
└────┬───────────────┘
     │ Yes
     ↓
[Execute Display Function]
     ↓
[Return to Sub-Menu Display]
     │
     │ No
     ↓
[Show Error Message]
     ↓
[Return to Sub-Menu Display]
```

---

## Validation Rules

### Main Menu Input

- **Valid inputs**: "0", "q", or any other menu option keys
- **Invalid inputs**: Empty string, whitespace-only, non-existent keys
- **Handling**: Trim whitespace, case-sensitive comparison

### Sub-Menu Input

- **Valid inputs**: "0" through "10", "q"
- **Invalid inputs**: 
  - Numbers outside range (e.g., "11", "-1")
  - Non-numeric characters (except "q")
  - Empty string
  - Whitespace-only
- **Handling**: Trim whitespace, case-sensitive comparison, display error and re-prompt

---

## Relationships

```
App (main.go)
  │
  ├── menu: map[string]MenuItem
  │     └── "0" → MenuItem{
  │           Description: "Lexical elements",
  │           Action: lexical_elements.DisplayMenu
  │         }
  │
  └── Run() → calls MenuItem.Action(stdin, stdout, stderr)
                    ↓
            lexical_elements.DisplayMenu(stdin, stdout, stderr)
                    ↓
            Creates internal map[string]SubMenuItem
                    ↓
            Menu loop → calls Display* functions
```

---

## Data Flow

### User Selects Main Menu Option "0"

```
1. User types "0" + Enter
2. main.App.Run() reads input
3. Looks up menu["0"]
4. Calls menu["0"].Action(app.stdin, app.stdout, app.stderr)
5. lexical_elements.DisplayMenu() executes
6. DisplayMenu() shows sub-menu
7. User interacts with sub-menu
8. User types "q" to return
9. DisplayMenu() returns
10. main.App.Run() continues main menu loop
```

### User Selects Sub-Menu Option "3" (Identifiers)

```
1. User types "3" + Enter (in sub-menu)
2. lexical_elements.DisplayMenu() reads input
3. Looks up subMenu["3"]
4. Calls DisplayIdentifiers()
5. DisplayIdentifiers() writes to stdout
6. Returns to DisplayMenu()
7. DisplayMenu() shows sub-menu again
```

---

## Assumptions

1. **No persistence**: Menu state is not saved between runs
2. **Single user**: No concurrent access considerations
3. **Synchronous**: All operations are blocking/synchronous
4. **UTF-8**: All text is UTF-8 encoded (for Chinese characters)

---

## Non-Functional Considerations

- **Performance**: Menu display and navigation should be instant (<10ms)
- **Memory**: Minimal overhead (menu maps are small, ~100 bytes)
- **Testability**: All menu logic testable with mock I/O (bytes.Buffer)

---

## Summary

The data model is intentionally minimal, consisting primarily of:
1. Modified `MenuItem` struct with I/O-aware Action signature
2. Two-level menu hierarchy (main → lexical elements sub-menu)
3. Simple state machines for navigation flow
4. No persistent data or complex relationships

This aligns with YAGNI (Principle V) and Simplicity (Principle I) principles.
