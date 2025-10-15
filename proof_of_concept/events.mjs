// channels.js
const Channels = {
  events: Object.create(null),

  subscribe(event, handler) {
    (this.events[event] ??= new Set()).add(handler);
    return () => this.unsubscribe(event, handler); // allow quick unlisten
  },

  unsubscribe(event, handler) {
    this.events[event]?.delete(handler);
  },

  publish(event, detail) {
    this.events[event]?.forEach(handler => {
      try { handler(detail); }
      catch (err) { console.error(`Error in handler for ${event}:`, err); }
    });
  }
};

export default Channels;
