/**
 * chat.js - Chat functionality for embedding in applications that handle authentication via cookies
 */

// Global variables
let myUserId = null;
let selectedUser = null;
let ws = null; // WebSocket connection
let userNotifications = {}; // Track notifications per user
let lastLoadedMessages = []; // Track the last batch of loaded messages
let isLoading = false;
let offset = 0;
const limit = 10;

/**
 * Initialize the chat functionality
 * @param {string} userId - The ID of the current logged-in user
 */
function initChat(userId) {
  // Store user ID
  myUserId = userId;

  // Connect WebSocket
  connectWebSocket();

  // Setup UI components
  setupUI();

  // Fetch initial user list
  fetchAllUsers();
}

/**
 * Connect to WebSocket server
 */
function connectWebSocket() {
  // Create WebSocket connection (cookies will be sent automatically)
  const wsURL = `ws://${window.location.host}/ws`;

  try {
    ws = new WebSocket(wsURL);

    ws.onopen = () => {
      console.log('WebSocket connected successfully');
      // Clear any previous connection error messages
      showToast('Connected to chat server', 'success');
      // Fetch users after successful connection
      fetchAllUsers();
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log('WebSocket message received:', data);

        if (data.type === 'auth_success') {
          console.log('Authentication successful, user ID:', data.userID);
          // Ensure myUserId is set correctly
          if (data.userID && !myUserId) {
            myUserId = data.userID;
          }
        } else if (data.type === 'message') {
          // If this is the currently selected chat, display the message
          if (selectedUser && (selectedUser.id == data.message.sender_id ||
              selectedUser.id == data.message.receiver_id)) {
            displayMessage(data.message);
            // Scroll to bottom when new message arrives
            const messageArea = document.getElementById('message-area');
            if (messageArea) {
              messageArea.scrollTop = messageArea.scrollHeight;
            }
          }

          // Show notification for new message if it's from someone else
          if (data.message.sender_id != myUserId) {
            showNotification(data.message);
          }

          // Refresh the user list to update the latest message preview
          fetchAllUsers();
        } else if (data.type === 'user_status_change') {
          // When receiving a user status change, refresh the user list
          console.log('User status changed:', data.user);
          fetchAllUsers();
        }
      } catch (error) {
        console.error('Error processing WebSocket message:', error);
      }
    };

    ws.onclose = (event) => {
      console.log('WebSocket disconnected. Code:', event.code, 'Reason:', event.reason);

      // Show a message to the user
      if (event.code !== 1000) { // 1000 is normal closure
        showToast('Disconnected from chat server. Reconnecting...', 'error');
      }

      // Try to reconnect in 3 seconds
      setTimeout(connectWebSocket, 3000);
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      showToast('Error connecting to chat server', 'error');
    };
  } catch (error) {
    console.error('Failed to create WebSocket connection:', error);
    showToast('Failed to connect to chat server', 'error');
    // Try to reconnect in 5 seconds
    setTimeout(connectWebSocket, 5000);
  }
}

/**
 * Set up the UI components and event listeners
 */
function setupUI() {
  // Set up send button
  const sendButton = document.getElementById('send-button');
  if (sendButton) {
    sendButton.addEventListener('click', sendMessage);
  }

  // Set up message input
  const messageInput = document.getElementById('message-input');
  if (messageInput) {
    messageInput.addEventListener('keypress', (e) => {
      if (e.key === 'Enter') sendMessage();
    });
  }

  // Add event listener for search input
  const searchInput = document.querySelector('.search-input');
  if (searchInput) {
    searchInput.addEventListener('input', (e) => {
      const searchTerm = e.target.value.toLowerCase();
      const userItems = document.querySelectorAll('.user-item');

      userItems.forEach(item => {
        const userName = item.textContent.toLowerCase();
        if (userName.includes(searchTerm) || searchTerm === '') {
          item.style.display = 'flex';
        } else {
          item.style.display = 'none';
        }
      });
    });
  }

  // Set up scroll handler for message area
  setupScrollHandler();
}

/**
 * Set up scroll handler for loading older messages
 */
