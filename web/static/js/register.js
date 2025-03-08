import { fetchErrorMessage } from '/static/js/form_error_message.js'

export function renderRegistrationPage() {
    document.title = 'sign up - real-time-forum'; // Set document title

    let formContainer = document.CreateElement('div');
    formContainer.classList.add('form-container');

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
    
    fetchErrorMessage(errorMessageContainer);

    // Create registration form
    let registrationForm = document.createElement('form');
    registrationForm.method = 'POST';
    registrationForm.action = '/register';
    registrationForm.enctype = 'multipart/form-data'

    // Create first and lastname fields
    let firstLastName = document.createElement('div');
    firstLastName.classList.add('shared-field');

    let firstName = document.createElement('input');
    firstName.type = 'text';
    firstName.id = 'first-name';
    firstName.name = 'first-name';
    firstName.placeholder = 'first name';
    firstName.required = true;
    firstLastName.appendChild(firstName);

    let lastName = document.createElement('input');
    lastName.type = 'text';
    lastName.id = 'last-name';
    lastName.name = 'last-name';
    lastName.placeholder = 'last name';
    lastName.required = true;
    firstLastName.appendChild(lastName);

    registrationForm.appendChild(firstLastName);

    // Create nickname/username field
    let username = document.createElement('input');
    username.type = 'text';
    username.id = 'username';
    username.name = 'username';
    username.placeholder = 'nickname / username';
    username.required = true;
    registrationForm.appendChild(username);

     // Create age and gender fields
     let ageGender = document.createElement('div');
     ageGender.classList.add('shared-field');
 
     let age = document.createElement('input');
     age.type = 'number';
     age.id = 'age';
     age.name = 'age';
     age.placeholder = 'age';
     age.required = true;
     ageGender.appendChild(age);
 
     let gender = document.createElement('select');
     gender.id = 'gender';
     gender.name = 'gender';
     
     let none = document.createElement('option');
     none.value = 'none';
     none.textContent = 'none';
     gender.appendChild(none);

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
     email.id = 'email';
     email.name = 'email';
     email.placeholder = 'email';
     email.required = 'true';
     registrationForm.appendChild(email);

     // Create password & confirm password fields
     let passwordConfirmPassword = document.createElement('div');
     passwordConfirmPassword.classList.add('shared-field');

     let password = document.createElement('input');
     password.type = 'password';
     password.id = 'password';
     password.name = 'password';
     password.placeholder = 'password';
     password.required = true
     passwordConfirmPassword.appendChild(password);

     let confirmPassword = document.createElement('input');
     confirmPassword.type = 'password';
     confirmPassword.id = 'confirm-password';
     confirmPassword.name = 'confirm-password';
     confirmPassword.placeholder = 'confirm password';
     confirmPassword.required = true
     passwordConfirmPassword.appendChild(confirmPassword);

     registrationForm.appendChild(passwordConfirmPassword);

     // Create bio field
     let bio = document.createElement('textarea');
     bio.id = 'bio';
     bio.name = 'bio';
     bio.placeholder = 'bio | about me';
     bio.required = true
     registrationForm.appendChild(bio);

     // Create image upload container
     let imageUpload = document.createElement('div');
     imageUpload.classList.add('custom-file-upload');

     // Create image upload field
     let image = document.createElement('input');
     image.type = 'file';
     image.id = 'image';
     image.name = 'image';
     image.accept = 'image/*'
     imageUpload.appendChild(image);

     // Create upload status section
     let fileName = document.createElement('p');
     fileName.id = 'file-name';
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

    document.body.appendChild(formContainer);

    // Create script tag and link to js
    let scriptTag = document.createElement('script');
    scriptTag.src = '/static/js/onboarding.js';
    scriptTag.defer = true;
    document.body.appendChild(scriptTag);

    // Attach event listener to handle registration via AJAX
    registrationForm.addEventListener("submit", async function (event) {
        event.preventDefault(); // Prevent full-page reload

        let formData = new FormData(registrationForm);
        let response = await fetch("/register", {
            method: "POST",
            body: formData,
        });

        let data = await response.json().catch(() => null);

        if (response.ok && !data?.error_message) {
            // Redirect to dashboard on success
            navigateTo("/");
        } else {
            // Show error message
            errorText.textContent = data?.error_message || "Registration failed";
        }
    });
}