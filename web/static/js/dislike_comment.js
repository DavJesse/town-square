import { navigateTo } from "/static/js/routes.js";
import { renderErrorPage } from "/static/js/error.js";
import { updateCommentEngagement } from '/static/js/like_comment.js';

export function handleDislikeComment(commentID, likeContainer, dislikeContainer) {
    let fetchLink = '/comments/dislike';
    let payload = JSON.stringify({
        commentID: commentID
    });

    fetch(fetchLink, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
        },
        body: payload,
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
        } else {
            updateCommentEngagement(fetchLink, likeContainer, dislikeContainer);            
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

