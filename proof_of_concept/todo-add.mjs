// todo-form.js
import todo  from "./todo.mjs";   // your localStorage API

export class TodoForm extends HTMLElement {
  connectedCallback() {
    this.attachShadow({ mode: "open" });
    this.render();
  }

  render() {
    this.shadowRoot.innerHTML = `
      <style>
        :host {
          display: block;
          border-top: 1px solid var(--color-primary, hotpink);
          padding: 1ch 0;
        }
        form {
          display: flex;
          flex-direction: column;
          gap: .5ch;
        }
        button {
          cursor: pointer;
          border: none;
          border-radius: 1ch;
          padding: .5ch 1ch;
        }
        button.primary { background: var(--color-primary, hotpink); color: white; }
        button.secondary { background: transparent; }
      </style>
      <form popover>
        <label>
          Short Headline
          <input name="short" type="text" required>
        </label>
        <div>
          <button class="primary" type="submit">Add Item</button>
          <button class="secondary" type="button" id="cancel">Cancel</button>
        </div>
      </form>
    `;

    const form = this.shadowRoot.querySelector("form");
    const cancel = this.shadowRoot.querySelector("#cancel");

    form.addEventListener("submit", e => this.onSubmit(e));
    cancel.addEventListener("click", () => this.reset());
  }

  onSubmit(event) {
    event.preventDefault();
    const data = new FormData(event.target);
    const short = data.get("short");
    const key = todo.save({ short });
    this.reset();

    // notify whoever cares that a new todo exists
    this.dispatchEvent(
      new CustomEvent("todo-added", {
        detail: { key },
        bubbles: true,
        composed: true
      })
    );
  }

  reset() {
    this.shadowRoot.querySelector("form").reset();
  }
}
