# Developer Notebook

An attempt to scratch my own itch.
I need one place to focus on the goal of a software development task at hand.
I need more structure than a text file but less features (i.e., distractions) than the "second brain" tools.

## Goals

- Optimize for information capture
  - Lowest friction to entering rows as possible
- Always visible: the problem statement
- Focus, concentrate, eliminate distractions
- Intentionally limit features to avoid having to make distracting decisions
- Simple, simple, simple implementation (already failed? lol)

## Steps

- [x] Restructure the model and initModel to group by #heading
- [x] Change the rendering to group by #heading
- [ ] Implement indentation
- [ ] Header
- [ ] LARGE: Update/delete (UX undefined)
- [ ] Save to disk
- [ ] Commands (e.g. /goto, /find, etc.)
- [ ] Footer (stats? current line?)
- [ ] Metadata?

## TODOs

- [ ] tab or shift+tab in the entry field to indent/dedent
- [ ] Display header with title
- [ ] Up/down arrows for selecting for edit or delete?
  - I might accidentally lose changes while just trying to move the cursor?
- [x] BUG: No notes display on open (cut off by extra `\n`'s)
- [x] REFACTOR: Model data structure -- map with key of tags/headers, with line numbers (indicated by slice index?) and indents
  - Question: How do I indent when a row is under multiple headings?
    - Only allow one tag per row??
- [x] Clear the text entry after enter
- [x] BUG: Initial text entry does not show, but it does from the second entry
- [x] BUG: Empty tag hash shows (off by 1 error from the string split?)
- [x] Show indentation in the text entry
- [x] Remember the previous row's heading so that the next ones can go under the same