 // MessagesView - Main component for the messaging interface
 class MessagesView {
    constructor(messageStore, wsManager, notificationManager) {
        this.messageStore = messageStore;
        this.wsManager = wsManager;
        this.notificationManager = notificationManager;
        
        this.currentPage = 1;
        this.isLoading = false;
        this.hasMoreMessages = true;
        this.scrollPositionToMaintain = null;
        this.pendingConversation = null;
        
        // Bind methods
        this.loadMoreMessages = this.throttle(this.loadMoreMessages.bind(this), 1000);
        this.handleIncomingMessage = this.handleIncomingMessage.bind(this);
        this.selectConversation = this.selectConversation.bind(this);
        this.sendMessage = this.sendMessage.bind(this);
        
        // Register message handler
        wsManager.registerMessageHandler(this.handleIncomingMessage);
        
        // Listen for the custom event for starting a new conversation
        document.addEventListener('startNewConversation', this.handleStartNewConversation.bind(this));
    }

    // Throttle function to prevent too many calls
    throttle(func, wait) {
        let timeout = null;
        let previous = 0;

        return function throttled(...args) {
            const now = Date.now();

            if (!previous) {
                previous = now;
            }

            const remaining = wait - (now - previous);

            if (remaining <= 0 || remaining > wait) {
                if (timeout) {
                    clearTimeout(timeout);
                    timeout = null;
                }

                previous = now;
                func.apply(this, args);
            } else if (!timeout) {
                timeout = setTimeout(() => {
                    previous = Date.now();
                    timeout = null;
                    func.apply(this, args);
                }, remaining);
            }
        };
    }

    handleStartNewConversation(event) {
        if (event.detail && event.detail.userId) {
            this.selectConversation(event.detail.userId);
        }
    }

    async init() {
        // Setup event listeners
        this.setupEventListeners();
        
        // Load conversations
        await this.loadConversations();
        
        // Check for pending conversations
        this.checkPendingConversations();
    }

    async loadConversations() {
        try {
            const contactsList = document.getElementById('contacts-list');
            contactsList.innerHTML = '<div class="loading">Loading conversations...</div>';
                        
            const response = await fetch('/api/messages/conversations', {
                method: 'GET',
                credentials: 'include',
                headers: {
                    'Accept': 'application/json',
                }
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            
            const conversations = data.conversations || [];
            this.messageStore.setConversations(conversations);
            this.renderConversationList();
        } catch (error) {
            console.error('Error loading conversations:', error);
            const contactsList = document.getElementById('contacts-list');
            contactsList.innerHTML = '<div class="error-message">Failed to load conversations</div>';
        }
    }

    renderConversationList() {
        const contactsList = document.getElementById('contacts-list');
        
        // Sort conversations by the timestamp of the last message (most recent first)
        const sortedConversations = [...this.messageStore.conversations].sort((a, b) => {
            // If no last_message_time, use last_seen as fallback
            const timeA = a.last_message_time ? new Date(a.last_message_time) : new Date(a.last_seen);
            const timeB = b.last_message_time ? new Date(b.last_message_time) : new Date(b.last_seen);
            // Sort in descending order (most recent first)
            return timeB - timeA;
        });
        
        if (sortedConversations.length === 0) {
            // Show a placeholder with a "Start New Conversation" button
            contactsList.innerHTML = `
                <div class="no-conversations">
                    <p>No conversations yet.</p>
                    <button id="start-new-conversation">Start a New Conversation</button>
                </div>
            `;
            document.getElementById('start-new-conversation')
                .addEventListener('click', () => this.showUserSearchModal());
        } else {
            contactsList.innerHTML = sortedConversations.map(conv => {
                // Format the timestamp for display
                const timestamp = conv.last_message_time 
                    ? this.formatTimestamp(conv.last_message_time, false) 
                    : this.formatTimestamp(conv.last_seen, false);
                
                return `
                    <div class="contact-item" data-user-id="${conv.other_user_id}">
                        <div class="sidebar-user-avatar">
                            ${conv.username.charAt(0).toUpperCase()}
                            <span class="sidebar-status-indicator ${conv.is_online ? 'online' : 'offline'}"></span>
                        </div>
                        <div class="sidebar-user-info">
                            <div class="sidebar-header-row">
                                <span class="sidebar-user-name">${conv.username}</span>
                                <span class="sidebar-last-time">${timestamp}</span>
                            </div>
                            <div class="sidebar-last-message">${conv.last_message || 'No messages yet'}</div>
                        </div>
                    </div>
                `;
            }).join('');
        }
    }

    async showUserSearchModal() {
        try {
            // Fetch registered users
            const response = await fetch('/api/users', { credentials: 'include' });
            if (!response.ok) throw new Error("Failed to fetch users");
            const data = await response.json();

            const currentUser = window.State.getCurrentUser();
            const users = data.users.filter(user => user.id !== currentUser.id);

            // Create a modal container
            const modal = document.createElement('div');
            modal.classList.add('user-search-modal');
            modal.innerHTML = `
                <div class="modal-overlay"></div>
                <div class="modal-content">
                    <h2>Select a user to chat with</h2>
                    <input type="text" id="user-search-input" placeholder="Search users..." />
                    <div id="user-search-results">
                        ${users.length > 0 ? 
                            users.map(user => `
                                <div class="user-search-item" data-user-id="${user.id}">
                                    ${user.nickname || user.username}
                                </div>
                            `).join('') :
                            '<div class="no-users">No other users found</div>'
                        }                   
                    </div>
                    <button id="close-user-search">Close</button>
                </div>
            `;
            document.body.appendChild(modal);

            // Close modal button
            modal.querySelector('#close-user-search').addEventListener('click', () => {
                this.closeUserSearchModal();
            });

            // Filter users as you type
            const searchInput = modal.querySelector('#user-search-input');
            searchInput.addEventListener('input', () => {
                const filter = searchInput.value.toLowerCase();
                const items = modal.querySelectorAll('.user-search-item');
                items.forEach(item => {
                    if (item.textContent.toLowerCase().includes(filter)) {
                        item.style.display = '';
                    } else {
                        item.style.display = 'none';
                    }
                });
            });

            // When a user is selected, start a conversation
            modal.querySelectorAll('.user-search-item').forEach(item => {
                item.addEventListener('click', (e) => {
                    const userId = e.target.dataset.userId;
                    this.selectConversation(userId);
                    this.closeUserSearchModal();
                });
            });

        } catch (error) {
            console.error('Error fetching users:', error);
            this.notificationManager.show('Failed to load users. Please try again later.');
        }
    }

    closeUserSearchModal() {
        const modal = document.querySelector('.user-search-modal');
        if (modal) {
            modal.remove();
        }
    }

    async loadMoreMessages(userId) {
        if (!userId || this.isLoading || !this.hasMoreMessages) return;
        
        const nextPage = this.currentPage + 1;
        const messagesList = document.getElementById('messages-list');
        
        // Remember scroll height before adding new messages
        const prevScrollHeight = messagesList.scrollHeight;
        this.scrollPositionToMaintain = prevScrollHeight;
        
        await this.loadMessages(userId, nextPage);
    }

    async loadMessages(userId, page = 1) {
        if (this.isLoading) return;
        this.isLoading = true;

        try {
            const messagesList = document.getElementById('messages-list');
            if (page === 1) {
                messagesList.innerHTML = '<div class="loading">Loading messages...</div>';
            }

            const response = await fetch(`/api/messages/${userId}?page=${page}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            const newMessages = data.messages || [];
            
            // If we got fewer messages than expected, we've reached the end
            if (newMessages.length === 0) {
                this.hasMoreMessages = false;
            }
            
            let updatedMessages;
            if (page === 1) {
                // First page - just set the messages
                updatedMessages = newMessages;
            } else {
                // Additional pages - prepend to existing messages
                const existing = this.messageStore.messages.get(userId) || [];
                updatedMessages = [...newMessages, ...existing];
            }
            
            // Sort all messages by timestamp to ensure correct order
            updatedMessages.sort((a, b) => new Date(a.created_at) - new Date(b.created_at));
            
            // Update the message store
            this.messageStore.setMessages(userId, updatedMessages);
            
            // Render the messages
            this.renderMessages(userId, page > 1);
            this.currentPage = page;
        } catch (error) {
            console.error('Error loading messages:', error);
            const messagesList = document.getElementById('messages-list');
            messagesList.innerHTML = '<div class="error-message">Failed to load messages</div>';
            this.hasMoreMessages = false; // Prevent further loading attempts on error
        } finally {
            this.isLoading = false;
        }
    }

    renderMessages(userId, maintainScrollPosition = false) {
        const messagesList = document.getElementById('messages-list');
        const messages = this.messageStore.messages.get(userId) || [];
        const currentUser = window.State.getCurrentUser();
        
        // Group messages by date for better visual separation
        const groupedMessages = this.groupMessagesByDate(messages);
        
        let messagesHtml = '';
        
        // Generate HTML for each date group
        Object.keys(groupedMessages).forEach(date => {
            messagesHtml += `<div class="message-date-separator">${date}</div>`;
            
            groupedMessages[date].forEach(msg => {
                // Check if this message is part of a sequence from the same sender
                const isContinuation = this.isMessageContinuation(msg, groupedMessages[date]);
                const isCurrentUser = msg.sender_id === currentUser?.id;
                
                // Get the username for the message
                let username;
                if (isCurrentUser && currentUser) {
                    username = currentUser.nickname || 
                             currentUser.username || 
                             `User ${currentUser.id}`;
                } else {
                    // For other users, first try to get from the message if it has sender info
                    username = msg.sender_nickname || msg.sender_username;
                    
                    // If not found in message, try to get from conversations
                    if (!username) {
                        const conversation = this.messageStore.conversations.find(conv => 
                            conv.other_user_id === msg.sender_id
                        );
                        username = conversation?.username || `User ${msg.sender_id}`;
                    }
                }
                
                messagesHtml += `
                    <div class="message ${isCurrentUser ? 'sent' : 'received'} ${isContinuation ? 'continuation' : ''}">
                        ${!isContinuation ? `<div class="message-username">${username}</div>` : ''}
                        <div class="message-content">${msg.content}</div>
                        <div class="message-info">
                            <span class="message-time">${this.formatTimestamp(msg.created_at, true)}</span>
                        </div>
                    </div>
                `;
            });
        });
        
        messagesList.innerHTML = messagesHtml;
        
        // Handle scroll position based on context
        if (maintainScrollPosition && this.scrollPositionToMaintain) {
            // When loading older messages, maintain scroll position
            const newScrollHeight = messagesList.scrollHeight;
            messagesList.scrollTop = newScrollHeight - this.scrollPositionToMaintain;
            this.scrollPositionToMaintain = null;
        } else {
            // For new messages or initial load, scroll to bottom
            messagesList.scrollTop = messagesList.scrollHeight;
        }
    }
    
    // Helper method to group messages by date
    groupMessagesByDate(messages) {
        const groups = {};
        
        messages.forEach(msg => {
            const date = new Date(msg.created_at).toLocaleDateString();
            if (!groups[date]) {
                groups[date] = [];
            }
            groups[date].push(msg);
        });
        
        return groups;
    }
    
    // Helper method to determine if a message is part of a sequence
    isMessageContinuation(message, messagesInGroup) {
        const messageIndex = messagesInGroup.findIndex(m => 
            m.id === message.id || m.temp_id === message.temp_id
        );
        
        if (messageIndex <= 0) return false;
        
        const previousMessage = messagesInGroup[messageIndex - 1];
        const timeDiff = new Date(message.created_at) - new Date(previousMessage.created_at);
        const isFromSameSender = message.sender_id === previousMessage.sender_id;
        
        // If less than 2 minutes between messages from same sender, consider it a continuation
        return isFromSameSender && timeDiff < 120000;
    }

    // Format timestamp for display
    formatTimestamp(timestamp, includeTime = true) {
        if (!timestamp) return '';
        
        const date = new Date(timestamp);
        const now = new Date();
        const diffMs = now - date;
        const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
        
        // Today
        if (diffDays === 0) {
            if (!includeTime) return 'Today';
            return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
        }
        
        // Yesterday
        if (diffDays === 1) {
            if (!includeTime) return 'Yesterday';
            return `Yesterday ${date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
        }
        
        // This week (within 7 days)
        if (diffDays < 7) {
            const options = includeTime 
                ? { weekday: 'short', hour: '2-digit', minute: '2-digit' }
                : { weekday: 'short' };
            return date.toLocaleDateString([], options);
        }
        
        // Older
        const options = includeTime 
            ? { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' }
            : { month: 'short', day: 'numeric' };
        return date.toLocaleDateString([], options);
    }

    setupEventListeners() {
        const contactsList = document.getElementById('contacts-list');
        const messageForm = document.getElementById('message-form');
        const messagesList = document.getElementById('messages-list');
        const messageInput = document.getElementById('message-input');
        const newConversationBtn = document.getElementById('new-conversation-btn');
        const themeToggle = document.getElementById('theme-toggle');

        // Theme toggle
        themeToggle.addEventListener('click', () => {
            document.body.classList.toggle('dark');
            const isDark = document.body.classList.contains('dark');
            themeToggle.innerHTML = isDark ? 
                '<i class="fas fa-sun"></i>' : 
                '<i class="fas fa-moon"></i>';
            localStorage.setItem('theme', isDark ? 'dark' : 'light');
        });

        // New conversation button
        newConversationBtn.addEventListener('click', () => {
            this.showUserSearchModal();
        });

        // Auto-resize textarea
        messageInput.addEventListener('input', () => {
            messageInput.style.height = 'auto';
            messageInput.style.height = messageInput.scrollHeight + 'px';
        });

        // Handle enter key (send on Enter, new line on Shift+Enter)
        messageInput.addEventListener('keydown', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                this.sendMessage();
            }
        });

        // Contact selection
        contactsList.addEventListener('click', (e) => {
            const contactItem = e.target.closest('.contact-item');
            if (contactItem) {
                const userId = contactItem.dataset.userId;
                this.selectConversation(userId);
            }
        });

        // Message form submission
        messageForm.addEventListener('submit', (e) => {
            e.preventDefault();
            this.sendMessage();
        });

        // Scroll handling for loading more messages
        messagesList.addEventListener('scroll', () => {
            // Consider "at top" when within 50px of the top
            if (messagesList.scrollTop <= 50 && 
                this.messageStore.currentConversation && 
                this.hasMoreMessages) {
                this.loadMoreMessages(this.messageStore.currentConversation);
            }
        });
    }

    async selectConversation(userId) {
        const currentUser = window.State.getCurrentUser();

        if (userId === currentUser.id) {
            return;
        }
        
        // Get user info from API if not in conversations
        let selectedUser = this.messageStore.conversations.find(conv => conv.other_user_id.toString() === userId);
        
        if (!selectedUser) {
            try {
                const response = await fetch(`/api/users/${userId}`);
                const userData = await response.json();
                
                // Create the user object but DON'T add to conversations list yet
                selectedUser = {
                    other_user_id: parseInt(userId),
                    username: userData.nickname || userData.username,
                    is_online: userData.is_online || false,
                    last_seen: userData.last_seen || new Date(),
                    last_message: 'No messages yet'
                };
                
                // Store as pending conversation
                this.pendingConversation = selectedUser;
            } catch (error) {
                console.error('Error fetching user data:', error);
                return;
            }
        }

        // Update active contact - only for existing conversations
        const contacts = document.querySelectorAll('.contact-item');
        contacts.forEach(contact => {
            contact.classList.remove('active');
            if (contact.dataset.userId === userId) {
                contact.classList.add('active');
            }
        });

        // Update chat header
        const chatHeader = document.getElementById('chat-header');
        chatHeader.innerHTML = `
            <div class="chat-user-info">
                <div class="chat-avatar">${selectedUser.username.charAt(0).toUpperCase()}</div>
                <div class="chat-user-details">
                    <div class="chat-username">${selectedUser.username}</div>
                    <div class="chat-status ${selectedUser.is_online ? 'online' : ''}">
                        ${selectedUser.is_online ? 'Online' : 'Last seen ' + this.formatTimestamp(selectedUser.last_seen)}
                    </div>
                </div>
            </div>
        `;

        // Show message input container
        const messageInputContainer = document.getElementById('message-input-container');
        messageInputContainer.style.display = 'block';

        this.messageStore.currentConversation = userId;
        await this.loadMessages(userId);
    }

    sendMessage() {
        const input = document.getElementById('message-input');
        const content = input.value.trim();
        if (!content || !this.messageStore.currentConversation) return;
    
        const receiverId = parseInt(this.messageStore.currentConversation);
        const currentUser = window.State.getCurrentUser();
    
        if (receiverId === currentUser.id) {
            this.notificationManager.show('You cannot send messages to yourself.');
            return;
        }
    
        // Create a temporary ID for immediate feedback
        const tempId = 'temp-' + Date.now();
        const timestamp = new Date();
        
        const message = {
            type: 'message',
            content,
            receiver_id: receiverId,
            sender_id: currentUser.id,
            timestamp,
            temp_id: tempId
        };
    
        try {
            // Send the message via WebSocket
            const sent = this.wsManager.send(message);
            
            if (!sent) {
                this.notificationManager.show('Could not connect to the messaging server. Please try again later.');
                return;
            }
    
            input.value = '';
            input.style.height = 'auto';
    
            // Add message to store immediately for instant feedback
            this.messageStore.addMessage(receiverId.toString(), {
                ...message,
                created_at: message.timestamp
            });
            this.renderMessages(receiverId.toString());
    
            // Check if this is a pending conversation that needs to be added to the list
            if (this.pendingConversation && this.pendingConversation.other_user_id.toString() === receiverId.toString()) {
                // Add the pending conversation to the conversations list with the message
                const newConversation = {
                    ...this.pendingConversation,
                    last_message: content.substring(0, 30) + (content.length > 30 ? '...' : ''),
                    last_message_time: timestamp
                };
                
                // Add to existing conversations
                const updatedConversations = [...this.messageStore.conversations, newConversation];
                
                // Clear the pending conversation
                this.pendingConversation = null;
                
                // Sort and update
                updatedConversations.sort((a, b) => {
                    const timeA = a.last_message_time ? new Date(a.last_message_time) : new Date(a.last_seen);
                    const timeB = b.last_message_time ? new Date(b.last_message_time) : new Date(b.last_seen);
                    return timeB - timeA;
                });
                
                this.messageStore.setConversations(updatedConversations);
                this.renderConversationList();
            } else {
                // Update existing conversation
                const updatedConversations = this.messageStore.conversations.map(conv => {
                    if (conv.other_user_id.toString() === receiverId.toString()) {
                        return {
                            ...conv,
                            last_message: content.substring(0, 30) + (content.length > 30 ? '...' : ''),
                            last_message_time: timestamp
                        };
                    }
                    return conv;
                });
                
                // Sort conversations by timestamp (most recent first)
                updatedConversations.sort((a, b) => {
                    const timeA = a.last_message_time ? new Date(a.last_message_time) : new Date(a.last_seen);
                    const timeB = b.last_message_time ? new Date(b.last_message_time) : new Date(b.last_seen);
                    return timeB - timeA;
                });
                
                this.messageStore.setConversations(updatedConversations);
                this.renderConversationList();
            }
        } catch (error) {
            console.error('Error sending message:', error);
            this.notificationManager.show('Failed to send message. Please try again.');
        }
    }

    handleIncomingMessage(message) {
        // First, check the message type
        if (message.type === 'message') {
            // Handle regular chat messages
            const conversationId = message.sender_id.toString();
            const currentUser = window.State.getCurrentUser();
            
            // Add message to store
            this.messageStore.addMessage(conversationId, {
                ...message,
                created_at: message.timestamp
            });
            
            // If this is the current conversation, render the messages
            if (this.messageStore.currentConversation === conversationId) {
                this.renderMessages(conversationId);
            }
            
            // Update the conversation with the latest message
            const conversations = this.messageStore.conversations;
            const conversationExists = conversations.some(conv => 
                conv.other_user_id.toString() === (message.sender_id === currentUser.id ? 
                    message.receiver_id.toString() : message.sender_id.toString())
            );
            
            if (conversationExists) {
                const updatedConversations = conversations.map(conv => {
                    // For messages we send, we need to check by receiver_id
                    // For messages we receive, we check by sender_id
                    const isTargetConversation = message.sender_id === currentUser.id ?
                        conv.other_user_id.toString() === message.receiver_id.toString() :
                        conv.other_user_id.toString() === message.sender_id.toString();
                    
                    if (isTargetConversation) {
                        return {
                            ...conv,
                            last_message: message.content.substring(0, 30) + (message.content.length > 30 ? '...' : ''),
                            last_message_time: message.timestamp
                        };
                    }
                    return conv;
                });
                
                // Sort conversations by timestamp (most recent first)
                updatedConversations.sort((a, b) => {
                    const timeA = a.last_message_time ? new Date(a.last_message_time) : new Date(a.last_seen);
                    const timeB = b.last_message_time ? new Date(b.last_message_time) : new Date(b.last_seen);
                    return timeB - timeA;
                });
                
                this.messageStore.setConversations(updatedConversations);
                this.renderConversationList();
            } else {
                // If the conversation doesn't exist in our list, refresh the list from server
                this.loadConversations();
            }
            
            // Only show notification if you're NOT the sender AND you're not currently viewing that conversation
            if (message.sender_id !== currentUser.id && this.messageStore.currentConversation !== conversationId) {
                // Find the sender's name
                const senderConversation = this.messageStore.conversations.find(
                    conv => conv && conv.other_user_id && conv.other_user_id.toString() === message.sender_id.toString()
                );
                
                const senderName = senderConversation ? senderConversation.username : 'Someone';
                
                this.notificationManager.newMessage({
                    sender: senderName
                });
            }
        }
        else if (message.type === 'status_update') {
            // Handle status update
            const userId = message.user_id.toString();
            
            // Update conversations list to reflect the new status
            const conversations = this.messageStore.conversations;
            const conversationExists = conversations.some(conv => conv.other_user_id.toString() === userId);
            
            // Update existing conversations
            if (conversationExists) {
                const updatedConversations = conversations.map(conv => {
                    if (conv.other_user_id.toString() === userId) {
                        return {
                            ...conv,
                            is_online: message.is_online,
                            last_seen: message.last_seen
                        };
                    }
                    return conv;
                });
                
                this.messageStore.setConversations(updatedConversations);
                this.renderConversationList();
            }
            
            // If this user is currently in an active conversation, update the header
            if (this.messageStore.currentConversation === userId) {
                const chatHeader = document.getElementById('chat-header');
                if (chatHeader) {
                    const statusElement = chatHeader.querySelector('.chat-status');
                    if (statusElement) {
                        statusElement.className = `chat-status ${message.is_online ? 'online' : ''}`;
                        statusElement.textContent = message.is_online ? 
                            'Online' : 
                            'Last seen ' + this.formatTimestamp(message.last_seen);
                    }
                }
            }
        }
    }

    checkPendingConversations() {
        // Check for pending conversation requests from sessionStorage
        const openConversationWith = sessionStorage.getItem('openConversationWith');
        if (openConversationWith) {
            // Clear the stored user ID
            sessionStorage.removeItem('openConversationWith');
            
            // Open the conversation
            this.selectConversation(openConversationWith);
            return;
        }
        
        // Also check for startConversationWithUser
        const startWithUser = sessionStorage.getItem('startConversationWithUser');
        if (startWithUser) {
            // Clear the stored value
            sessionStorage.removeItem('startConversationWithUser');
            
            // Open the conversation
            this.selectConversation(startWithUser);
        }
    }

    destroy() {
        // Unregister the message handler when the view is destroyed
        this.wsManager.unregisterMessageHandler(this.handleIncomingMessage);
    }
}
