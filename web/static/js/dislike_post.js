import { navigateTo } from "/static/js/routes.js";
import { renderErrorPage } from "/static/js/error.js";

// Add event listener when the document loads
document.addEventListener('DOMContentLoaded', () => {
    // Use event delegation to handle like button clicks
    document.addEventListener('click', (e) => {
        const dislikeBtn = e.target.closest('#dislike_link');
        if (dislikeBtn) {
            e.preventDefault();
            const postId = dislikeBtn.dataset.postId;
            handleDislikePost(postId, dislikeBtn);
        }
    });
});

export function handleDislikePost(postId, dislikeContainer) {
    fetch('/posts/dislike', {
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
            renderErrorPage(response.statusText, response.status);            
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
    })
    .then(data => {
        if (data.success) {
            // Update the likes count in the container
            const countSpan = dislikeContainer.querySelector('#engagement_count');
            if (countSpan) {
                countSpan.textContent = data.likesCount;
            }
            
            // Optional: Add visual feedback
            const likeLink = dislikeContainer.querySelector('#ldisike_link');
            if (likeLink) {
                likeLink.classList.add('disliked');
                setTimeout(() => likeLink.classList.remove('disliked'), 200);
            }
        } else if (data.message === "User is not logged in. Please log in to try again") {
            navigateTo('/login');
        } else {
            console.error('Dislike failed:', data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
}