function setupScrollHandler() {
  const messageArea = document.getElementById('message-area');
  if (!messageArea) return;

  // Use a very small threshold to detect when user is near the top
  const SCROLL_THRESHOLD = 30;

  // Use a smaller throttle time for more responsive loading
  const handleScroll = throttle(async function () {
    // Only proceed if we have a selected user and we're not already loading
    if (isLoading || !selectedUser) return;

    // Check if we're near the top of the message area
    if (messageArea.scrollTop <= SCROLL_THRESHOLD && lastLoadedMessages.length >= limit) {
      isLoading = true;
      offset += limit;

      // Add loading indicator at the top of the message area
      const loadingIndicator = createLoadingIndicator();
      messageArea.insertBefore(loadingIndicator, messageArea.firstChild);

      // Save current scroll position and height
      const prevHeight = messageArea.scrollHeight;
      const prevScrollTop = messageArea.scrollTop;

      try {
        const messages = await fetchMessageHistory(selectedUser.id, offset, limit);

        // Remove loading indicator
        removeLoadingIndicator(loadingIndicator);

        if (messages.length > 0) {
          // Calculate new scroll position to maintain the same view
          const newHeight = messageArea.scrollHeight;
          const heightDifference = newHeight - prevHeight;

          // Set scroll position to maintain the same view
          messageArea.scrollTop = prevScrollTop + heightDifference;
        }

        if (messages.length < limit) {
          console.log('Reached the beginning of the conversation history');
        }
      } catch (error) {
        console.error('Error loading more messages:', error);
        removeLoadingIndicator(loadingIndicator);
      } finally {
        isLoading = false;
      }
    }
  }, 200);

  messageArea.addEventListener('scroll', handleScroll);

  // Also check if we need to load more messages when the chat is first opened
  function checkInitialScroll() {
    if (messageArea.scrollHeight <= messageArea.clientHeight &&
        !isLoading &&
        selectedUser &&
        lastLoadedMessages.length >= limit) {
      // If the message area isn't filled, load more messages
      handleScroll();
    }
  }

  // Check after initial messages are loaded
  setTimeout(checkInitialScroll, 300);

  // Also check when window is resized
  window.addEventListener('resize', checkInitialScroll);
}

/**
 * Send a message to another user
 */
function sendMessage() {
  const messageInput = document.getElementById('message-input');
  const messageText = messageInput.value.trim();
  if (!messageText || !selectedUser) return;

  const messageData = {
    type: 'message',
    message: {
      sender_id: myUserId,
      receiver_id: selectedUser.id,
      content: messageText
    }
  };

  // Send via WebSocket if available, otherwise send HTTP request
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(messageData));
    messageInput.value = '';
  } else {
    // Fallback to HTTP
    fetch('/send_message', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        // Include credentials to ensure cookies are sent
        'credentials': 'same-origin'
      },
      body: JSON.stringify(messageData)
    })
    .then(response => {
      if (response.ok) {
        messageInput.value = '';
        return response.json();
      }
      // If unauthorized, try to refresh the page to re-establish session
      if (response.status === 401) {
        showToast('Session expired. Please log in again.', 'error');
        return null;
      }
      throw new Error(`Failed to send message: ${response.status}`);
    })
    .then(data => {
      // Only process if we got valid data
      if (data && data.message) {
        // Display the sent message directly instead of refetching all messages
        displayMessage(data.message);
        const messageArea = document.getElementById('message-area');
        if (messageArea) {
          messageArea.scrollTop = messageArea.scrollHeight;
        }
      }
    })
    .catch(error => {
      showToast('Error sending message. Please try again.', 'error');
      console.error('Error sending message:', error);
    });
  }
}

/**
 * Display a message in the message area
 * @param {Object} msg - Message object
 */
function displayMessage(msg) {
  const messageArea = document.getElementById('message-area');
  if (!messageArea) return;

  const messageElement = createMessageElement(msg);

  // Check if we should auto-scroll after adding the message
  // (if user is already at the bottom)
  const shouldScroll = messageArea.scrollTop + messageArea.clientHeight >= messageArea.scrollHeight - 20;

  // Add the message to the message area
  messageArea.appendChild(messageElement);

  // Scroll to bottom if we were already at the bottom
  if (shouldScroll) {
    // Use a small delay to ensure rendering is complete
    setTimeout(() => {
      messageArea.scrollTop = messageArea.scrollHeight;
    }, 10);
  }
}

/**
 * Show notification for new message
 * @param {Object} msg - Message object
 */
function showNotification(msg) {
  // Don't notify for our own messages
  if (parseInt(msg.sender_id) === parseInt(myUserId)) return;

  // Update notification counter for this user
  if (!userNotifications[msg.sender_id]) {
    userNotifications[msg.sender_id] = 0;
  }
  userNotifications[msg.sender_id]++;

  // Update UI to show notification
  updateNotificationBadges();

  // Show popup notification
  showMiniNotification(msg.sender_nickname);
}

