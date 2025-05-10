import { renderRegistrationPage } from "/static/js/register.js";
import { renderLoginPage } from "/static/js/login.js";
import { renderErrorPage } from "/static/js/error.js";
import { renderIndexPage  } from "/static/js/index.js";
import { renderUsersList } from "/static/js/chat.js";

document.addEventListener("DOMContentLoaded", () => {
    handleRouteChange(); // Run when the page loads

    // Listen for back/forward navigation
    window.addEventListener("popstate", handleRouteChange);

    // Attach event listeners to links for single-page navigation
    document.body.addEventListener("click", (event) => {
        const target = event.target.closest("a"); // Check if an <a> was clicked
        if (target && target.getAttribute("href") && !target.target) {
            event.preventDefault(); // Prevent full page reload
            navigateTo(target.getAttribute("href")); // Handle navigation
        }
    });
});

// Function to handle navigation updates
export function navigateTo(url) {
    history.pushState(null, "", url); // Change URL without reloading
    handleRouteChange(); // Trigger re-rendering logic
}

// Function to detect the current route and update the DOM
function handleRouteChange() {
    const path = window.location.pathname;
    const app = document.getElementById("app");
    app.innerHTML = ""; // Clear only the app div

    // Extract chat ID if it's a chat route
      if (path.startsWith('/chat/')) {
        const userId = path.split('/').pop();
        if (userId && userId !== 'chat') {
          // Fetch user details and render chat interface
          fetch(`/users/${userId}`)
            .then(response => response.json())
            .then(user => {
              navigateToChat(user);
            })
            .catch(error => {
              console.error('Error fetching user:', error);
              // Fallback to users list if user not found
              renderUsersList();
            });
        }
        return;
      }

    switch (path) {
        case "/":
            renderIndexPage();
            break;
        case "/login":
            renderLoginPage();
            break;
        case "/register":
            renderRegistrationPage();
            break;
        default:
            window.history.pushState({}, "", "/404") // Re-route page URL without reloading
            renderErrorPage();
            break;
    }
}