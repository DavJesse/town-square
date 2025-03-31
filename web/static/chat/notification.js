// Notification Manager - Handles notifications
class NotificationManager {
    constructor() {
        this.container = this.createContainer();
        this.notifications = [];
    }

    createContainer() {
        let container = document.getElementById('notification-container');
        if (!container) {
            container = document.createElement('div');
            container.id = 'notification-container';
            document.body.appendChild(container);
        }
        return container;
    }

    show(message) {
        const notification = this.createNotification(message);
        this.container.appendChild(notification);
        this.notifications.push(notification);
        
        // Auto remove after 10 seconds
        setTimeout(() => {
            this.remove(notification);
        }, 10000);
    }

    createNotification(message) {
        const notification = document.createElement('div');
        notification.className = 'notification';
        
        notification.innerHTML = `
            <div class="notification-content">
                ${message}
            </div>
            <button class="notification-close">&times;</button>
        `;
        
        // Add click to navigate to messages if it's a message notification
        if (message.startsWith('New message from')) {
            notification.addEventListener('click', (e) => {
                // Don't navigate if clicking the close button
                if (!e.target.closest('.notification-close')) {
                    this.remove(notification);
                }
            });
        }
        
        // Add close button functionality
        const closeButton = notification.querySelector('.notification-close');
        closeButton.addEventListener('click', (e) => {
            e.stopPropagation(); // Prevent notification click handler from firing
            this.remove(notification);
        });
        
        return notification;
    }

    remove(notification) {
        if (notification && notification.parentElement) {
            notification.remove();
            this.notifications = this.notifications.filter(n => n !== notification);
        }
    }

    newMessage(data) {
        this.show(`New message from ${data.sender}`);
    }
}