/**
 * Fetch all users from the server
 */
async function fetchAllUsers() {
  console.log("TRY FETCH USERS");
  try {
    // Add a small delay to ensure the session is properly established
    await new Promise(resolve => setTimeout(resolve, 100));

    const response = await fetch("/users", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        // Include credentials to ensure cookies are sent
        "credentials": "same-origin"
      }
    });

    if (!response.ok) {
      console.error(`Failed to fetch users: ${response.status}`);
      // If unauthorized, try to refresh the page to re-establish session
      if (response.status === 401) {
        showToast('Session expired. Please log in again.', 'error');
        return;
      }
      throw new Error(`Failed to fetch users: ${response.status}`);
    }

    const users = await response.json();
    console.log("Fetched users:", users);
    updateUsersList(users);

  } catch (error) {
    console.error('Error fetching users:', error);
    showToast('Error loading users. Please try refreshing the page.', 'error');
  }
}

/**
 * Fetch message history between current user and selected user
 * @param {string|number} receiverId - ID of the user to fetch messages with
 * @param {number} offset - Number of messages to skip (for pagination)
 * @param {number} limit - Number of messages to fetch
 * @returns {Promise<Array>} - Promise resolving to array of messages or empty array
 */
async function fetchMessageHistory(receiverId, offset = 0, limit = 10) {
  try {
    // Add a small delay to ensure the session is properly established
    await new Promise(resolve => setTimeout(resolve, 100));

    const response = await fetch(`/messages?receiver_id=${receiverId}&offset=${offset}&limit=${limit}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        // Include credentials to ensure cookies are sent
        "credentials": "same-origin"
      }
    });

    if (!response.ok) {
      console.error(`Failed to fetch messages: ${response.status}`);
      // If unauthorized, try to refresh the page to re-establish session
      if (response.status === 401) {
        showToast('Session expired. Please log in again.', 'error');
        return [];
      }
      throw new Error(`Failed to fetch messages: ${response.status}`);
    }

    const messages = await response.json();
    console.log("Fetched messages:", messages);

    // Make sure we handle null/undefined properly
    const messageArray = Array.isArray(messages) ? messages : [];

    // Sort messages by timestamp to ensure chronological order (oldest first)
    messageArray.sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp));

    // Track the last loaded batch
    lastLoadedMessages = messageArray;

    // Clear message area if this is the first batch
    const messageArea = document.getElementById('message-area');
    if (!messageArea) return messageArray;

    if (offset === 0) {
      messageArea.innerHTML = '';

      // For first batch, append messages in chronological order
      const fragment = document.createDocumentFragment();
      messageArray.forEach(msg => {
        const messageElement = createMessageElement(msg);
        fragment.appendChild(messageElement);
      });
      messageArea.appendChild(fragment);

      // Scroll to bottom after loading initial messages
      // Use a small delay to ensure rendering is complete
      setTimeout(() => {
        messageArea.scrollTop = messageArea.scrollHeight;
      }, 50);
    } else {
      // For subsequent batches, prepend messages
      const fragment = document.createDocumentFragment();

      // Add messages in chronological order
      messageArray.forEach(msg => {
        const messageElement = createMessageElement(msg);
        fragment.appendChild(messageElement);
      });

      // Prepend all messages at once for better performance
      if (messageArea.firstChild) {
        messageArea.insertBefore(fragment, messageArea.firstChild);
      } else {
        messageArea.appendChild(fragment);
      }
    }

    return messageArray;
  } catch (error) {
    console.error('Error fetching messages:', error);
    showToast('Error loading messages', 'error');
    // Return empty array instead of null to prevent errors
    return [];
  }
}

/**
 * Update the users list in the UI
 * @param {Array} users - Array of user objects
 */
function updateUsersList(users) {
  const onlineUsersList = document.getElementById('online-users-list');
  if (!onlineUsersList) return;

  // Clear current list
  onlineUsersList.innerHTML = '';

  // Filter out current user
  const filteredUsers = users.filter(user => user.id !== myUserId);

  if (filteredUsers.length === 0) {
    const emptyItem = document.createElement('li');
    emptyItem.textContent = 'No users available';
    emptyItem.className = 'user-item';
    onlineUsersList.appendChild(emptyItem);
    return;
  }

  // Create section headers
  const createSection = (title) => {
    const section = document.createElement('li');
    section.className = 'section-header';
    section.textContent = title;
    return section;
  };

  // Separate users into three categories
  const onlineUsers = filteredUsers.filter(user => user.is_online && user.has_chat_history);
  const offlineUsers = filteredUsers.filter(user => !user.is_online && user.has_chat_history);
  const noMessageUsers = filteredUsers.filter(user => !user.has_chat_history);

  // Sort users with messages by last message time (most recent first)
  onlineUsers.sort((a, b) => new Date(b.last_message?.timestamp || 0) - new Date(a.last_message?.timestamp || 0));
  offlineUsers.sort((a, b) => new Date(b.last_message?.timestamp || 0) - new Date(a.last_message?.timestamp || 0));

  // Sort users with no messages alphabetically
  noMessageUsers.sort((a, b) => a.nickname.localeCompare(b.nickname));

  // Create user item
  const createUserItem = (user) => {
    const userItem = document.createElement('li');
    userItem.className = 'user-item';
    userItem.dataset.userId = user.id;

    // Create online status indicator
    const statusIndicator = document.createElement('span');
    statusIndicator.className = `status-indicator ${user.is_online ? 'status-online' : 'status-offline'}`;

    // Create user nickname element
    const userNickname = document.createElement('span');
    userNickname.textContent = user.nickname;

    // Add notification badge if needed
    if (userNotifications[user.id] && userNotifications[user.id] > 0) {
      const badge = document.createElement('span');
      badge.className = 'notification-badge';
      badge.textContent = userNotifications[user.id];
      userItem.appendChild(badge);
    }

    // Add last message preview if available
    if (user.last_message && user.last_message.content) {
      const messagePreview = document.createElement('div');
      messagePreview.className = 'message-preview';
      messagePreview.textContent = truncateText(user.last_message.content, 30);

      // Add timestamp
      const timestamp = document.createElement('span');
      timestamp.className = 'message-time';
      timestamp.textContent = formatMessageTime(new Date(user.last_message.timestamp));

      userItem.appendChild(messagePreview);
      userItem.appendChild(timestamp);
    }

    // Assemble user item
    userItem.insertAdjacentElement('afterbegin', statusIndicator);
    userItem.insertAdjacentElement('afterbegin', userNickname);

    // Handle click to select user
    userItem.addEventListener('click', () => {
      // Navigate to the chat page with this user
      navigateToChat(user);
    });

    return userItem;
  };

  // Add online users with chat history section
  if (onlineUsers.length > 0) {
    onlineUsersList.appendChild(createSection('Online Users'));
    onlineUsers.forEach(user => {
      onlineUsersList.appendChild(createUserItem(user));
    });
  }

  // Add offline users with chat history section
  if (offlineUsers.length > 0) {
    onlineUsersList.appendChild(createSection('Offline Users'));
    offlineUsers.forEach(user => {
      onlineUsersList.appendChild(createUserItem(user));
    });
  }

  // Add users with no chat history section
  if (noMessageUsers.length > 0) {
    onlineUsersList.appendChild(createSection('New Users'));
    noMessageUsers.forEach(user => {
      onlineUsersList.appendChild(createUserItem(user));
    });
  }
}

/**
 * Navigate to chat with a specific user
 * @param {Object} user - The user to chat with
 */
function navigateToChat(user) {
  // Set selected user
  selectedUser = {
    id: user.id,
    nickname: user.nickname,
    isOnline: user.is_online
  };

  // Clear notifications for this user
  userNotifications[user.id] = 0;
  updateNotificationBadges();

  // Reset offset for fresh load
  offset = 0;

  // Update chat title
  const chatTitle = document.getElementById('chat-title');
  if (chatTitle) {
    chatTitle.textContent = `Chat with ${user.nickname}`;
  }

  // Update the URL without reloading the page
  history.pushState(null, "", `/chat/${user.id}`);

  // Render the chat interface
  renderChatInterface(user);

  // Fetch message history
  fetchMessageHistory(user.id, 0, limit);
}

/**
 * Render the chat interface for a specific user
 * @param {Object} user - The user to chat with
 */
function renderChatInterface(user) {
  const app = document.getElementById('app');
  if (!app) return;

  app.innerHTML = `
    <div class="h-screen flex">
      <!-- Mini notification for new messages -->
      <div id="mini-notification" class="fixed top-4 right-4 z-50 p-4 bg-transparent text-transparent rounded transition-all duration-300"></div>

      <div class="flex w-full h-full bg-white app-container">
        <!-- Main content area -->
        <div class="flex-grow flex flex-col">
          <!-- Chat Area -->
          <div id="chat-area" class="flex-grow flex flex-col h-full">
            <!-- App header with logo and user info -->
            <div class="app-header flex items-center justify-between p-4 border-b">
              <div class="flex items-center">
                <div class="app-logo">
                  <span class="logo-icon">💬</span>
                </div>
                <h3 id="chat-title" class="text-lg font-semibold ml-3">Chat with ${user.nickname}</h3>
              </div>
              <button id="back-button" class="back-btn">
                <span class="back-text">Back to Users</span>
                <span class="back-icon">←</span>
              </button>
            </div>

            <!-- Message area with improved styling -->
            <div id="message-area" class="flex-grow overflow-y-auto p-4 bg-gray-50 message-container">
              <!-- Messages will appear here -->
            </div>

            <!-- Message input area with improved styling -->
            <div class="message-input-container p-4 border-t flex">
              <input id="message-input" class="message-input flex-grow px-4 py-3" type="text" placeholder="Type your message...">
              <button id="send-button" class="send-button">
                <span>Send</span>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  `;

  // Add back button handler
  const backButton = document.getElementById('back-button');
  if (backButton) {
    backButton.addEventListener('click', () => {
      // Navigate back to the users list
      history.pushState(null, "", "/chat");
      renderUsersList();
    });
  }

  // Re-setup UI components after rendering
  setupUI();
}

/**
 * Render the user list view
 */
function renderUsersList() {
  const app = document.getElementById('app');
  if (!app) return;

  app.innerHTML = `
    <div class="h-screen flex">
      <!-- Mini notification for new messages -->
      <div id="mini-notification" class="fixed top-4 right-4 z-50 p-4 bg-transparent text-transparent rounded transition-all duration-300"></div>

      <div class="flex w-full h-full bg-white app-container">
        <!-- Sidebar for users -->
        <div id="sidebar" class="sidebar w-80 flex flex-col h-full">
          <div class="sidebar-header p-4 border-b border-gray-200">
            <h2 class="text-lg font-semibold">Contacts</h2>
            <div class="search-container mt-2">
              <input type="text" class="search-input" placeholder="Search users...">
              <span class="search-icon">🔍</span>
            </div>
          </div>
          <div class="overflow-y-auto flex-grow">
            <ul id="online-users-list" class="user-list divide-y divide-gray-200">
              <!-- User categories and lists will be populated here -->
            </ul>
          </div>
        </div>

        <!-- Main content area -->
        <div class="flex-grow flex flex-col">
          <!-- Empty state for chat area -->
          <div class="flex-grow flex flex-col items-center justify-center">
            <div class="text-center p-8">
              <div class="text-6xl mb-4">💬</div>
              <h2 class="text-2xl font-semibold mb-2">Chat Application</h2>
              <p class="text-gray-500">Select a user from the sidebar to start chatting</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  `;

  // Fetch users after rendering the list view
  fetchAllUsers();
}

/**
 * Create a message element with the proper format
 * @param {Object} msg - Message object
 * @returns {HTMLElement} Message element
 */
function createMessageElement(msg) {
  const isSender = parseInt(msg.sender_id) === parseInt(myUserId);
  const messageDiv = document.createElement('div');
  messageDiv.classList.add('message', isSender ? 'sent' : 'received');

  // Create username element (only for received messages)
  if (!isSender && msg.sender_nickname) {
    const usernameElement = document.createElement('div');
    usernameElement.classList.add('message-username');
    usernameElement.textContent = msg.sender_nickname;
    messageDiv.appendChild(usernameElement);
  }

  // Create message content
  const messageContent = document.createElement('div');
  messageContent.classList.add('message-content');
  messageContent.textContent = msg.content;
  messageDiv.appendChild(messageContent);

  // Create timestamp
  const messageTime = document.createElement('div');
  messageTime.classList.add('message-time');

  // Format date and time
  let timeText = "Time unavailable";
  if (msg.timestamp) {
    try {
      const msgDate = new Date(msg.timestamp);
      // Check if date is valid
      if (!isNaN(msgDate.getTime())) {
        timeText = formatFullDate(msgDate);
      }
    } catch (e) {
      console.error("Error formatting date:", e);
    }
  }
  messageTime.textContent = timeText;
  messageDiv.appendChild(messageTime);

  return messageDiv;
}

/**
 * Create loading indicator for message loading
 * @returns {HTMLElement} Loading indicator element
 */
function createLoadingIndicator() {
  const loadingIndicator = document.createElement('div');
  loadingIndicator.className = 'message-loading';
  loadingIndicator.textContent = 'Loading earlier messages...';
  return loadingIndicator;
}

/**
 * Remove loading indicator from DOM
 * @param {HTMLElement} indicator - Loading indicator element
 */
function removeLoadingIndicator(indicator) {
  if (indicator && indicator.parentNode) {
    indicator.parentNode.removeChild(indicator);
  }
}

/**
 * Update notification badges in the user list
 */
function updateNotificationBadges() {
  const userItems = document.querySelectorAll('.user-item');
  userItems.forEach(li => {
    const userId = li.dataset.userId;
    if (userId && userNotifications[userId]) {
      // Get or create badge
      let badge = li.querySelector('.notification-badge');
      if (!badge) {
        badge = document.createElement('span');
        badge.classList.add('notification-badge');
        li.appendChild(badge);
      }
      badge.textContent = userNotifications[userId];
    }
  });
}

/**
 * Show mini notification popup
 * @param {string} senderName - Name of the sender
 */
function showMiniNotification(senderName) {
  const miniNotification = document.getElementById('mini-notification');
  if (!miniNotification) return;

  miniNotification.textContent = `New message from ${senderName}`;
  miniNotification.style.backgroundColor = '#ef4444';
  miniNotification.style.color = 'white';

  // Clear after 3 seconds with fade effect
  setTimeout(() => {
    miniNotification.style.backgroundColor = 'transparent';
    miniNotification.style.color = 'transparent';

    // Clear text after fade completes
    setTimeout(() => {
      miniNotification.textContent = '';
    }, 300);
  }, 3000);
}

/**
 * Display a toast notification
 * @param {string} message - Message to display
 * @param {string} type - Type of toast (success or error)
 */
function showToast(message, type = 'success') {
  // Create toast element
  const toast = document.createElement('div');
  toast.className = `toast toast-${type}`;
  toast.textContent = message;

  // Add to document
  document.body.appendChild(toast);

  // Remove after 3 seconds
  setTimeout(() => {
    toast.classList.add('toast-hide');
    setTimeout(() => {
      document.body.removeChild(toast);
    }, 300);
  }, 3000);
}

/**
 * Creates a throttled function that only invokes the provided function
 * at most once per every 'limit' milliseconds.
 *
 * @param {Function} func - The function to throttle
 * @param {number} limit - The number of milliseconds to throttle invocations to
 * @returns {Function} The throttled function
 */
function throttle(func, limit) {
  let inThrottle;

  return function(...args) {
    const context = this;

    if (!inThrottle) {
      func.apply(context, args);
      inThrottle = true;
      setTimeout(() => {
        inThrottle = false;
      }, limit);
    }
  };
}

/**
 * Format date for message display
 * @param {Date} date - Date to format
 * @returns {string} Formatted date string
 */
function formatFullDate(date) {
  if (!date) return '';

  const now = new Date();
  const yesterday = new Date(now);
  yesterday.setDate(yesterday.getDate() - 1);

  // Format the time part
  const timeStr = date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });

  // If today, just show time
  if (date.toDateString() === now.toDateString()) {
    return `Today at ${timeStr}`;
  }
  // If yesterday, show Yesterday + time
  else if (date.toDateString() === yesterday.toDateString()) {
    return `Yesterday at ${timeStr}`;
  }
  // Otherwise show full date + time
  else {
    const dateStr = date.toLocaleDateString([], {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
    return `${dateStr} at ${timeStr}`;
  }
}

/**
 * Format time for messages list
 * @param {Date} date - Date to format
 * @returns {string} Formatted time string
 */
function formatMessageTime(date) {
  if (!date) return '';

  const now = new Date();
  const yesterday = new Date(now);
  yesterday.setDate(yesterday.getDate() - 1);

  // If today, show time
  if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }
  // If yesterday, show "Yesterday"
  else if (date.toDateString() === yesterday.toDateString()) {
    return 'Yesterday';
  }
  // Otherwise show date
  else {
    return date.toLocaleDateString([], { month: 'short', day: 'numeric' });
  }
}

/**
 * Helper function to truncate text
 * @param {string} text - Text to truncate
 * @param {number} maxLength - Maximum length
 * @returns {string} Truncated text
 */
function truncateText(text, maxLength) {
  if (!text) return '';
  if (text.length <= maxLength) return text;
  return text.substring(0, maxLength - 3) + '...';
}

// Export the public API
export { initChat, renderUsersList, renderChatInterface };