/**
 * chat.js - Chat functionality for embedding in applications that handle authentication via cookies
 */
import { renderLogoutButton } from '/static/js/logout_button.js';

// Global variables
let myUserId = null;
let selectedUser = null;
let ws = null; // WebSocket connection
let userNotifications = {}; // Track notifications per user
let lastLoadedMessages = []; // Track the last batch of loaded messages
let isLoading = false;
let firstMessageId = null; // ID of the first (oldest) message loaded
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
  // 500ms is a good balance between responsiveness and preventing spam
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

      // Add a more noticeable loading indicator at the top of the message area
      const loadingIndicator = document.createElement('div');
      loadingIndicator.className = 'loading-indicator';
      loadingIndicator.innerHTML = `
        <div class="loading-spinner"></div>
        <div class="loading-text">Loading older messages...</div>
      `;
      messageArea.insertBefore(loadingIndicator, messageArea.firstChild);

      try {
        console.log(`Loading older messages: firstMessageId=${firstMessageId}, limit=${limit}`);

        // Add a deliberate delay to make loading more noticeable (800ms)
        await new Promise(resolve => setTimeout(resolve, 800));

        // Use the ID of the first (oldest) message as the cursor for pagination
        const messages = await fetchMessageHistory(selectedUser.id, firstMessageId, limit);

        // Remove loading indicator with a slight delay for better UX
        setTimeout(() => {
          if (loadingIndicator && loadingIndicator.parentNode) {
            loadingIndicator.parentNode.removeChild(loadingIndicator);
          }
        }, 200);

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
        // Remove loading indicator if there's an error
        if (loadingIndicator && loadingIndicator.parentNode) {
          loadingIndicator.parentNode.removeChild(loadingIndicator);
        }
        showToast('Error loading older messages', 'error');
      } finally {
        isLoading = false;
      }
    }
  }, 500); // 500ms throttle time

  messageArea.addEventListener('scroll', handleScroll);

  // We're not going to auto-load more messages on initial load
  // This ensures users see the loading process when they scroll up

  // Still listen for window resize events
  window.addEventListener('resize', () => {
    // Only check if the window is resized significantly
    if (!isLoading && selectedUser && messageArea.scrollHeight <= messageArea.clientHeight) {
      // If the message area is too small, we might need to load more
      // But we'll let the user scroll to trigger loading
    }
  });
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

      // We only need to track the oldest message ID for pagination
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

  // Reset message ID for fresh load
  firstMessageId = null;

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

  // Clear the app container
  app.innerHTML = "";
  
  // Import and render the navbar
  import('/static/js/navbar.js').then(module => {
    const { renderNavBar } = module;
    renderNavBar();

    // Create the chat page container
    const chatPageContainer = document.createElement('div');
    chatPageContainer.id = 'chat_page';
    chatPageContainer.classList.add('chat-page-container');
    app.appendChild(chatPageContainer);

    // Create the chat container
    const chatContainer = document.createElement('div');
    chatContainer.classList.add('chat-container');
    chatPageContainer.appendChild(chatContainer);

    // Add mini notification
    const miniNotification = document.createElement('div');
    miniNotification.id = 'mini-notification';
    miniNotification.classList.add('mini-notification');
    chatContainer.appendChild(miniNotification);

    // Create chat sidebar
    const chatSidebar = document.createElement('div');
    chatSidebar.classList.add('chat-sidebar');

    // Create sidebar header
    const sidebarHeader = document.createElement('div');
    sidebarHeader.classList.add('sidebar-header');

    const sidebarTitle = document.createElement('h2');
    sidebarTitle.classList.add('sidebar-title');
    sidebarTitle.textContent = 'Messages';

    const searchContainer = document.createElement('div');
    searchContainer.classList.add('search-container');

    const searchInput = document.createElement('input');
    searchInput.type = 'text';
    searchInput.classList.add('search-input');
    searchInput.placeholder = 'Search users...';

    searchContainer.appendChild(searchInput);
    sidebarHeader.appendChild(sidebarTitle);
    sidebarHeader.appendChild(searchContainer);
    chatSidebar.appendChild(sidebarHeader);

    // Create users container
    const usersContainer = document.createElement('div');
    usersContainer.classList.add('users-container');

    const usersList = document.createElement('ul');
    usersList.id = 'online-users-list';
    usersList.classList.add('user-list');

    usersContainer.appendChild(usersList);
    chatSidebar.appendChild(usersContainer);
    chatContainer.appendChild(chatSidebar);

    // Create main chat area
    const chatMain = document.createElement('div');
    chatMain.classList.add('chat-main');

    // Create chat header
    const chatHeader = document.createElement('div');
    chatHeader.classList.add('chat-header');

    const mobileToggle = document.createElement('button');
    mobileToggle.id = 'mobile-toggle';
    mobileToggle.classList.add('mobile-toggle');

    const mobileToggleIcon = document.createElement('span');
    mobileToggleIcon.classList.add('mobile-toggle-icon');
    mobileToggleIcon.textContent = '👥';

    mobileToggle.appendChild(mobileToggleIcon);

    const chatTitle = document.createElement('h3');
    chatTitle.id = 'chat-title';
    chatTitle.classList.add('chat-title');
    chatTitle.textContent = `Chat with ${user.nickname}`;

    const backButton = document.createElement('button');
    backButton.id = 'back-button';
    backButton.classList.add('back-button', 'navbar__button');

    const backIcon = document.createElement('span');
    backIcon.classList.add('back-icon');
    backIcon.textContent = '←';

    const backText = document.createElement('span');
    backText.textContent = 'Back';

    backButton.appendChild(backIcon);
    backButton.appendChild(backText);

    // Create container for right-side buttons
    const headerButtons = document.createElement('div');
    headerButtons.classList.add('header-buttons');
    
    // Add back button to the header buttons container
    headerButtons.appendChild(backButton);
    
    chatHeader.appendChild(mobileToggle);
    chatHeader.appendChild(chatTitle);
    chatHeader.appendChild(headerButtons);
    chatMain.appendChild(chatHeader);

    // Create message container
    const messageContainer = document.createElement('div');
    messageContainer.classList.add('message-container');

    const messageArea = document.createElement('div');
    messageArea.id = 'message-area';
    messageArea.classList.add('message-area');

    const newMessageNotification = document.createElement('div');
    newMessageNotification.id = 'new-message-notification';
    newMessageNotification.classList.add('new-message-notification');
    newMessageNotification.textContent = 'New messages ↓';

    messageContainer.appendChild(messageArea);
    messageContainer.appendChild(newMessageNotification);
    chatMain.appendChild(messageContainer);

    // Create message input container
    const messageInputContainer = document.createElement('div');
    messageInputContainer.classList.add('message-input-container');

    const messageInput = document.createElement('input');
    messageInput.id = 'message-input';
    messageInput.classList.add('message-input');
    messageInput.type = 'text';
    messageInput.placeholder = 'Type your message...';

    const sendButton = document.createElement('button');
    sendButton.id = 'send-button';
    sendButton.classList.add('send-button', 'navbar__button');
    sendButton.textContent = 'Send';

    const sendIcon = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    sendIcon.classList.add('send-icon');
    sendIcon.setAttribute('width', '16');
    sendIcon.setAttribute('height', '16');
    sendIcon.setAttribute('viewBox', '0 0 24 24');
    sendIcon.setAttribute('fill', 'none');
    sendIcon.setAttribute('stroke', 'currentColor');
    sendIcon.setAttribute('stroke-width', '2');
    sendIcon.setAttribute('stroke-linecap', 'round');
    sendIcon.setAttribute('stroke-linejoin', 'round');

    const line = document.createElementNS('http://www.w3.org/2000/svg', 'line');
    line.setAttribute('x1', '22');
    line.setAttribute('y1', '2');
    line.setAttribute('x2', '11');
    line.setAttribute('y2', '13');

    const polygon = document.createElementNS('http://www.w3.org/2000/svg', 'polygon');
    polygon.setAttribute('points', '22 2 15 22 11 13 2 9 22 2');

    sendIcon.appendChild(line);
    sendIcon.appendChild(polygon);
    sendButton.appendChild(sendIcon);

    messageInputContainer.appendChild(messageInput);
    messageInputContainer.appendChild(sendButton);
    chatMain.appendChild(messageInputContainer);

    chatContainer.appendChild(chatMain);

    // Add back button handler
    if (backButton) {
      backButton.addEventListener('click', () => {
        // Navigate back to the users list
        history.pushState(null, "", "/chat");
        renderUsersList();
      });
    }

    // Add mobile toggle button handler
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
  });
}

