import { renderRegistrationPage } from "/static/js/register.js";
import { renderLoginPage } from "/static/js/login.js";

// routes.js
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
function navigateTo(url) {
    if (path.startsWith("/posts/")) {
        const postId = path.split("/")[2]; // Extract post ID
        app.innerHTML = `<h1>Post ${postId}</h1><p>Details of the post...</p>`;
        return;
    }
    
    history.pushState(null, "", url); // Change URL without reloading
    handleRouteChange(); // Handle the new route
}

// Function to detect the current route and update the DOM
function handleRouteChange() {
    const path = window.location.pathname;
    console.log(`Current path: ${path}`);

    document.body.innerHTML = ""

    switch (path) {
        case "/":
            document.body.innerHTML = "<h1>Home Page</h1><p>Welcome to the forum!</p>";
            break;
        case "/login":
            renderLoginPage();
            break;
        case "/register":
            renderRegistrationPage();
            break;
        case "/posts":
            document.body.innerHTML = "<h1>Posts</h1><p>Here are some posts.</p>";
            break;
        default:
            document.body.innerHTML = "<h1>404</h1><p>Page not found.</p>";
            break;
    }
}
