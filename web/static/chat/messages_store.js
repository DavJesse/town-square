// Message Store - Manages state for messages and conversations
class MessageStore {
    constructor() {
        this.messages = new Map();
        this.conversations = [];
        this.currentConversation = null;
    }

    setConversations(conversations) {
        this.conversations = conversations;
    }

    setMessages(userId, messages) {
        this.messages.set(userId, messages);
    }

    addMessage(userId, message) {
        const messages = this.messages.get(userId) || [];
        messages.push(message);
        this.messages.set(userId, messages);
    }

    getCurrentMessages() {
        return this.currentConversation ? 
            this.messages.get(this.currentConversation) || [] : [];
    }
}