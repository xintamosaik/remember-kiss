// todo-edit.js
import todo from "./todo.mjs";   // your localStorage API
import "./button-primary.mjs";
export class TodoEdit extends HTMLElement {
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
        form { display: flex; flex-direction: column; gap: .5ch; }
        button {
          cursor: pointer;
          border: none;
          border-radius: 1ch;
          padding: .5ch 1ch;
        }
  
      </style>
      <form>
        <input name="key" type="hidden" />
        <label>
          Short Headline
          <input name="short" type="text" required />
        </label>
        <label>
          Long Description
          <input name="long" type="text" />
        </label>
        <div>
          <button class="primary" type="submit">Update Item</button>
          <button-primary>Add Item</button-primary>
          <button class="secondary" type="button" id="cancel">Cancel</button>
        </div>
      </form>
    `;

        const form = this.shadowRoot.querySelector("form");
        const cancel = this.shadowRoot.querySelector("#cancel");
        form.addEventListener("submit", e => this.onSubmit(e));
        cancel.addEventListener("click", () => this.hide());
    }

    /**
     * Opens the form and fills it with todo data.
     * @param {string} key
     */
    open(key) {
        const item = todo.load(key);
        if (!item) return;

        const form = this.shadowRoot.querySelector("form");
        form.key.value = key;
        form.short.value = item.short;
        form.long.value = item.long ?? "";

        this.show();
    }

    /**
     * Handle form submission.
     */
    onSubmit(e) {
        e.preventDefault();
        const data = new FormData(e.target);
        const key = data.get("key");
        const short = data.get("short");
        const long = data.get("long");
        todo.update({ short, long }, key);
        this.hide();
        this.dispatchEvent(
            new CustomEvent("todo-updated", {
                detail: { key, short, long },
                bubbles: true,
                composed: true,
            })
        );
    }

    show() { this.style.display = "block"; }
    hide() {
        const form = this.shadowRoot.querySelector("form");
        form.reset();
        this.style.display = "none";
    }
}




customElements.define("todo-edit", TodoEdit);
