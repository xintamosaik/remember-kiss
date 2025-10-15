// todo-list.js
import todo  from "./todo.mjs";   // your localStorage API

export class TodoItem extends HTMLElement {
    static get observedAttributes() { return ["key"]; }

    /**
     * Called when element is added to the DOM
     */
    connectedCallback() {
        // shadow DOM for style encapsulation
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
          display: list-item;
          cursor: pointer;
        }
        :host(:hover) {
          color: var(--color-primary, hotpink);
        }
      </style>
      <span>${item.short}</span>
    `;

        // remove any previous listener
        this.shadowRoot.onclick = null;
        this.shadowRoot.addEventListener("click", () => this.onClick(key, item));
    }

    /**
     * Emit a decoupled event when clicked
     */
    onClick(key, item) {
        this.dispatchEvent(new CustomEvent("todo-open", {
            bubbles: true,
            composed: true,
            detail: { key, item }
        }));
    }
}


customElements.define("todo-item", TodoItem);