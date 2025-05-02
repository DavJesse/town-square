export function populateComments(postCard, comments) {
    let commentsSection = document.createElement('div');
    commentsSection.id = 'comments_section';
    postCard.appendChild(commentsSection);

    // Add comment form
    let commentForm = document.createElement('form');
    let inputField = document.createElement('textarea');
    let submitButton = document.createElement('button');
    inputField.id = 'comment_input_field';
    submitButton.id = 'comment_submit_button';
    inputField.type = 'text';
    submitButton.type = 'submit';
    inputField.placeholder = 'Write a comment...';
    inputField.maxLength = 750;
    inputField.name = 'comment';
    inputField.required = true;
    commentForm.appendChild(inputField);
    commentForm.appendChild(submitButton);

    comments.forEach(comment => {
        // Create comment container
        let commmentContainer = document.createElement('div');
        commmentContainer.id = 'comment_container';
        commentsSection.appendChild(commmentContainer);

        //
    });
}