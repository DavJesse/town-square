async function fetchErrorMessage(errorMessageContainer) {
    try {
        let response =  await fetch('/login');
        if (!response.ok) {
            let data = response.json();
            setErrorMessage(errorMessageContainer, data.error_message);
        }
    } catch (error) {
        console.error(`Error fetcing message: ${error}`);
    }
}

function setErrorMessage(errorMessageContainer, message) {
    errorMessageContainer.value = message;
    if (message) {
        let errorText = document.createElement('p');
        errorText.textContent = message;
        errorText.style.color = 'red';
        errorText.insertAdjacentElement('afterend', errorText);
    }
}

export function renderLoginPage() {
    // update document title
    document.title = 'login - real-time-forum';

    // create login container
    let loginContainer = document.createElement('div');
    loginContainer.classList.add('login-container');
    
    // create logo
    let heading = document.createElement('h1');
    heading.classList.add('logo');
    heading.textContent = 'real-time-forum';
    loginContainer.appendChild(heading);

    // create login text
    let loginText = document.createElement('h1');
    heading.textContent = 'real-time-forum';
    loginContainer.appendChild(loginText);

    // create login sub-container
    let loginSubContainer = document.createElement('div');
    loginSubContainer.classList.add('login');

    // create login tagline
    let loginTagline = document.createElement('h2');
    loginTagline.classList.add('login-tagline');
    loginTagline.textContent = 'Welcome back, sign in to continue';
    loginSubContainer.appendChild(loginTagline);

    // create error message container
    let errorMessageContainer = document.createElement('input');
    errorMessageContainer.classList.add('error-message');
    errorMessageContainer.id = 'error_message';
    errorMessageContainer.type = 'hidden';
    errorMessageContainer.value = '';
    fetchErrorMessage(errorMessageContainer);
    loginSubContainer.appendChild(errorMessageContainer);

    // create login form
    let loginForm = document.createElement('form');
    loginForm.method = 'POST';
    loginForm.action = '/login';

    // create email/username input
    let emailUsername = document.createElement('input');
    emailUsername.type = 'text';
    emailUsername.id = 'email_username';
    emailUsername.placeholder = 'username/email';
    emailUsername.name = 'email_username';
    emailUsername.required = true;
    loginForm.appendChild(emailUsername);

    // create password input
    let password = document.createElement('input');
    password.type = 'text';
    password.id = 'password';
    password.placeholder = 'password';
    password.name = 'password';
    password.required = true;
    loginForm.appendChild(password)
    
    // create submit button
    let loginButton = document.createElement('button');
    loginButton.type ='submit';
    loginButton.textContent = 'login';
    loginForm.appendChild(loginButton)
    
    loginSubContainer.appendChild(loginForm);
    loginContainer.appendChild(loginSubContainer); // add login sub-container

    // create account issues
    let accountIssues = document.createElement('div');
    accountIssues.classList.add('account-issues');

    let registerOption = document.createElement('p');
    registerOption.textContent = 'Don\'t have an account? ';

    let registerLink = document.createElement('a');
    registerLink.classList.add('register-link');
    registerLink.textContent = 'Sign up';
    registerLink.href = '/register';

    registerOption.appendChild(registerLink);
    accountIssues.appendChild(registerOption);
    loginContainer.appendChild(accountIssues);

    document.body.appendChild(loginContainer);
    
}

