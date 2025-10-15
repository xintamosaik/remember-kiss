
import { TodoItem } from "./todo-item.mjs";
import { TodoList } from "./todo-list.mjs";
import { TodoForm } from "./todo-add.mjs";
import { TodoEdit } from "./todo-edit.mjs";
customElements.define("todo-item", TodoItem);
customElements.define("todo-list", TodoList);
customElements.define("todo-form", TodoForm);
customElements.define("todo-edit", TodoEdit);
