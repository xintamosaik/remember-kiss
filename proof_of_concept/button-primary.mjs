// button-primary.js
export class ButtonPrimary extends HTMLElement {
  connectedCallback() {
    this.attachShadow({ mode: "open" });
    this.render();
  }

  render() {
    const label = this.textContent.trim() || "OK";
    const type = this.getAttribute("type") || "button";

    this.shadowRoot.innerHTML = `
      <style>
        :host {
          display: inline-block;
        }
        button {
          cursor: pointer;
          border: none;
          border-radius: 1ch;
          padding: .5ch 1ch;
          font: inherit;
          background: hotpink;
          color: white;
        }
        button:hover {
          opacity: 0.9;
        }
        button:focus {
          outline: 2px solid var(--color-primary, hotpink);
          outline-offset: 1px;
        }
      </style>
      <button type="${type}">
        ${label}
      </button>
    `;

    // re-emit click so parent elements can listen normally
    const btn = this.shadowRoot.querySelector("button");
    btn.addEventListener("click", e =>
      this.dispatchEvent(new Event("click", { bubbles: true }))
    );
  }
}

customElements.define("button-primary", ButtonPrimary);