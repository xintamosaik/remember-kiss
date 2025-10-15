/**
 * @typedef {Object} Todo
 * @property {string} short A short headline for the todo item
 * @property {string} [long] A longer description of the todo item
 * @property {number} [update] A timestamp in milliseconds from 1970-01-01
 * 
 */

/**
 * Creates a new timestamp in milliseconds from 1970-01-01
 * 
 * @return {number} The current timestamp in milliseconds
 */
const new_timestamp = () => new Date().getTime();

const todo = {
    /**
     * Update an existing item in localStorage
     * 
     * @param {Todo} item The item to update
     * @param {string} key The key of the item to update
     */
    update: function update_item(item, key) {
        item.update = new_timestamp();
        const json = JSON.stringify(item)
        localStorage.setItem(key, json);
    },
    /**
     * Save a new item to localStorage
     * 
     * @param {Todo} item The item to save
     */
    save: function save_item(item) {
        const linux_timestamp = new_timestamp()
        const json = JSON.stringify(item)
        localStorage.setItem(linux_timestamp, json);
    },
    /**
     * Load an item from localStorage
     * 
     * @param {string} key The key of the item to load
     * @return {Todo|undefined} The loaded item or undefined if not found
     */
    load: function load_item(key) {
        const data = localStorage.getItem(key);
        if (!data) return undefined;
        return JSON.parse(data)
    }
}

export default todo;