import { navigateTo } from "/static/js/routes.js";
import { renderErrorPage } from "/static/js/error.js";

// Add event listener when the document loads
document.addEventListener('DOMContentLoaded', () => {
    // Use event delegation to handle like button clicks
    document.addEventListener('click', (e) => {
        const likeBtn = e.target.closest('#like_link');
        if (likeBtn) {
            e.preventDefault();
            const postId = likeBtn.dataset.postId;
            handleLikePost(postId, likeBtn);
        }
    });
});

export function handleLikePost(postId, likeContainer) {
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
            }
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
    })
    .then(data => {
        if (data.success) {
            // Update the likes count in the container
            const countSpan = likeContainer.querySelector('#engagement_count');
            if (countSpan) {
                countSpan.textContent = data.likesCount;
            }
            
            // Optional: Add visual feedback
            const likeLink = likeContainer.querySelector('#like_link');
            if (likeLink) {
                likeLink.classList.add('liked');
                setTimeout(() => likeLink.classList.remove('liked'), 200);
            }
        } else if (data.message === "User is not logged in. Please log in to try again") {
            navigateTo('/login');
        } else {
            console.error('Like failed:', data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
}
