export function renderCreatePostForm() {
    document.title = 'Create Post';
    let app = document.getElementById('app');
    
    // Create a smoke screen over page comtent
    let smokeScreen = document.createElement('div');
    smokeScreen.id = 'smoke_screen';
    app.appendChild(smokeScreen);

    // Create form container
    let formContainer = document.createElement('div');
    formContainer.id = 'create_post_form_container';
    smokeScreen.appendChild(formContainer);

    // Create forum logo
    let logo = document.createElement('h1');
    logo.textContent = 'Create New Post';
    formContainer.appendChild(logo);

    // Create text to capture user-generated errors
    let errorMessage = document.createElement('p');
    errorMessage.id = 'error_text';
    errorMessage.style.color = 'red';
    formContainer.appendChild(errorMessage);

    // Create form to capture user input
    let form = document.createElement('form');
    form.method = 'POST';
    form.action = '/posts/create';
    form.enctype = 'multipart/form-data';
    formContainer.appendChild(form);

    // Create title field
    let postTitle = document.createElement('input');
    postTitle.type = 'text';
    postTitle.id = 'post_title';
    postTitle.name = 'title';
    postTitle.placeholder = 'subject of your post';
    postTitle.required = true;
    form.appendChild(postTitle);

    // Create content field
    let postContent = document.createElement('textarea');
    postContent.id = 'post_content';
    postContent.name = 'content';
    postContent.placeholder = 'share your thoughts';
    postContent.required = true;
    form.appendChild(postContent);

    // Create image upload container
    let imageUpload = document.createElement('div');
    imageUpload.classList.add('custom-file-upload');
    form.appendChild(imageUpload);

    // Add image upload label
    let uploadButton = document.createElement('label');
    uploadButton.htmlFor = 'image';
    uploadButton.textContent = 'Upload Image (PNG, GIF or JPG, max 20MB)';
    uploadButton.classList.add('upload-btn');
    imageUpload.appendChild(uploadButton);

    // Create image upload field
    let postImage = document.createElement('input');
    postImage.type = 'file';
    postImage.id = 'image';
    postImage.name = 'image';
    postImage.accept = 'image/gif, image/jpeg, image/png';
    imageUpload.appendChild(postImage);

    // Create upload status section
    let fileName = document.createElement('p');
    fileName.id = 'file_name';
    fileName.textContent = 'No file chosen';
    imageUpload.appendChild(fileName);

    // Create submit button
    let removeButton = document.createElement('button');
    removeButton.id = "remove-image";
    removeButton.classList.add("remove-btn");
    removeButton.style.display = "none";
    removeButton.textContent = "Remove Image";
    imageUpload.appendChild(removeButton);

    // Create categories container
    let categoriesContainer = document.createElement('div');
    categoriesContainer.id = 'categories_container';
    form.appendChild(categoriesContainer);

    // Create categories title
    let categoriesTitle = document.createElement('h2');
    categoriesTitle.textContent = 'Select Categories';
    categoriesTitle.id = 'categories_title';
    categoriesContainer.appendChild(categoriesTitle);

    // Create left cluster of categories
    let leftCluster = document.createElement('div');
    leftCluster.id = 'left_cluster';
    categoriesContainer.appendChild(leftCluster);

    // Create right cluster of categories
    let rightCluster = document.createElement('div');
    rightCluster.id = 'right_cluster';
    categoriesContainer.appendChild(rightCluster);

    // Create submit button
    let submitButton = document.createElement('button');
    submitButton.type = 'submit';
    submitButton.textContent = 'Create Post';
    form.appendChild(submitButton);

    // Link js
    let scriptTag = document.createElement('script');
    scriptTag.src = '/static/js/onboarding.js';
    scriptTag.defer = true;
    app.appendChild(scriptTag);
}