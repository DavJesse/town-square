import { renderRegistrationPage } from "/static/js/register.js";
import { renderLoginPage } from "/static/js/login.js";
import { renderErrorPage } from "/static/js/error.js";
import { renderProfilePage } from "/static/js/profile.js";
import { renderIndexPage  } from "/static/js/index.js";
import { handleLikePost } from "/static/js/like_post.js";

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
        case "/posts":
            app.innerHTML = "<h1>Posts</h1><p>Here are some posts.</p>";
            break;
        case "/profile":
            console.log("PROFILE")
            renderProfilePage();
            break;
        case "/posts/like":
            handleLikePost();
            break;
        default:
            window.history.pushState({}, "", "/404") // Re-route page URL without reloading
            renderErrorPage();
            break;
    }
}
