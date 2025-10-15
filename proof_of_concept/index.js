const new_timestamp = () => new Date().getTime();

function update_item(item, key) {
    item.update = new_timestamp();
    const json = JSON.stringify(item)
    localStorage.setItem(key, json);
}

function save_item(item) {
    const linux_timestamp = new_timestamp()
    const json = JSON.stringify(item)
    localStorage.setItem(linux_timestamp, json);
}

function load_item(key) {
    const data = localStorage.getItem(key);
    if (!data) return
    return JSON.parse(data)
}