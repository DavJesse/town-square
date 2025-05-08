import { navigateTo } from "/static/js/routes.js";
import { renderErrorPage } from "/static/js/error.js";

export function handleLikeComment(commentID, likeContainer, dislikeContainer) {
    let payload = JSON.stringify({
        commentID: commentID
    });

    fetch('/comments/like', {
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
            updateCommentEngagement(likeContainer, dislikeContainer);            
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

function updateCommentEngagement(likeContainer, dislikeContainer) {
    fetch('/comments/like', {
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
        const likeCountSpan = likeContainer.querySelector('#comment_like_count');
        if (likeCountSpan) {
            likeCountSpan.textContent = data.like_count;
        }
        
        // Update the dislike count in the container
        const dislikeCountSpan = dislikeContainer.querySelector('#comment_dislike_count');
        if (dislikeCountSpan) {
            dislikeCountSpan.textContent = data.dislike_count;
        }
    })

    .catch(error => {
        console.error('Error:', error);
    });
}
