# remember-kiss
this file was generated with LLM. I was too lazy to write it down at 1am.

A **vanilla full-stack web app** built in pure Go, HTML, and CSS —  
with the goal of showing how little you actually need.

No framework. No dependencies. No node_modules.  
Just code you can read, understand, and extend.

---

## 🧭 Philosophy

> **KISS** — Keep It Server-Side.

Most web apps don’t need a framework.  
If your backend and frontend are fast, an MPA (multi-page app) can **beat** an SPA in simplicity, reliability, and even speed.

`remember-kiss` is a tiny experiment in what the web could still be if we design for:

- **DIY simplicity** over external dependencies  
- **Local-first** apps that run with a single `go run .`  
- **Understandability** from top to bottom  
- **Ease of extension** — add a page, handler, or form, and it just works  
- **Performance through clarity**, not caching layers

---

## ⚙️ Features

- Written entirely in Go (`net/http`, `html/template`)
- Server-side rendering for all pages (MPA style)
- No build step — regenerates static HTML files directly
- Persistent storage via a single `data.json` file
- Fully functional without JavaScript
- Optional progressive enhancement (transition API coming soon)

---

## 🚀 Run It

```bash
go run .
# then open http://localhost:3000


