import { navigateTo } from "/static/js/routes.js";
import { renderErrorPage } from "/static/js/error.js";

export function handleLikePost(postId, likeContainer) {
    let fetchLink = '/posts/like';
    fetch(fetchLink, {
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

export function updatePostEngagement(fetchLink, likeContainer, dislikeContainer) {
    fetch(fetchLink, {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
        },
        credentials: 'include'
    })

    .then(response => {
        if (!response.ok) {
            if (response.status === 401) {
                navigateTo('/login');
                return;
            } else {
                renderErrorPage(response.statusText, response.status);
            }   
            throw new Error(`Unable To Like Comment`);
        }
        return response.json();
    })

    .then(data => {
         // Update the like count in the container
        const likeCountSpan = likeContainer.querySelector('#post_like_count');
        if (likeCountSpan) {
            likeCountSpan.textContent = data.like_count;
        }
        
        // Update the dislike count in the container
        const dislikeCountSpan = dislikeContainer.querySelector('#post_dislike_count');
        if (dislikeCountSpan) {
            dislikeCountSpan.textContent = data.dislike_count;
        }
    })

    .catch(error => {
        console.error('Error:', error);
    });
}
