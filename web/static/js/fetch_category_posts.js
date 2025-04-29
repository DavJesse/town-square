import { navigateTo } from '/static/js/routes.js';
import { renderErrorPage } from '/static/js/error.js';
import { populatePosts, setToggleEventListeners } from '/static/js/profile.js';

// Function to fetch posts by category
export function fetchPostsPerCategogy(categoryID) {
    fetch(`/categories/${categoryID}`, {
        method: 'GET',
        credentials: 'include',
        headers: {
            'Accept': 'application/json'
        }
    })

    .then(response => {
        if (!response.ok) {
                    // Check for unauthorized resposes
                    if (response.status === 401) {
                        navigateTo('/login');
                        return Promise.reject('Unauthorized');
                    } else {
                        // throw new Error(`HTTP error! status: ${response.status}`);
                        return response.json().then(errorData => {
                            renderErrorPage(errorData.Issue || response.statusText, response.status);
                            return Promise.reject(errorData.Issue || 'Error occurred');
                        });
                    }
                }
                return response.json();
    })

    .then(data => {
        // Clear current posts for repopulation
        let mainContent = document.getElementById('posts_container');
        mainContent.innerHTML = '';

        let categoryButton = document.getElementById(categoryID);
        categoryButton.style.backgroundColor = '#0000EE';

        const posts = data.data.posts
        const likedPosts = data.data.liked_posts

        // Update posts
        populatePosts(posts);
        setToggleEventListeners(likedPosts);
    })
    .catch(error => {
        console.error('Error fetching category posts:', error);
    });
}