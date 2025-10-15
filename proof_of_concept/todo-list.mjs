// todo-list.js
import todo from "./todo.mjs";   // your localStorage API
import "./todo-item.mjs";            // ensure <todo-item> is defined

export class TodoList extends HTMLElement {
  connectedCallback() {
    // attach a shadow root for encapsulation
    this.attachShadow({ mode: "open" });
    this.render();
  }

  /**
   * Render all todos from localStorage
   */
  render() {
    // clear current shadow contents
    this.shadowRoot.innerHTML = `
      <style>
        :host {
          display: block;
        }
        ol {
          padding-inline-start: 2ch;
        }
      </style>
      <ol></ol>
    `;

    const ol = this.shadowRoot.querySelector("ol");

    // iterate through localStorage keys
    Object.keys(localStorage).forEach(key => {
      const item = todo.load(key);
      if (!item) return;

      const li = document.createElement("todo-item");
      li.setAttribute("key", key);
      ol.appendChild(li);
    });

    // bubble up todo-open events from children
    ol.addEventListener("todo-open", e => {
      this.dispatchEvent(
        new CustomEvent("todo-open", {
          detail: e.detail,
          bubbles: true,
          composed: true,
        })
      );
    });
  }
}

