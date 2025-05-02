import { handleLikeComment } from '/static/js/like_comment.js';
import { handleDislikeComment } from '/static/js/dislike_comment.js';

export function populateComments(postCard, postID, comments) {
    let commentsSection = document.createElement('div');
    commentsSection.id = 'comments_section';
    postCard.appendChild(commentsSection);

    // Add comment form
    let commentForm = document.createElement('form');
    let inputField = document.createElement('textarea');
    let submitButton = document.createElement('button');
    commentForm.method = 'POST';
    commentForm.action = '/comment';
    inputField.id = 'comment_input_field';
    submitButton.id = 'comment_submit_button';
    submitButton.type = 'submit';
    inputField.placeholder = 'Write a comment...';
    submitButton.textContent = 'comment';
    inputField.maxLength = 750;
    inputField.name = 'comment';
    inputField.required = true;
    commentForm.appendChild(inputField);
    commentForm.appendChild(submitButton);
    commentsSection.appendChild(commentForm);

    comments.forEach(comment => {
        // Create comment container
        let commentContainer = document.createElement('div');
        commentContainer.id = 'comment_container';
        commentsSection.appendChild(commentContainer);

        // Create creator image container
        let creatorImageContainer = document.createElement('div');
        let creatorImage = document.createElement('img');
        creatorImageContainer.id = 'comment_creator_image_container';
        creatorImage.id = 'comment_creator_image';
        creatorImage.src = comment.creator_image? `/static/images/${comment.creator_image}` : '/static/user-circle-svgrepo-com.svg';
        creatorImageContainer.appendChild(creatorImage);
        commentContainer.appendChild(creatorImageContainer);

        // Create info info container
        let creatorInfoContainer = document.createElement('div');
        let creatorName = document.createElement('h3');
        let creatorUsername = document.createElement('h3');
        creatorInfoContainer.id = 'comment_creator_info_container';
        creatorName.id = 'comment_creator_name';
        creatorUsername.id = 'comment_creator_username';
        creatorName.textContent = `${comment.creator_first_name.charAt(0).toUpperCase()}${comment.creator_first_name.slice(1)} ${comment.creator_last_name.charAt(0).toUpperCase()}${comment.creator_last_name.slice(1)}`;
        creatorUsername.textContent = `@${comment.creator_username}`;
        creatorInfoContainer.appendChild(creatorName);
        creatorInfoContainer.appendChild(creatorUsername);
        commentContainer.appendChild(creatorInfoContainer);

        // Create comment content container
        let commentContentContainer = document.createElement('div');
        let commentContent = document.createElement('p');
        commentContentContainer.appendChild(commentContent)
        commentContainer.appendChild(commentContentContainer);

        // Create engagement container
        let commentEngagement = document.createElement('div');
        let likeLink = document.createElement('button');
        let dislikeLink = document.createElement('button');
        let likeCount = document.createElement('p');
        let dislikeCount = document.createElement('p');
        let likeIcon = document.createElement('span');
        let dislikeIcon = document.createElement('span');

        // Set IDs and classes
        likeLink.id = 'comment_like_link';
        dislikeLink.id = 'comment_dislike_link';
        likeCount.id = 'comment_engagement_count';
        dislikeCount.id = 'comment_engagement_count';
        likeIcon.classList.add('material-symbols-outlined');
        dislikeIcon.classList.add('material-symbols-outlined');
        likeIcon.id = 'comment_engagement_icon';
        dislikeIcon.id = 'comment_engagement_icon';
        commentIcon.id = 'comment_engagement_icon';

        // Set content
        likeCount.textContent = comment.likes_count;
        dislikeCount.textContent = comment.dislikes_count;            
        likeIcon.textContent = 'thumb_up';
        dislikeIcon.textContent = 'thumb_down';

        // Add event listeners
        likeLink.addEventListener('click', (e) => {
            e.preventDefault();
            handleLikeComment(comment.uuid, likeLink);
        });

        dislikeLink.addEventListener('click', (e) => {
            e.preventDefault();
            handleDislikeComment(postID, dislikeLink);
        });

        // Append elements to engagement containers
        likeLink.appendChild(likeCount);
        likeLink.appendChild(likeIcon);
        dislikeLink.appendChild(dislikeCount);
        dislikeLink.appendChild(dislikeIcon);
        commentEngagement.appendChild(likeLink);
        commentEngagement.appendChild(dislikeLink);
        commentContainer.appendChild(commentEngagement);

        // Add 'hide comments' button
        let hideComments = document.createElement('button');
        hideComments.id = 'hide_comments';
        hideComments.textContent = 'hide comments...';
        hideComments.addEventListener('click', () => {
            commentsSection.style.display = 'none';
        });
        commentContainer.appendChild(hideComments);
    });
}