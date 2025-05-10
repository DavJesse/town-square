import { navigateTo } from '/static/js/routes.js';
import { renderErrorPage } from '/static/js/error.js';
import { populatePosts, setToggleEventListeners } from '/static/js/populate_posts.js';

// Function to fetch posts by category
export function fetchPostsPerCategory(categoryID) {
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
                            renderErrorPage(response.statusText, response.status);
                            return Promise.reject(errorData.Issue || 'Error occurred');
                        });
                    }
                }
                return response.json();
    })

    .then(data => {
        let postsContainer = document.getElementById('posts_container');
        let allPostsButton = document.getElementById('all_posts_button');
        let likedPostsButton = document.getElementById('liked_posts_button');
        let myPostsButton = document.getElementById('my_posts_button');

        // Clear current posts for repopulation
        postsContainer.innerHTML = '';

        let categoryContainer = document.getElementById('category_content_container');
        let activeButton = categoryContainer.querySelector(`.active`);
        
        // Restore active button to white
        if (activeButton) {
            activeButton.style.backgroundColor = '#FFFFFF';
            activeButton.style.color = '#0000EE';
            
            // Update class names for correct identification
            activeButton.classList.remove('active');
            activeButton.classList.add('inactive');
        }

        // Set category button to blue
        let categoryButton = document.getElementById('category_button_' + categoryID);
        categoryButton.style.backgroundColor = '#0000EE';
        categoryButton.classList.remove('inactive');
        categoryButton.classList.add('active');

         // Set category button text to white
        categoryButton.style.color = '#FFFFFF';

        // Assign data to variables for easy acces
        const posts = data.data.posts
        const likedPosts = data.data.liked_posts

        // Update posts
        populatePosts(posts);

        // Set active button to all posts
        if (likedPostsButton.classList.contains('active')) {
            likedPostsButton.classList.remove('active');
            allPostsButton.classList.add('active');
        } else if (myPostsButton.classList.contains('active')) {
            myPostsButton.classList.remove('active');
            allPostsButton.classList.add('active');
        }
        setToggleEventListeners(posts, likedPosts, []);
    })

    .catch(error => {
        console.error('Error fetching category posts:', error);
    });
}