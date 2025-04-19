import { navigateTo } from "/static/js/routes.js";
import { renderErrorPage } from "/static/js/error.js";

export function handleLikePost(postId, likeButton) {
    fetch('/posts/like', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            'Accept': 'application/json',
        },
        body: `post-id=${postId}`,
        credentials: 'include'
    })

    .then(response => {
        if (!response.ok) {
            if (response.status === 401) {
                navigateTo('/login');
                return;
            } else {
                renderErrorPage(response.statusText, response.status);
                return;
            }
        }
        
        return response.json()
    })

    .then(data => {
        if (data.success) {
            // Update the likes count in the button
            const countSpan = likeButton.querySelector('.likes-count');
            countSpan.textContent = data.likesCount;
            
            // Optional: Add visual feedback
            likeButton.classList.add('liked');
            setTimeout(() => likeButton.classList.remove('liked'), 200);
        } else if (data.message === "User is not logged in. Please log in to try again") {
            navigateTo('/login');
            return
        } else {
            console.error('Like failed:', data.message);
        }
    })

    .catch(error => {
        console.error('Error:', error);
    });
}