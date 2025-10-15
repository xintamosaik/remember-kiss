// todo-edit.js
import todo from "./todo.mjs";   // your localStorage API
import "./button-primary.mjs";
// In todo-edit.mjs
import Channels from "./events.mjs";

export class TodoEdit extends HTMLElement {
    connectedCallback() {
        this.attachShadow({ mode: "open" });
        this.render();
        /**
         * Listen for todo-edit events on the custom element itself
         */
        Channels.subscribe("todo-edit", ({ key }) => {
            this.open(key);
        });
    }

    render() {
        this.shadowRoot.innerHTML = `
      <style>
      </style>
      <form id="edit" popover>
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
          <button-primary type="submit">Update Item</button-primary>
          <button class="secondary" type="button" id="cancel">Cancel</button>
        </div>
      </form>
    `;
        const form = this.shadowRoot.querySelector("form");

        const submit = this.shadowRoot.querySelector("button-primary");
        submit.addEventListener("click", e => this.onSubmit(e));

        const cancel = this.shadowRoot.querySelector("#cancel");
        cancel.addEventListener("click", () => this.hidePopover());

    }

    /**
     * Opens the form and fills it with todo data.
     * @param {string} key
     */
    open(key) {
        const item = todo.load(key);
        if (!item) return;
        console.log(this.edit)
        form.key.value = key;
        form.short.value = item.short;
        form.long.value = item.long ?? "";

        form.showPopover();
    }

    /**
     * Handle form submission.
     */
    onSubmit(e) {
        console.log(this)
        e.preventDefault();
 
        const data = new FormData(form);
        const key = data.get("key");
        const short = data.get("short");
        const long = data.get("long");
        todo.update({ short, long }, key);
        form.reset();
        form.hidePopover();
        Channels.publish("todo-updated", { key });
    }
}

customElements.define("todo-edit", TodoEdit);
