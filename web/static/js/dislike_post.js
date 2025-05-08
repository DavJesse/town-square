import { navigateTo } from "/static/js/routes.js";
import { renderErrorPage } from "/static/js/error.js";
import { updatePostEngagement } from '/static/js/like_post.js';

export function handleDislikePost(postId, likeContainer, dislikeContainer) {
    let fetchLink = '/posts/dislike';
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
            }
            renderErrorPage(response.statusText, response.status);            
            throw new Error(`HTTP error! status: ${response.status}`);
        } else {
            updatePostEngagement(fetchLink, likeContainer, dislikeContainer);
        }        
    })

    .catch(error => {
        console.error('Error:', error);
    });
}
