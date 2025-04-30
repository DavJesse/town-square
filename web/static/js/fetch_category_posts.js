import { navigateTo } from '/static/js/routes.js';
import { renderErrorPage } from '/static/js/error.js';
import { populatePosts, setToggleEventListeners } from '/static/js/profile.js';

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
        // Clear current posts for repopulation
        let mainContent = document.getElementById('posts_container');
        mainContent.innerHTML = '';

        // Restore active button to white
        let activeButton = document.querySelector(`.active`);
        if (activeButton) {
            let activeText = activeButton.querySelector('p');
            activeButton.style.backgroundColor = '#FFFFFF';
            if (activeText) {
                activeText.style.color = '#0000EE';
            }
            activeButton.classList.remove('active');
            activeButton.classList.add('inactive');
        }

        // Set category button to blue
        let categoryButton = document.getElementById('category_button_' + categoryID);
        categoryButton.style.backgroundColor = '#0000EE';

         // Set category button text to white
         let categoryText = document.getElementById(`category_text_${categoryID}`);
         categoryText.style.color = '#FFFFFF';

        // Assign data to variables for easy acces
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