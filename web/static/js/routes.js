import { renderRegistrationPage } from "/static/js/register.js";
import { renderLoginPage } from "/static/js/login.js";
import { renderErrorPage } from "/static/js/error.js";
import { renderNavBar } from "/static/js/navbar.js";

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
export async function navigateTo(url) {
    try {
        if (url.startsWith("/posts/")) {
            const postId = url.split("/")[2]; // Extract post ID
            app.innerHTML = `<h1>Post ${postId}</h1><p>Details of the post...</p>`;
            return;
        } 
    
        let response = await fetch(url);

        if (!response.ok) {
            console.error(`Error: ${response.status} - ${response.statusText}`)
            renderErrorPage();
            return;
        }
        
        let content = await response.text();
        document.body.innerHTML = content;

        history.pushState(null, "", url); // Change URL without reloading
        handleRouteChange(); // Handle the new route

    } catch (error) {
        console.error(`Navigation Error: ${error}`)
    }
}

// Function to detect the current route and update the DOM
function handleRouteChange() {
    const path = window.location.pathname;
    const app = document.getElementById("app");
    app.innerHTML = ""; // Clear only the app div

    switch (path) {
        case "/":
            app.innerHTML = "<h1>Home Page</h1><p>Welcome to the forum!</p>";
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
            renderNavBar();
            break;
        default:
            window.history.pushState({}, "", "/404") // Re-route page URL without reloading
            renderErrorPage();
            break;
    }
}
