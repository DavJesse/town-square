import { fetchErrorMessage } from '/static/js/form_error_message.js'
import { navigateTo } from '/static/js/routes.js'
import { setupImageUpload } from '/static/js/onboarding.js';

export function renderRegistrationPage() {
    document.title = 'sign up - real-time-forum'; // Set document title

    let formContainer = document.createElement('div');
    formContainer.classList.add('registration-form-container');

    // create logo
    let heading = document.createElement('h1');
    heading.classList.add('logo');
    heading.textContent = 'real-time-forum';
    formContainer.appendChild(heading);

    // create error message container
    let errorMessageContainer = document.createElement('input');
    errorMessageContainer.id = 'error_message';
    errorMessageContainer.type = 'hidden';
    errorMessageContainer.value = '';

    // Create error message placeholder
    let errorText = document.createElement('p');
    errorText.id = 'error_text';
    errorText.style.color = 'red';
    formContainer.appendChild(errorText);
    formContainer.appendChild(errorMessageContainer);
    
    fetchErrorMessage(errorMessageContainer, '/register');

    // Create registration form
    let registrationForm = document.createElement('form');
    registrationForm.method = 'POST';
    registrationForm.action = '/register';
    registrationForm.enctype = 'multipart/form-data'

    // Create first and lastname fields
    let firstLastName = document.createElement('div');
    firstLastName.id ='shared_input_containers';

    let firstName = document.createElement('input');
    firstName.type = 'text';
    firstName.id = 'register_first_name';
    firstName.name = 'first-name';
    firstName.placeholder = 'first name';
    firstName.autocomplete = 'given-name';
    firstName.autocapitalize = 'words';
    firstName.required = true;
    firstLastName.appendChild(firstName);

    let lastName = document.createElement('input');
    lastName.type = 'text';
    lastName.id = 'register_last_name';
    lastName.name = 'last-name';
    lastName.placeholder = 'last name';
    lastName.autocapitalize = 'words';
    lastName.required = true;
    firstLastName.appendChild(lastName);

    registrationForm.appendChild(firstLastName);

    // Create nickname/username field
    let username = document.createElement('input');
    username.type = 'text';
    username.id = 'register_username';
    username.name = 'username';
    username.placeholder = 'nickname / username';
    username.autocomplete = 'username';
    username.required = true;
    registrationForm.appendChild(username);

     // Create age and gender fields
     let ageGender = document.createElement('div');
     ageGender.id ='shared_input_containers';
 
     let age = document.createElement('input');
     age.type = 'number';
     age.id = 'register_age';
     age.name = 'age';
     age.placeholder = 'age';
     age.min = '13';
     age.max = '120';
     age.required = true;
     ageGender.appendChild(age);

     // Add gender label
     let genderLabel = document.createElement('label');
     genderLabel.htmlFor = 'register_gender';
     genderLabel.id = 'register_gender_label';
     genderLabel.textContent = 'gender:';
     ageGender.appendChild(genderLabel);
 
     let gender = document.createElement('select');
     gender.id = 'register_gender';
     gender.name = 'gender';
     gender.textContent = 'gender';
     
     let other = document.createElement('option');
     other.value = 'other';
     other.textContent = 'other';
     gender.appendChild(other);

     let male = document.createElement('option');
     male.value ='male';
     male.textContent ='male';
     gender.appendChild(male);

     let female = document.createElement('option');
     female.value = 'female';
     female.textContent = 'female';
     gender.appendChild(female);

     ageGender.appendChild(gender); 
     registrationForm.appendChild(ageGender);

     // Create email field
     let email = document.createElement('input');
     email.type = 'email';
     email.id = 'register_email';
     email.name = 'email';
     email.placeholder = 'email';
     email.autocomplete = 'email';
     email.required = true;
     registrationForm.appendChild(email);

     // Create password & confirm password fields
     let passwordConfirmPassword = document.createElement('div');
     passwordConfirmPassword.id = 'shared_input_containers';

     let password = document.createElement('input');
     password.type = 'password';
     password.id = 'register_password';
     password.name = 'password';
     password.placeholder = 'password';
     password.required = true
     passwordConfirmPassword.appendChild(password);

     let confirmPassword = document.createElement('input');
     confirmPassword.type = 'password';
     confirmPassword.id = 'register_confirm_password';
     confirmPassword.name = 'confirm-password';
     confirmPassword.placeholder = 'confirm password';
     confirmPassword.required = true
     passwordConfirmPassword.appendChild(confirmPassword);

     registrationForm.appendChild(passwordConfirmPassword);

     // Create bio field
     let bio = document.createElement('textarea');
     let charCounter = document.createElement('p');
     bio.id = 'register_bio';
     charCounter.id = 'char_counter';
     bio.name = 'bio';
     bio.placeholder = 'bio | about me (300 characters max)';
     bio.required = true
     bio.maxLength = 300;
     charCounter.textContent = '300 characters remaining';
     registrationForm.appendChild(bio);
     registrationForm.appendChild(charCounter);

     bio.addEventListener('input', () => {
        let remaining = 300 - bio.value.length;
        if (remaining < 50) {
            charCounter.style.color = '#f95615';
        }
        charCounter.textContent = `${remaining} characters left`;
     })

     // Create image upload container
     let imageUpload = document.createElement('div');
     imageUpload.classList.add('custom-file-upload');

     // Add image upload label
     let uploadButton = document.createElement('label');
     uploadButton.htmlFor = 'image';
     uploadButton.textContent = 'choose image';
     uploadButton.classList.add('upload-btn');
     imageUpload.appendChild(uploadButton);

     // Create image upload field
     let image = document.createElement('input');
     image.type = 'file';
     image.id = 'image';
     image.name = 'image';
     image.accept = 'image/*'
     imageUpload.appendChild(image);

     // Create upload status section
     let fileName = document.createElement('p');
     fileName.id = 'file_name';
     fileName.textContent = 'no file chosen';
     imageUpload.appendChild(fileName);

     registrationForm.appendChild(imageUpload);

     // Create submit button
     let submitButton = document.createElement('button');
     submitButton.type ='submit';
     submitButton.textContent = 'sign up';
     registrationForm.appendChild(submitButton);

     formContainer.appendChild(registrationForm);

     // Create account issues section
     let accountIssues = document.createElement('div');
     accountIssues.classList.add('account-issues');

     // Create login text and login link
     let loginText = document.createElement('p');
     loginText.textContent = 'Already have an account? ';

     // Create login link
     let loginLink = document.createElement('a');
     loginLink.classList.add('register-link');
     loginLink.textContent = 'Login';
     loginLink.href = '/login';
     loginText.appendChild(loginLink);

     accountIssues.appendChild(loginText)
     formContainer.appendChild(accountIssues);

    // Add login container to #app div
    const app = document.getElementById("app");
    app.appendChild(formContainer);

    // Create script tag and link to js
    let scriptTag = document.createElement('script');
    scriptTag.src = '/static/js/onboarding.js';
    scriptTag.defer = true;
    scriptTag.onload = function() {
        setupImageUpload();
    };
    app.appendChild(scriptTag);

    // Attach event listener to handle registration via AJAX
    registrationForm.addEventListener("submit", async function (event) {
        event.preventDefault(); // Prevent full-page reload

        if (event.target !== registrationForm) return; // Ensure accurate form submision

        let formData = new FormData();
        formData.append("first_name", document.getElementById('register_first_name').value);
        formData.append("last_name", document.getElementById('register_last_name').value);
        formData.append("username", document.getElementById('register_username').value);
        formData.append("age", document.getElementById('register_age').value);
        formData.append("gender", document.getElementById('register_gender').value);
        formData.append("email", document.getElementById('register_email').value);
        formData.append("password", document.getElementById('register_password').value);
        formData.append("confirm_password", document.getElementById('register_confirm_password').value);
        formData.append("bio", document.getElementById('register_bio').value);

        let imageFile = document.getElementById('image').files[0];
        if (imageFile)  {
            formData.append("image", imageFile);
        }
        
        try {
        let response = await fetch("/register", {
            method: "POST",
            body: formData,
        });

        let data = await response.json().catch(() => null);

        if (response.ok && !data?.error_message) {
            // Redirect to login page on success
            navigateTo("/login");
        } else {
            // Show error message
            errorText.textContent = data?.error_message || "Registration failed";
        }
    } catch (error) {
        console.log(`Fetch Error: ${error}`);
        errorText.textContent = "Network Error, please try again";
    }
    });

     // Prevent 'sign up' link from being blocked
     loginLink.addEventListener("click", function(event) {
        event.stopPropagation();
        navigateTo("/login");
    });
}