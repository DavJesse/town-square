
        // WebSocket Manager - Handles WebSocket connections
        class WebSocketManager {
            constructor(messageStore) {
                this.ws = null;
                this.messageStore = messageStore;
                this.messageHandlers = [];
                this.reconnectAttempts = 0;
                this.maxReconnectAttempts = 5;
                this.reconnectDelay = 2000;
            }

            connect() {
                if (this.ws && this.ws.readyState === WebSocket.OPEN) {
                    return;
                }

                const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
                const wsUrl = `${protocol}//${window.location.host}/ws`;

                this.ws = new WebSocket(wsUrl);

                this.ws.onopen = () => {
                    console.log('WebSocket connection established');
                    this.reconnectAttempts = 0;
                };

                this.ws.onmessage = (event) => {
                    console.log("Raw WebSocket message received:", event.data);
                    try {
                        const data = JSON.parse(event.data);
                        console.log("Parsed WebSocket data:", data);
                        this.handleMessage(data);
                    } catch (error) {
                        console.error('Error parsing WebSocket message:', error);
                    }
                };

                this.ws.onclose = (event) => {
                    console.log('WebSocket connection closed:', event.code, event.reason);
                    this.attemptReconnect();
                };

                this.ws.onerror = (error) => {
                    console.error('WebSocket error:', error);
                };
            }

            attemptReconnect() {
                if (this.reconnectAttempts >= this.maxReconnectAttempts) {
                    console.log('Max reconnect attempts reached');
                    return;
                }

                this.reconnectAttempts++;
                console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);
                
                setTimeout(() => {
                    this.connect();
                }, this.reconnectDelay * this.reconnectAttempts);
            }

            send(message) {
                if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
                    console.error('WebSocket is not connected');
                    return false;
                }
                
                try {
                    this.ws.send(JSON.stringify(message));
                    return true;
                } catch (error) {
                    console.error('Error sending message via WebSocket:', error);
                    return false;
                }
            }

            registerMessageHandler(handler) {
                if (typeof handler === 'function' && !this.messageHandlers.includes(handler)) {
                    this.messageHandlers.push(handler);
                }
            }

            unregisterMessageHandler(handler) {
                this.messageHandlers = this.messageHandlers.filter(h => h !== handler);
            }

            // In your WebSocketManager class, add a specific method to handle initialization
            handleMessage(message) {
                // Check specifically for the init message type
                if (message.type === 'init') {
                    console.log('Received init message with userID:', message.userID);
                    // Store userID where needed
                    window.userID = message.userID;
                    window.State.setCurrentUser(message.user)
                    let user = window.State.getCurrentUser()
                    console.log("GETTED_CURR: ", user)
                    // You could also store it on the WebSocketManager instance
                    this.userID = message.userID;
                    // Dispatch an event for other parts of your application
                    window.dispatchEvent(new CustomEvent('userInitialized', { 
                        detail: { userID: message.userID } 
                    }));
                }
                
                // Continue with notifying all registered handlers
                this.messageHandlers.forEach(handler => {
                    try {
                        handler(message);
                    } catch (error) {
                        console.error('Error in message handler:', error);
                    }
                });
            }

            close() {
                if (this.ws) {
                    this.ws.close();
                    this.ws = null;
                }
            }
        }