/**
 * Render the user list view
 */
function renderUsersList() {
  const app = document.getElementById('app');
  if (!app) return;

  // Clear the app container
  app.innerHTML = "";

  // Import and render the navbar
  import('/static/js/navbar.js').then(module => {
    const { renderNavBar } = module;
    renderNavBar();

    // Create the chat page container
    const chatPageContainer = document.createElement('div');
    chatPageContainer.id = 'chat_page';
    chatPageContainer.classList.add('chat-page-container');
    app.appendChild(chatPageContainer);

    // Create the chat container
    const chatContainer = document.createElement('div');
    chatContainer.classList.add('chat-container');
    chatPageContainer.appendChild(chatContainer);

    // Add mini notification
    const miniNotification = document.createElement('div');
    miniNotification.id = 'mini-notification';
    miniNotification.classList.add('mini-notification');
    chatContainer.appendChild(miniNotification);

    // Create chat sidebar
    const chatSidebar = document.createElement('div');
    chatSidebar.classList.add('chat-sidebar');

    // Create sidebar header
    const sidebarHeader = document.createElement('div');
    sidebarHeader.classList.add('sidebar-header');

    const sidebarTitle = document.createElement('h2');
    sidebarTitle.classList.add('sidebar-title');
    sidebarTitle.textContent = 'Messages';

    const searchContainer = document.createElement('div');
    searchContainer.classList.add('search-container');

    const searchInput = document.createElement('input');
    searchInput.type = 'text';
    searchInput.classList.add('search-input');
    searchInput.placeholder = 'Search users...';

    searchContainer.appendChild(searchInput);
    sidebarHeader.appendChild(sidebarTitle);
    sidebarHeader.appendChild(searchContainer);
    chatSidebar.appendChild(sidebarHeader);

    // Create users container
    const usersContainer = document.createElement('div');
    usersContainer.classList.add('users-container');

    const usersList = document.createElement('ul');
    usersList.id = 'online-users-list';
    usersList.classList.add('user-list');

    usersContainer.appendChild(usersList);
    chatSidebar.appendChild(usersContainer);
    chatContainer.appendChild(chatSidebar);

    // Create main chat area with empty state
    const chatMain = document.createElement('div');
    chatMain.classList.add('chat-main');

    const emptyState = document.createElement('div');
    emptyState.classList.add('empty-state');

    const emptyStateIcon = document.createElement('div');
    emptyStateIcon.classList.add('empty-state-icon');
    emptyStateIcon.textContent = '💬';

    const emptyStateTitle = document.createElement('h2');
    emptyStateTitle.classList.add('empty-state-title');
    emptyStateTitle.textContent = 'Select a conversation';

    const emptyStateText = document.createElement('p');
    emptyStateText.classList.add('empty-state-text');
    emptyStateText.textContent = 'Choose a user from the sidebar to start chatting';

    emptyState.appendChild(emptyStateIcon);
    emptyState.appendChild(emptyStateTitle);
    emptyState.appendChild(emptyStateText);
    chatMain.appendChild(emptyState);

    chatContainer.appendChild(chatMain);

    // Fetch users after rendering the list view
    fetchAllUsers();
  });
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

// Removed unused loading indicator functions

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

  // Filter out current user
  const filteredUsers = users.filter(user => user.id !== myUserId);

  if (filteredUsers.length === 0) {
    const emptyItem = document.createElement('div');
    emptyItem.textContent = 'No users available';
    emptyItem.className = 'empty-user-item';
    container.appendChild(emptyItem);
    return;
  }

  // Split users into those with and without message history
  const usersWithMessages = filteredUsers.filter(user => user.has_chat_history);
  const usersWithoutMessages = filteredUsers.filter(user => !user.has_chat_history);

  // Sort users with messages by last message timestamp
  usersWithMessages.sort((a, b) => {
    if (a.last_message?.timestamp && b.last_message?.timestamp) {
      return new Date(b.last_message.timestamp) - new Date(a.last_message.timestamp);
    }
    return 0;
  });

  // Sort users without messages alphabetically
  usersWithoutMessages.sort((a, b) => a.nickname.localeCompare(b.nickname));

  // Create a compact user item
  const createCompactUserItem = (user) => {
    const userItem = document.createElement('div');
    userItem.className = 'compact-user-item';
    userItem.dataset.userId = user.id;

    // Create status indicator
    const userHandlerContainer = document.createElement('div');
    const statusIndicator = document.createElement('span');
    userHandlerContainer.id = 'user_handler_container';
    statusIndicator.className = `status-indicator ${user.is_online ? 'status-online' : 'status-offline'}`;

    // Create user nickname
    const userNickname = document.createElement('span');
    userNickname.className = 'user-nickname';
    userNickname.textContent = user.nickname;

    // Add notification badge if needed
    let chatIcon = document.getElementById('chat_container');
    if (userNotifications[user.id] && userNotifications[user.id] > 0) {
      const badge = document.createElement('span');
      badge.className = 'notification-badge';
      badge.textContent = userNotifications[user.id];
      chatIcon.appendChild(badge);
    }

    // Assemble user item
    userHandlerContainer.appendChild(statusIndicator);
    userHandlerContainer.appendChild(userNickname);
    userItem.appendChild(userHandlerContainer);

    // Add last message preview if available
    if (user.has_chat_history && user.last_message) {
      const messagePreviewContainer = document.createElement('div');
      const messagePreview = document.createElement('div');
      messagePreviewContainer.id = 'message_preview_container';
      messagePreview.className = 'message-preview';
      messagePreview.textContent = truncateText(user.last_message.content, 25);
      
      const timestamp = document.createElement('div');
      timestamp.className = 'message-time';
      timestamp.textContent = formatMessageTime(new Date(user.last_message.timestamp));
      
      messagePreviewContainer.appendChild(messagePreview);
      messagePreviewContainer.appendChild(timestamp);
      userItem.appendChild(messagePreviewContainer);
    }

    // Add click handler to open chat
    userItem.addEventListener('click', () => {
      // Navigate to chat with this user
      history.pushState(null, "", `/chat/${user.id}`);
      navigateToChat(user);
    });

    return userItem;
  };

  // Add users with message history first
  if (usersWithMessages.length > 0) {
    const historyHeader = document.createElement('div');
    historyHeader.className = 'user-section-header';
    historyHeader.textContent = 'Recent Chats';
    container.appendChild(historyHeader);

    usersWithMessages.forEach(user => {
      container.appendChild(createCompactUserItem(user));
    });
  }

  // Add users without message history
  if (usersWithoutMessages.length > 0) {
    const newUsersHeader = document.createElement('div');
    newUsersHeader.className = 'user-section-header';
    newUsersHeader.textContent = 'Other Users';
    container.appendChild(newUsersHeader);

    // Show first 5 users without message history
    const displayedUsers = usersWithoutMessages.slice(0, 5);
    displayedUsers.forEach(user => {
      container.appendChild(createCompactUserItem(user));
    });

    // Add "Show more" button if there are more users
    if (usersWithoutMessages.length > 5) {
      const showMoreBtn = document.createElement('div');
      showMoreBtn.className = 'show-more-btn';
      showMoreBtn.textContent = `Show ${usersWithoutMessages.length - 5} more...`;
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