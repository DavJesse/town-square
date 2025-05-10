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
let firstMessageId = null; // ID of the first (oldest) message loaded
let lastMessageId = null; // ID of the last (newest) message loaded
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

            // No need to update pagination IDs for new messages
            // We only care about the oldest message ID for loading older messages

            displayMessage(data.message);

            // Check if we should auto-scroll after adding the message
            const messageArea = document.getElementById('message-area');
            if (messageArea) {
              // We'll auto-scroll if the user is already near the bottom (within 100px)
              const isNearBottom = messageArea.scrollTop + messageArea.clientHeight >= messageArea.scrollHeight - 100;
              if (isNearBottom) {
                // Use a small delay to ensure rendering is complete
                setTimeout(() => {
                  messageArea.scrollTop = messageArea.scrollHeight;
                }, 10);
              } else {
                // If the user has scrolled up to read older messages, show a notification
                const newMessageNotification = document.getElementById('new-message-notification');
                if (newMessageNotification) {
                  newMessageNotification.classList.add('active');

                  // Add click handler to scroll to bottom
                  newMessageNotification.onclick = () => {
                    messageArea.scrollTop = messageArea.scrollHeight;
                    newMessageNotification.classList.remove('active');
                  };

                  // Auto-hide after 5 seconds
                  setTimeout(() => {
                    newMessageNotification.classList.remove('active');
                  }, 5000);
                }
              }
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
 * Set up scroll handler for loading older messages using ID-based pagination
 */
function setupScrollHandler() {
  const messageArea = document.getElementById('message-area');
  if (!messageArea) return;

  // Use a threshold to detect when user is near the top
  const SCROLL_THRESHOLD = 100;

  // Use throttling to prevent excessive API calls
  // 300ms is a good balance between responsiveness and preventing spam
  const handleScroll = throttle(async function () {
    // Only proceed if we have a selected user and we're not already loading
    if (isLoading || !selectedUser) return;

    // Check if we're near the top of the message area
    if (messageArea.scrollTop <= SCROLL_THRESHOLD && lastLoadedMessages.length >= limit) {
      isLoading = true;

      // Check if we already have the "Beginning of conversation" message
      const existingEndOfHistoryMsg = messageArea.querySelector('.end-of-history-message');
      if (existingEndOfHistoryMsg) {
        // We've already reached the beginning, no need to load more
        isLoading = false;
        return;
      }

      // Add loading indicator at the top of the message area
      const loadingIndicator = createLoadingIndicator();
      messageArea.insertBefore(loadingIndicator, messageArea.firstChild);

      try {
        console.log(`Loading older messages: firstMessageId=${firstMessageId}, limit=${limit}`);

        // Use the ID of the first (oldest) message as the cursor for pagination
        const messages = await fetchMessageHistory(selectedUser.id, firstMessageId, limit);

        // Remove loading indicator
        removeLoadingIndicator(loadingIndicator);

        // Check if we've reached the beginning of the conversation
        if (messages.length === 0 || messages.length < limit) {
          console.log('Reached the beginning of the conversation history');

          // Only add the end-of-history message if it doesn't already exist
          if (!messageArea.querySelector('.end-of-history-message')) {
            const endOfHistoryMsg = document.createElement('div');
            endOfHistoryMsg.className = 'end-of-history-message';
            endOfHistoryMsg.textContent = 'Beginning of conversation';
            messageArea.insertBefore(endOfHistoryMsg, messageArea.firstChild);
          }
        }
      } catch (error) {
        console.error('Error loading more messages:', error);
        removeLoadingIndicator(loadingIndicator);
        showToast('Error loading older messages', 'error');
      } finally {
        isLoading = false;
      }
    }
  }, 300); // 300ms throttle time

  messageArea.addEventListener('scroll', handleScroll);

  // Check if we need to load more messages when the chat is first opened
  function checkInitialScroll() {
    // Make sure we have messages and the message area isn't filled
    if (messageArea.scrollHeight <= messageArea.clientHeight &&
        !isLoading &&
        selectedUser &&
        lastLoadedMessages &&
        lastLoadedMessages.length >= limit) {

      // Check if we already have the "Beginning of conversation" message
      const existingEndOfHistoryMsg = messageArea.querySelector('.end-of-history-message');
      if (existingEndOfHistoryMsg) {
        // We've already reached the beginning, no need to load more
        return;
      }

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
  // We'll auto-scroll if the user is already near the bottom (within 100px)
  const isNearBottom = messageArea.scrollTop + messageArea.clientHeight >= messageArea.scrollHeight - 100;

  // Add the new message to the message area (always at the bottom since it's the newest)
  messageArea.appendChild(messageElement);

  // Scroll to bottom if we were already near the bottom
  if (isNearBottom) {
    // Use a small delay to ensure rendering is complete
    setTimeout(() => {
      messageArea.scrollTop = messageArea.scrollHeight;
    }, 10);
  } else {
    // If the user has scrolled up to read older messages, show a notification
    // that new messages have arrived
    const newMessageNotification = document.getElementById('new-message-notification');
    if (newMessageNotification) {
      newMessageNotification.classList.add('active');

      // Add click handler to scroll to bottom if not already there
      newMessageNotification.onclick = () => {
        messageArea.scrollTop = messageArea.scrollHeight;
        newMessageNotification.classList.remove('active');
      };

      // Auto-hide after 5 seconds
      setTimeout(() => {
        newMessageNotification.classList.remove('active');
      }, 5000);
    }
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

    // Update the main chat users list if it exists
    updateUsersList(users);

    // Also update the homepage online users list if it exists
    const onlineUsersContent = document.getElementById('online_users_content');
    if (onlineUsersContent) {
      populateOnlineUsersList(users, onlineUsersContent);
    }

    return users;
  } catch (error) {
    console.error('Error fetching users:', error);
    showToast('Error loading users. Please try refreshing the page.', 'error');
    return [];
  }
}

/**
 * Fetch message history between current user and selected user using ID-based pagination
 * @param {string|number} receiverId - ID of the user to fetch messages with
 * @param {number|null} oldestMessageId - ID of the oldest message loaded (for cursor-based pagination)
 * @param {number} limit - Number of messages to fetch
 * @returns {Promise<Array>} - Promise resolving to array of messages or empty array
 */
async function fetchMessageHistory(receiverId, oldestMessageId = null, limit = 10) {
  try {
    console.log(`Fetching messages with receiverId=${receiverId}, oldestMessageId=${oldestMessageId}, limit=${limit}`);

    // Build the query URL with ID-based pagination parameters
    let url = `/messages?receiver_id=${receiverId}&limit=${limit}`;

    // Add oldestMessageId if provided (for loading older messages)
    if (oldestMessageId) {
      url += `&last_message_id=${oldestMessageId}&direction=older`;
    }

    const response = await fetch(url, {
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

    // Make sure we handle null/undefined properly
    const messageArray = Array.isArray(messages) ? messages : [];
    console.log(`Fetched ${messageArray.length} messages:`, messageArray);

    // Messages should already be in chronological order (oldest first) from the server

    // Track the last loaded batch
    lastLoadedMessages = messageArray;

    // Update the first message ID for pagination (oldest message)
    if (messageArray.length > 0) {
      // Update firstMessageId if this is the first batch or if we found older messages
      if (!firstMessageId || messageArray[0].id < firstMessageId) {
        firstMessageId = messageArray[0].id;
      }

      // Update lastMessageId if this is the first batch
      if (!oldestMessageId && messageArray.length > 0) {
        lastMessageId = messageArray[messageArray.length - 1].id;
      }
    }

    // Get the message area element
    const messageArea = document.getElementById('message-area');
    if (!messageArea) return messageArray;

    if (!oldestMessageId) {
      // For the first batch (initial load), clear the message area
      messageArea.innerHTML = '';

      // Create a document fragment for better performance
      const fragment = document.createDocumentFragment();

      // Add all messages to the fragment in chronological order (oldest first)
      messageArray.forEach(msg => {
        const messageElement = createMessageElement(msg);
        fragment.appendChild(messageElement);
      });

      // Append all messages at once
      messageArea.appendChild(fragment);

      // Scroll to bottom after loading initial messages
      // Use a small delay to ensure rendering is complete
      setTimeout(() => {
        messageArea.scrollTop = messageArea.scrollHeight;
      }, 50);
    } else {
      // For older messages, prepend them at the top
      const fragment = document.createDocumentFragment();

      // Add older messages to the fragment
      messageArray.forEach(msg => {
        const messageElement = createMessageElement(msg);
        fragment.appendChild(messageElement);
      });

      // Save current scroll position and height
      const prevHeight = messageArea.scrollHeight;
      const prevScrollTop = messageArea.scrollTop;

      // Prepend all messages at once for better performance
      if (messageArea.firstChild) {
        messageArea.insertBefore(fragment, messageArea.firstChild);
      } else {
        messageArea.appendChild(fragment);
      }

      // Calculate new scroll position to maintain the same view
      const newHeight = messageArea.scrollHeight;
      const heightDifference = newHeight - prevHeight;

      // Set scroll position to maintain the same view
      messageArea.scrollTop = prevScrollTop + heightDifference;
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
    emptyItem.className = 'user-item empty-state';

    const emptyIcon = document.createElement('div');
    emptyIcon.className = 'empty-state-icon';
    emptyIcon.textContent = '👥';

    const emptyText = document.createElement('div');
    emptyText.className = 'empty-state-text';
    emptyText.textContent = 'No users available';

    emptyItem.appendChild(emptyIcon);
    emptyItem.appendChild(emptyText);

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
  const usersWithMessages = filteredUsers.filter(user => user.has_chat_history);
  const usersWithoutMessages = filteredUsers.filter(user => !user.has_chat_history);

  // Further separate by online status
  const onlineUsersWithMessages = usersWithMessages.filter(user => user.is_online);
  const offlineUsersWithMessages = usersWithMessages.filter(user => !user.is_online);
  const onlineUsersWithoutMessages = usersWithoutMessages.filter(user => user.is_online);
  const offlineUsersWithoutMessages = usersWithoutMessages.filter(user => !user.is_online);

  // Sort users with messages by last message time (most recent first)
  onlineUsersWithMessages.sort((a, b) => {
    // Check if both users have last_message
    if (a.last_message && b.last_message) {
      return new Date(b.last_message.timestamp) - new Date(a.last_message.timestamp);
    }
    // If only a has last_message, a comes first
    else if (a.last_message) {
      return -1;
    }
    // If only b has last_message, b comes first
    else if (b.last_message) {
      return 1;
    }
    // If neither has last_message, sort alphabetically
    return a.nickname.localeCompare(b.nickname);
  });

  // Sort offline users with messages the same way
  offlineUsersWithMessages.sort((a, b) => {
    if (a.last_message && b.last_message) {
      return new Date(b.last_message.timestamp) - new Date(a.last_message.timestamp);
    }
    else if (a.last_message) {
      return -1;
    }
    else if (b.last_message) {
      return 1;
    }
    return a.nickname.localeCompare(b.nickname);
  });

  // Sort users with no messages alphabetically
  onlineUsersWithoutMessages.sort((a, b) => a.nickname.localeCompare(b.nickname));
  offlineUsersWithoutMessages.sort((a, b) => a.nickname.localeCompare(b.nickname));

  // Create user item
  const createUserItem = (user) => {
    const userItem = document.createElement('li');
    userItem.className = 'user-item';
    if (selectedUser && selectedUser.id === user.id) {
      userItem.classList.add('selected');
    }
    userItem.dataset.userId = user.id;

    // Create user item header (status + nickname)
    const userItemHeader = document.createElement('div');
    userItemHeader.className = 'user-item-header';

    // Create online status indicator
    const statusIndicator = document.createElement('span');
    statusIndicator.className = `status-indicator ${user.is_online ? 'status-online' : 'status-offline'}`;
    userItemHeader.appendChild(statusIndicator);

    // Create user nickname element
    const userNickname = document.createElement('span');
    userNickname.className = 'user-nickname';
    userNickname.textContent = user.nickname;
    userItemHeader.appendChild(userNickname);

    userItem.appendChild(userItemHeader);

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
      userItem.appendChild(messagePreview);

      // Add timestamp
      const timestamp = document.createElement('div');
      timestamp.className = 'message-time';
      timestamp.textContent = formatMessageTime(new Date(user.last_message.timestamp));
      userItem.appendChild(timestamp);
    }

    // Handle click to select user
    userItem.addEventListener('click', () => {
      // Navigate to the chat page with this user
      navigateToChat(user);
    });

    return userItem;
  };

  // Add users with chat history first (online)
  if (onlineUsersWithMessages.length > 0) {
    onlineUsersList.appendChild(createSection('Recent Conversations (Online)'));
    onlineUsersWithMessages.forEach(user => {
      onlineUsersList.appendChild(createUserItem(user));
    });
  }

  // Add users with chat history (offline)
  if (offlineUsersWithMessages.length > 0) {
    onlineUsersList.appendChild(createSection('Recent Conversations (Offline)'));
    offlineUsersWithMessages.forEach(user => {
      onlineUsersList.appendChild(createUserItem(user));
    });
  }

  // Add users without chat history (online)
  if (onlineUsersWithoutMessages.length > 0) {
    onlineUsersList.appendChild(createSection('New Users (Online)'));
    onlineUsersWithoutMessages.forEach(user => {
      onlineUsersList.appendChild(createUserItem(user));
    });
  }

  // Add users without chat history (offline)
  if (offlineUsersWithoutMessages.length > 0) {
    onlineUsersList.appendChild(createSection('New Users (Offline)'));
    offlineUsersWithoutMessages.forEach(user => {
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

  // Reset message IDs for fresh load
  firstMessageId = null;
  lastMessageId = null;

  // Update chat title
  const chatTitle = document.getElementById('chat-title');
  if (chatTitle) {
    chatTitle.textContent = `Chat with ${user.nickname}`;
  }

  // Update the URL without reloading the page
  history.pushState(null, "", `/chat/${user.id}`);

  // Render the chat interface
  renderChatInterface(user);

  // Fetch initial message history (most recent messages)
  fetchMessageHistory(user.id, null, limit);
}

/**
 * Render the chat interface for a specific user
 * @param {Object} user - The user to chat with
 */
function renderChatInterface(user) {
  const app = document.getElementById('app');
  if (!app) return;

  app.innerHTML = `
    <div class="chat-container">
      <!-- Mini notification for new messages -->
      <div id="mini-notification" class="fixed top-4 right-4 z-50 p-4 bg-transparent text-transparent rounded transition-all duration-300"></div>

      <!-- Chat Sidebar -->
      <div class="chat-sidebar">
        <div class="sidebar-header">
          <h2 class="sidebar-title">Messages</h2>
          <div class="search-container">
            <input type="text" class="search-input" placeholder="Search users...">
          </div>
        </div>
        <div class="users-container">
          <ul id="online-users-list" class="user-list">
            <!-- User categories and lists will be populated here -->
          </ul>
        </div>
      </div>

      <!-- Main Chat Area -->
      <div class="chat-main">
        <div class="chat-header">
          <button id="mobile-toggle" class="mobile-toggle">
            <span class="mobile-toggle-icon">👥</span>
          </button>
          <h3 id="chat-title" class="chat-title">Chat with ${user.nickname}</h3>
          <button id="back-button" class="back-button">
            <span class="back-icon">←</span>
            <span>Back</span>
          </button>
        </div>

        <div class="message-container">
          <div id="message-area" class="message-area">
            <!-- Messages will appear here -->
          </div>

          <!-- New message notification that appears when user has scrolled up -->
          <div id="new-message-notification" class="new-message-notification">
            New messages ↓
          </div>
        </div>

        <div class="message-input-container">
          <input id="message-input" class="message-input" type="text" placeholder="Type your message...">
          <button id="send-button" class="send-button">
            Send
            <svg class="send-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="22" y1="2" x2="11" y2="13"></line>
              <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
            </svg>
          </button>
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

  // Add mobile toggle button handler
  const mobileToggle = document.getElementById('mobile-toggle');
  const chatSidebar = document.querySelector('.chat-sidebar');
  if (mobileToggle && chatSidebar) {
    mobileToggle.addEventListener('click', () => {
      // Toggle the active class on the sidebar
      chatSidebar.classList.toggle('active');
    });

    // Close sidebar when a user is selected on mobile
    const userItems = chatSidebar.querySelectorAll('.user-item');
    userItems.forEach(item => {
      item.addEventListener('click', () => {
        chatSidebar.classList.remove('active');
      });
    });
  }

  // Fetch users to populate the sidebar
  fetchAllUsers();

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
    <div class="chat-container">
      <!-- Mini notification for new messages -->
      <div id="mini-notification" class="fixed top-4 right-4 z-50 p-4 bg-transparent text-transparent rounded transition-all duration-300"></div>

      <!-- Chat Sidebar -->
      <div class="chat-sidebar">
        <div class="sidebar-header">
          <h2 class="sidebar-title">Messages</h2>
          <div class="search-container">
            <input type="text" class="search-input" placeholder="Search users...">
          </div>
        </div>
        <div class="users-container">
          <ul id="online-users-list" class="user-list">
            <!-- User categories and lists will be populated here -->
          </ul>
        </div>
      </div>

      <!-- Main Chat Area (Empty State) -->
      <div class="chat-main">
        <div class="empty-state">
          <div class="empty-state-icon">💬</div>
          <h2 class="empty-state-title">Select a conversation</h2>
          <p class="empty-state-text">Choose a user from the sidebar to start chatting</p>
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

  // Add data attribute for message ID to help with debugging
  messageDiv.dataset.messageId = msg.id;

  // Create message header with username and timestamp
  const messageHeader = document.createElement('div');
  messageHeader.classList.add('message-header');

  // Create username element (only for received messages)
  if (!isSender && msg.sender_nickname) {
    const usernameElement = document.createElement('span');
    usernameElement.classList.add('message-username');
    usernameElement.textContent = msg.sender_nickname;
    messageHeader.appendChild(usernameElement);
  }

  // Create timestamp
  const messageTime = document.createElement('span');
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
  messageHeader.appendChild(messageTime);

  // Add header to message
  messageDiv.appendChild(messageHeader);

  // Create message content
  const messageContent = document.createElement('div');
  messageContent.classList.add('message-content');
  messageContent.textContent = msg.content;
  messageDiv.appendChild(messageContent);

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
  let lastFunc;
  let lastRan;

  return function(...args) {
    const context = this;

    if (!lastRan) {
      // First time being called, execute immediately
      func.apply(context, args);
      lastRan = Date.now();
    } else {
      // Clear previous scheduled execution
      clearTimeout(lastFunc);

      // Schedule a new execution after the throttle period
      lastFunc = setTimeout(function() {
        // Only execute if enough time has passed since last execution
        if ((Date.now() - lastRan) >= limit) {
          func.apply(context, args);
          lastRan = Date.now();
        }
      }, limit - (Date.now() - lastRan));
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

  // Format the time part (more compact)
  const timeStr = date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });

  // If today, just show time
  if (date.toDateString() === now.toDateString()) {
    return timeStr;
  }
  // If yesterday, show Yesterday
  else if (date.toDateString() === yesterday.toDateString()) {
    return `Yesterday ${timeStr}`;
  }
  // Otherwise show compact date + time
  else {
    const dateStr = date.toLocaleDateString([], {
      month: 'short',
      day: 'numeric'
    });
    return `${dateStr} ${timeStr}`;
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

/**
 * Create a compact user list for the homepage
 * @param {Array} users - Array of user objects
 * @param {HTMLElement} container - Container element to append users to
 */
function populateOnlineUsersList(users, container) {
  if (!container) return;

  // Clear current list
  container.innerHTML = '';

  // Filter out current user and get only online users
  const onlineUsers = users.filter(user => user.id !== myUserId && user.is_online);
  const offlineUsers = users.filter(user => user.id !== myUserId && !user.is_online);

  if (onlineUsers.length === 0 && offlineUsers.length === 0) {
    const emptyItem = document.createElement('div');
    emptyItem.textContent = 'No users available';
    emptyItem.className = 'empty-user-item';
    container.appendChild(emptyItem);
    return;
  }

  // Create a compact user item
  const createCompactUserItem = (user) => {
    const userItem = document.createElement('div');
    userItem.className = 'compact-user-item';
    userItem.dataset.userId = user.id;

    // Create status indicator
    const statusIndicator = document.createElement('span');
    statusIndicator.className = `status-indicator ${user.is_online ? 'status-online' : 'status-offline'}`;

    // Create user nickname
    const userNickname = document.createElement('span');
    userNickname.className = 'user-nickname';
    userNickname.textContent = user.nickname;

    // Add notification badge if needed
    if (userNotifications[user.id] && userNotifications[user.id] > 0) {
      const badge = document.createElement('span');
      badge.className = 'notification-badge';
      badge.textContent = userNotifications[user.id];
      userItem.appendChild(badge);
    }

    // Assemble user item
    userItem.appendChild(statusIndicator);
    userItem.appendChild(userNickname);

    // Add click handler to open chat
    userItem.addEventListener('click', () => {
      // Navigate to chat with this user
      history.pushState(null, "", `/chat/${user.id}`);
      navigateToChat(user);
    });

    return userItem;
  };

  // Add online users first
  if (onlineUsers.length > 0) {
    const onlineHeader = document.createElement('div');
    onlineHeader.className = 'user-section-header';
    onlineHeader.textContent = 'Online';
    container.appendChild(onlineHeader);

    // Sort online users by last message time or alphabetically
    onlineUsers.sort((a, b) => {
      if (a.has_chat_history && b.has_chat_history) {
        if (a.last_message?.timestamp && b.last_message?.timestamp) {
          return new Date(b.last_message.timestamp) - new Date(a.last_message.timestamp);
        }
      }
      return a.nickname.localeCompare(b.nickname);
    });

    onlineUsers.forEach(user => {
      container.appendChild(createCompactUserItem(user));
    });
  }

  // Add offline users
  if (offlineUsers.length > 0) {
    const offlineHeader = document.createElement('div');
    offlineHeader.className = 'user-section-header';
    offlineHeader.textContent = 'Offline';
    container.appendChild(offlineHeader);

    // Sort offline users by last message time or alphabetically
    offlineUsers.sort((a, b) => {
      if (a.has_chat_history && b.has_chat_history) {
        if (a.last_message?.timestamp && b.last_message?.timestamp) {
          return new Date(b.last_message.timestamp) - new Date(a.last_message.timestamp);
        }
      }
      return a.nickname.localeCompare(b.nickname);
    });

    // Only show first 5 offline users to save space
    const displayedOfflineUsers = offlineUsers.slice(0, 5);
    displayedOfflineUsers.forEach(user => {
      container.appendChild(createCompactUserItem(user));
    });

    // Add "Show more" button if there are more offline users
    if (offlineUsers.length > 5) {
      const showMoreBtn = document.createElement('div');
      showMoreBtn.className = 'show-more-btn';
      showMoreBtn.textContent = `Show ${offlineUsers.length - 5} more...`;
      showMoreBtn.addEventListener('click', () => {
        // Navigate to full chat page
        history.pushState(null, "", "/chat");
        renderUsersList();
      });
      container.appendChild(showMoreBtn);
    }
  }
}

// Export the public API
export { initChat, renderUsersList, renderChatInterface, navigateToChat, populateOnlineUsersList, fetchAllUsers };