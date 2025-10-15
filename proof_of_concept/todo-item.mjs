// todo-list.js
import todo from "./todo.mjs";   // your localStorage API
// In todo-item.mjs
import Channels  from "./events.mjs";

export class TodoItem extends HTMLElement {
  static get observedAttributes() { return ["key"]; }

  /**
   * Called when element is added to the DOM
   */
  connectedCallback() {
    // shadow DOM for style encapsulation
    Channels.subscribe("todo-updated", ({ key }) => {
      if (key === this.getAttribute("key")) {
        this.render();
      }
    });
    this.attachShadow({ mode: "open" });
    this.render();
  }

  /**
   * Called when the "key" attribute changes
   */
  attributeChangedCallback(name, oldValue, newValue) {
    if (name === "key" && oldValue !== newValue) {
      this.render();
    }
  }

  /**
   * Render the item
   */
  render() {
    if (!this.shadowRoot) return;
    const key = this.getAttribute("key");
    if (!key) return;

    const item = todo.load(key);
    if (!item) return;

    this.shadowRoot.innerHTML = `
      <style>
        :host {
     
          cursor: pointer;
        }
        :host(:hover) {
          color: var(--color-primary, hotpink);
        }
      </style>
      <li>${item.short}</li>
    `;

    // remove any previous listener
    this.shadowRoot.onclick = null;
    this.shadowRoot.addEventListener("click", () => this.onClick(key, item));
  }

  /**
   * Emit a decoupled event when clicked
   */
  onClick(key, item) {
    Channels.publish("todo-edit", { key });
  }
}


customElements.define("todo-item", TodoItem);