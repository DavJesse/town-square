import { navigateTo } from "/static/js/routes.js";
import { renderErrorPage } from "/static/js/error.js";

export function handleLikeComment(commentID, postID, likeContainer) {
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
        if (response.ok) {
            UpdateCommentLikesAndDislikes(commentID, postID);            
        } else {
            throw new Error(`Unable To Like Comment`);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

// .then(data => {
//     if (data.success) {
//         // Update the likes count in the container
//         const countSpan = likeContainer.querySelector('#comment_engagement_count');
//         if (countSpan) {
//             countSpan.textContent = data.likesCount;
//         }
        
//     } else {
//         console.error('Like failed:', data.message);
//     }
// })