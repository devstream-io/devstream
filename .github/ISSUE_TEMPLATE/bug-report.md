---
name: Bug Report
description: Report a bug encountered while operating DevStream
labels: bug
---

body:
- type: textarea
  id: problem
  attributes:
  label: What happened?
  description: Please provide as much info as possible.
  validations:
  required: true

- type: textarea
  id: repro
  attributes:
  label: How to reproduce it?
  validations:
  required: true

- type: textarea
  id: dtmVersion
  attributes:
  label: DevStream version (command: dtm version)
  required: true
