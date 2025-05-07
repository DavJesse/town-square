import { navigateTo } from "/static/js/routes.js";
import { renderErrorPage } from "/static/js/error.js";

export function handleLikeComment(commentID, postID, likeContainer, dislikeContainer) {
    let payload = JSON.stringify({
        commentID: commentID,
        postID: postID
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
        if (countSpan) {
            countSpan.textContent = data.likes_count;
        }
        
        // Update the dislike count in the container
        const dislikeCountSpan = likeContainer.querySelector('#comment_dislike_count');
        if (countSpan) {
            countSpan.textContent = data.dislikes_count;
        }
    })

    .catch(error => {
        console.error('Error:', error);
    });
}
